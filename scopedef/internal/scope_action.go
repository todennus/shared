package internal

import "github.com/todennus/x/scope"

type Action struct {
	*scope.BaseAction

	Read  *scope.BaseAction `action:"read"`
	Write *WriteAction      `action:"write"`
}

type WriteAction struct {
	*scope.BaseAction

	Create *scope.BaseAction `action:"create"`
	Update *scope.BaseAction `action:"update"`
	Delete *scope.BaseAction `action:"delete"`
}
