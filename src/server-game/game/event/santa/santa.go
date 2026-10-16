package santa

import (
	"fmt"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/event"
	"github.com/xujintao/balgass/src/server-game/game/event/internal/herd"
)

func init() {
	cfg, err := herd.Load(conf.PathCommon, "Events/IGC_XMasAttack.xml", "XMasAttack", "圣诞活动")
	if err != nil {
		panic(fmt.Errorf("santa: %w", err))
	}
	if cfg.Enabled {
		if err := event.EventManager.Register("santa", herd.New(cfg), cfg.Schedules); err != nil {
			panic(err)
		}
	}
}
