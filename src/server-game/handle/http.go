package handle

import (
	"encoding/json"
	"log/slog"
	"reflect"
	"strconv"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game"
	"github.com/xujintao/balgass/src/server-game/game/model"
)

func init() {
	HTTPHandle.init()
}

var HTTPHandle httpHandle

type httpHandle struct {
	*gin.Engine
	validate *validator.Validate
	commands map[string]*command
}

func (h *httpHandle) init() {
	gin.SetMode(gin.ReleaseMode)
	if conf.ServerEnv.Debug {
		gin.SetMode(gin.DebugMode)
	}
	h.Engine = gin.Default()
	h.validate = validator.New()
	h.commands = make(map[string]*command)
	for _, cmd := range commands {
		h.commands[cmd.Action] = cmd
	}
	h.GET("/", h.handleHome)
	h.POST("/api/accounts", h.CreateAccount, h.handleErr)
	h.GET("/api/accounts", h.GetAccountList, h.handleErr)
	h.DELETE("/api/accounts/:id", h.DeleteAccount, h.handleErr)
	h.POST("/api/bots", h.AddBot, h.handleErr)
	h.DELETE("/api/bots", h.DeleteBot, h.handleErr)
	h.GET("/api/game", h.handleGame)
	h.POST("/api/command", h.handleCommand, h.handleErr)
}

func (h *httpHandle) setErr(c *gin.Context, service int, err error) {
	defer c.Next()
	err = MakeError(service, err)
	c.Set("err", err)
	c.Set("path", c.FullPath())
}

func (h *httpHandle) handleErr(c *gin.Context) {
	defer c.Next()
	err, ok := c.Get("err")
	if ok {
		ce := err.(*ConfigError)
		path, _ := c.Get("path")
		slog.Error("http", "path", path.(string), "description", ce.Description, "err", ce.Error())
		c.JSON(ce.Code, gin.H{
			"message": ce.Description,
		})
	}
}

func (h *httpHandle) handleHome(c *gin.Context) {
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
			var drawPoint = (x, y, size) => {
			  ctx.fillRect(x * 2, y * 2, size, size);
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
					bg?.map((p) => {
					  drawPoint(p.x, p.y, 2);
					});
					ctx.fillStyle = "red";
					out.data.monsters?.map((p) => {
					  drawPoint(p.x, p.y, 4);
					});
					ctx.fillStyle = "green";
					out.data.npcs?.map((p) => {
					  drawPoint(p.x, p.y, 4);
					});
					ctx.fillStyle = "blue";
					out.data.players?.map((p) => {
					  drawPoint(p.x, p.y, 4);
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
	homeTemplate.Execute(c.Writer, "ws://"+c.Request.Host+"/api/game")
}

func (h *httpHandle) CreateAccount(c *gin.Context) {
	in := model.Account{}

	// bind
	if err := c.ShouldBind(&in); err != nil {
		h.setErr(c, CreateAccountBind, err)
		return
	}

	// validate
	if err := h.validate.Struct(&in); err != nil {
		h.setErr(c, CreateAccountValidate, err)
		return
	}

	// command
	if _, err := game.Game.Command("CreateAccount", &in); err != nil {
		h.setErr(c, CreateAccountDB, err)
		return
	}

	c.JSON(200, in)
}

func (h *httpHandle) GetAccountList(c *gin.Context) {
	// get param
	email := c.Query("user_email")

	// command
	accs, err := game.Game.Command("GetAccountList", &model.Account{UserEmail: email})
	if err != nil {
		h.setErr(c, GetAccountListDB, err)
		return
	}

	c.JSON(200, accs)
}

func (h *httpHandle) DeleteAccount(c *gin.Context) {
	// get param
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.setErr(c, DeleteAccountMissingParam, err)
		return
	}

	// command
	if _, err := game.Game.Command("DeleteAccount", &model.Account{ID: id}); err != nil {
		h.setErr(c, DeleteAccountDB, err)
		return
	}

	c.JSON(200, gin.H{})
}

func (h *httpHandle) AddBot(c *gin.Context) {
	in := model.MsgAddBot{}
	if err := c.ShouldBind(&in); err != nil {
		h.setErr(c, AddBotBind, err)
		return
	}
	if err := h.validate.Struct(&in); err != nil {
		h.setErr(c, AddBotValidate, err)
		return
	}
	out, err := game.Game.Command("AddBot", &in)
	if err != nil {
		h.setErr(c, AddBotCommand, err)
		return
	}
	c.JSON(200, out)
}

func (h *httpHandle) DeleteBot(c *gin.Context) {
	in := model.MsgDeleteBot{}
	if err := c.ShouldBind(&in); err != nil {
		h.setErr(c, DeleteBotBind, err)
		return
	}
	if err := h.validate.Struct(&in); err != nil {
		h.setErr(c, DeleteBotValidate, err)
		return
	}
	out, err := game.Game.Command("DeleteBot", &in)
	if err != nil {
		h.setErr(c, DeleteBotCommand, err)
		return
	}
	c.JSON(200, out)
}

func (h *httpHandle) handleGame(c *gin.Context) {
	handleGame(c.Writer, c.Request)
}

func (h *httpHandle) handleCommand(c *gin.Context) {
	type commandRequest struct {
		Action string          `json:"action" validate:"required"`
		In     json.RawMessage `json:"in"`
	}
	req := commandRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.setErr(c, CommandBind, err)
		return
	}
	if err := h.validate.Struct(&req); err != nil {
		h.setErr(c, CommandValidate, err)
		return
	}
	var cmd *command
	var ok bool
	if cmd, ok = h.commands[req.Action]; !ok {
		h.setErr(c, CommandActionInvalid, nil)
		return
	}
	in := reflect.New(reflect.TypeOf(cmd.In).Elem()).Interface()
	data := req.In
	if len(data) == 0 {
		data = []byte("{}")
	}
	if err := json.Unmarshal(data, in); err != nil {
		h.setErr(c, CommandInInvalid, err)
		return
	}
	out, err := game.Game.Command(req.Action, in)
	if err != nil {
		h.setErr(c, CommandExec, err)
		return
	}
	c.JSON(200, gin.H{
		"action": req.Action,
		"out":    out,
	})
}

type command struct {
	Action string `json:"action"`
	In     any    `json:"in"`
}

var commands = [...]*command{
	{"CreateAccount", (*model.Account)(nil)},
	{"GetAccountList", (*model.Account)(nil)},
	{"DeleteAccount", (*model.Account)(nil)},
	{"AddBot", (*model.MsgAddBot)(nil)},
	{"DeleteBot", (*model.MsgDeleteBot)(nil)},
}
