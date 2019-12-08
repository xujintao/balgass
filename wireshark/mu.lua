local path_common = "D:\\github.com\\xujintao\\balgass\\wireshark\\"
package.path = path_common .. "aes\\?.lua;"
            .. path_common .. "aes\\?\\init.lua;"
            .. path_common .. "xor\\?.lua;"
            .. path_common .. "xor\\?\\init.lua;"
            .. package.path

local aes = require("aes")
local xor = require("xor")
local array = require("util.array")

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
function do_dissector(tvb, pinfo, tree)
    -- protocol tree
    local subtree
    local flag = tvb(0,1):uint()
    local len = 0
    local code = 0
    local data = "null"
    local offset = 0
    local text = ""
    local disC1C2 = function()
        if pinfo.dst_port==56900 then -- only packets send to 56900 need xor
            -- xor.dec
            if flag==0xC1 then
                len = tvb(1,1):uint()
                local xortext = array.fromHex(tvb(0, len):bytes():tohex())
                xor.dec(xortext, 4, #xortext)
                tvb = ByteArray.new(array.toHex(xortext)):tvb("tvb_xor")
            else
                len = tvb(1,2):uint()       
                local xortext = array.fromHex(tvb(0, len):bytes():tohex())
                xor.dec(xortext, 5, #xortext)
                tvb = ByteArray.new(array.toHex(xortext)):tvb("tvb_xor")        
            end
        else
            if flag==0xC1 then
                len = tvb(1,1):uint()
            else
                len = tvb(1,2):uint()
            end
        end
        
        -- flag
        subtree:add(f_flag, tvb(offset,1))
        offset = offset + 1
        
        -- len
        if flag==0xC1 then          
            subtree:add(f_len, tvb(offset,1))
            offset = offset + 1
        else
            subtree:add(f_len2, tvb(offset,2))
            offset = offset + 2         
        end
        
        -- code
        code = tvb(offset,2):uint()
        subtree:add(f_code, tvb(offset,2))
        offset = offset + 2
            
        -- data 
        if offset<len then
            subtree:add(f_data, tvb(offset,len-offset))
            if len-offset > 10 then
                data = tvb(offset,10):bytes():tohex(true, " ") .. " ..."
            else
                data = tvb(offset,len-offset):bytes():tohex(true, " ")
            end
            offset = len
        end
        
        -- prepare subtree text
        if flag==0xC1 then
            text = text .. string.format("<flag:%02x + len:%02x + code:%04x + data:%s>", 0xC1, len, code, data)
        else
            text = text .. string.format("<flag:%02x + len:%04x + code:%04x + data:%s>", 0xC2, len, code, data)
        end
    end
    
    local flag_table = {
        [0xC1] = function()
            disC1C2()
        end,
        [0xC2] = function()
            disC1C2()
        end,
        [0xC3] = function()
            -- aes.decrypt
            local lenC3 = tvb(1,1):uint()
            text = string.format("<flag:%x + len:%02x>", 0xC3, lenC3)
            local ciphertext = array.fromHex(tvb(2, lenC3-2):bytes():tohex())
            local plaintext,err = aes.decrypt(ciphertext)
            if string.len(err)~=0 then
                subtree:set_text(text) -- save field
                error(string.format("decrypt failed! %s", err))
            end
            table.insert(plaintext, 1, #plaintext + 2)
            table.insert(plaintext, 1, 0xC1)
            tvb = ByteArray.new(array.toHex(plaintext)):tvb("tvb_aes")
            flag = tvb(0,1):uint()
            disC1C2()
        end,
        [0xC4] = function()
            -- aes.decrypt
            local lenC4 = tvb(1,2):uint()
            text = string.format("<flag:%x + len:%04x>", 0xC4, lenC4)
            local ciphertext = array.fromHex(tvb(3, lenC4-3):bytes():tohex())
            local plaintext,err = aes.decrypt(ciphertext)
            if string.len(err)~=0 then
                subtree:set_text(text) -- save field
                error(string.format("decrypt failed! %s", err))
            end
            table.insert(plaintext, 1, (#plaintext + 3)%256)
            table.insert(plaintext, 1, math.floor((#plaintext + 3)/256))
            table.insert(plaintext, 1, 0xC2)
            tvb = ByteArray.new(array.toHex(plaintext)):tvb("tvb_aes")
            flag = tvb(0,1):uint()
            disC1C2()
        end,
    }

    local dis = flag_table[flag]
    if (dis == nil) then
        data_dis:call(tvb, pinfo, tree)
        return
    end
    
    -- prepare subtree
    subtree = tree:add(p_mu, tvb(0,tvb:len()))
    
    -- dissect
    dis()

    -- set subtree text
    subtree:set_text(text)
end

function get_message_len(tvb, pinfo, offset)
    local flag = tvb(offset,1):uint()
    if flag==0xC1 or flag==0xC3 then
        return tvb(offset+1,1):uint()
    elseif flag==0xC2 or flag==0xC4 then
        return tvb(offset+1,2):uint()
    else
        return 1
    end
end

--- bind dissect handle
function p_mu.dissector(tvb, pinfo, tree)
    -- Reassemble tcp fragment
    dissect_tcp_pdus(tvb, tree, 3, get_message_len, do_dissector)
end

local tcp_encap_table = DissectorTable.get("tcp.port")
tcp_encap_table:add(44405, p_mu)
tcp_encap_table:add(56900, p_mu)
tcp_encap_table:add(56970, p_mu)
tcp_encap_table:add(56960, p_mu)
tcp_encap_table:add(56906, p_mu)