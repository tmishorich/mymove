package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/auth"
	moveop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/moves"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"go.uber.org/zap"
)

func payloadForMoveModel(user models.User, move models.Move) internalmessages.MovePayload {
	movePayload := internalmessages.MovePayload{
		CreatedAt:        fmtDateTime(move.CreatedAt),
		SelectedMoveType: move.SelectedMoveType,
		ID:               fmtUUID(move.ID),
		UpdatedAt:        fmtDateTime(move.UpdatedAt),
		UserID:           fmtUUID(user.ID),
	}
	return movePayload
}

// CreateMoveHandler creates a new move via POST /move
type CreateMoveHandler HandlerContext

// Handle ... creates a new Move from a request payload
func (h CreateMoveHandler) Handle(params moveop.CreateMoveParams) middleware.Responder {
	var response middleware.Responder
	// User should always be populated by middleware
	user, _ := auth.GetUser(params.HTTPRequest.Context())

	// Create a new move for an authenticated user
	newMove := models.Move{
		UserID:           user.ID,
		SelectedMoveType: params.CreateMovePayload.SelectedMoveType,
	}
	if verrs, err := h.db.ValidateAndCreate(&newMove); verrs.HasAny() || err != nil {
		if verrs.HasAny() {
			h.logger.Error("DB Validation", zap.Error(verrs))
		} else {
			h.logger.Error("DB Insertion", zap.Error(err))
		}
		response = moveop.NewCreateMoveBadRequest()
	} else {
		movePayload := payloadForMoveModel(user, newMove)
		response = moveop.NewCreateMoveCreated().WithPayload(&movePayload)
	}
	return response
}

// IndexMovesHandler returns a list of all moves
type IndexMovesHandler HandlerContext

// Handle retrieves a list of all moves in the system belonging to the logged in user
func (h IndexMovesHandler) Handle(params moveop.IndexMovesParams) middleware.Responder {
	var response middleware.Responder
	// User should always be populated by middleware
	user, _ := auth.GetUser(params.HTTPRequest.Context())

	moves, err := models.GetMovesForUserID(h.db, user.ID)
	if err != nil {
		h.logger.Error("DB Query", zap.Error(err))
		response = moveop.NewIndexMovesBadRequest()
	} else {
		movePayloads := make(internalmessages.IndexMovesPayload, len(moves))
		for i, move := range moves {
			movePayload := payloadForMoveModel(user, move)
			movePayloads[i] = &movePayload
		}
		response = moveop.NewIndexMovesOK().WithPayload(movePayloads)
	}
	return response
}

// ShowMoveHandler returns a move for a user and move ID
type ShowMoveHandler HandlerContext

// Handle retrieves a move in the system belonging to the logged in user given move ID
func (h ShowMoveHandler) Handle(params moveop.ShowMoveParams) middleware.Responder {
	// User should always be populated by middleware
	user, _ := auth.GetUser(params.HTTPRequest.Context())
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.db, user, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	movePayload := payloadForMoveModel(user, *move)
	return moveop.NewShowMoveOK().WithPayload(&movePayload)

}

// PatchMoveHandler patches a move via PATCH /moves/{moveId}
type PatchMoveHandler HandlerContext

// Handle ... patches a new Move from a request payload
func (h PatchMoveHandler) Handle(params moveop.PatchMoveParams) middleware.Responder {
	// User should always be populated by middleware
	user, _ := auth.GetUser(params.HTTPRequest.Context())
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.db, user, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	payload := params.PatchMovePayload
	newSelectedMoveType := payload.SelectedMoveType

	if newSelectedMoveType != nil {
		move.SelectedMoveType = newSelectedMoveType
	}

	verrs, err := h.db.ValidateAndUpdate(move)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}
	movePayload := payloadForMoveModel(user, *move)
	return moveop.NewPatchMoveCreated().WithPayload(&movePayload)
}
