package handlers

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/auth"
	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *HandlerSuite) TestCreatePPMHandler() {
	t := suite.T()

	user1 := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "whoever@example.com",
	}

	user2 := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "someoneelse@example.com",
	}

	verrs, err := suite.db.ValidateAndCreate(&user1)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}
	verrs, err = suite.db.ValidateAndCreate(&user2)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}
	var selectedType = internalmessages.SelectedMoveTypeCOMBO
	move := models.Move{
		UserID:           user1.ID,
		SelectedMoveType: &selectedType,
	}
	verrs, err = suite.db.ValidateAndCreate(&move)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}

	fmt.Println(user1.ID, user2.ID, move.UserID)

	request := httptest.NewRequest("POST", "/fake/path", nil)
	ctx := request.Context()
	ctx = auth.PopulateUserModel(ctx, user1)
	request = request.WithContext(ctx)

	newPPMPayload := internalmessages.CreatePersonallyProcuredMovePayload{WeightEstimate: swag.Int64(12), PickupZip: swag.String("00112"), DaysInStorage: swag.Int64(3)}

	newPPMParams := ppmop.CreatePersonallyProcuredMoveParams{
		MoveID: strfmt.UUID(move.ID.String()),
		CreatePersonallyProcuredMovePayload: &newPPMPayload,
		HTTPRequest:                         request,
	}

	handler := CreatePersonallyProcuredMoveHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(newPPMParams)
	// assert we got back the 201 response
	createdResponse := response.(*ppmop.CreatePersonallyProcuredMoveCreated)
	createdIssuePayload := createdResponse.Payload

	fmt.Println(createdIssuePayload)

	// Next try the wrong user
	ctx = auth.PopulateUserModel(ctx, user2)
	request = request.WithContext(ctx)
	newPPMParams.HTTPRequest = request

	badUserResponse := handler.Handle(newPPMParams)
	suite.checkResponseForbidden(badUserResponse)

	// Now try a bad move
	newPPMParams.MoveID = strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	badMoveResponse := handler.Handle(newPPMParams)
	suite.checkResponseNotFound(badMoveResponse)

}

func (suite *HandlerSuite) TestIndexPPMHandler() {
	t := suite.T()

	user1 := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "whoever@example.com",
	}

	verrs, err := suite.db.ValidateAndCreate(&user1)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}
	var selectedType = internalmessages.SelectedMoveTypeCOMBO
	move := models.Move{
		UserID:           user1.ID,
		SelectedMoveType: &selectedType,
	}
	verrs, err = suite.db.ValidateAndCreate(&move)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}

	move2 := models.Move{
		UserID:           user1.ID,
		SelectedMoveType: &selectedType,
	}
	verrs, err = suite.db.ValidateAndCreate(&move2)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}

	ppm1 := models.PersonallyProcuredMove{
		MoveID:             move.ID,
		WeightEstimate:     swag.Int64(1),
		EstimatedIncentive: swag.String("$2681.25 - $4111.25"),
	}
	ppm2 := models.PersonallyProcuredMove{
		MoveID:         move.ID,
		WeightEstimate: swag.Int64(2),
	}
	otherPPM := models.PersonallyProcuredMove{
		MoveID:         move2.ID,
		WeightEstimate: swag.Int64(4),
	}

	verrs, err = suite.db.ValidateAndCreate(&ppm1)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}

	verrs, err = suite.db.ValidateAndCreate(&ppm2)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}

	verrs, err = suite.db.ValidateAndCreate(&otherPPM)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}

	request := httptest.NewRequest("GET", "/fake/path", nil)
	ctx := request.Context()
	ctx = auth.PopulateUserModel(ctx, user1)
	request = request.WithContext(ctx)

	indexPPMParams := ppmop.IndexPersonallyProcuredMovesParams{
		MoveID:      strfmt.UUID(move.ID.String()),
		HTTPRequest: request,
	}

	handler := IndexPersonallyProcuredMovesHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(indexPPMParams)

	// assert we got back the 201 response
	okResponse := response.(*ppmop.IndexPersonallyProcuredMovesOK)
	indexPPMPayload := okResponse.Payload

	fmt.Println(indexPPMPayload)
	for _, ppm := range indexPPMPayload {
		if *ppm.ID == *fmtUUID(otherPPM.ID) {
			t.Error("We should only have got back ppms associated with this move")
		}
	}
	if len(indexPPMPayload) != 2 {
		t.Error("We should have gotten back two good ones. ")
	}

}

func (suite *HandlerSuite) TestPatchPPMHandler() {
	t := suite.T()
	initialSize := internalmessages.TShirtSize("S")
	newSize := internalmessages.TShirtSize("L")
	initialWeight := swag.Int64(1)
	newWeight := swag.Int64(5)
	initialMoveDate := time.Now().Add(-2 * 24 * time.Hour)
	newMoveDate := time.Now()
	destinationZip := swag.String("00112")

	user1 := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "whoever@example.com",
	}
	suite.mustSave(&user1)
	var selectedType = internalmessages.SelectedMoveTypeCOMBO
	move := models.Move{
		UserID:           user1.ID,
		SelectedMoveType: &selectedType,
	}
	suite.mustSave(&move)

	ppm1 := models.PersonallyProcuredMove{
		MoveID:          move.ID,
		Size:            &initialSize,
		WeightEstimate:  initialWeight,
		PlannedMoveDate: &initialMoveDate,
		DestinationZip:  destinationZip,
	}
	suite.mustSave(&ppm1)

	request := httptest.NewRequest("GET", "/fake/path", nil)
	ctx := request.Context()
	ctx = auth.PopulateUserModel(ctx, user1)
	request = request.WithContext(ctx)

	payload := internalmessages.PatchPersonallyProcuredMovePayload{
		Size:            &newSize,
		WeightEstimate:  newWeight,
		PlannedMoveDate: fmtDatePtr(&newMoveDate),
	}

	patchPPMParams := ppmop.PatchPersonallyProcuredMoveParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(move.ID.String()),
		PersonallyProcuredMoveID:           strfmt.UUID(ppm1.ID.String()),
		PatchPersonallyProcuredMovePayload: &payload,
	}

	handler := PatchPersonallyProcuredMoveHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(patchPPMParams)

	// assert we got back the 201 response
	okResponse := response.(*ppmop.PatchPersonallyProcuredMoveCreated)
	patchPPMPayload := okResponse.Payload

	fmt.Println(patchPPMPayload)
	if *patchPPMPayload.Size != newSize {
		t.Error("Size should have been updated.")
	}

	if patchPPMPayload.WeightEstimate != newWeight {
		t.Error("Weight should have been updated.")
	}

	if !(*time.Time)(patchPPMPayload.PlannedMoveDate).Equal(newMoveDate) {
		t.Error("MoveDate should have been updated.")
	}

	if *patchPPMPayload.DestinationZip != *destinationZip {
		t.Error("DestinationZip should have been updated.")
	}
}

func (suite *HandlerSuite) TestPatchPPMHandlerWrongUser() {
	initialSize := internalmessages.TShirtSize("S")
	newSize := internalmessages.TShirtSize("L")
	initialWeight := swag.Int64(1)
	newWeight := swag.Int64(5)
	initialMoveDate := time.Now().Add(-2 * 24 * time.Hour)
	newMoveDate := time.Now()

	user1 := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "whoever@example.com",
	}
	suite.mustSave(&user1)

	user2 := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "whoever@example.com",
	}
	suite.mustSave(&user2)
	var selectedType = internalmessages.SelectedMoveTypeCOMBO
	move := models.Move{
		UserID:           user1.ID,
		SelectedMoveType: &selectedType,
	}
	suite.mustSave(&move)

	ppm1 := models.PersonallyProcuredMove{
		MoveID:          move.ID,
		Size:            &initialSize,
		WeightEstimate:  initialWeight,
		PlannedMoveDate: &initialMoveDate,
	}
	suite.mustSave(&ppm1)

	request := httptest.NewRequest("GET", "/fake/path", nil)
	ctx := request.Context()
	ctx = auth.PopulateUserModel(ctx, user2)
	request = request.WithContext(ctx)

	payload := internalmessages.PatchPersonallyProcuredMovePayload{
		Size:            &newSize,
		WeightEstimate:  newWeight,
		PlannedMoveDate: fmtDatePtr(&newMoveDate),
	}

	patchPPMParams := ppmop.PatchPersonallyProcuredMoveParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(move.ID.String()),
		PersonallyProcuredMoveID:           strfmt.UUID(ppm1.ID.String()),
		PatchPersonallyProcuredMovePayload: &payload,
	}

	handler := PatchPersonallyProcuredMoveHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(patchPPMParams)

	suite.checkResponseForbidden(response)
}

func (suite *HandlerSuite) TestPatchPPMHandlerWrongMoveID() {
	initialSize := internalmessages.TShirtSize("S")
	newSize := internalmessages.TShirtSize("L")
	initialWeight := swag.Int64(1)
	newWeight := swag.Int64(5)

	user1 := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "whoever@example.com",
	}
	suite.mustSave(&user1)

	user2 := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "whoever@example.com",
	}
	suite.mustSave(&user2)
	var selectedType = internalmessages.SelectedMoveTypeCOMBO
	move := models.Move{
		UserID:           user1.ID,
		SelectedMoveType: &selectedType,
	}
	suite.mustSave(&move)

	move2 := models.Move{
		UserID:           user2.ID,
		SelectedMoveType: &selectedType,
	}
	suite.mustSave(&move2)

	ppm1 := models.PersonallyProcuredMove{
		MoveID:         move2.ID,
		Size:           &initialSize,
		WeightEstimate: initialWeight,
	}
	suite.mustSave(&ppm1)

	request := httptest.NewRequest("GET", "/fake/path", nil)
	ctx := request.Context()
	ctx = auth.PopulateUserModel(ctx, user1)
	request = request.WithContext(ctx)

	payload := internalmessages.PatchPersonallyProcuredMovePayload{
		Size:           &newSize,
		WeightEstimate: newWeight,
	}

	patchPPMParams := ppmop.PatchPersonallyProcuredMoveParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(move.ID.String()),
		PersonallyProcuredMoveID:           strfmt.UUID(ppm1.ID.String()),
		PatchPersonallyProcuredMovePayload: &payload,
	}

	handler := PatchPersonallyProcuredMoveHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(patchPPMParams)

	suite.checkResponseForbidden(response)

}

func (suite *HandlerSuite) TestPatchPPMHandlerNoMove() {
	t := suite.T()

	initialSize := internalmessages.TShirtSize("S")
	newSize := internalmessages.TShirtSize("L")
	initialWeight := swag.Int64(1)
	newWeight := swag.Int64(5)

	user1 := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "whoever@example.com",
	}
	suite.mustSave(&user1)

	user2 := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "whoever@example.com",
	}
	suite.mustSave(&user2)

	var selectedType = internalmessages.SelectedMoveTypeCOMBO
	move := models.Move{
		UserID:           user1.ID,
		SelectedMoveType: &selectedType,
	}
	suite.mustSave(&move)

	badMoveID := uuid.Must(uuid.NewV4())

	ppm1 := models.PersonallyProcuredMove{
		MoveID:         move.ID,
		Size:           &initialSize,
		WeightEstimate: initialWeight,
	}
	suite.mustSave(&ppm1)

	request := httptest.NewRequest("GET", "/fake/path", nil)
	ctx := request.Context()
	ctx = auth.PopulateUserModel(ctx, user1)
	request = request.WithContext(ctx)

	payload := internalmessages.PatchPersonallyProcuredMovePayload{
		Size:           &newSize,
		WeightEstimate: newWeight,
	}

	patchPPMParams := ppmop.PatchPersonallyProcuredMoveParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(badMoveID.String()),
		PersonallyProcuredMoveID:           strfmt.UUID(ppm1.ID.String()),
		PatchPersonallyProcuredMovePayload: &payload,
	}

	handler := PatchPersonallyProcuredMoveHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(patchPPMParams)

	// assert we got back the badrequest response
	_, ok := response.(*ppmop.PatchPersonallyProcuredMoveBadRequest)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}

}
