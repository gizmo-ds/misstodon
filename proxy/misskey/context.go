package misskey

type Context interface {
	ProxyServer() string
	Token() *string
	UserID() *string
	HOST() *string
}
