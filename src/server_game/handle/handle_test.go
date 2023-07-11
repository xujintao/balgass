package handle

import (
	"testing"

	"github.com/xujintao/balgass/src/server_game/model"
)

func TestMarshal(t *testing.T) {
	tm := model.MsgConnectFailed{Result: 1}
	_, err := APIHandleDefault.Marshal(&tm)
	if err != nil {
		t.Error(err)
	}
}
