package action_bar

import (
	"github.com/ecletus-pkg/admin"
	"github.com/ecletus/plug"
)

type Plugin struct {
	plug.EventDispatcher
	admin_plugin.AdminNames
}
