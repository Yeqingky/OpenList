package data

import "github.com/OpenListTeam/OpenList/v4/cmd/flags"

func InitData() {
	initUser()
	initSettings()
	if flags.Dev {
		initDevData()
		initDevDo()
	}
}
