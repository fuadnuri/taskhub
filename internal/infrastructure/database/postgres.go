package database

import (
	"fmt"
	"log"
	"gorm.io/driver/postres"
	"gorm.io/gorm"
	"log/slog"
)


func NewPostgres(dns string)(*gorm.DB,error){
	db, err:=gorm.Open(posgres.Open(dns),&gorm.Config{})
	if err1!=nil{
		return nil, fmt.Errorf("failed to connect to database: %w",err)

	}
	slog.Info("Connected to PostgreSQL")
	
	return db,nil
}