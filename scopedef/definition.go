package scopedef

import (
	"github.com/todennus/shared/scopedef/internal"
	"github.com/todennus/x/scope"
)

var Actions, actionMap = scope.DefineAction[internal.Action]()
var Resources, resourceMap = scope.DefineResource[internal.Resource]()
var Engine = scope.NewEngine("todennus", actionMap, resourceMap)
