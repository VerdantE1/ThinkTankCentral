/**
 * Created by GoLand.
 * User: buzzlight.frank@qq.com
 * Date: 2025/3/14
 * Time: 16:26
 */

/**
全局变量实例
*/

package global

import (
	"ThinkTankCentral/config"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config     *config.Config
	Log        *zap.Logger
	DB         *gorm.DB
	ESClient   *elasticsearch.TypedClient
	Redis      redis.Client
	BlackCache local_cache.Cache
)
