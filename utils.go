/**
 * @author [Double]
 * @email [2637309949@qq.com.com]
 * @create date 2019-01-12 22:46:31
 * @modify date 2019-01-12 22:46:31
 * @desc [bulrush Limit plugin]
 */

package limit

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

// Some get or a default value
func Some(t interface{}, i interface{}) interface{} {
	if t != nil && t != "" && t != 0 {
		return t
	}
	return i
}

// DefaultFailureHandler default error handler
var DefaultFailureHandler ErrorHandler = func(ctx *gin.Context) {
	rushLogger.Warn("rate limited access, pease check again later")
	ctx.JSON(http.StatusBadRequest, gin.H{"message": "rate limited access, pease check again later"})
	ctx.Abort()
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

func hashPath(ctx *gin.Context) string {
	raw := ctx.Request.URL.RawQuery
	path := ctx.Request.URL.Path
	if raw != "" {
		path = path + "?" + raw
	}
	return path
}
