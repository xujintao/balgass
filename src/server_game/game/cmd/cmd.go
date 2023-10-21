package cmd

import (
	"github.com/xujintao/balgass/src/server_game/game/bot"
	"github.com/xujintao/balgass/src/server_game/game/model"
	"github.com/xujintao/balgass/src/server_game/game/object"
)

var Command command

type command struct{}

func (*command) AddBot(msg *model.MsgAddBot) (any, error) {
	bot.BotManager.AddBot(msg.Name)
	return nil, nil
}

func (*command) DeleteBot(msg *model.MsgDeleteBot) (any, error) {
	bot.BotManager.DeleteBot(msg.Name)
	return nil, nil
}

func (*command) OfflineAllObjects(msg any) (any, error) {
	object.ObjectManager.OfflineAllObjects()
	return nil, nil
}

func (*command) GetOnlineObjectsNumber(msg any) (*model.MsgGetOnlineObjectNumberReply, error) {
	return object.ObjectManager.GetOnlineObjectsNumber(), nil
}

// func (*command) GetObjectsByMapNumber(msg *model.MsgSubscribeMap) (*model.MsgSubscribeMapReply, error) {
// 	number := msg.Number
// 	pots := object.ObjectManager.GetObjectsByMapNumber(number)
// 	return &model.MsgSubscribeMapReply{
// 		Name: "object",
// 		Data: pots,
// 	}, nil
// }

func (*command) CreateAccount(msg *model.Account) (*model.Account, error) {
	return model.DB.CreateAccount(msg)
}

func (*command) DeleteAccount(account string) (any, error) {
	return nil, model.DB.DeleteAccount(account)
}
