package service

var (
	ServerManager  *serverManager
	AccountManager *accountManager
)

func init() {
	ServerManager = &serverManager{}
	AccountManager = &accountManager{}
}
