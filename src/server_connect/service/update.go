package service

import (
	"github.com/xujintao/balgass/cmd/server_connect/conf"

	"github.com/xujintao/balgass/cmd/server_connect/model"
)

type update struct{}

func (*update) CheckVersion(v *model.Version) interface{} {
	// var major, minor, patch uint8
	// fmt.Sscanf(conf.Version, "%u.%u.%u", &major, &minor, &patch)
	// if v.Major == major && v.Minor == minor && v.Patch == patch {
	// 	return true
	// }
	if *v == conf.AutoUpdate.Ver {
		return true
	}
	return &conf.AutoUpdate // copy
}
