package tests

import (
	"log"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jojiejo/user-point/api/models"
	"github.com/stretchr/testify/assert"
)

func TestFindAllUsers(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	_, err = seedUsers()
	if err != nil {
		log.Fatal(err)
	}

	users, err := userInstance.FindAllUsers(server.DB)
	if err != nil {
		t.Errorf("Error fetching users: %v\n", err)
		return
	}

	assert.Equal(t, len(*users), 2)
}

func TestSaveUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	newUser := models.User{
		Email: "djodi@example.com",
	}

	savedUser, err := newUser.CreateUser(server.DB)
	if err != nil {
		t.Errorf("Error fetching users: %v\n", err)
		return
	}

	assert.Equal(t, newUser.ID, savedUser.ID)
	assert.Equal(t, newUser.Email, savedUser.Email)
}

func TestFindUserByID(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedUser()
	if err != nil {
		log.Fatalf("Cannot seed user table: %v", err)
	}

	foundUser, err := userInstance.FindUserByID(server.DB, user.ID)
	if err != nil {
		t.Errorf("Error fetching users: %v\n", err)
		return
	}

	assert.Equal(t, foundUser.ID, user.ID)
	assert.Equal(t, foundUser.Email, user.Email)
}

func TestUpdateUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	_, err = seedUser()
	if err != nil {
		log.Fatalf("Cannot seed user table: %v\n", err)
	}

	userUpdate := models.User{
		ID:           1,
		Email:        "djodiupdate@example.com",
		CurrentPoint: 50,
	}

	updatedUser, err := userUpdate.UpdateUserPoint(server.DB)
	if err != nil {
		t.Errorf("Error updating users: %v\n", err)
		return
	}

	assert.Equal(t, updatedUser.ID, userUpdate.ID)
	assert.Equal(t, updatedUser.Email, userUpdate.Email)
}

func TestDeleteAUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedUser()
	if err != nil {
		log.Fatalf("Cannot seed user table: %v\n", err)
	}

	isDeleted, err := user.DeleteUser(server.DB)
	if err != nil {
		t.Errorf("Error deleting users: %v\n", err)
		return
	}

	assert.Equal(t, isDeleted, int64(1))
}
