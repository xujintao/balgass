
-- protocol
local p_mu = Proto("mu", "mu protocol")

-- field
local f_flag = ProtoField.uint8("flag", "mu.flag", base.HEX)
local f_len = ProtoField.uint8("len", "mu.len", base.DEC)
local f_len2 = ProtoField.uint16("len2", "mu.len2", base.DEC)
local f_code = ProtoField.uint16("code", "mu.code", base.HEX)
local f_data = ProtoField.string("data", "mu.data")
local f_data_enc = ProtoField.string("enc data", "mu.dataenc")

--- bind field with protocol
p_mu.fields = {f_flag, f_len, f_len2, f_code, f_data, f_data_enc}

local data_dis = Dissector.get("data")

--- bind dissect handle
function p_mu.dissector(buf, pkt, tree)

    while(true)
    do
        -- protocol tree
        local subtree = tree:add(p_mu, buf(0,1))

        local len = 0
        local offset = 0
        local flag_table = {
            [0xC1] = function()
                -- len
                len = buf(offset,1):uint()
                subtree:add(f_len, buf(offset,1))
                offset = offset + 1

                -- code
                subtree:add(f_code, buf(offset,2))
                offset = offset + 2

                -- data
                if len > offset then
                    subtree:add(f_data, buf(offset,len-offset))
                    offset = len
                end
            end,
            [0xC2] = function()
                -- len
                len = buf(offset,2):uint()
                subtree:add(f_len2, buf(offset,2))
                offset = offset + 2

                -- code
                subtree:add(f_code, buf(offset,2))
                offset = offset + 2

                -- data
                if len > offset then
                    subtree:add(f_data, buf(offset,len-offset))
                    offset = len
                end
            end,
            [0xC3] = function()
                -- len
                len = buf(offset,1):uint()
                subtree:add(f_len, buf(offset,1))
                offset = offset + 1

                -- data
                if len > offset then
                    subtree:add(f_data_enc, buf(offset,len-offset))
                    offset = len
                end
            end,
            [0xC4] = function()
                -- len
                len = buf(offset,2):uint()
                subtree:add(f_len2, buf(offset,2))
                offset = offset + 2

                -- data
                if len > offset then
                    subtree:add(f_data_enc, buf(offset,len-offset))
                    offset = len
                end
            end,
        }

        local flag = buf(offset,1):uint()
        local dis = flag_table[flag]
        if (dis == nil) then
            data_dis:call(buf, pkt, tree)
            return
        end

        -- flag
        subtree:add(f_flag, buf(offset,1))
        offset = offset + 1

        -- other field
        dis()

        if (buf:len() == len) then
            return
        end

        -- next
        buf = buf(len):tvb()
    end
end

local tcp_encap_table = DissectorTable.get("tcp.port")
tcp_encap_table:add(44405, p_mu)
tcp_encap_table:add(56900, p_mu)