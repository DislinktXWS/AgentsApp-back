package company_store

import (
	"modules/dto"
	"modules/utils"
)

func UserMapper(userReq *dto.RequestUser) User {
	user := User{
		Email:       userReq.Email,
		Username:    userReq.Username,
		Password:    utils.HashPassword(userReq.Password),
		DateOfBirth: userReq.DateOfBirth,
		Name:        userReq.Name,
		Gender:      Gender(userReq.Gender),
		Surname:     userReq.Surname,
		Phone:       userReq.Phone,
		Role:        Role(userReq.Role),
	}
	return user
}
