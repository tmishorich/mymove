package handlers

import (
	"github.com/gobuffalo/uuid"

	shipmentop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/shipments"
	internalmessages "github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *HandlerSuite) TestIndexShipmentsHandler() {
	t := suite.T()

	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	tsp := models.TransportationServiceProvider{
		StandardCarrierAlphaCode: "scac",
		Name: "Transportation Service Provider 1",
	}
	suite.mustSave(&tsp)

	tdl := models.TrafficDistributionList{
		CodeOfService:     "cos",
		DestinationRegion: "dr",
		SourceRateArea:    "sra",
	}
	suite.mustSave(&tdl)

	selectedType := internalmessages.SelectedMoveTypeHHG
	move := models.Move{
		UserID:           user.ID,
		SelectedMoveType: &selectedType,
	}
	suite.mustSave(&move)

	avs := models.Shipment{
		TrafficDistributionListID: tdl.ID,
		SourceGBLOC:               "AGFM",
		MoveID:                    move.ID,
	}
	suite.mustSave(&avs)

	aws := models.Shipment{
		TrafficDistributionListID: tdl.ID,
		SourceGBLOC:               "AGFM",
		MoveID:                    move.ID,
	}
	suite.mustSave(&aws)

	offer := models.ShipmentOffer{
		ShipmentID:                      aws.ID,
		TransportationServiceProviderID: tsp.ID,
	}
	suite.mustSave(&offer)

	params := shipmentop.NewIndexShipmentsParams()
	indexHandler := IndexShipmentsHandler(NewHandlerContext(suite.db, suite.logger))
	indexResponse := indexHandler.Handle(params)

	okResponse, ok := indexResponse.(*shipmentop.IndexShipmentsOK)
	if !ok {
		t.Fatalf("Request failed: %#v", indexResponse)
	}
	shipments := okResponse.Payload

	if len(shipments) != 2 {
		t.Errorf("expected %d shipments, got %d", 2, len(shipments))
	}

	offeredCount := 0
	availableCount := 0
	for _, shipment := range shipments {
		if shipment.TransportationServiceProviderID != nil {
			offeredCount++
			if shipment.TransportationServiceProviderID.String() != tsp.ID.String() {
				t.Errorf("got wrong tsp id, expected %s, got %s", tsp.ID.String(), shipment.TransportationServiceProviderID.String())

			}
		} else {
			availableCount++
		}
	}

	if offeredCount != 1 {
		t.Errorf("expected %d offered shipments, got %d", 1, offeredCount)
	}

	if availableCount != 1 {
		t.Errorf("expected %d available shipments, got %d", 1, availableCount)
	}
}
