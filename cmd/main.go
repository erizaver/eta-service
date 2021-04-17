package main

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/erizaver/eta-service/internal/pkg/app"
)

func main() {
	ctx := context.Background()

	log.Info("creating application")
	application := app.NewApp()

	log.Info("starting application")
	if err := application.Run(ctx); err != nil {
		log.Fatal("error running application", err)
	}
}
