package action_bar

import (
	"github.com/aghape-pkg/admin"
	"github.com/aghape/plug"
)

type Plugin struct {
	plug.EventDispatcher
	admin_plugin.AdminNames
}
