package enumdef

import (
	"github.com/xybor-x/enum"
)

type subjectType any
type SubjectType = enum.WrapEnum[subjectType]

const (
	SubjectUser SubjectType = iota
	SubjectClient
)

func init() {
	enum.Map(SubjectUser, "user")
	enum.Map(SubjectClient, "client")
	enum.Finalize[SubjectType]()
}
