package svc

import (
	"context"
	"demo-webapp-lottery/internal/app"
	"demo-webapp-lottery/internal/repos"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Draw interface {
	SetTicketSvc(ticket Ticket)
	GetCurrentDrawId() int64
	Draw() (string, error)
	IsStartDraw() bool
	StartNewDraw() error
	CreateDrawDbRecord(drawId int64, drewAt time.Time) error
	GetDraw(drawId int64) (repos.GetDrawRow, error)
}

type DrawImpl struct {
	db  *repos.DbConn
	mu  sync.RWMutex
	cfg *app.LotteryConfig

	currDrawId  int64
	scheduledAt time.Time

	ticketSvc Ticket
}

func NewDrawImpl(db *repos.DbConn, cfg *app.LotteryConfig) Draw {
	return &DrawImpl{db: db, cfg: cfg}
}

func (svc *DrawImpl) SetTicketSvc(ticket Ticket) {
	svc.ticketSvc = ticket
}

func (svc *DrawImpl) GetCurrentDrawId() int64 {
	svc.mu.RLock()
	return svc.currDrawId
}

func (svc *DrawImpl) Draw() (string, error) {
	drawId := svc.currDrawId
	if err := svc.StartNewDraw(); err != nil {
		return "", err
	}

	cnt, err := svc.db.Querier.CountTicket(context.TODO(), drawId)
	if err != nil {
		return "", err
	}

	ticketId, err := svc.SetDrawWinner(svc.currDrawId, rand.Int31n(int32(cnt)))
	if err != nil {
		return "", err
	}

	return ticketId, nil
}

func (svc *DrawImpl) SetDrawWinner(drawId int64, offset int32) (string, error) {
	arg := repos.SetDrawWinnerParams{DrawID: drawId, Offset: offset}
	ticketId, err := svc.db.Querier.SetDrawWinner(context.Background(), arg)
	if err != nil {
		return "", err
	}
	return ticketId.String, nil
}

func (svc *DrawImpl) StartNewDraw() error {
	if !svc.mu.TryLock() {
		return fmt.Errorf("draw already started by another process")
	}
	defer svc.mu.Unlock()

	now, interval := time.Now(), time.Duration(svc.cfg.Interval)*time.Second
	err := svc.CreateDrawDbRecord(now.Unix(), now.Add(interval))
	if err != nil {
		return err
	}

	svc.currDrawId, svc.scheduledAt = now.Unix(), now.Add(interval)
	svc.ticketSvc.ResetTicketIndex()
	return nil
}

func (svc *DrawImpl) IsStartDraw() bool {
	svc.mu.RLock()
	return svc.scheduledAt.Compare(time.Now()) <= 0
}

func (svc *DrawImpl) CreateDrawDbRecord(drawId int64, drewAt time.Time) error {
	arg := repos.CreateDrawParams{DrawID: drawId, DrewAt: drewAt}
	return svc.db.Querier.CreateDraw(context.Background(), arg)
}

func (svc *DrawImpl) GetDraw(drawId int64) (repos.GetDrawRow, error) {
	return svc.db.Querier.GetDraw(context.Background(), drawId)
}
