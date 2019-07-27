// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package limit

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

// DefaultFailureHandler default error handler
var DefaultFailureHandler ErrorHandler = func(ctx *gin.Context) {
	rushLogger.Error("rate limited access, pease check again later")
	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "rate limited access, pease check again later"})
}

func matchRule(rules []Rule, params struct {
	path   string
	method string
}) (bool, Rule) {
	rule := funk.Find(rules, func(rule Rule) bool {
		r, _ := regexp.Compile(rule.Match)
		ruleMatch := r.MatchString(params.path)
		methodMatch := false
		if rule.Match == "" {
			ruleMatch = true
		}
		if m := funk.Find(rule.Methods, func(method string) bool {
			return strings.ToUpper(params.method) == strings.ToUpper(method)
		}); m != nil {
			methodMatch = true
		}
		if len(rule.Methods) == 0 {
			methodMatch = true
		}
		return methodMatch && ruleMatch
	})
	if rule == nil {
		return false, Rule{}
	}
	return true, rule.(Rule)
}
