package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	//1. connect to postgresql
	dns:="host:localhost user=postgres password=password dbname=myapp port=5432 sslmode=disabled"
	db,err:=database.NewPostgres(dns)
	if err!=nil{
		slog.Error("%s %s",time.Now(),err)
	}
	//run migration
	if err:=db.Automigrate(&entities.UserEntity{});err!=nil{
		slog.Error("%s : %s",time.Now(),err)

	}
	//create repository
	userRepository:=repository.NewUserRepositoryImp{db}

	// create application service
	userService:=service.NewUserService(userRepository)
	
	//create handler
	userHandler:=httpHandler.NewUserHandler(userService)

	//create route
	route:=gin.Default()
	router.Post("/users",userHandler.Create)
	router.Get("/users/:id",userHandler.GetAll)

	
	godotenv.Load()
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})
	router.Run(":8080")
}
