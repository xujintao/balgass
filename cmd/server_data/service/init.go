package service

var (
	ServerManager  *serverManager
	AccountManager *accountManager
	VIPManager     *vipManager
)

func init() {
	ServerManager = &serverManager{}
	AccountManager = &accountManager{}
	VIPManager = &vipManager{}
}
