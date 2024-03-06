package main

import (
	"flag"
	"fmt"
	"go-pentor-bank/internal/batch"
	"go-pentor-bank/internal/clog"
	"go-pentor-bank/internal/common"
	"go-pentor-bank/internal/config"
	"go-pentor-bank/internal/consumer"
	"go-pentor-bank/internal/createIndex"
	"go-pentor-bank/internal/infra"
	"go-pentor-bank/internal/onetimebatch"
	"go-pentor-bank/internal/routes/api"
	"go-pentor-bank/internal/routes/doc"
	"go-pentor-bank/internal/utils"
	"os"
	"strings"
	_ "time/tzdata"
)

func main() {
	port := flag.String("p", "12345", "port number")
	//socketPort := flag.String("scPort", "8000", "socket port number")
	socketPortV2 := flag.String("scPortV2", "3011", "socket port number")
	socketPortV4 := flag.String("scPortV4", "3012", "socket port number")
	state := flag.String("state", "local", "set working environment")
	configPath := flag.String("config", "configs/default", "set configs path, default as: 'configs'")
	errorPath := flag.String("error", "configs/common", "set configs path, default as: 'configs'")
	app := flag.String("app", "api", "batch, api, and bo")
	logFormat := flag.String("logFormat", "text", "text, json")
	debug := flag.Bool("debug", false, "sets log level to debug")
	jobs := flag.String("job", "", "Job for run batch")
	prefixPath := flag.String("prefixPath", "", "Job for run batch")
	flag.Parse()

	osPort := os.Getenv("PRT")
	if osPort != "" {
		port = &osPort
	}
	osStage := os.Getenv("STATE")
	if osStage != "" {
		state = &osStage
	}
	osApp := os.Getenv("APP")
	if osApp != "" {
		app = &osApp
	}
	osLogFormat := os.Getenv("LOGFORMAT")
	if osLogFormat != "" {
		logFormat = &osLogFormat
	}
	osJobs := os.Getenv("JOBS")
	if osJobs != "" {
		jobs = &osJobs
	}
	//osSocketPort := os.Getenv("SOCKETPRT")
	//if osSocketPort != "" {
	//	socketPort = &osSocketPort
	//}
	osSocketPortV2 := os.Getenv("SOCKETPRTV2")
	if osSocketPortV2 != "" {
		socketPortV2 = &osSocketPortV2
	}
	osSocketPortV4 := os.Getenv("SOCKETPRTV4")
	if osSocketPortV4 != "" {
		socketPortV4 = &osSocketPortV4
	}

	clog.New(*logFormat, *debug)

	log := clog.GetLog()
	log.Info().Msg(fmt.Sprintf("Running on %v environment, app: %v", *state, *app))

	if err := config.Conf.InitViperWithStage(*state, *configPath); err != nil {
		log.Panic().Err(err).Msg("Read config file err:")
		return
	}

	if err := config.InitDefaultValidators(); err != nil {
		log.Error().Err(err).Msg("InitDefaultValidators err")
		return
	}

	if err := config.EM.Init(*errorPath); err != nil {
		log.Panic().Err(err).Msg("Read Error file")
		return
	}

	// Set time location
	err := config.LoadTimeLocation()
	if err != nil {
		log.Panic().Err(err).Msg("parse location error")
	}

	infra.EstablishInfraConnection(*app)

	defer func() {
		log.Print("Running cleanup tasks...")
		infra.ShutdownInfra()
		log.Print("Server down")
	}()

	log.Info().Msgf("running app %v, job %v", *app, strings.Split(*jobs, ","))
	switch *app {
	case common.AppAPI:
		//socketIO.RunBroadcast(*socketPortV2, *socketPortV4)
		api.App(*app, *port, *prefixPath)

	//case common.AppBO:
	//	socketIO.RunBroadcast(*socketPortV2, *socketPortV4)
	//	bo.App(*app, *port)
	//
	//case common.AppOpenAPI:
	//	socketIO.RunBroadcast(*socketPortV2, *socketPortV4)
	//	openapi.App(*app, *port)
	//
	//case common.AppCDN:
	//	cdnImg.App(*app, *port)

	case "doc":
		doc.App(*app, *port)

	case common.AppBatch:
		log.Info().Msgf("Batch: running at port %v", *port)
		batch.Batch(strings.Split(*jobs, ","), *port)

	//case common.AppSocketIOV2:
	//	socketIO.RunBroadcast(*socketPortV2, *socketPortV4)
	//	appv2.NewSocketServerConnector(*socketPort, false)
	//
	//case common.AppSocketio:
	//	socketIO.RunBroadcast(*socketPortV2, *socketPortV4)
	//	appv4.NewSocketServerConnector(*socketPort, false)

	case "index":
		createIndex.Run()

	case common.AppConsumer:
		consumer.Run(*jobs, *port, *socketPortV2, *socketPortV4)

	case "onetimebatch":
		onetimebatch.Run(*jobs)

	default:
		log.Panic().Msg("invalid app")
	}
	utils.WaitGoRoutines()
}
