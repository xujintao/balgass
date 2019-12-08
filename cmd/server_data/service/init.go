package service

var (
	ServerManager *serverManager
)

func init() {
	ServerManager = &serverManager{}
}
