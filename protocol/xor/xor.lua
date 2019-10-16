local bit32 = require("bit32")

local xor = {}

xor.enc = function(buf, begin, foot)
    local key = {
        0xAB, 0x11, 0xCD, 0xFE,
        0x18, 0x23, 0xC5, 0xA3,
        0xCA, 0x33, 0xC1, 0xCC,
        0x66, 0x67, 0x21, 0xF3,
        0x32, 0x12, 0x15, 0x35,
        0x29, 0xFF, 0xFE, 0x1D,
        0x44, 0xEF, 0xCD, 0x41,
        0x26, 0x3C, 0x4E, 0x4D,
    }
    for i=begin,foot,1 do
        buf[i] = bit32.bxor(buf[i], buf[i-1], key[i%32])
    end
end

xor.dec = function(buf, begin, foot)
    local key = {
        0xAB, 0x11, 0xCD, 0xFE,
        0x18, 0x23, 0xC5, 0xA3,
        0xCA, 0x33, 0xC1, 0xCC,
        0x66, 0x67, 0x21, 0xF3,
        0x32, 0x12, 0x15, 0x35,
        0x29, 0xFF, 0xFE, 0x1D,
        0x44, 0xEF, 0xCD, 0x41,
        0x26, 0x3C, 0x4E, 0x4D,
    }
    for i=foot,begin,-1 do
        buf[i] = bit32.bxor(buf[i], buf[i-1], key[i%32])
    end
end
return xor