// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package limit

// Option defined implement of option
type (
	Option func(*Limit) interface{}
)

// option defined implement of option
func (o Option) apply(r *Limit) *Limit { return o(r).(*Limit) }

// option defined implement of option
func (o Option) check(r *Limit) interface{} { return o(r) }
