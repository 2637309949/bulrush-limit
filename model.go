// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package limit

import (
	"strconv"
	"time"

	redisext "github.com/2637309949/bulrush-addition/redis"
	gredis "github.com/go-redis/redis"
)

// RedisModel adapter for redis
type RedisModel struct {
	Model
	Redis *redisext.Redis
}

// Save save a token
func (model *RedisModel) Save(ip string, url string, method string, rate int) {
	value, err := model.Redis.Client.Get("LIMIT:" + ip + url + method).Result()
	if value != "" && err != gredis.Nil {
		model.Redis.Client.Incr("LIMIT:" + ip + url + method).Result()
	} else {
		model.Redis.Client.Set("LIMIT:"+ip+url+method, "1", time.Duration(10)*time.Second).Result()
		model.Redis.Client.Get("LIMIT:" + ip + url + method).Result()
	}
}

// Find find a token
func (model *RedisModel) Find(ip string, url string, method string, rate int) interface{} {
	value, err := model.Redis.Client.Get("LIMIT:" + ip + url + method).Result()
	if err == gredis.Nil {
		return nil
	}
	if num, err := strconv.Atoi(value); err == nil {
		if num < rate {
			return nil
		}
		return num
	}
	return nil
}
