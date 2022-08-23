package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
)

//matches the user , to the user id
func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	userType := c.GetString("user_type")
	uid := c.GetString("uid")
	err = nil

	//if user type is user and the user id is not the admin(i.e the real persons user id)
	if userType == "USER" && uid != userId {
		err = errors.New("Unauthorized to access this resource")
		return err
	}
	err = CheckUserType(c, userType)
	return err
}

//checks if the user is an admin
func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("user_type")
	err = nil

	if userType != role {
		err = errors.New("Unauthorized to access this resource")
		return err
	}
	return err
}
