package notifier

import "github.com/Maltide/notification_bot_tg/pkg/types"

type TGNotifier struct {
	store  types.Store
	parser types.Parser
	// TODO: tg_bot
}
