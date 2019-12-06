package handle

import (
	"log"
	"math/rand"

	"github.com/xujintao/balgass/cmd/server_data/model"
	"github.com/xujintao/balgass/network"
)

const (
	chanNum  = 3
	chanSize = 1000
	gNum     = 10
)

type data struct {
	v   interface{}
	req *network.Request
}

type handleData struct {
	handleBase
	queue []chan data
	exit  chan struct{}
}

// Handle *CMDHandle implements network.Handler
func (h *handleData) Handle(v interface{}, req *network.Request) {

	if h.queue == nil {
		h.exit = make(chan struct{}, 1)
		h.queue = make([]chan data, chanNum)
		for i := range h.queue {
			h.queue[i] = make(chan data, chanSize)
		}
		for i := 0; i < gNum; i++ {
			go func() {
				for {
					select {
					case d := <-h.queue[i%chanNum]:
						// id := v.(string)
						// user := service.ServerManager.GetServer(id)
						_ = d.req
					// code := req.Code
					// if h, ok := h.cmds[int(code)]; ok {
					// 	h(user, req)
					// 	return
					// }
					// subcode := req.Body[0]
					// codes := []byte{code, subcode}
					// code16 := binary.BigEndian.Uint16(codes)
					// if h, ok := h.cmds[int(code16)]; ok {
					// 	req.Body = req.Body[1:]
					// 	h(user, req)
					// 	return
					// }
					// log.Printf("[%s], invalid cmd, code:[%02x], body:[%v]", uid, code, network.Hex2string(req.Body))
					case <-h.exit:
						log.Println("goroutine exit")
					}

				}

			}()
		}
	}
	h.queue[rand.Int()%chanNum] <- data{v, req}
}

func (h *handleData) Exit() {
	h.exit <- struct{}{}
}

var cmdsData = map[int]func(server *model.Server, req *network.Request){
	// 0x05:   checkVersion,
	// 0xf406: getServerList,
	// 0xf403: getServerInfo,
}
