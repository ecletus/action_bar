package action_bar

import (
	"fmt"

	"github.com/moisespsena-go/xroute"

	"github.com/ecletus/admin"
	"github.com/ecletus/auth"
	"github.com/ecletus/core/utils"
	"github.com/ecletus/core/utils/url"
	"github.com/moisespsena/template/html/template"
)

// ActionBar stores configuration about a action bar.
type ActionBar struct {
	Admin         *admin.Admin
	GlobalActions []ActionInterface
	actions       []ActionInterface
}

// New will create an ActionBar object
func New(admin *admin.Admin) *ActionBar {
	bar := &ActionBar{Admin: admin}
	ctr := &controller{ActionBar: bar}
	admin.OnRouter(func(r xroute.Router) {
		r.Get("/action_bar/switch_mode", ctr.SwitchMode)
		r.Get("/action_bar/inline_edit", ctr.InlineEdit)
	})
	return bar
}

// RegisterAction register global action
func (bar *ActionBar) RegisterAction(action ActionInterface) {
	bar.GlobalActions = append(bar.GlobalActions, action)
	bar.actions = bar.GlobalActions
}

// Actions register actions
func (bar *ActionBar) Actions(actions ...ActionInterface) *ActionBar {
	newBar := &ActionBar{Admin: bar.Admin, actions: bar.GlobalActions}
	newBar.actions = append(newBar.actions, actions...)
	return newBar
}

// Render will return the HTML of the bar, used this function to render the bar in frontend page's template or layout
func (bar *ActionBar) Render(context *admin.Context) template.HTML {
	var (
		actions, inlineActions []ActionInterface
	)
	for _, action := range bar.actions {
		if action.InlineAction() {
			inlineActions = append(inlineActions, action)
		} else {
			actions = append(actions, action)
		}
	}
	result := map[string]interface{}{
		"EditMode":      bar.EditMode(context),
		"Auth":          bar.Admin.Auth,
		"CurrentUser":   context.Context.CurrentUser(),
		"Actions":       actions,
		"InlineActions": inlineActions,
		"RouterPrefix":  bar.Admin.Config.MountPath,
	}
	defer context.PushI18nGroup(I18NGROUP)()
	return context.RenderHtml("action_bar/action_bar", result)
}

// FuncMap will return helper to render inline edit button
func (bar *ActionBar) FuncMap(context *admin.Context) template.FuncMap {
	funcMap := template.FuncMap{}

	funcMap["render_edit_button"] = func(value interface{}, resources ...*admin.Resource) template.HTML {
		return bar.RenderEditButtonWithResource(context, value, resources...)
	}

	return funcMap
}

// EditMode return whether current mode is `Preview` or `Edit`
func (bar *ActionBar) EditMode(context *admin.Context) bool {
	return isEditMode(context)
}

func isEditMode(context *admin.Context) bool {
	if Auth := context.Admin.Auth; Auth != nil {
		if user, err := Auth.GetCurrentUser(context); user == nil {
			return false
		} else if err != nil && err != auth.ErrNoSession {
			context.AddError(err)
		}
	}
	if cookie, err := context.Request.Cookie("qor-action-bar"); err == nil {
		return cookie.Value == "true"
	}
	return false
}

func (bar *ActionBar) RenderEditButtonWithResource(context *admin.Context, value interface{}, resources ...*admin.Resource) template.HTML {
	var (
		admin        = bar.Admin
		resourceName = utils.ModelType(value).String()
		editURL, _   = url.JoinURL(context.URLFor(value, resources...), "edit")
	)

	if res := admin.GetResourceByID(resourceName); res != nil {
		resourceName = string(admin.T(context.Context, fmt.Sprintf("%v.name", res.ToParam()), res.Name))
	}

	for _, res := range resources {
		resourceName = string(admin.T(context.Context, fmt.Sprintf("%v.name", res.ToParam()), res.Name))
	}

	title := string(admin.T(context.Context, I18NGROUP+".action.edit_resource", "Edit {{$1}}", resourceName))
	return bar.RenderEditButton(context, title, editURL)
}

func (bar *ActionBar) RenderEditButton(context *admin.Context, title string, link string) template.HTML {
	if bar.EditMode(context) {
		var (
			prefix   = bar.Admin.Router.Prefix()
			jsURL    = fmt.Sprintf("<script data-prefix=\"%v\" src=\"%v/assets/javascripts/action_bar_check.js?theme=action_bar\"></script>", prefix, prefix)
			frameURL = fmt.Sprintf("%v/action_bar/inline_edit", prefix)
		)

		return template.HTML(fmt.Sprintf(`%v<a target="blank" data-iframe-url="%v" data-url="%v" href="#" class="qor-actionbar-button">%v</a>`, jsURL, frameURL, link, title))
	}
	return template.HTML("")
}

func (bar *ActionBar) PrepareAdmin(Admin *admin.Admin) {
}
