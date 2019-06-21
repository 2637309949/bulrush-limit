// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package limit

import (
	"regexp"

	"github.com/2637309949/bulrush"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

type (
	// Rule defined
	Rule struct {
		Methods []string
		Match   string
		Rate    int16
	}
	// Model token store
	Model interface {
		Save(string, string, string, int16)
		Find(string, string, string, int16) interface{}
	}
	// ErrorHandler handler error
	ErrorHandler func(ctx *gin.Context)
	// Frequency limit
	Frequency struct {
		Passages       []string
		Rules          []Rule
		Model          Model
		FailureHandler ErrorHandler
	}
	// Limit plugin
	Limit struct {
		bulrush.PNBase
		Frequency *Frequency
	}
)

// Plugin for Limit
func (l *Limit) Plugin() bulrush.PNRet {
	return func(router *gin.RouterGroup) {
		router.Use(func(ctx *gin.Context) {
			path := hashPath(ctx)
			ip := ctx.ClientIP()
			method := ctx.Request.Method
			pass := funk.Find(l.Frequency.Passages, func(regex string) bool {
				r, _ := regexp.Compile(regex)
				return r.MatchString(path)
			})
			if pass != nil {
				ctx.Next()
				return
			}
			ruleMatch, rule := matchRule(l.Frequency.Rules, struct {
				path   string
				method string
			}{path: path, method: method})
			if ruleMatch {
				rate := Some(rule.Rate, 1).(int16)
				line := l.Frequency.Model.Find(ip, path, method, rate)
				if line != nil {
					var errorHandler ErrorHandler
					if l.Frequency.FailureHandler == nil {
						errorHandler = DefaultFailureHandler
					}
					errorHandler(ctx)
					return
				}
				l.Frequency.Model.Save(ip, path, method, rule.Rate)
			}
			ctx.Next()
		})
	}
}
