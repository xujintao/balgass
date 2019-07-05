package Init

type profile struct {
	isneed bool
	Handle func()
}

var profiles = [...]profile{
	// 224个
	{true, foo},
	{false, bar},
}

func init() {
	for _, p := range profiles {
		if p.isneed {
			p.Handle()
		}
	}
}

// 初始化模块，224个
func foo() {
	println("foo")
}

func bar() {
	println("bar")
}
