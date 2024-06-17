package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gospotify.com/db"
	"gospotify.com/models"
	"gospotify.com/utils"
)

func TestPasswordHash(t *testing.T) {
	password := "123"
	hash, err := utils.HashSha256(password)
	if err != nil {
		t.Error(err)
	}
	t.Log(hash)
	assert.Nil(t, err)
}

func TestPasswordHashAreSame(t *testing.T) {
	passwordOne := "3022003"
	hashOne, err := utils.HashSha256(passwordOne)
	if err != nil {
		t.Error(err)
	}
	passwordTwo := "3022003"
	hashTwo, err := utils.HashSha256(passwordTwo)
	if err != nil {
		t.Error(err)
	}
	t.Logf(fmt.Sprintf("Hash1: %s, Hash2: %s", hashOne, hashTwo))
	assert.Equal(t, hashTwo, hashOne)
}

func TestUserPasswordHashWithSalt(t *testing.T) {
	password := "3022003"
	salt := make([]byte, 32)
	rand.Read(salt)
	passwordSalt, err := utils.BytesToHex(salt)
	if err != nil {
		t.Error(err)
	}
	passwordHash, err := utils.HashSha256(password + passwordSalt)
	if err != nil {
		t.Error(err)
	}

	testUser := models.RegisterUserDb{
		Username: "test1",
		Password: passwordHash,
		Salt:     passwordSalt,
	}
	client, err := db.DbClient()
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	if err != nil {
		t.Error(err)
	}

	cursor, err := db.Db.Collection("users").InsertOne(context.TODO(), testUser)
	if err != nil {
		t.Error(err)
	}

	t.Log(cursor.InsertedID)
	assert.NotNil(t, cursor)
}
