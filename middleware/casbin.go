package middleware

import (
	"auction-website/conf"
	"auction-website/middleware/casbin"
	"path/filepath"
	"time"

	sqlxadapter "github.com/Blank-Xu/sqlx-adapter"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type Claims struct {
	ID uint32 `json:"id"`
	jwt.RegisteredClaims
}

func (c *Config) Casbin() *casbin.CasbinMiddleware {
	configPath := conf.GetConfigPath() + string(filepath.Separator)
	auth, err := casbin.NewCasbinMiddleware(configPath+"rbac_model.conf", c.getAdapter("casbin_rule", "rbac_model.conf"), subjectFromJWT)
	if err != nil {
		panic(err)
	}
	return auth
}

func (c *Config) getAdapter(tableName, rbacModelName string) *sqlxadapter.Adapter {
	// Initialize a Sqlx adapter and use it in a Casbin enforcer:
	// The adapter will use the Sqlite3 table name "casbin_rule_test",
	// the default table name is "casbin_rule".
	// If it doesn't exist, the adapter will create it automatically.
	a, err := sqlxadapter.NewAdapter(c.db, tableName)
	if err != nil {
		panic(err)
	}
	return a
}

// TODO jwt 过期时间考虑一下
func subjectFromJWT(c *gin.Context) (string, error) {
	// Check if access_token exists in both header and cookie
	tokenString := c.Request.Header.Get("Access-Token")
	if tokenString == "" {
		var err error
		tokenString, err = c.Cookie("Access-Token")
		if err != nil {
			return "", err
		}
	}

	cl, err := VerifyToken(tokenString)
	if err != nil {
		return "", err
	}

	// Set "uid" key in Gin context
	c.Set("uid", cl.ID)
	//c.Set("username", cl.Username)
	return cl.Subject, nil
}

func GetSecret() []byte {
	jwtKey := []byte(viper.GetString("jwt.secret"))
	return jwtKey
}

func GenerateToken(uid uint32, roleName string) (string, error) {
	mySigningKey := GetSecret()
	claims := Claims{
		ID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   roleName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), //TODO jwt过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return ss, nil
}

// VerifyToken verifies the JWT token and returns the claims if valid
func VerifyToken(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return GetSecret(), nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, err
	} else {
		return nil, err
	}
}
