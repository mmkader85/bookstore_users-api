package users

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mmkader85/bookstore_users-api/domain/users"
	"github.com/mmkader85/bookstore_users-api/services"
	"github.com/mmkader85/bookstore_users-api/utils"
)

func parseUserIDParam(userIDParam string) (int64, *utils.RestErr) {
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		parseErr := utils.RestErrUtils.BadRequestErr(fmt.Sprintf("user_id '%s' should be a number", userIDParam))
		return 0, parseErr
	}

	return userID, nil
}

func Create(ctx *gin.Context) {
	var user *users.User

	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		restErr := utils.RestErrUtils.BadRequestErr(fmt.Sprintf("Invalid JSON. %s", err.Error()))
		ctx.JSON(restErr.Status, restErr)
		return
	}

	result, createUserErr := services.UserService.Create(user)
	log.Println("Result", result)
	if createUserErr != nil {
		ctx.JSON(createUserErr.Status, createUserErr)
		return
	}

	var isPublic = ctx.GetHeader("X-Private") != "true"
	ctx.JSON(http.StatusCreated, result.Marshal(isPublic))
}

func Get(ctx *gin.Context) {
	userID, err := parseUserIDParam(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	user, getUserErr := services.UserService.Get(userID)
	if getUserErr != nil {
		ctx.JSON(getUserErr.Status, getUserErr)
		return
	}

	var isPublic = ctx.GetHeader("X-Private") != "true"
	ctx.JSON(http.StatusOK, user.Marshal(isPublic))
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
		restErr := utils.RestErrUtils.BadRequestErr(fmt.Sprintf("Invalid JSON. %s", err.Error()))
		ctx.JSON(restErr.Status, restErr)
		return
	}

	isPartial := ctx.Request.Method == http.MethodPatch
	updatedUser, updateErr := services.UserService.Update(isPartial, &user)
	if updateErr != nil {
		ctx.JSON(updateErr.Status, updateErr)
		return
	}

	var isPublic = ctx.GetHeader("X-Private") != "true"
	ctx.JSON(http.StatusOK, updatedUser.Marshal(isPublic))
}

func Delete(ctx *gin.Context) {
	userID, err := parseUserIDParam(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	deleteUserErr := services.UserService.Delete(userID)
	if deleteUserErr != nil {
		ctx.JSON(deleteUserErr.Status, deleteUserErr)
		return
	}

	response := map[string]string{"status": "ok", "message": "successfully deleted"}
	ctx.JSON(http.StatusOK, &response)
}

func Search(ctx *gin.Context) {
	status := ctx.Query("status")

	results, searchUserErr := services.UserService.Search(status)
	if searchUserErr != nil {
		ctx.JSON(searchUserErr.Status, searchUserErr)
		return
	}

	var isPublic = ctx.GetHeader("X-Private") != "true"
	var output = make([]interface{}, len(results))
	for index, user := range results {
		output[index] = user.Marshal(isPublic)
	}

	ctx.JSON(http.StatusOK, output)
}
