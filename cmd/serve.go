package cmd

import (
	"context"
	"fmt"

	"github.com/delivc/identity/api"
	"github.com/delivc/identity/conf"
	"github.com/delivc/identity/storage"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serveCmd = cobra.Command{
	Use:  "serve",
	Long: "Start API server",
	Run: func(cmd *cobra.Command, args []string) {
		execWithConfig(cmd, serve)
	},
}

func serve(globalConfig *conf.GlobalConfiguration, config *conf.Configuration) {
	db, err := storage.Dial(globalConfig)
	if err != nil {
		logrus.Fatalf("Error opening database: %+v", err)
	}
	defer db.Close()

	ctx, err := api.WithInstanceConfig(context.Background(), config, uuid.Nil)
	if err != nil {
		logrus.Fatalf("Error loading instance config: %+v", err)
	}

	api := api.NewAPIWithVersion(ctx, globalConfig, db, Version)

	l := fmt.Sprintf("%v:%v", globalConfig.API.Host, globalConfig.API.Port)
	logrus.Infof("Delivc Identiy API started on: %s", l)
	api.ListenAndServe(l)

	logrus.Info("Hello from serveCmd")
}
