local string = require("string")
local lockbox = require("lockbox")
lockbox.ALLOW_INSECURE = true
local array = require("lockbox.util.array")
local stream = require("lockbox.util.stream")
local ecbMode = require("lockbox.cipher.mode.ecb")
local aes256 = require("lockbox.cipher.aes256")
local zeroPadding = require("lockbox.padding.zero")


local iv = array.fromHex("")

local aes = {}

aes.encrypter = function()
    local public = {}
    local cipher = ecbMode.Cipher()

    public.setKey = function(key)
        cipher.setKey(key)
        return public
    end
    
    public.setBlockCipher = function(blockCipher)
        cipher.setBlockCipher(blockCipher)
        return public
    end
    
    public.setPadding = function(paddingMode)
        cipher.setPadding(paddingMode)
        return public
    end

    public.encrypt = function(src)
        if #src == 0 then
            return nil, "src len is zero"
        end
        local dst = cipher.init()
                        .update(stream.fromArray(iv))
                        .update(stream.fromArray(src))
                        .finish()
                        .asBytes()
        local padsize = 0
        if rawlen(src)%aes256.blockSize~=0 then
            padsize = aes256.blockSize - rawlen(src)%aes256.blockSize
        end
        table.insert(dst, padsize)
        return dst, ""
    end

    return public
end

aes.decrypter = function()
    local public = {}
    local cipher = ecbMode.Decipher()

    public.setKey = function(key)
        cipher.setKey(key)
        return public
    end
    
    public.setBlockCipher = function(blockCipher)
        cipher.setBlockCipher(blockCipher)
        return public
    end
    
    public.setPadding = function(paddingMode)
        cipher.setPadding(paddingMode)
        return public
    end

    public.decrypt = function(src)
        if (#src <= aes256.blockSize) or (#src % aes256.blockSize ~= 1) then
            return nil, "src len invalid"
        end
        local padsize = table.remove(src)
        local dst = cipher.init()
                        .update(stream.fromArray(iv))
                        .update(stream.fromArray(src))
                        .finish()
                        .asBytes()
        while padsize>0 do
            table.remove(dst)
            padsize = padsize - 1
        end
        return dst, ""
    end

    return public
end

-- 默认key
local key = {
    0x7A, 0x2C, 0x74, 0x6D, 0xB5, 0x4F, 0xF7, 0xAF, 0x4A, 0x18, 0x8D, 0x94, 0x7A, 0xE4, 0x71, 0x01,
	0x44, 0x19, 0xE6, 0x83, 0x68, 0x46, 0x86, 0xDB, 0xBE, 0x6D, 0xD9, 0x9C, 0x8C, 0x3C, 0x08, 0x40,
}

--默认加密器
local default_encrypter = aes.encrypter()
                            .setKey(key)
                            .setBlockCipher(aes256)
                            .setPadding(zeroPadding)

--默认解密器
local default_decrypter = aes.decrypter()
                            .setKey(key)
                            .setBlockCipher(aes256)
                            .setPadding(zeroPadding)

-- 默认加密
aes.encrypt = function(src)
    return default_encrypter.encrypt(src)
end

-- 默认解密
aes.decrypt = function(src)
    return default_decrypter.decrypt(src)
end

return aes