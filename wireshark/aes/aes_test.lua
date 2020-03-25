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
    },
    { -- 0e
        plaintext = "" ..
        "0E372F36F3509AA968A4C29C9355496A" ..
        "51554FB04E",
        ciphertext = "" ..
        "01E8D4F82A281F279A1479C148B007BA" ..
        "87D74286A7310CD98B243A0BBF6AF1F1" ..
        "0B",
    },
    { -- f303
        plaintext = "" ..
        "5CF30333E1330000000003777A8E6100" ..
        "00000378DF3E5203001807E200680046" ..
        "01C9059A076B02FD02D539D539510125" ..
        "0300009CD67A0003007F007F00000000" ..
        "007F000000000000000000C90500009A" ..
        "0700006B020000FD020000",
        ciphertext = "" ..
        "43EA94326417E574E67553761EE3677A" ..
        "2CE4846B7D7D6694F11F8FD8A2410E09" ..
        "19548FFC87372F934CB39A868471046E" ..
        "B2D9720F6ED1C8BEFBC1CCC64F22C511" ..
        "31DCD20F483880023F1D05F98049D413" ..
        "693F522679BE2130A6CF97197940A241" ..
        "05",
    }
}
for k,v in pairs(tests) do
    local plaintext = array.fromHex(v.plaintext)
    local dstcipher, err = aes.encrypt(plaintext)
    if string.len(err)~=0 then
        error(string.format("encrypt failed! %s", err))
    end
    assert(v.ciphertext == array.toHex(dstcipher),
        string.format("encrypt failed! expected(%s) got(%s)", v.ciphertext, array.toHex(dstcipher)))
    
    local ciphertext = array.fromHex(v.ciphertext)
    local dstplain,err = aes.decrypt(ciphertext)
    if string.len(err)~=0 then
        error(string.format("decrypt failed! %s", err))
    end
    assert(v.plaintext == array.toHex(dstplain), 
        string.format("decrypt failed! expected(%s) got(%s)", v.plaintext, array.toHex(dstplain)))
end

