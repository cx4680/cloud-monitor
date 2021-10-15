package auth

/*
import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Get params from request and to authenticate
func HandleAuthenticateRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetHeader(configuration.UIN)
		log.Debug("HandleAuthenticateRequest,user id is ", userId)
		identity, err := GetCasBinAuthInstance().HandleAuthenticate(userId, nil)
		if identity == "" || err != nil {
			c.Status(http.StatusUnauthorized)
			c.Abort()
		}
		c.Next()
	}
}

// Get params from request and to  authorize
func HandleAuthorizeRequest(ctx *gin.Context, objectType string, operation string) error {
	log.Debug("HandleAuthorizeRequest, object and operation is ", objectType, operation)
	userId := getUserIdFromRequest(ctx)
	if userId == "" {
		return errs.New(errs.UNAUTHORIZED, "user authorize failed,can not get user id")
	}
	merchantId := getMerchantIdFromRequest(ctx)
	if merchantId == "" {
		return errs.New(errs.UNAUTHORIZED, "user authorize failed,can not get merchant id")
	}
	permission := getPermission(objectType, operation, merchantId)
	success, err := GetCasBinAuthInstance().HandleAuthorize(userId, permission)
	if err != nil || !success {
		log.Warn("user authorize failed, userId=", userId, ", err=", err)
		return errs.New(errs.UNAUTHORIZED, "user authorize failed")
	}
	return nil
}
*/
