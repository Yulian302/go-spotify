package auth

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
