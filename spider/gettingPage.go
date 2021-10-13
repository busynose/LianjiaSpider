package spider

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"lianjiaSpider/common"

	"github.com/gocolly/colly"
)

//定义page结构体用来处理json
type Page struct {
	TotalPage int `json:"totalPage"`
	CurPage   int `json:"curPage"`
}

func GetSellingPageSpider(districtName string, area string) int {
	var totalPage int
	c := colly.NewCollector(
		//colly.Async(true),并发
		colly.AllowURLRevisit(),
		colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"),
	)
	c.SetRequestTimeout(time.Duration(300) * time.Second)
	c.Limit(&colly.LimitRule{DomainGlob: common.ErshoufangUrl, Parallelism: 1}) //Parallelism代表最大并发数
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})
	//获取不同地区的总页数
	c.OnHTML(".contentBottom .house-lst-page-box", func(e *colly.HTMLElement) {
		page := Page{}
		err := json.Unmarshal([]byte(e.Attr("page-data")), &page)
		if err != nil {
			log.Fatalln(err)
		}
		totalPage = page.TotalPage
	})
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
		c.Visit(common.ErshoufangUrl + districtName + "/" + area)
	})
	c.Visit(common.ErshoufangUrl + districtName + "/" + area)
	c.Wait()
	return totalPage
}

func GetSoldPageSpider(districtName string, area string) int {
	var totalPage int
	c := colly.NewCollector(
		//colly.Async(true),并发
		colly.AllowURLRevisit(),
		colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"),
	)
	c.SetRequestTimeout(time.Duration(300) * time.Second)
	c.Limit(&colly.LimitRule{DomainGlob: common.ChengjiaoUrl, Parallelism: 1}) //Parallelism代表最大并发数
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})
	//获取不同地区的总页数
	c.OnHTML(".contentBottom .house-lst-page-box", func(e *colly.HTMLElement) {
		page := Page{}
		err := json.Unmarshal([]byte(e.Attr("page-data")), &page)
		if err != nil {
			log.Fatalln(err)
		}
		totalPage = page.TotalPage
	})
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
		c.Visit(common.ChengjiaoUrl + districtName + "/" + area)
	})
	c.Visit(common.ChengjiaoUrl + districtName + "/" + area)
	c.Wait()
	return totalPage
}
