package authorization

import "gospotify.com/models"

type User = models.UserDb

func IsUserAdmin(user *User) bool {
	return user.IsAdmin
}
