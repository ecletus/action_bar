package action_bar

import (
	"net/http"

	"github.com/ecletus/admin"
	"github.com/ecletus/core/utils"
	"github.com/moisespsena-go/xroute"
)

type controller struct {
	ActionBar *ActionBar
}

// SwitchMode is handle to store switch status in cookie
func (controller) SwitchMode(rctx *xroute.RouteContext) {
	context := admin.ContextFromRouteContext(rctx)
	utils.SetCookie(http.Cookie{Name: "qor-action-bar", Value: context.Request.URL.Query().Get("checked")}, context.Context)

	referrer := context.Request.Referer()
	if referrer == "" {
		referrer = "/"
	}

	http.Redirect(context.Writer, context.Request, referrer, http.StatusFound)
}

// InlineEdit using to make inline edit resource shown as slideout
func (controller) InlineEdit(rctx *xroute.RouteContext) {
	context := admin.ContextFromRouteContext(rctx)
	context.Writer.Write([]byte(context.Render("action_bar/inline_edit")))
}
