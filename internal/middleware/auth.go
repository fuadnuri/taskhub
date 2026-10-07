package middleware 
import (
	"github.com/gin-gonic/gin"
)

type User struct{
	Token string `json:"token"`
}
func AuthMiddleware(ctx gin.Context){
	var token string
	token := json.Marshal(ctx.Request)
	if token.Isvalid{
		return ctx.Next()
	}
	return ctx.(401,gin)

}

//add auth middleware 