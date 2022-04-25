package config

import "medium_com/pkg/config"

func init() {
	config.Add("spider", config.StrMap{
		// 应用名称
		"name": config.Env("SPIDER_NAME", "spiders"),
		// 当前环境，用以区分多环境
		"env": config.Env("SPIDER_ENV", "develop"),
		// 是否开启调试模式
		"debug": config.Env("SPIDER_DEBUG", false),
		// APP 安全密钥，务必去创建一个自己的 GUID 作为密钥：https://www.guidgen.com
		"key": config.Env("SPIDER_KEY", "b2581f25-99a2-4dd2-826d-753a6702903e"),
		// 爬虫域名
		"domain": config.Env("SPIDER_DOMAIN"),
		// 是否开启异步
		"async": config.Env("SPIDER_ASYNC", false),
		// user_agent
		"user_agent": config.Env("SPIDER_USER_AGENT", "Mozilla/5.0 (Windows NT 11.0; Win64; x64) "+
			"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4515.131 Safari/537.36"),
		// Socks5 代理
		"socks5": config.Env("SPIDER_SOCKS5"),
		// 消息队列数量
		"queue_count": config.Env("SPIDER_QUEUE_COUNT", 10),
		// 缓存目录
		"cache_dir": config.Env("SPIDER_CACHE_DIR", "./runtime/cache"),
		// 日志存放目录
		"logger_dir": config.Env("SPIDER_LOGGER_DIR", "./runtime/logs"),
	})
}
