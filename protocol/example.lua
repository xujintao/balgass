local aes = require("aes")

local plaintext = {0x12, 0x34}
local ciphertext = aes.encrypt(plaintext)
for k,v in pairs(ciphertext) do
    print(v)
end
