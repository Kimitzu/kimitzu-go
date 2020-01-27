package mobile

import (
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/gorilla/mux"
	"github.com/nokusukun/particles/config"
	"github.com/nokusukun/particles/roggy"

	"github.com/kimitzu/kimitzu-services/api"
	"github.com/kimitzu/kimitzu-services/configs"
	"github.com/kimitzu/kimitzu-services/p2p"

	"github.com/kimitzu/kimitzu-services/location"
	"github.com/kimitzu/kimitzu-services/servicestore"
	"github.com/kimitzu/kimitzu-services/voyager"
)

var (
	logger  = roggy.Printer("services")
	confSat = config.Satellite{
		Port:        9009,
		Host:        "0.0.0.0",
		DisableUPNP: false,
	}
)

type Services struct {
	daemon *configs.Daemon
}

func NewServices(dataPath string, testnet bool, logLvl int) *Services {
	daemon := &configs.Daemon{
		// ArgLogLevel:           "log",
		DataPath:              dataPath,
		LogLevel:              logLvl,
		Version:               "0.2.0-alpha.1",
		DialTo:                "",
		ApiListen:             "0.0.0.0:8109",
		KeyPath:               path.Join(dataPath, "p2pkeys"),
		GenerateNewKeys:       true,
		ShowHelp:              false,
		DatabasePath:          path.Join(dataPath, "p2p"),
		BootstrapNodeIdentity: "",
		Testnet:               testnet,
	}
	roggy.LogLevel = logLvl
	roggy.Simple = true
	return &Services{daemon: daemon}
}

func (s *Services) Start() error {
	log := logger

	time.Sleep(time.Second * 1)

	log.Info(fmt.Sprintf("Kimitzu Services Daemon (%v)", s.daemon.Version))
	log.Info(" --- --- --- --- --- ")
	log.Infof("Log Level: %v, Simple Mode: %v", roggy.LogLevel, roggy.Simple)
	log.Info("Starting Services")

	if s.daemon.Testnet {
		log.Info("Network: Testnet...")
	} else {
		log.Info("Network: Mainnet...")
	}
	store := servicestore.InitializeManagedStorage(s.daemon.DataPath)
	p2pKillSig := make(chan int, 1)

	ratingManager, err := p2p.InitializeRatingManager(s.daemon.DatabasePath)
	if err != nil {
		log.Error("Opening database failed")
		roggy.Wait()
		panic(err)
	}

	apiRouter := mux.NewRouter()

	time.Sleep(time.Second * 10)
	go p2p.Bootstrap(s.daemon, &confSat, ratingManager, p2pKillSig)
	go voyager.RunVoyagerService(log.Sub("voyager"), store)
	location.RunLocationService(log.Sub("location"))

	p2p.AttachAPI(p2p.Sat, apiRouter, ratingManager)
	api.AttachStore(store)
	api.AttachAPI(log.Sub("api"), apiRouter)

	log.Infof("Running API on %v", s.daemon.ApiListen)
	go http.ListenAndServe(s.daemon.ApiListen, apiRouter)
	log.Infof("Kimitzu services is running...")

	return nil
}
