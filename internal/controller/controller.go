package controller

import (
	"demo-webapp-lottery/internal/svc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LotteryRestImpl struct {
	ticketSvc svc.Ticket
	drawSvc   svc.Draw
}

var _ (RestController) = (*LotteryRestImpl)(nil)

func NewLottoRestImpl(ticketSvc svc.Ticket, drawSvc svc.Draw) RestController {
	return &LotteryRestImpl{
		ticketSvc: ticketSvc, drawSvc: drawSvc,
	}
}

func (ctr *LotteryRestImpl) RegisterRoute(root *gin.RouterGroup) {
	r := root.Group("/")
	{
		r.POST("/ticket", ctr.CreateTicket)
		r.GET("/ticket/:id", ctr.GetTicket)
		r.GET("/draw/:id")
	}
}

func (ctr *LotteryRestImpl) CreateTicket(c *gin.Context) {
	id, err := ctr.ticketSvc.CreateTicket()
	if err != nil {
		c.JSON(
			http.StatusInternalServerError, gin.H{"error": err.Error()},
		)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ticketId": id})
}

func (ctr *LotteryRestImpl) GetTicket(c *gin.Context) {
	res, err := ctr.ticketSvc.GetTicket(c.Param("id"))
	if err != nil {
		c.JSON(
			http.StatusInternalServerError, gin.H{"error": err.Error()},
		)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (ctr *LotteryRestImpl) GetDraw(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 0 {
		c.JSON(
			http.StatusBadRequest, gin.H{"error": "invalid draw id"},
		)
		return
	}
	res, err := ctr.drawSvc.GetDraw(id)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError, gin.H{"error": err.Error()},
		)
		return
	}
	c.JSON(http.StatusOK, res)
}
