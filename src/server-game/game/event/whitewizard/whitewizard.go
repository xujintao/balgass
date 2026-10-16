package whitewizard

import (
	"fmt"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/event"
	"github.com/xujintao/balgass/src/server-game/game/event/internal/herd"
)

func init() {
	cfg, err := herd.Load(conf.PathCommon, "Events/IGC_WhiteWizardAttack.xml", "WhiteWizard", "白色魔法师活动")
	if err != nil {
		panic(fmt.Errorf("whitewizard: %w", err))
	}
	if cfg.Enabled {
		if err := event.EventManager.Register("whitewizard", herd.New(cfg), cfg.Schedules); err != nil {
			panic(err)
		}
	}
}
