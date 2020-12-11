package users

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mmkader85/bookstore_users-api/domain/users"
	"github.com/mmkader85/bookstore_users-api/services"
	"github.com/mmkader85/bookstore_users-api/utils/errors"
)

func parseUserIDParam(userIDParam string) (int64, *errors.RestErr) {
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		parseErr := errors.BadRequestErr(fmt.Sprintf("user_id '%s' should be a number", userIDParam))
		return 0, parseErr
	}

	return userID, nil
}

func Create(ctx *gin.Context) {
	var user *users.User

	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		restErr := errors.BadRequestErr(fmt.Sprintf("Invalid JSON. %s", err.Error()))
		ctx.JSON(restErr.Status, restErr)
		return
	}

	result, createUserErr := services.CreateUser(user)
	log.Println("Result", result)
	if createUserErr != nil {
		ctx.JSON(createUserErr.Status, createUserErr)
		return
	}

	ctx.JSON(http.StatusCreated, result.CleanResponse())
}

func Get(ctx *gin.Context) {
	userID, err := parseUserIDParam(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	user, getUserErr := services.GetUser(userID)
	if getUserErr != nil {
		ctx.JSON(getUserErr.Status, getUserErr)
		return
	}

	ctx.JSON(http.StatusOK, user.CleanResponse())
}

func Update(ctx *gin.Context) {
	userID, userIDErr := parseUserIDParam(ctx.Param("user_id"))
	if userIDErr != nil {
		ctx.JSON(userIDErr.Status, userIDErr)
		return
	}

	var user users.User
	user.ID = userID

	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		restErr := errors.BadRequestErr(fmt.Sprintf("Invalid JSON. %s", err.Error()))
		ctx.JSON(restErr.Status, restErr)
		return
	}

	isPartial := ctx.Request.Method == http.MethodPatch
	updatedUser, updateErr := services.UpdateUser(isPartial, &user)
	if updateErr != nil {
		ctx.JSON(updateErr.Status, updateErr)
		return
	}

	ctx.JSON(http.StatusOK, updatedUser.CleanResponse())
}

func Delete(ctx *gin.Context) {
	userID, err := parseUserIDParam(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	deleteUserErr := services.DeleteUser(userID)
	if deleteUserErr != nil {
		ctx.JSON(deleteUserErr.Status, deleteUserErr)
		return
	}

	response := map[string]string{"status": "ok", "message": "successfully deleted"}
	ctx.JSON(http.StatusOK, &response)
}

func Search(ctx *gin.Context) {
	status := ctx.Query("status")

	users, searchUserErr := services.Search(status)
	if searchUserErr != nil {
		ctx.JSON(searchUserErr.Status, searchUserErr)
		return
	}

	ctx.JSON(http.StatusOK, &users)
}
