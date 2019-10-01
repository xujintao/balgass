
-- protocol
local p_mu = Proto("mu", "mu protocol")

-- field
local f_flag = ProtoField.uint8("mu.flag", "flag", base.HEX)
local f_len = ProtoField.uint8("mu.len", "len", base.DEC)
local f_len2 = ProtoField.uint16("mu.len2", "len2", base.DEC)
local f_code = ProtoField.uint16("mu.code", "code", base.HEX)
local f_data = ProtoField.bytes("mu.data", "data")
local f_data_enc = ProtoField.bytes("mu.dataenc", "dataenc")

--- bind field with protocol
p_mu.fields = {f_flag, f_len, f_len2, f_code, f_data, f_data_enc}

local data_dis = Dissector.get("data")

--- bind dissect handle
function do_dissector(buf, pkt, tree)

    -- protocol tree
    local subtree = tree:add(p_mu, buf(0,1))

    local len = 0
    local offset = 0
    local code = 0
    local flag_table = {
        [0xC1] = function()
            -- len
            len = buf(offset,1):uint()
            subtree:add(f_len, buf(offset,1))
            offset = offset + 1

            -- code
            code = buf(offset,2):uint()
            subtree:add(f_code, buf(offset,2))
            offset = offset + 2

            -- data 
            subtree:add(f_data, buf(offset,len-offset))
            offset = len            
        end,
        [0xC2] = function()
            -- len
            len = buf(offset,2):uint()
            subtree:add(f_len2, buf(offset,2))
            offset = offset + 2

            -- code
            code = buf(offset,2):uint()
            subtree:add(f_code, buf(offset,2))
            offset = offset + 2

            -- data
            subtree:add(f_data, buf(offset,len-offset))
            offset = len
        end,
        [0xC3] = function()
            -- len
            len = buf(offset,1):uint()
            subtree:add(f_len, buf(offset,1))
            offset = offset + 1

            -- data
            subtree:add(f_data_enc, buf(offset,len-offset))
            offset = len
        end,
        [0xC4] = function()
            -- len
            len = buf(offset,2):uint()
            subtree:add(f_len2, buf(offset,2))
            offset = offset + 2

            -- data
            subtree:add(f_data_enc, buf(offset,len-offset))
            offset = len
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

    -- set subtree the real length
    subtree:set_len(len)
    local text = string.format("<flag:0x%x + len:%04d + code:0x%04x>", flag, len, code)
    subtree:set_text(text)
end

function get_message_len(buf, pinfo, offset)
    local flag = buf(offset,1):uint()
    if flag==0xC1 or flag==0xC3 then
        return buf(offset+1,1):uint()
    elseif flag==0xC2 or flag==0xC4 then
        return buf(offset+1,2):uint()
    else
        return 1
    end
end

--- bind dissect handle
function p_mu.dissector(buf, pinfo, tree)
    -- Reassemble tcp fragment
    dissect_tcp_pdus(buf, tree, 3, get_message_len, do_dissector)
end

local tcp_encap_table = DissectorTable.get("tcp.port")
tcp_encap_table:add(44405, p_mu)
tcp_encap_table:add(56900, p_mu)