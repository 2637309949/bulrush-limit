// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

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
