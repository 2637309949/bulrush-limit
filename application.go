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

const defaultRate = 1

// DefaultFailureHandler default error handler
var DefaultFailureHandler ErrorHandler = func(ctx *gin.Context) {
	rushLogger.Warn("Rate Limited access, Pease Check Again Later")
	ctx.JSON(http.StatusBadRequest, gin.H{
		"message": "Rate Limited access, Pease Check Again Later.",
	})
	ctx.Abort()
}

// Plugin for Limit
func (l *Limit) Plugin() bulrush.PNRet {
	return func(router *gin.RouterGroup) {
		router.Use(func(ctx *gin.Context) {
			raw := ctx.Request.URL.RawQuery
			path := ctx.Request.URL.Path
			if raw != "" {
				path = path + "?" + raw
			}
			ip := ctx.ClientIP()
			method := ctx.Request.Method
			pass := funk.Find(l.Frequency.Passages, func(regex string) bool {
				r, _ := regexp.Compile(regex)
				return r.MatchString(path)
			})
			if pass != nil {
				ctx.Next()
			} else {
				rule := funk.Find(l.Frequency.Rules, func(rule Rule) bool {
					r, _ := regexp.Compile(rule.Match)
					ruleMatch := r.MatchString(path)
					methodMatch := false
					if rule.Match == "" {
						ruleMatch = true
					}
					m := funk.Find(rule.Methods, func(method string) bool {
						return strings.ToUpper(method) == strings.ToUpper(method)
					})
					if m != nil {
						methodMatch = true
					}
					if len(rule.Methods) == 0 {
						methodMatch = true
					}
					return methodMatch && ruleMatch
				})
				if rule != nil {
					rate := Some((rule.(Rule)).Rate, defaultRate).(int16)
					line := l.Frequency.Model.Find(ip, path, method, rate)
					if line != nil {
						var handler ErrorHandler
						if l.Frequency.FailureHandler == nil {
							handler = DefaultFailureHandler
						}
						handler(ctx)
						ctx.Abort()
					} else {
						l.Frequency.Model.Save(ip, path, method, (rule.(Rule)).Rate)
					}
				} else {
					ctx.Next()
				}
			}
		})
	}
}
