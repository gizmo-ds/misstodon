package misskey

type Context interface {
	Server() string
	Token() *string
}
