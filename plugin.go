package action_bar

import (
	"github.com/aghape/admin/adminplugin"
	"github.com/aghape/plug"
)

type Plugin struct {
	plug.EventDispatcher
	adminplugin.AdminNames
}
