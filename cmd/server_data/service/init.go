package service

var (
	ServerManager *serverManager
	UserManager   *userManager
	VIPManager    *vipManager
)

func init() {
	ServerManager = &serverManager{}
	UserManager = &userManager{}
	VIPManager = &vipManager{}
}
