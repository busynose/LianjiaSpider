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

func GetSoldInfoSpider(esClient *elasticsearch.Client, districtName string, page int, filter string) {

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
	//访问所有info 访问前20页采用goroutine
	c.OnHTML(".listContent>li", func(e *colly.HTMLElement) {
		re, _ := regexp.Compile(`\d+`)                                                                             //正则表达式用来匹配数字
		houseId := string(re.Find([]byte(strings.Split(e.ChildAttr("div.info > div.title > a", "href"), "/")[4]))) //获取房子ID，可根据ID直接访问房子详情主页
		name := strings.Split(e.ChildText("div.info > div.title > a"), " ")[0]                                     //获取小区名
		area := 0
		if len(strings.Split(e.ChildText("div.info > div.title > a"), " ")) == 3 {
			area, _ = strconv.Atoi(string(re.Find([]byte(strings.Split(e.ChildText("div.info > div.title > a"), " ")[2])))) //获取总面积
		}
		totalPrice, _ := strconv.Atoi(e.DOM.Find(".info .address .totalPrice span").Eq(0).Text())                      //获取总价
		unitPrice, _ := strconv.Atoi(string(re.Find([]byte(e.DOM.Find(".info .flood .unitPrice span").Eq(0).Text())))) //获取单价
		dealDate := e.DOM.Find(".info .address .dealDate").Eq(0).Text()                                                //获取成交年月日
		soldYear := strings.Split(dealDate, ".")[0]                                                                    //分离出成交年份
		soldMonth := strings.Split(dealDate, ".")[1]                                                                   //分离出成交月
		t, _ := time.Parse("2006.01.02 15:04:05", dealDate+" 00:00:00")
		if houseId != "" {
			// fmt.Println("start save", houseId, page)
			soldInfo := model.Sold{HouseId: houseId, Name: name, TotalPrice: totalPrice, UnitPrice: unitPrice, District: districtName, SoldYear: soldYear, SoldMonth: soldMonth, Area: area, TimeStamp: t.UTC()}
			body, err := json.Marshal(soldInfo)
			if err != nil {
				fmt.Println(err)
				return
			}

			req := esapi.IndexRequest{
				Index:      common.SoldIndex,
				DocumentID: houseId,
				Body:       strings.NewReader(string(body)),
			}

			rsp, err := req.Do(context.Background(), esClient)
			if err != nil {
				fmt.Println(err)
			}

			_ = rsp.Body.Close()
		}
	})
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
		time.Sleep(time.Second * 10)
		c.Visit(common.ChengjiaoUrl + districtName + "/pg" + strconv.Itoa(page) + filter + "/")
	})
	c.Visit(common.ChengjiaoUrl + districtName + "/pg" + strconv.Itoa(page) + filter + "/")
	c.Wait()
}
