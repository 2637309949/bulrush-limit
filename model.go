/**
 * @author [Double]
 * @email [2637309949@qq.com.com]
 * @create date 2019-01-12 22:46:31
 * @modify date 2019-01-12 22:46:31
 * @desc [bulrush Limit plugin]
 */

package limit

import (
	"time"

	"github.com/2637309949/bulrush-addition/redis"
)

// RedisModel adapter for redis
type RedisModel struct {
	Model
	Redis *redis.Redis
}

// Save save a token
func (group *RedisModel) Save(ip string, url string, method string, rate int16) {
	group.Redis.Client.Set("LIMIT:"+ip+url+method, 1, time.Duration(rate)*time.Second)
}

// Find find a token
func (group *RedisModel) Find(ip string, url string, method string, rate int16) interface{} {
	if value, err := group.Redis.Client.Get("LIMIT:" + ip + url + method).Result(); err == nil {
		return value
	}
	return nil
}
