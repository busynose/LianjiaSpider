package spider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gocolly/colly"
	"xietong.me/LianjiaSpider/common"
	"xietong.me/LianjiaSpider/model"
)

func GetSellingInfoSpider(esClient *elasticsearch.Client, districtName string, page int) {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"),
	)
	c.SetRequestTimeout(time.Duration(120) * time.Second)
	c.Limit(&colly.LimitRule{DomainGlob: common.ErshoufangUrl, Parallelism: 1}) //Parallelism代表最大并发数
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})
	//访问所有info 访问前20页采用goroutine
	c.OnHTML(".sellListContent>li", func(e *colly.HTMLElement) {
		re, _ := regexp.Compile(`\d+`)                                                                                                    //正则表达式用来匹配数字
		houseId := e.Attr("data-lj_action_housedel_id")                                                                                   //获取房子ID，可根据ID直接访问房子详情主页
		nameRegion := e.ChildText("div.info > div.flood > div.positionInfo > a")                                                          //同时获取小区名和详细地区
		name := strings.Split(nameRegion, " ")[0]                                                                                         //将同时获取的小区名和详细地区分离，取其中的小区名字
		region := strings.Split(nameRegion, " ")[1]                                                                                       //将同时获取的小区名和详细地区分离，取其中的详细地区
		totalPrice, _ := strconv.Atoi(string(re.Find([]byte(e.DOM.Find(".info .priceInfo .totalPrice span").Eq(0).Text()))))              //根据页面元素获取总价，正则匹配数字，转换成int类型
		unitPrice, _ := strconv.Atoi(string(re.Find([]byte(e.DOM.Find(".info .priceInfo .unitPrice span").Eq(0).Text()))))                //读取页面元素获取单价,正则匹配单价的数字，转换成int类型
		area, _ := strconv.Atoi(string(re.Find([]byte(strings.Split(e.ChildText("div.info > div.address > div.houseInfo "), " | ")[1])))) // //读取页面元素获取面积,正则匹配单价的数字，转换成int类型
		info := string([]byte(e.ChildText("div.info > div.address > div.houseInfo ")))
		if houseId != "" {
			// fmt.Println("start save", houseId, page)
			sellingInfo := model.Selling{Id: houseId, Name: name, TotalPrice: totalPrice, UnitPrice: unitPrice, District: districtName, Region: region, Area: area, Info: info}

			body, err := json.Marshal(sellingInfo)
			if err != nil {
				fmt.Println(err)
				return
			}

			req := esapi.IndexRequest{
				Index:        common.SellingIndex,
				DocumentType: "selling",
				DocumentID:   houseId,
				Body:         strings.NewReader(string(body)),
				Refresh:      "true",
				Timeout:      time.Second * 60,
			}

			_, err = req.Do(context.Background(), esClient)
			if err != nil {
				fmt.Println(err)
			}
		}
	})
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
		c.Visit(common.ErshoufangUrl + districtName + "/pg" + strconv.Itoa(page))
	})
	c.Visit(common.ErshoufangUrl + districtName + "/pg" + strconv.Itoa(page))
	c.Wait()

}
