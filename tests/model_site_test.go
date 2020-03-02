package tests

import (
	"strconv"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/mssql" //Ms. SQL driver
	"github.com/stretchr/testify/assert"
)

func TestFindAllSites(t *testing.T) {
	sites, err := siteInstance.FindAllSites(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the sites: %v\n", err)
		return
	}

	assert.Equal(t, len(*sites), 2)
}

func TestFindUserByID(t *testing.T) {
	inputSiteID := "14"
	convertedSiteID, err := strconv.ParseUint(inputSiteID, 10, 64)
	foundSite, err := siteInstance.FindSiteByID(server.DB, convertedSiteID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundSite.ID, convertedSiteID)
}
