package handle

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"text/template"

	"github.com/gorilla/websocket"
	"github.com/xujintao/balgass/src/server_game/game"
	"github.com/xujintao/balgass/src/server_game/game/model"
)

func init() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/api/game", handleGame)
	wsHandleDefault.init()
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/api/game")
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Document</title>
    <script>
      window.addEventListener("load", () => {
        var input = document.getElementById("input");
        var canvas = document.getElementById("canvas");
        var ctx = canvas.getContext("2d");
        ctx.transform(1, 0, 0, -1, 0, canvas.height);
        ctx.strokeRect(0, 0, 512, 512);
        var bg;
        var drawPoint = (x, y) => {
          ctx.fillRect(x * 2, y * 2, 2, 2);
        };
        var print = function (message) {
          var d = document.createElement("div");
          d.textContent = message;
          output.appendChild(d);
          output.scroll(0, output.scrollHeight);
        };
        ws = new WebSocket("{{.}}");
        ws.onopen = () => {
          print("websocket open");
          ws.send(
            JSON.stringify({
              action: "SubscribeMap",
              in: { number: 0 },
            })
          );
        };
        ws.onclose = () => {
          print("websocket close");
          ws = null;
        };
        ws.onmessage = (e) => {
          var resp = JSON.parse(e.data);
          var out = resp.out;
          if (out.err) {
            print(out.err);
            return;
          }
          switch (resp.action) {
            case "HandleErrorReply":
              break;
            case "SubscribeMapReply":
              if (out.name == "map") {
                bg = out.data;
              }
              if (out.name == "object") {
                ctx.clearRect(0, 0, 512, 512);
                ctx.fillStyle = "black";
                bg.map((p) => {
                  drawPoint(p.x, p.y);
                });
                ctx.fillStyle = "red";
                out.data.map((p) => {
                  drawPoint(p.x, p.y);
                });
              }
              break;
            default:
              print("unknown [action]" + resp.action);
          }
        };
        ws.onerror = (e) => {
          print("err:" + e.data);
        };
        document.getElementById("send").onclick = (e) => {
          if (!ws) {
            return false;
          }
          const number = parseInt(input.value, 10);
          ws.send(
            JSON.stringify({
              action: "SubscribeMap",
              in: {
                number,
              },
            })
          );
          return false;
        };
      });
    </script>
  </head>

  <body>
    <canvas id="canvas" width="512" height="512"></canvas>
    <table>
      <tr>
        <td valign="top" width="50%">
          <input id="input" type="text" value="0" />
          <button id="send">Send</button>
        </td>
        <td valign="top" width="50%">
          <div id="output" style="max-height: 20vh; overflow-y: scroll"></div>
        </td>
      </tr>
    </table>
  </body>
</html>
`))

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type wsconn struct {
	*websocket.Conn
	addr string
}

func (c *wsconn) Addr() string {
	return c.addr
}

func (c *wsconn) Write(msg any) error {
	reply, err := wsHandleDefault.addAction(msg)
	if err != nil {
		log.Printf("addAction failed [err]%v\n", err)
		return err
	}
	err = c.WriteJSON(reply)
	if err != nil {
		log.Printf("WriteJSON failed [err]%v\n", err)
		return err
	}
	return nil
}

func handleGame(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	addr := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		addr = realIP
	}
	conn := wsconn{Conn: c, addr: addr}
	writeErr := func(s string) {
		msg := model.MsgHandleErrorReply{
			Err: s,
		}
		conn.Write(&msg)
	}
	id, err := game.Game.UserConn(&conn)
	if err != nil {
		writeErr(err.Error())
		c.Close()
		return
	}
	defer game.Game.UserCloseConn(id)
	for {
		var req map[string]any
		err := c.ReadJSON(&req)
		if err != nil {
			log.Printf("ReadJSON failed [addr]%s [err]%v\n", c.RemoteAddr().String(), err)
			return
		}
		actionField, ok := req["action"]
		if !ok {
			s := "websocket request has no action field"
			log.Printf("%s [addr]%s\n", s, c.RemoteAddr().String())
			writeErr(s)
			continue
		}
		action, ok := actionField.(string)
		if !ok {
			s := "websocket request action field is not a string"
			log.Printf("%s [addr]%s\n", s, c.RemoteAddr().String())
			writeErr(s)
			continue
		}
		inField, ok := req["in"]
		if !ok {
			s := fmt.Sprintf("websocket request [action]%s has no in field", action)
			log.Printf("%s [addr]%s\n", s, c.RemoteAddr().String())
			writeErr(s)
			continue
		}
		data, err := json.Marshal(inField)
		if err != nil {
			s := fmt.Sprintf("websocket request [action]%s in field is not a valid json", action)
			log.Printf("%s [addr]%s [err]%v\n", s, c.RemoteAddr().String(), err)
			writeErr(s)
			continue
		}
		wsHandleDefault.Handle(id, action, data)
	}
}

var wsHandleDefault wsHandle

type wsHandle struct {
	wsIns  map[string]*wsIn
	wsOuts map[any]*wsOut
}

func (h *wsHandle) init() {
	// ingress
	h.wsIns = make(map[string]*wsIn)
	for _, v := range wsIns {
		h.wsIns[v.action] = v
	}
	// egress
	h.wsOuts = make(map[any]*wsOut)
	for _, v := range wsOuts {
		t := reflect.TypeOf(v.msg)
		if t.Kind() != reflect.Ptr {
			log.Printf("wsOuts [action]%s msg field must be a pointer\n", v.action)
		}
		h.wsOuts[t] = v
	}
}

func (h *wsHandle) Handle(id int, action string, data []byte) {
	var api *wsIn
	var ok bool
	if api, ok = h.wsIns[action]; !ok {
		log.Printf("invalid websocket request [action]%s\n", action)
		return
	}
	msg := reflect.New(reflect.TypeOf(api.msg).Elem()).Interface()
	if err := json.Unmarshal(data, msg); err != nil {
		log.Printf("json.Unmarshal failed [data]%v [msg]%v\n", data, msg)
		return
	}
	game.Game.UserAction(id, api.action, msg)
}

func (h *wsHandle) addAction(msg any) (any, error) {
	v := reflect.ValueOf(msg)
	t := v.Type()
	api, ok := h.wsOuts[t]
	if !ok {
		err := fmt.Errorf("%s has not yet be registered to wsOuts table", t.String())
		return nil, err
	}
	reply := map[string]any{
		"action": api.action,
		"out":    msg,
	}
	return &reply, nil
}

type wsIn struct {
	action string
	msg    any
}
type wsOut struct {
	action string
	msg    any
}

var wsIns = [...]*wsIn{
	{"SubscribeMap", (*model.MsgSubscribeMap)(nil)},
	// {"SubscribePlayer", nil},
}

var wsOuts = [...]*wsOut{
	{"HandleErrorReply", (*model.MsgHandleErrorReply)(nil)},
	{"SubscribeMapReply", (*model.MsgSubscribeMapReply)(nil)},
}
