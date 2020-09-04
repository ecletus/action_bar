package action_bar

import (
	"net/http"

	"github.com/ecletus/admin"
	"github.com/ecletus/core/utils"
)

type controller struct {
	ActionBar *ActionBar
}

// SwitchMode is handle to store switch status in cookie
func (controller) SwitchMode(w http.ResponseWriter, r *http.Request) {
	context := admin.ContextFromRequest(r)
	utils.SetCookie(http.Cookie{Name: "qor-action-bar", Value: r.URL.Query().Get("checked")}, context.Context)

	referrer := r.Referer()
	if referrer == "" {
		referrer = "/"
	}

	http.Redirect(context.Writer, r, referrer, http.StatusFound)
}

// InlineEdit using to make inline edit resource shown as slideout
func (controller) InlineEdit(w http.ResponseWriter, r *http.Request) {
	admin.ContextFromRequest(r).Include(w, "action_bar/inline_edit")
}
