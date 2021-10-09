package common

import "github.com/spf13/viper"

var ErshoufangUrl string
var ChengjiaoUrl string
var District []string

func InitSource() {
	ErshoufangUrl = viper.GetString("source.ershoufang")
	ChengjiaoUrl = viper.GetString("source.chengjiao")
	District = viper.GetStringSlice("source.district")
}
