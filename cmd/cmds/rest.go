package cmds

import (
	"github.com/spf13/cobra"
	"gitlab.company.com/projectname/configs"
	"gitlab.company.com/projectname/pkg/logx"
	"gitlab.company.com/projectname/protocol"
)

var restCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP API server",
	RunE:  runRest,
}

func runRest(_ *cobra.Command, _ []string) error {
	conf := configs.GetConfigs()
	logx.Init(conf.App.LogLevel)
	log := logx.GetLog()

	app, err := protocol.New(conf)
	if err != nil {
		log.WithError(err).Fatal("failed to initialize application")
		return err
	}

	app.StartCron()
	log.WithField("port", conf.App.Port).Info("starting server")
	return app.Start()
}
