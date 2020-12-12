package users

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/mmkader85/bookstore_users-api/utils/errors"
)

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" binding:"required"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
	Password  string `json:"password"`
}

type PublicUser struct {
	ID        int64  `json:"id"`
	Email     string `json:"email" binding:"required"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
}

type PrivateUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" binding:"required"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
}

func (u *User) Validate() *errors.RestErr {
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)

	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	if u.Email == "" {
		return errors.BadRequestErr("Empty user email.")
	}

	u.Status = strings.TrimSpace(strings.ToLower(u.Status))
	if u.Status != "active" && u.Status != "suspended" && u.Status != "deleted" {
		return errors.BadRequestErr("User status should be active/suspended/deleted only.")
	}

	u.Password = strings.TrimSpace(strings.ToLower(u.Password))
	if u.Password == "" {
		return errors.BadRequestErr("Empty user password.")
	}

	return nil
}

// Marshall can be done manually or thru json.Marshal().
// This method demonstrates both.
func (u User) Marshal(isPublic bool) interface{} {
	if isPublic {
		return &PublicUser{
			ID:        u.ID,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
			Status:    u.Status,
		}
	}

	user, err := json.Marshal(u)
	if err != nil {
		log.Fatal("Unable to marshall user.")
	}

	var privateUser PrivateUser
	if err = json.Unmarshal(user, &privateUser); err != nil {
		log.Fatal("Unable to unmarshall user.")
	}

	return privateUser
}

func (u *User) CleanResponse() map[string]string {
	var response = map[string]string{
		"ID":        strconv.FormatInt(u.ID, 10),
		"FirstName": u.FirstName,
		"LastName":  u.LastName,
		"Email":     u.Email,
		"CreatedAt": u.CreatedAt,
		"Status":    u.Status,
	}

	return response
}
