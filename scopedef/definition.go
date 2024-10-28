package scopedef

import (
	"github.com/todennus/shared/scopedef/internal"
	"github.com/todennus/x/scope"
)

const (
	TitleUser  = "user"
	TitleApp   = "app"
	TitleAdmin = "admin"
)

var Actions, actionMap = scope.DefineAction[internal.Action]()
var Resources, resourceMap = scope.DefineResource[internal.Resource]()
var Engine = scope.NewEngine("todennus", actionMap, resourceMap)

func init() {
	Engine.DefineTitle(TitleUser)
	Engine.DefineTitle(TitleApp)
	Engine.DefineTitle(TitleAdmin)

	Engine.DefineCoverable(TitleAdmin, TitleUser)
	Engine.DefineCoverable(TitleAdmin, TitleApp)
}
