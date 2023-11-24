package main

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"{0}/bot/scenario/lib/report"
	"{0}/bot/scenario/scenarios"

	"github.com/Diarkis/diarkis/client/go/udp"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/util"
)

var commonParams map[string]any
var scenarioParams map[string]any
var gp = scenarios.GlobalParams{
	Host:            "localhost:7000",
	ReceiveByteSize: 1400,
	UDPSendInterval: 100,
	// 20: sys (can be changed by configs)
	LogLevel: 20,
}

var config string = "config"

type ScenarioSettings struct {
	ScenarioName    string `json:"type"`
	ScenarioPattern string `json:"run"`
	HowMany         int    `json:"howmany"`
	Duration        int    `json:"duration"`
	Interval        int    `json:"interval"`
}

var scenarioSettings *ScenarioSettings

// var scenarioName, scenarioPattern string
// var howmany, duration, interval int

var scenarioFactory *func() scenarios.Scenario
var logger = log.New("BOT/MAIN")

func load() error {
	// read flags
	// var serverMode int
	// flag.StringVar(&config, "config", "config", "Directory for config files.")
	flag.StringVar(&scenarioSettings.ScenarioName, "type", "Connect", "Scenario name that you implement and defined in ScenarioList.")
	flag.StringVar(&scenarioSettings.ScenarioPattern, "run", "ConnectUDP", "Scenario instance that is defined as 'hint' in Json file.")
	flag.IntVar(&scenarioSettings.HowMany, "howmany", -1, "The number of clients to join matching.")
	flag.IntVar(&scenarioSettings.Duration, "duration", -1, "Duration to run the scenario in seconds.")
	flag.IntVar(&scenarioSettings.Interval, "interval", -1, "Interval to create scenario clients in millisecond.")
	// flag.BoolVar(&serverMode, "serverMode", false, "Set true to run this as a server mode to wait for the scenario input from http.")
	flag.Parse()
	return nil
}
func setup() error {
	// get scenario
	scenarioFactory_, ok := scenarios.ScenarioFactoryList[scenarioSettings.ScenarioName]
	if !ok {
		return util.NewError("No Scenario named \"%v\". Please check 'ScenarioList' in bot/scenario/scenarios/main.go", scenarioSettings.ScenarioName)
	}
	scenarioFactory = &scenarioFactory_

	// read config files
	filepath.Walk(config, func(path string, info os.FileInfo, err error) error {
		filename := info.Name()
		logger.Debug("Reading file \"%s\" ...", filename)
		if strings.HasSuffix(filename, ".json") {
			raw, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			logger.Debug("File Contents: %v", string(raw))
			json.Unmarshal(raw, &commonParams)
			scenarioParamsIf, ok := commonParams[scenarioSettings.ScenarioPattern]
			if ok {
				scenarioParams = scenarioParamsIf.(map[string]any)
			}
		}
		return nil
	})
	if scenarioParams == nil {
		logger.Warn("No Scenario Parameter [%s] found in json files. Using only global parameters.", scenarioSettings.ScenarioPattern)
	}

	// set global params
	globalParamsBytes, _ := scenarios.GenerateParams(0, commonParams, scenarioParams)
	err := json.Unmarshal(globalParamsBytes, &gp)
	if err != nil {
		return util.StackError(util.NewError("Failed to pars common params. %v", err.Error()))
	}

	// set log level
	udp.LogLevel(gp.LogLevel)

	// value from args will be prioritized
	// todo: review
	if scenarioSettings.HowMany >= 0 {
		gp.HowMany = scenarioSettings.HowMany
	}
	if scenarioSettings.Duration >= 0 {
		gp.Duration = scenarioSettings.Duration
	}
	if scenarioSettings.Interval >= 0 {
		gp.Interval = scenarioSettings.Interval
	}
	gp.RawConfigs.CommonParams = commonParams
	gp.RawConfigs.ScenarioParams = scenarioParams

	logger.Info("Setup done. commonParams:%+v scenarioParams:%+v", commonParams, scenarioParams)
	return nil
}

func start() error {
	clients := make([]*scenarios.Scenario, gp.HowMany)
	// create bot clients
	var wg sync.WaitGroup
	wg.Add(gp.HowMany)
	for i := 0; i <= gp.HowMany-1; i++ {
		go func(i int) {
			defer wg.Done()
			scenarioClient := (*scenarioFactory)()
			clients[i] = &scenarioClient

			// set scenario params
			scenarioParamsBytes, _ := scenarios.GenerateParams(i, commonParams, scenarioParams)
			scenarioClient.ParseParam(i, scenarioParamsBytes)

			// execute scenario
			err := scenarioClient.Run(&gp)
			if err != nil {
				logger.Error(util.StackError(util.NewError("Scenario execution failed. %v", err.Error())))
				// todo: report.IncrementScenarioError()
			}
		}(i)
		time.Sleep(time.Millisecond * time.Duration(gp.Interval))
	}

	// wait for "duration" time if it's set. mainly used for looping scenario
	if gp.Duration == 0 {
		wg.Wait()
	} else {
		time.Sleep(time.Second * time.Duration(gp.Duration))
	}

	// loop again to call scenario end callback
	for i := 0; i <= gp.HowMany-1; i++ {
		(*clients[i]).OnScenarioEnd()
		globalReport, report := (*clients[i]).WriteReport()
		// todo: merge
		logger.Notice("global:%+v individual:%+v", globalReport, report)
	}

	// logger.Notice("Report: %+v", report.GetAllCommandMetrics())
	// todo: merge
	report.PrintCustomMetrics()
	report.PrintAllCommandMetrics()
	return nil
}

func run() error {
	err := setup()
	if err != nil {
		logger.Fatal("\x1b[0;91m%v\x1b[0m", err.Error())
		return err
	}

	err = start()
	if err != nil {
		logger.Fatal("\x1b[0;91m%v\x1b[0m", err.Error())
		return err
	}
	return nil
}

func main() {
	scenarioSettings = &ScenarioSettings{HowMany: -1, Interval: -1, Duration: -1}
	isServerMode := util.GetEnv("BOT_SEVER_MODE")
	configPath := util.GetEnv("BOT_CONFIG")
	if configPath != "" {
		config = configPath
	}

	if isServerMode == "true" {
		err := listen()
		if err != nil {
			logger.Fatal("\x1b[0;91m%v\x1b[0m", err.Error())
			os.Exit(1)
		}
	} else {
		err := load()
		if err != nil {
			logger.Fatal("\x1b[0;91m%v\x1b[0m", err.Error())
			os.Exit(1)
		}

		err = run()
		if err != nil {
			logger.Fatal("\x1b[0;91m%v\x1b[0m", err.Error())
			os.Exit(1)
		}
	}
}
