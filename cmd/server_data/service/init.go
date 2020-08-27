package service

var (
	ServerManager *serverManager
	UserManager   *userManager
	VIPManager    *vipManager
	CharManager   *charManager
)

func init() {
	ServerManager = &serverManager{}
	UserManager = &userManager{}
	VIPManager = &vipManager{}
	CharManager = &charManager{}
}
