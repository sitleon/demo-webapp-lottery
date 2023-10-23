package main

import (
	"context"
	"database/sql"
	"demo-webapp-lottery/internal/app"
	"demo-webapp-lottery/internal/controller"
	"demo-webapp-lottery/internal/repos"
	"demo-webapp-lottery/internal/schd"
	"demo-webapp-lottery/internal/server"
	"demo-webapp-lottery/internal/svc"
	"flag"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var (
	cfg *app.AppConfig
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339, FullTimestamp: true,
	})
}

func main() {
	// configuration
	cfgPath := flag.String("config", "", "configuration file")
	if flag.Parse(); len(*cfgPath) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	cfg = app.LoadCfg(*cfgPath)

	// establish db connection
	db, err := sql.Open("postgres", cfg.Db.ConnURI)
	if err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}
	dbConn := repos.NewDbConn(db)

	// service & deps injection
	var (
		drawSvc   = svc.NewDrawImpl(dbConn, &cfg.Lottery)
		ticketSvc = svc.NewTicketSvc(dbConn)
	)
	drawSvc.SetTicketSvc(ticketSvc)
	ticketSvc.SetDrawSvc(drawSvc)

	// server
	srv := server.NewServer(cfg, []controller.RestController{
		controller.NewLottoRestImpl(ticketSvc, drawSvc),
	})
	go srv.Serve()

	// run schd task
	task := schd.NewDrawTask(drawSvc)
	go task.Run()

	// graceful shut down
	<-app.Sigterm() // wait for sigterm
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatalf("Server forced to shutdown: %v", err)
	}

	logrus.Info("Server exiting")
}
