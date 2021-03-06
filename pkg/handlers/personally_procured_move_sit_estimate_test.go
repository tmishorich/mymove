package handlers

import (
	"net/http/httptest"

	"github.com/gobuffalo/uuid"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *HandlerSuite) TestShowPPMSitEstimateHandlerWithDcos() {
	t := suite.T()

	// Given: a TDL, TSP and TSP performance with SITRate for relevant location and date
	tdl, _ := testdatagen.MakeTDL(suite.db, "US68", "5", "D") // Victoria, TX to Salina, KS
	tsp, _ := testdatagen.MakeTSP(suite.db, testdatagen.RandomSCAC())

	suite.mustSave(&models.Tariff400ngZip3{Zip3: "779", RateArea: "US68", BasepointCity: "Victoria", State: "TX", ServiceArea: 748, Region: 6})
	suite.mustSave(&models.Tariff400ngZip3{Zip3: "674", Region: 5, BasepointCity: "Salina", State: "KS", RateArea: "US58", ServiceArea: 320})

	originServiceArea := models.Tariff400ngServiceArea{
		Name:               "Victoria, TX",
		ServiceArea:        748,
		LinehaulFactor:     39,
		ServiceChargeCents: 350,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		SIT185ARateCents:   unit.Cents(1402),
		SIT185BRateCents:   unit.Cents(53),
		SITPDSchedule:      3,
	}
	suite.mustSave(&originServiceArea)

	destServiceArea := models.Tariff400ngServiceArea{
		Name:               "Salina, KS",
		ServiceArea:        320,
		LinehaulFactor:     43,
		ServiceChargeCents: 350,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		SIT185ARateCents:   unit.Cents(1292),
		SIT185BRateCents:   unit.Cents(51),
		SITPDSchedule:      2,
	}
	suite.mustSave(&destServiceArea)

	tspPerformance := models.TransportationServiceProviderPerformance{
		PerformancePeriodStart:          testdatagen.PerformancePeriodStart,
		PerformancePeriodEnd:            testdatagen.PerformancePeriodEnd,
		RateCycleStart:                  testdatagen.PeakRateCycleStart,
		RateCycleEnd:                    testdatagen.PeakRateCycleEnd,
		TrafficDistributionListID:       tdl.ID,
		TransportationServiceProviderID: tsp.ID,
		QualityBand:                     models.IntPointer(1),
		BestValueScore:                  90,
		LinehaulRate:                    50.5,
		SITRate:                         50,
	}
	suite.mustSave(&tspPerformance)

	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/estimates/ppm_sit", nil)
	req = suite.authenticateRequest(req, user)

	params := ppmop.ShowPPMSitEstimateParams{
		HTTPRequest:     req,
		PlannedMoveDate: *fmtDate(testdatagen.DateInsidePeakRateCycle),
		DaysInStorage:   4,
		OriginZip:       "77901",
		DestinationZip:  "67401",
		WeightEstimate:  3000,
	}
	// And: show Queue is queried
	showHandler := ShowPPMSitEstimateHandler(NewHandlerContext(suite.db, suite.logger))
	showResponse := showHandler.Handle(params)

	// Then: Expect a 200 status code
	okResponse := showResponse.(*ppmop.ShowPPMSitEstimateOK)
	sitCost := okResponse.Payload

	// And: Returned SIT cost to be as expected
	expectedSitCost := int64(3060)
	if *sitCost.Estimate != expectedSitCost {
		t.Errorf("Expected move ppm SIT cost to be '%v', instead is '%v'", expectedSitCost, *sitCost.Estimate)
	}
}

func (suite *HandlerSuite) TestShowPPMSitEstimateHandler2cos() {
	t := suite.T()

	// Given: a TDL, TSP and TSP performance with SITRate for relevant location and date
	tdl, _ := testdatagen.MakeTDL(suite.db, "US68", "5", "2") // Victoria, TX to Salina, KS
	tsp, _ := testdatagen.MakeTSP(suite.db, testdatagen.RandomSCAC())

	suite.mustSave(&models.Tariff400ngZip3{Zip3: "779", RateArea: "US68", BasepointCity: "Victoria", State: "TX", ServiceArea: 748, Region: 6})
	suite.mustSave(&models.Tariff400ngZip3{Zip3: "674", Region: 5, BasepointCity: "Salina", State: "KS", RateArea: "US58", ServiceArea: 320})

	originServiceArea := models.Tariff400ngServiceArea{
		Name:               "Victoria, TX",
		ServiceArea:        748,
		LinehaulFactor:     39,
		ServiceChargeCents: 350,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		SIT185ARateCents:   unit.Cents(1402),
		SIT185BRateCents:   unit.Cents(53),
		SITPDSchedule:      3,
	}
	suite.mustSave(&originServiceArea)

	destServiceArea := models.Tariff400ngServiceArea{
		Name:               "Salina, KS",
		ServiceArea:        320,
		LinehaulFactor:     43,
		ServiceChargeCents: 350,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		SIT185ARateCents:   unit.Cents(1292),
		SIT185BRateCents:   unit.Cents(51),
		SITPDSchedule:      2,
	}
	suite.mustSave(&destServiceArea)

	tspPerformance := models.TransportationServiceProviderPerformance{
		PerformancePeriodStart:          testdatagen.PerformancePeriodStart,
		PerformancePeriodEnd:            testdatagen.PerformancePeriodEnd,
		RateCycleStart:                  testdatagen.PeakRateCycleStart,
		RateCycleEnd:                    testdatagen.PeakRateCycleEnd,
		TrafficDistributionListID:       tdl.ID,
		TransportationServiceProviderID: tsp.ID,
		QualityBand:                     models.IntPointer(1),
		BestValueScore:                  90,
		LinehaulRate:                    50.5,
		SITRate:                         50,
	}
	suite.mustSave(&tspPerformance)

	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/estimates/ppm_sit", nil)
	req = suite.authenticateRequest(req, user)

	params := ppmop.ShowPPMSitEstimateParams{
		HTTPRequest:     req,
		PlannedMoveDate: *fmtDate(testdatagen.DateInsidePeakRateCycle),
		DaysInStorage:   4,
		OriginZip:       "77901",
		DestinationZip:  "67401",
		WeightEstimate:  3000,
	}
	// And: show Queue is queried
	showHandler := ShowPPMSitEstimateHandler(NewHandlerContext(suite.db, suite.logger))
	showResponse := showHandler.Handle(params)

	// Then: Expect a 200 status code
	okResponse := showResponse.(*ppmop.ShowPPMSitEstimateOK)
	sitCost := okResponse.Payload

	// And: Returned SIT cost to be as expected
	expectedSitCost := int64(3060)
	if *sitCost.Estimate != expectedSitCost {
		t.Errorf("Expected move ppm SIT cost to be '%v', instead is '%v'", expectedSitCost, *sitCost.Estimate)
	}
}
