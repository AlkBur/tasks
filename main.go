package main

import (
	"log"
	"net/http"
	"tasks/db"
	"tasks/lib"
	"tasks/server"
	"tasks/service"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world\n"))
}

func main() {
	config, err := lib.ConfigLoad(".")
	if err != nil {
		log.Fatalf("Could not load config. %v", err)
	}
	l := lib.NewLogger(config.LogLevel)
	sqldb, err := db.NewSQL("app.db", &l)
	if err != nil {
		l.Fatal().Err(err).Send()
	}
	defer sqldb.Close()

	s := service.NewTask(sqldb)

	r, err := server.NewChiRouter(s, config.PASETOSecret, config.AccessTokenDuration, &l)
	if err != nil {
		l.Fatal().Err(err).Send()
	}

	httpServer, err := server.NewHTTP(r, config.HTTPServerAddress, &l)
	if err != nil {
		l.Fatal().Err(err).Send()
	}

	err = httpServer.Shutdown()
	if err != nil {
		l.Fatal().Err(err).Send()
	}
}
