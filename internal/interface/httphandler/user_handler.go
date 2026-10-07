package httpHandler


import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"

	"github.com/fuadroid/taskhub/internal/core/domain/entities"
	"github.com/fuadroid/taskhub/internal/core/application/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler{
	return &UserHandler{
		service:service,
	}
}

func (h *UserHandler) Create(c *gin.Context){
	var user entities.UserEntity
	if err:=c.ShouldBind(&user); err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":err.Error(),
		})
		return 
	}

	if err:=h.service.CreateUser(&user);err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error",err.Error(),
		})
		return
	}
}


func (h *UserHandler)GetAllUsers(c *gin.Context){
	
	if users,err:=h.service.GetAll().Error; err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"internal server error",
		})
		return 
	}else{
		c.JSON(http.StatusOk,user)
	}
}

//the rest of the user related endpionts goes here