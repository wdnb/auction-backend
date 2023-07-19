package casbin

import (
	"auction-website/utils/resp"
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"regexp"
	"strings"
)

type CasbinMiddleware struct {
	enforcer *casbin.Enforcer
	subFn    SubjectFn
}

// SubjectFn is used to look up current subject in runtime.
// If it can not find anything, just return an empty string.
type SubjectFn func(c *gin.Context) (string, error)

// Logic is the logical operation (AND/OR) used in permission checks
// in case multiple permissions or roles are specified.
type Logic int

const (
	AND Logic = iota
	OR
)

var (
	ErrSubFnNil = errors.New("subFn is nil")
)

// NewCasbinMiddleware returns a new CasbinMiddleware using Casbin's Enforcer internally.
// modelFile is the file path to Casbin model file e.g. path/to/rbac_model.conf.
// policyAdapter can be a file or a Mysql adapter.
// File: path/to/basic_policy.csv
// MySQL Mysql: mysqladapter.NewDBAdapter("mysql", "mysql_username:mysql_password@tcp(127.0.0.1:3306)/")
// subFn is a function that looks up the current subject in runtime and returns an empty string if nothing found.
func NewCasbinMiddleware(modelFile string, policyAdapter interface{}, subFn SubjectFn) (*CasbinMiddleware, error) {
	e, err := casbin.NewEnforcer(modelFile, policyAdapter)
	if err != nil {
		return nil, err
	}

	return NewCasbinMiddlewareFromEnforcer(e, subFn)
}

// Create from given Enforcer.
func NewCasbinMiddlewareFromEnforcer(e *casbin.Enforcer, subFn SubjectFn) (*CasbinMiddleware, error) {
	if subFn == nil {
		return nil, ErrSubFnNil
	}

	return &CasbinMiddleware{
		enforcer: e,
		subFn:    subFn,
	}, nil
}

// Option is used to change some default behaviors.
type Option interface {
	apply(*options)
}

type options struct {
	logic Logic
}

type logicOption Logic

func (lo logicOption) apply(opts *options) {
	opts.logic = Logic(lo)
}

// WithLogic sets the logical operator used in permission or role checks.
func WithLogic(logic Logic) Option {
	return logicOption(logic)
}

// CheckPermissions tries to find the current subject by calling SubjectFn
func (am *CasbinMiddleware) CheckPermissions(opts ...Option) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := am.RequiresPermissions(c, am.AutoGenPermissions(c), opts...)
		if err != nil {
			resp.ErrorResponse(c, http.StatusForbidden, resp.ERROR, err)
			return
		}
	}
}

// AutoGenPermissions 自动生成权限字符串
func (am *CasbinMiddleware) AutoGenPermissions(c *gin.Context) []string {
	permissions := make([]string, 0)
	// 获取请求方法
	method := c.Request.Method
	// 获取请求URL
	URL := c.Request.URL.Path
	// 去除ID字段
	URL = regexp.MustCompile(`/\d+`).ReplaceAllString(URL, "/*")
	// 拼接权限字符串
	permission := URL + ":" + method
	permissions = append(permissions, permission)
	// 处理RESTful风格的URL中含有参数的情况
	return permissions
}

// RequiresPermissions tries to find the current subject by calling SubjectFn
// and determine if the subject has the required permissions according to predefined Casbin policies.
// permissions are formatted strings. For example, "file:read" represents the permission to read a file.
// opts is some optional configurations such as the logical operator (default is AND) in case multiple permissions are specified.
func (am *CasbinMiddleware) RequiresPermissions(c *gin.Context, permissions []string, opts ...Option) error {
	if len(permissions) == 0 {
		return nil
		//c.Next()
		//return
	}
	// Here we provide default options.
	actualOptions := options{
		logic: AND,
	}
	//fmt.Println(actualOptions)
	// Apply actual options.
	for _, opt := range opts {
		opt.apply(&actualOptions)
	}
	// Look up current subject.
	sub, err := am.subFn(c)
	if err != nil {
		return err
	}

	if sub == "" {
		//return nil
		//c.AbortWithStatus(401)
		//resp.ErrorResponse(c, http.StatusUnauthorized, "sub is empty", "Unauthorized")
		return errors.New("sub is empty")
	}
	// Enforce Casbin policies.
	if actualOptions.logic == AND {
		//fmt.Println(permissions)
		// Must pass all tests.
		for _, permission := range permissions {
			obj, act := parsePermissionStrings(permission)
			if obj == "" || act == "" {
				// Can not handle any illegal permission strings.
				return errors.New("illegal permission string:" + permission)
			}
			//sub = "seller"
			fmt.Println(sub)
			fmt.Println(obj)
			fmt.Println(act)
			ok, err := am.enforcer.Enforce(sub, obj, act)
			//fmt.Println(err)
			if err != nil {
				return err
				// 处理err
			}

			if ok == true {
				return nil
			} else {
				return errors.New("权限不足")
			}
			//!ok || err != nil {
			//	fmt.Println(ok)
			//	fmt.Println(err)
			//	return err
			//}
		}
	} else {
		// Need to pass at least one test.
		for _, permission := range permissions {
			obj, act := parsePermissionStrings(permission)
			if obj == "" || act == "" {
				zap.L().Error("illegal permission string:" + permission)
				continue
			}

			if ok, err := am.enforcer.Enforce(sub, obj, act); ok && err == nil {
				return nil
			}
		}
		return errors.New("不存在的权限")
	}
	return nil
}

func parsePermissionStrings(str string) (string, string) {
	if !strings.Contains(str, ":") {
		return "", ""
	}
	vals := strings.Split(str, ":")
	return vals[0], vals[1]
}
