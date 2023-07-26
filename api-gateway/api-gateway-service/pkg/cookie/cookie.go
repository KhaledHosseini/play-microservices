package cookie

import (
	"github.com/KhaledHosseini/play-microservices/api-gateway/api-gateway-service/config"
	"github.com/gin-gonic/gin"
)

const (
	AUTH_HEADER          = "authorization"
	REFRESH_TOKEN_HEADER = "x-refresh-token"
)

type CookieManager struct {
	cfg *config.Config
}

func NewCookie(cfg *config.Config) *CookieManager {
	return &CookieManager{cfg: cfg}
}

func (cm *CookieManager) setCookies(c *gin.Context, key string, value string, age int64) {
	c.SetCookie(key, value, int(age), "/", cm.cfg.ClientDomain, false, true)
}

func (cm *CookieManager) SetAccessToken(c *gin.Context, token string, age int64) {
	cm.setCookies(c, AUTH_HEADER, token, age)
}

func (cm *CookieManager) SetRefreshToken(c *gin.Context, token string, age int64) {
	cm.setCookies(c, REFRESH_TOKEN_HEADER, token, age)
}

func getValueFromCookie(c *gin.Context, key string) (string, error) {
	value, err := c.Cookie(key)
	if err != nil {
		return "", err
	}
	return value, nil
}

func GetAccessToken(c *gin.Context) (string, error) {
	return getValueFromCookie(c, AUTH_HEADER)

}
func GetRefreshToken(c *gin.Context) (string, error) {
	return getValueFromCookie(c, REFRESH_TOKEN_HEADER)
}
