package bootstrap

import (
	"github.com/olivere/elastic/v7"
	"medium_com/pkg/config"
	pkgelastic "medium_com/pkg/elastic"
	"time"
)

// SetupElastic 初始化 Elastic.
func SetupElastic() {
	pkgelastic.Options = []elastic.ClientOptionFunc{
		elastic.SetURL(config.GetString("elastic.host")),
		elastic.SetBasicAuth(
			config.GetString("elastic.username"),
			config.GetString("elastic.password"),
		),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(5 * time.Second),
	}
}
