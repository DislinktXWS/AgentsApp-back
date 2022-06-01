package company_store

import (
	"modules/dto"
	"modules/utils"
)

func UserMapper(userReq *dto.RequestUser) User {
	user := User{
		ID:          userReq.ID,
		Email:       userReq.Email,
		Username:    userReq.Username,
		Password:    utils.HashPassword(userReq.Password),
		Birthday:    userReq.Birthday,
		FirstName:   userReq.FirstName,
		Gender:      Gender(userReq.Gender),
		LastName:    userReq.LastName,
		PhoneNumber: userReq.PhoneNumber,
		Role:        Role(userReq.Role),
	}
	return user
}
