package common

import "github.com/spf13/viper"

var ErshoufangUrl string
var ChengjiaoUrl string
var District []string
var SellingIndex string
var SoldIndex string

func InitSource() {
	ErshoufangUrl = viper.GetString("source.ershoufang")
	ChengjiaoUrl = viper.GetString("source.chengjiao")
	District = viper.GetStringSlice("source.district")
	SellingIndex = viper.GetString("elastic.selling")
	SoldIndex = viper.GetString("elastic.sold")
}
