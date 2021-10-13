# LianjiaSpider
## source
https://github.com/xietongMe/LianjiaSpider

## fix
- 区域/城市做配置项
- 修改存储方式为elastic，方便kibana做图形化分析
- lianjia接口限制最大返回100页码，通过面积分类分散查找以完成采集所有信息

## bug
- 多线程采集遇到限流问题，采用单线程+周期休眠
