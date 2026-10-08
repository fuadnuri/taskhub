package userRepositoryImpl

import (
	"gorm.io/gorm"
)

type UserRepositoryImp struct {
	db *gorm.DB
}

func NewUserRepositoryImp(db *gorm.DB) *UserRepositoryImp {
	return &UserRepositoryImp{db: db}
}


func (r *userRepositoryImpl) Create(user *Entities.UserEntity) error{
	return r.db.Create(user).Error
}

func (r *userRepositoryImpl) GetByID(id uint)(*Entities.UserEntity,error){
	var user Entities.UserEntity
	err:=r.db.First(&user,id).Error
	if erro!=nil{
		return nil, err
	}
	return &user,nil
}

func (r *UserRepositoryImp) GetAll() ([]Entities.UserEntity,error){
	var users []Entities.UserEntity

	if err:= r.db.Find(&users).Error; err!=nil{
		return nil, err
	}else{
		return users, nil
	}
}