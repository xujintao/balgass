package service

var (
	ServerManager *serverManager
	Update        *update
	UserManager   *userManager
)

func init() {
	ServerManager = &serverManager{}
	Update = &update{}
	UserManager = &userManager{}
}
