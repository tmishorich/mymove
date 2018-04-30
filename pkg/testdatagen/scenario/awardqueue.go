package scenario

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// RunAwardQueueScenario1 creates 17 shipments and 5 TSPs in 1 TDL. This allows testing against
// award queue to ensure it behaves as expected. This doesn't track blackout dates.
func RunAwardQueueScenario1(db *pop.Connection) {
	shipmentsToMake := 17

	// Make a TDL to contain our tests
	tdl, _ := testdatagen.MakeTDL(db, "california", "90210", "2")

	// Make a market
	market := "dHHG"

	// Make a source GBLOC
	sourceGBLOC := "OHAI"

	// Make shipments in this TDL
	for i := 0; i < shipmentsToMake; i++ {
		testdatagen.MakeShipment(db, time.Now(), time.Now(), time.Now(), tdl, sourceGBLOC, &market)
	}

	// Make TSPs in the same TDL to handle these shipments
	tsp1, _ := testdatagen.MakeTSP(db, "Excellent TSP", testdatagen.RandomSCAC())
	tsp2, _ := testdatagen.MakeTSP(db, "Pretty Good TSP", testdatagen.RandomSCAC())
	tsp3, _ := testdatagen.MakeTSP(db, "Good TSP", testdatagen.RandomSCAC())
	tsp4, _ := testdatagen.MakeTSP(db, "OK TSP", testdatagen.RandomSCAC())
	tsp5, _ := testdatagen.MakeTSP(db, "Bad TSP", testdatagen.RandomSCAC())

	// TSPs should be ordered by offer_count first, then BVS.
	testdatagen.MakeTSPPerformance(db, tsp1, tdl, swag.Int(1), 5, 0, 4.2, 4.2)
	testdatagen.MakeTSPPerformance(db, tsp2, tdl, swag.Int(1), 4, 0, 3.3, 3.3)
	testdatagen.MakeTSPPerformance(db, tsp3, tdl, swag.Int(2), 3, 0, 2.1, 2.1)
	testdatagen.MakeTSPPerformance(db, tsp4, tdl, swag.Int(3), 2, 0, 1.1, 1.1)
	testdatagen.MakeTSPPerformance(db, tsp5, tdl, swag.Int(4), 1, 0, .5, .5)
}

// RunAwardQueueScenario2 creates 9 shipments to be divided between 5 TSPs in 1 TDL and 10 shipments to be divided among 4 TSPs in TDL 2.
// This allows testing against award queue to ensure it behaves as expected. Two TSPs in TDL1 and one TSP in TDL 2 have blackout dates.
func RunAwardQueueScenario2(db *pop.Connection) {
	shipmentsToMake := 9
	shipmentDate := time.Now()

	// Make a TDL to contain our tests
	tdl, _ := testdatagen.MakeTDL(db, "california", "90210", "2")
	tdl2, _ := testdatagen.MakeTDL(db, "New York", "10024", "2")

	// Make a market
	market := "dHHG"

	// Make a source GBLOC
	sourceGBLOC := "OHAI"

	// Make shipments in first TDL
	for i := 0; i < shipmentsToMake; i++ {
		testdatagen.MakeShipment(db, shipmentDate, shipmentDate, shipmentDate, tdl, sourceGBLOC, &market)
	}
	// Make shipments in second TDL
	for i := 0; i <= shipmentsToMake; i++ {
		testdatagen.MakeShipment(db, shipmentDate, shipmentDate, shipmentDate, tdl2, sourceGBLOC, &market)
	}

	// Make TSPs
	tsp1, _ := testdatagen.MakeTSP(db, "Excellent TSP with Blackout Date", testdatagen.RandomSCAC())
	tsp2, _ := testdatagen.MakeTSP(db, "Very Good TSP", testdatagen.RandomSCAC())
	tsp3, _ := testdatagen.MakeTSP(db, "Pretty Good TSP", testdatagen.RandomSCAC())
	tsp4, _ := testdatagen.MakeTSP(db, "OK TSP with Blackout Date", testdatagen.RandomSCAC())
	tsp5, _ := testdatagen.MakeTSP(db, "Are you even trying TSP", testdatagen.RandomSCAC())
	tsp6, _ := testdatagen.MakeTSP(db, "Excellent TSP", testdatagen.RandomSCAC())
	tsp7, _ := testdatagen.MakeTSP(db, "Pretty Good TSP with Blackout Date", testdatagen.RandomSCAC())
	tsp8, _ := testdatagen.MakeTSP(db, "OK TSP", testdatagen.RandomSCAC())
	tsp9, _ := testdatagen.MakeTSP(db, "Going out of business TSP", testdatagen.RandomSCAC())

	// Put TSPs in 2 TDLs to handle these shipments
	testdatagen.MakeTSPPerformance(db, tsp1, tdl, swag.Int(1), 5, 0, 4.2, 4.4)
	testdatagen.MakeTSPPerformance(db, tsp2, tdl, swag.Int(1), 4, 0, 3.1, 3.2)
	testdatagen.MakeTSPPerformance(db, tsp3, tdl, swag.Int(2), 3, 0, 2.4, 2.5)
	testdatagen.MakeTSPPerformance(db, tsp4, tdl, swag.Int(3), 2, 0, 1.1, 1.3)
	testdatagen.MakeTSPPerformance(db, tsp5, tdl, swag.Int(4), 1, 0, .5, .8)

	testdatagen.MakeTSPPerformance(db, tsp6, tdl2, swag.Int(1), 5, 0, 4.2, 4.4)
	testdatagen.MakeTSPPerformance(db, tsp7, tdl2, swag.Int(2), 4, 0, 3.1, 3.2)
	testdatagen.MakeTSPPerformance(db, tsp8, tdl2, swag.Int(3), 2, 0, 1.1, 1.3)
	testdatagen.MakeTSPPerformance(db, tsp9, tdl2, swag.Int(4), 1, 0, .5, .8)
	// Add blackout dates
	blackoutStart := shipmentDate.AddDate(0, 0, -3)
	blackoutEnd := shipmentDate.AddDate(0, 0, 3)

	gbloc := "BKAS"
	testdatagen.MakeBlackoutDate(db,
		tsp1,
		blackoutStart,
		blackoutEnd,
		&tdl,
		&gbloc,
		&market)
	testdatagen.MakeBlackoutDate(db,
		tsp4,
		blackoutStart,
		blackoutEnd,
		&tdl,
		&gbloc,
		&market)
	testdatagen.MakeBlackoutDate(db,
		tsp7,
		blackoutStart,
		blackoutEnd,
		&tdl,
		&gbloc,
		&market)
}
