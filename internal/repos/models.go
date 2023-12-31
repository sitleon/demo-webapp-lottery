// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package repos

import (
	"database/sql"
	"time"
)

type Draw struct {
	DrawID       int64          `json:"drawId"`
	WinnerTicket sql.NullString `json:"winnerTicket"`
	DrewAt       time.Time      `json:"drewAt"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdateAt     sql.NullTime   `json:"updateAt"`
}

type Ticket struct {
	TicketID  string       `json:"ticketId"`
	DrawID    int64        `json:"drawId"`
	CreatedAt sql.NullTime `json:"createdAt"`
}
