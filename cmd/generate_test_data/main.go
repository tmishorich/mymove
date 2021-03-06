package main

import (
	"log"

	"github.com/gobuffalo/pop"
	"github.com/namsral/flag"
	"github.com/transcom/mymove/pkg/testdatagen"
	tdgs "github.com/transcom/mymove/pkg/testdatagen/scenario"
)

// Hey, refactoring self: you can pull the UUIDs from the objects rather than
// querying the db for them again.
func main() {
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, configures the database, presently.")
	rounds := flag.String("rounds", "none", "If not using premade scenarios: Specify none (no awards), full (1 full round of awards), or half (partial round of awards)")
	numTSP := flag.Int("numTSP", 15, "If not using premade scenarios: Specify the number of TSPs you'd like to create")
	scenario := flag.Int("scenario", 0, "Specify which scenario you'd like to run. Current options: 1, 2.")
	flag.Parse()

	//DB connection
	pop.AddLookupPaths(*config)
	db, err := pop.Connect(*env)
	if err != nil {
		log.Panic(err)
	}

	if *scenario == 1 {
		tdgs.RunAwardQueueScenario1(db)
	} else if *scenario == 2 {
		tdgs.RunAwardQueueScenario2(db)
	} else if *scenario == 3 {
		tdgs.RunDutyStationScenario3(db)
	} else {
		// Can this be less repetitive without being overly clever?
		testdatagen.MakeTDLData(db)
		testdatagen.MakeTSPs(db, *numTSP)
		testdatagen.MakeShipmentData(db)
		testdatagen.MakeShipmentOfferData(db)
		testdatagen.MakeTSPPerformanceData(db, *rounds)
		testdatagen.MakeBlackoutDateData(db)
		testdatagen.MakeMoveData(db)
		testdatagen.MakeServiceMember(db)
	}
}
