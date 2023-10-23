package schd

import (
	"demo-webapp-lottery/internal/app"
	"demo-webapp-lottery/internal/svc"
	"time"

	"github.com/sirupsen/logrus"
)

type DrawTaskImpl struct {
	drawSvc svc.Draw
}

func NewDrawTask(drawSvc svc.Draw) *DrawTaskImpl {
	return &DrawTaskImpl{drawSvc: drawSvc}
}

func (task *DrawTaskImpl) Run() {
	ticker := time.NewTicker(1 * time.Second)
	logrus.Infoln("start draw scheduled task")

	for {
		select {
		case <-app.Sigterm():
			logrus.Infoln("stop draw scheduled task")
			return
		case <-ticker.C:
			if task.drawSvc.IsStartDraw() {
				logrus.Infoln("start draw lottery")
				if err := task.drawSvc.StartNewDraw(); err != nil {
					logrus.Errorln(err)
				}
			}
		}
	}
}
