package enumdef

type SubjectType int

const (
	SubjectUser SubjectType = iota
	SubjectClient
)
