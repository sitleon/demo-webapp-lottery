package svc

import (
	"context"
	"demo-webapp-lottery/internal/repos"
	"demo-webapp-lottery/internal/utils"
	"sync"
)

type Ticket interface {
	SetDrawSvc(draw Draw)
	CreateTicket() (id string, err error)
	GetNextTicketIdSet() (ticketId string, drawId int64)
	ResetTicketIndex()
	GetTicket(ticketId string) (repos.GetTicketRow, error)
}

type TicketImpl struct {
	db        *repos.DbConn
	incLock   sync.Mutex
	resetLock sync.RWMutex
	nxtTicIdx int32

	// service
	drawSvc Draw
}

var _ Ticket = (*TicketImpl)(nil)

func NewTicketSvc(db *repos.DbConn) Ticket {
	return &TicketImpl{db: db}
}

func (svc *TicketImpl) SetDrawSvc(draw Draw) {
	svc.drawSvc = draw
}

func (svc *TicketImpl) CreateTicket() (string, error) {
	id, d := svc.GetNextTicketIdSet()

	err := svc.db.Querier.CreateTicket(
		context.Background(), repos.CreateTicketParams{
			TicketID: id, DrawID: d,
		})
	if err != nil {
		return "", err
	}

	return id, nil
}

func (svc *TicketImpl) GetNextTicketIdSet() (string, int64) {
	svc.incLock.Lock()
	defer svc.incLock.Unlock()
	svc.resetLock.RLock()

	d, x := svc.drawSvc.GetCurrentDrawId(), svc.nxtTicIdx
	svc.nxtTicIdx += 1

	id, _ := utils.Int64ToUuid(uint64(d), uint64(x))
	return id, d
}

func (svc *TicketImpl) ResetTicketIndex() {
	svc.resetLock.Lock()
	defer svc.resetLock.Unlock()

	svc.nxtTicIdx = 0
}

func (svc *TicketImpl) GetTicket(ticketId string) (repos.GetTicketRow, error) {
	return svc.db.Querier.GetTicket(context.Background(), ticketId)
}
