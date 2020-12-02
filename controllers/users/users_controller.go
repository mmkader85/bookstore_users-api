package users

import (
	"github.com/gin-gonic/gin"
	"github.com/mmkader85/bookstore_users-api/domain/users"
	"github.com/mmkader85/bookstore_users-api/services"
	"github.com/mmkader85/bookstore_users-api/utils/errors"
	"net/http"
	"strconv"
)

func CreateUser(ctx *gin.Context) {
	var user users.User
	// body, err := ioutil.ReadAll(ctx.Request.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	//
	// if err := json.Unmarshal(body, &user); err != nil {
	// 	fmt.Println("Error during unmarshall:", err)
	// 	return
	// }

	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		restErr := errors.BadRequestErr("invalid json")
		ctx.JSON(restErr.Status, restErr)

		return
	}

	result, createUserErr := services.CreateUser(user)
	if createUserErr != nil {
		ctx.JSON(createUserErr.Status, createUserErr)
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

func GetUser(ctx *gin.Context) {
	userID, err := strconv.ParseInt(ctx.Param("user_id"), 10, 64)
	if err != nil {
		parseErr := errors.BadRequestErr("user_id should be a number")
		ctx.JSON(parseErr.Status, parseErr)

		return
	}

	user, getUserErr := services.GetUser(userID)
	if getUserErr != nil {
		ctx.JSON(getUserErr.Status, getUserErr)

		return
	}

	ctx.JSON(http.StatusOK, user)
}
