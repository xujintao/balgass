local array = require("util.array")
local xor = require("xor")

local tests = {
    {
        plaintext  = "C104F40C",
        ciphertext = "C104F406",
    },
    {
        plaintext  = "C10400FF",
        ciphertext = "C1040001",
    },
    {
        plaintext =  "C110FA044B15D2A7D92493A24D752580",
        ciphertext = "C110FA00536572766572204E65777300",
    },
    {
        plaintext  = "C109FA118BFF558BEC",
        ciphertext = "C109FA15865ACAE2C4",
    }
}

for k,v in pairs(tests) do
    local ciphertext = array.fromHex(v.plaintext)
    xor.enc(ciphertext, 4, #ciphertext)
    assert(v.ciphertext == array.toHex(ciphertext),
    string.format("xor failed! expected(%s) got(%s)", v.ciphertext, array.toHex(ciphertext)))

    local plaintext = array.fromHex(v.ciphertext)
    xor.dec(plaintext, 4, #plaintext)
    assert(v.plaintext == array.toHex(plaintext),
    string.format("xor failed! expected(%s) got(%s)", v.plaintext, array.toHex(plaintext)))
end