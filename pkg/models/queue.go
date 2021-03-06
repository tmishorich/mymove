package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// MoveQueueItem represents a single move queue item within a queue.
type MoveQueueItem struct {
	ID               uuid.UUID                           `json:"id" db:"id"`
	CreatedAt        time.Time                           `json:"created_at" db:"created_at"`
	Edipi            string                              `json:"edipi" db:"edipi"`
	Rank             *internalmessages.ServiceMemberRank `json:"rank" db:"rank"`
	CustomerName     string                              `json:"customer_name" db:"customer_name"`
	Locator          string                              `json:"locator" db:"locator"`
	Status           string                              `json:"status" db:"status"`
	OrdersType       string                              `json:"orders_type" db:"orders_type"`
	MoveDate         time.Time                           `json:"move_date" db:"move_date"`
	CustomerDeadline time.Time                           `json:"customer_deadline" db:"customer_deadline"`
	LastModifiedDate time.Time                           `json:"last_modified_date" db:"last_modified_date"`
	LastModifiedName string                              `json:"last_modified_name" db:"last_modified_name"`
}

// GetMoveQueueItems gets all moveQueueItems for a specific lifecycleState
func GetMoveQueueItems(db *pop.Connection, lifecycleState string) ([]MoveQueueItem, error) {
	moveQueueItems := []MoveQueueItem{}
	// TODO: Add clause `SELECT moves.ID, sm.edipi, sm.rank, CONCAT(sm.last_name, ', ', sm.first_name) AS customer_name`
	// TODO: add clause `JOIN personally_procured_moves AS ppm ON moves.id = ppm.move_id`
	// TODO: add `CONCAT(sm.last_name, ', ', sm.first_name)` back to query, once data is in place

	// TODO: replace hardcoded values with actual query values once data is available
	query := `
		SELECT moves.ID,
			COALESCE(sm.edipi, 'test ID') as edipi,
			COALESCE(sm.rank, 'major') as rank,
			'Tester, Telly' AS customer_name,
			'12345' as locator,
			COALESCE(moves.selected_move_type, 'COMBO') as orders_type,
			current_date as move_date,
			current_time as created_at,
			'Awaiting review' as status,
			current_time as last_modified_date
		FROM moves
		INNER JOIN service_members AS sm ON moves.user_id = sm.user_id
	`
	// TODO: add clause `WHERE moves.lifecycle_state = $1`
	// err = db.RawQuery(query, lifecycleState).All(&moveQueueItems)
	err := db.RawQuery(query).All(&moveQueueItems)
	return moveQueueItems, err
}
