package testdatagen

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/markbates/pop"

	"github.com/transcom/mymove/pkg/models"
)

// No test includes a zip3 or a volume_move value, as we're not
// currently fully implementing those.

// MakeBlackoutDate creates a test blackoutDate object to add to the database.
func MakeBlackoutDate(db *pop.Connection, tsp models.TransportationServiceProvider,
	tdl *models.TrafficDistributionList, cos *string, channel *string, gbloc *string,
	market *string) (models.BlackoutDate, error) {

	blackoutDates := models.BlackoutDate{
		TransportationServiceProviderID: tsp.ID,
		StartBlackoutDate:               time.Now(),
		EndBlackoutDate:                 time.Now(),
		TrafficDistributionListID:       &tdl.ID,
		CodeOfService:                   cos,
		Channel:                         channel,
		GBLOC:                           gbloc,
		Market:                          market,
	}

	_, err := db.ValidateAndSave(&blackoutDates)
	if err != nil {
		log.Panic(err)
	}

	return blackoutDates, err
}

// MakeBlackoutDateData creates three blackoutDate objects and commits them to the blackout_dates table.
func MakeBlackoutDateData(db *pop.Connection) {
	// These two queries duplicate ones in other testdatagen files; not optimal
	tspList := []models.TransportationServiceProvider{}
	err := db.All(&tspList)
	if err != nil {
		fmt.Println("TSP ID import failed.")
	}

	tdlList := []models.TrafficDistributionList{}
	err = db.All(&tdlList)
	if err != nil {
		fmt.Println("TDL ID import failed.")
	}

	conus := "CONUS_CONUS"
	market := "market"
	cos := "BKAS"
	channel := "dHHG"

	// Make a blackout date with market and channel.
	MakeBlackoutDate(db,
		tspList[rand.Intn(len(tspList))],
		&tdlList[rand.Intn(len(tdlList))],
		nil,
		&conus,
		nil,
		&market,
	)

	// Make a blackout date with a channel.
	MakeBlackoutDate(db,
		tspList[rand.Intn(len(tspList))],
		&tdlList[rand.Intn(len(tdlList))],
		nil,
		&conus,
		nil,
		nil,
	)

	// Make a blackout date with market, GBLOC, and channel.
	MakeBlackoutDate(db,
		tspList[rand.Intn(len(tspList))],
		&tdlList[rand.Intn(len(tdlList))],
		nil,
		&conus,
		&cos,
		&channel,
	)
}