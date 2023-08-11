package handle

import (
	"context"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xujintao/balgass/src/server_game/game"
	"github.com/xujintao/balgass/src/server_game/game/maps"
	"github.com/xujintao/balgass/src/server_game/game/model"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Map(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	ctx, cancel := context.WithCancel(context.Background())
	t1s := time.NewTicker(time.Second)
	t1s.Stop()
	numberChan := make(chan int)
	go func() {
		defer cancel()
		for {
			var req model.MsgGetObjectsByMapNumber
			err := c.ReadJSON(&req)
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNoStatusReceived) {
					log.Printf("[addr]%s closed.\n", c.RemoteAddr().String())
				} else {
					log.Printf("ReadJSON failed [err]%v\n", err)
				}
				return
			}
			number := req.Number
			if number < 0 || number >= maps.MAX_MAP_NUMBER {
				resp := model.MsgGetObjectsByMapNumberReply{Err: "invalid map number"}
				err = c.WriteJSON(&resp)
				if err != nil {
					log.Printf("WriteJSON failed [err]%v\n", err)
					return
				}
				continue
			}
			t1s.Stop()
			data := maps.MapManager.GetMapPots(number)
			resp := model.MsgGetObjectsByMapNumberReply{Name: "map", Data: data}
			err = c.WriteJSON(&resp)
			if err != nil {
				log.Printf("WriteJSON failed [err]%v\n", err)
				return
			}
			numberChan <- number
		}
	}()

	number := 0
	for {
		select {
		case number = <-numberChan:
			t1s.Reset(time.Second)
		case <-t1s.C:
			var req model.MsgGetObjectsByMapNumber
			req.Number = number
			data, err := game.Game.Command("GetObjectsByMapNumber", &req)
			if err != nil {
				log.Printf("Command failed [err]%v\n", err)
				return
			}
			err = c.WriteJSON(data)
			if err != nil {
				log.Printf("WriteJSON failed [err]%v\n", err)
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/api/map")
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
    <script>
        window.addEventListener("load", () => {
            var input = document.getElementById("input")
            var canvas = document.getElementById("canvas")
            var ctx = canvas.getContext("2d")
            ctx.transform(1, 0, 0, -1, 0, canvas.height)
            ctx.strokeRect(0, 0, 512, 512)
            var bg
            var drawPoint = (x, y) => {
                ctx.fillRect(x * 2, y * 2, 2, 2)
            }
            var print = function (message) {
                var d = document.createElement("div")
                d.textContent = message
                output.appendChild(d)
                output.scroll(0, output.scrollHeight)
            }
            ws = new WebSocket("{{.}}")
            ws.onopen = () => {
                print("websocket open")
                ws.send(JSON.stringify({ number: 0 }))
            }
            ws.onclose = () => {
                print("websocket close")
                ws = null;
            }
            ws.onmessage = (e) => {
                var resp = JSON.parse(e.data)
                if (resp.err) {
                    print(resp.err)
                    return
                }
                if (resp.name == "map") {
                    bg = resp.data
                }
                if (resp.name == "object") {
                    ctx.clearRect(0, 0, 512, 512)
                    ctx.fillStyle = "black"
                    bg.map((p) => {
                        drawPoint(p.x, p.y)
                    })
                    ctx.fillStyle = "red"
                    resp.data.map((p) => {
                        drawPoint(p.x, p.y)
                    })
                }
            }
            ws.onerror = (e) => {
                print("err:" + e.data)
            }
            document.getElementById("send").onclick = (e) => {
                if (!ws) {
                    return false
                }
                const number = parseInt(input.value, 10)
                ws.send(JSON.stringify({ number }))
                return false
            }
        })
    </script>
</head>

<body>
    <canvas id="canvas" width="512" height="512"></canvas>
    <table>
        <tr>
            <td valign="top" width="50%">
                <input id="input" type="text" value="0">
                <button id="send">Send</button>
            </td>
            <td valign="top" width="50%">
                <div id="output" style="max-height: 20vh;overflow-y: scroll;"></div>
            </td>
        </tr>
    </table>
</body>

</html>
`))
