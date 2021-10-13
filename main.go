package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
	"xietong.me/LianjiaSpider/common"
	"xietong.me/LianjiaSpider/spider"
)

func main() {
	//初始化配置
	InitConfig()
	// 初始化爬虫源站信息
	common.InitSource()
	// 初始化es实例
	esClient := common.InitElastic()
	district := common.District
	// 空间过滤条件 ba10ea50 - 大于10平米小于50平米
	filters := []string{"ba10ea50", "ba50ea60", "ba60ea70", "ba70ea80", "ba80ea90", "ba90ea100", "ba100ea110", "ba110ea120", "ba120ea130", "ba130ea140", "ba140ea150"}

	/********************************************************************************************
										获取在售房源信息
	********************************************************************************************/
	// 通过循环来爬取不同地区，同时获取不同地区的总分页数来爬取不同页面的数据
	for _, districtName := range district {
		// 链家每次搜索限制返回30页，通过不同平方数量过滤，获取全部房源信息，如果区域房源数量大，可调整filter粒度
		for _, filter := range filters {
			// 获取分页数量
			totalSellingPage := spider.GetSellingPageSpider(districtName, filter)
			fmt.Println("totalPage:", totalSellingPage)
			for page := 1; page <= totalSellingPage; page++ {
				// 20页休眠5秒
				if page%20 == 0 {
					time.Sleep(time.Second * 5)
				}
				func(page int) {
					spider.GetSellingInfoSpider(esClient, districtName, page, filter)
				}(page)
			}
		}
	}

	/********************************************************************************************
										获取历史交易房源信息
	********************************************************************************************/
	for _, districtName := range district {
		for _, filter := range filters {
			totalSoldPage := spider.GetSoldPageSpider(districtName, filter)
			fmt.Println("totalPage:", totalSoldPage)
			for page := 1; page <= totalSoldPage; page++ {
				// 20页休眠5秒
				if page%20 == 0 {
					time.Sleep(time.Second * 5)
				}
				func(page int) {
					spider.GetSoldInfoSpider(esClient, districtName, page, filter)
				}(page)
			}
		}
	}
}

//初始化配置函数
func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
