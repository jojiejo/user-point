package tests

import (
	"log"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jojiejo/user-point/api/models"
	"github.com/stretchr/testify/assert"
)

func TestSaveUserPoint(t *testing.T) {

	err := refreshUserAndUserPointTable()
	if err != nil {
		log.Fatalf("Error refreshing user & user point tables: %v\n", err)
	}

	user, err := seedUser()
	if err != nil {
		log.Fatalf("Error seeding user table: %v\n", err)
	}

	userPoint := models.UserPoint{
		Value:  20,
		UserID: user.ID,
	}

	createdUserPoint, err := userPoint.CreateUserPoint(server.DB)
	if err != nil {
		t.Errorf("Error fetching user point: %v\n", err)
		return
	}

	assert.Equal(t, userPoint.Value, createdUserPoint.Value)
	assert.Equal(t, userPoint.UserID, createdUserPoint.UserID)
}

func TestFindPostByID(t *testing.T) {
	err := refreshUserAndUserPointTable()
	if err != nil {
		log.Fatalf("Error refreshing user & user point tables: %v\n", err)
	}

	_, userPoint, err := seedUserAndUserPoint()
	if err != nil {
		log.Fatalf("Error seeding user & user point table: %v\n", err)
	}

	foundUserPoint, err := userPoint.FindPointHistoryByUserID(server.DB, userPoint.UserID)
	if err != nil {
		t.Errorf("Error fetching user point: %v\n", err)
		return
	}

	assert.Equal(t, 1, len(*foundUserPoint))
}
