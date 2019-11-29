package service

var (
	Server *server
	Update *update
)

func init() {
	Server = &server{}
	Update = &update{}
}
