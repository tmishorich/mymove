package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth/context"
	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForPPMModel(personallyProcuredMove models.PersonallyProcuredMove) internalmessages.PersonallyProcuredMovePayload {
	ppmPayload := internalmessages.PersonallyProcuredMovePayload{
		ID:             fmtUUID(personallyProcuredMove.ID),
		CreatedAt:      fmtDateTime(personallyProcuredMove.CreatedAt),
		UpdatedAt:      fmtDateTime(personallyProcuredMove.UpdatedAt),
		Size:           personallyProcuredMove.Size,
		WeightEstimate: personallyProcuredMove.WeightEstimate,
	}
	return ppmPayload
}

// CreatePersonallyProcuredMoveHandler creates a PPM
type CreatePersonallyProcuredMoveHandler HandlerContext

// Handle is the handler
func (h CreatePersonallyProcuredMoveHandler) Handle(params ppmop.CreatePersonallyProcuredMoveParams) middleware.Responder {
	// User should always be populated by middleware
	user, _ := context.GetUser(params.HTTPRequest.Context())
	// MoveID is validated as a UUID by the swagger validator
	moveID, _ := uuid.FromString(params.MoveID.String())

	move, err := models.FetchMove(h.db, user.ID, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	ppmParams := params.CreatePersonallyProcuredMovePayload
	newPPM, verrs, err := move.CreatePPM(h.db, ppmParams.Size, ppmParams.WeightEstimate)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	ppmPayload := payloadForPPMModel(newPPM)
	return ppmop.NewCreatePersonallyProcuredMoveCreated().WithPayload(&ppmPayload)

}

// IndexPersonallyProcuredMovesHandler returns a list of all the PPMs associated with this move.
type IndexPersonallyProcuredMovesHandler HandlerContext

// Handle handles the request
func (h IndexPersonallyProcuredMovesHandler) Handle(params ppmop.IndexPersonallyProcuredMovesParams) middleware.Responder {
	var response middleware.Responder
	userID, ok := context.GetUserID(params.HTTPRequest.Context())
	if !ok {
		h.logger.Fatal("No User ID, this should never happen.")
	}
	moveID, err := uuid.FromString(params.MoveID.String())
	if err != nil {
		h.logger.Fatal("Invalid MoveID, this should never happen.")
	}

	// Validate that this move belongs to the current user
	moveResult, err := models.GetMoveForUser(h.db, userID, moveID)
	if err != nil {
		h.logger.Error("DB Error checking on move validity", zap.Error(err))
		response = ppmop.NewCreatePersonallyProcuredMoveInternalServerError()
	} else if !moveResult.IsValid() {
		switch errCode := moveResult.ErrorCode(); errCode {
		case models.FetchErrorNotFound:
			response = ppmop.NewCreatePersonallyProcuredMoveNotFound()
		case models.FetchErrorForbidden:
			response = ppmop.NewCreatePersonallyProcuredMoveForbidden()
		default:
			h.logger.Fatal("An error type has occurred that is unaccounted for in this case statement.")
		}
	} else { // The given move does belong to the current user.
		ppms, err := models.GetPersonallyProcuredMovesForMoveID(h.db, moveID)
		if err != nil {
			h.logger.Error("DB Error checking on move validity", zap.Error(err))
			response = ppmop.NewCreatePersonallyProcuredMoveInternalServerError()
		} else {
			ppmsPayload := make(internalmessages.IndexPersonallyProcuredMovePayload, len(ppms))
			for i, ppm := range ppms {
				ppmPayload := payloadForPPMModel(ppm)
				ppmsPayload[i] = &ppmPayload
			}
			response = ppmop.NewIndexPersonallyProcuredMovesOK().WithPayload(ppmsPayload)
		}
	}

	return response
}

// PatchPersonallyProcuredMoveHandler Patchs a PPM
type PatchPersonallyProcuredMoveHandler HandlerContext

// Handle is the handler
func (h PatchPersonallyProcuredMoveHandler) Handle(params ppmop.PatchPersonallyProcuredMoveParams) middleware.Responder {
	var response middleware.Responder
	userID, ok := context.GetUserID(params.HTTPRequest.Context())
	if !ok {
		h.logger.Fatal("No User ID, this should never happen.")
	}
	moveID, err := uuid.FromString(params.MoveID.String())
	if err != nil {
		h.logger.Fatal("Invalid MoveID, this should never happen.")
	}
	ppmID, err := uuid.FromString(params.PersonallyProcuredMoveID.String())
	if err != nil {
		h.logger.Fatal("Invalid PersonallyProcuredMoveID, this should never happen.")
	}

	// Make sure the move exists and is owned by the user
	exists, userOwns := models.ValidateMoveOwnership(h.db, userID, moveID)
	if !exists {
		response = ppmop.NewPatchPersonallyProcuredMoveNotFound()
		return response
	} else if !userOwns {
		response = ppmop.NewPatchPersonallyProcuredMoveForbidden()
		return response
	}

	ppm, err := models.GetPersonallyProcuredMoveForID(h.db, ppmID)
	if err != nil {
		response = ppmop.NewPatchPersonallyProcuredMoveNotFound()
		return response
	} else if ppm.MoveID != moveID {
		// Saved move ID should match request move ID
		response = ppmop.NewPatchPersonallyProcuredMoveBadRequest()
		return response
	}

	// TODO: Is there a pattern for updating that doesn't require hardcoding fields?
	size := params.PatchPersonallyProcuredMovePayload.Size
	weightEstimate := params.PatchPersonallyProcuredMovePayload.WeightEstimate

	if size != nil {
		ppm.Size = size
	}
	if weightEstimate != nil {
		ppm.WeightEstimate = weightEstimate
	}

	if verrs, err := h.db.ValidateAndUpdate(&ppm); err != nil {
		h.logger.Error("DB Patch", zap.Error(err))
		response = ppmop.NewPatchPersonallyProcuredMoveInternalServerError()
	} else if verrs.HasAny() {
		h.logger.Error("We got verrs!", zap.String("verrs", verrs.String()))
		response = ppmop.NewPatchPersonallyProcuredMoveBadRequest()
	} else {
		ppmPayload := payloadForPPMModel(ppm)
		response = ppmop.NewPatchPersonallyProcuredMoveCreated().WithPayload(&ppmPayload)
	}

	return response
}
