local string = require("string")
local array = require("lockbox.util.array")
local aes = require("aes")

local tests = {
    {
        plaintext = {104, 101, 108, 108, 111},
        ciphertext = {
            122, 84, 253, 52, 133, 151, 96, 4, 121, 51, 80, 230, 132, 236, 126, 210,
            11 -- fill 11
        }
    },
    {
        plaintext = {101, 120, 97, 109, 112, 108, 101, 112, 108, 97, 105, 110, 116, 101, 120, 116},
        ciphertext = {
            54, 51, 193, 49, 72, 227, 229, 47, 249, 28, 98, 150, 219, 34, 197, 252,
            0 -- fill 0
        }
    }
}

for k,v in pairs(tests) do
    local dstcipher = aes.encrypt(v.plaintext)
    assert(array.toHex(v.ciphertext) == array.toHex(dstcipher),
        string.format("encrypt failed! expected(%s) got(%s)", array.toHex(v.ciphertext), array.toHex(dstcipher)))
    
    local dstplain = aes.decrypt(v.ciphertext)
    assert(array.toHex(v.plaintext) == array.toHex(dstplain), 
        string.format("decrypt failed! expected(%s) got(%s)", array.toHex(v.plaintext), array.toHex(dstplain)))
end

