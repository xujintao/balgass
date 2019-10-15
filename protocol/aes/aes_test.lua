local string = require("string")
local array = require("lockbox.util.array")
local aes = require("aes")

local tests = {
    {
        plaintext = "68656C6C6F",
        ciphertext = "7A54FD3485976004793350E684EC7ED20B",
    },
    {
        plaintext = "6578616D706C65706C61696E74657874",
        ciphertext = "3633C13148E3E52FF91C6296DB22C5FC00",
    }
}
for k,v in pairs(tests) do
    local plaintext = array.fromHex(v.plaintext)
    local dstcipher = aes.encrypt(plaintext)
    assert(v.ciphertext == array.toHex(dstcipher),
        string.format("encrypt failed! expected(%s) got(%s)", v.ciphertext, array.toHex(dstcipher)))
    
    local ciphertext = array.fromHex(v.ciphertext)
    local dstplain = aes.decrypt(ciphertext)
    assert(v.plaintext == array.toHex(dstplain), 
        string.format("decrypt failed! expected(%s) got(%s)", v.plaintext, array.toHex(dstplain)))
end

