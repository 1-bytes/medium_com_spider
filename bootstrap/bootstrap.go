package bootstrap

import "medium_com/config"

// Setup 初始化指定的服务，例如：Redis MySQL Logger 等模块.
func Setup() {
	autoLoader(
		config.Initialize, // 配置文件
		SetupElastic,
		SetupMySQL,
		SetupRedisStorage,
		SetupQueued,
	)
}

// autoLoader 自动加载初始化.
func autoLoader(funcName ...func()) {
	for _, v := range funcName {
		v()
	}
}
