// © 2019-2024 Diarkis Inc. All rights reserved.

// Used to wait for the scenario run online like to use in k8s.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/Diarkis/diarkis/util"

	"github.com/Diarkis/diarkis-server-template/bot/scenario/lib/report"
	"github.com/Diarkis/diarkis-server-template/bot/scenario/scenarios"
)

type String string

func (s String) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, s)
}

func handleRun(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")

	if r.Method == "OPTIONS" {
		// In case the request comes from a browser
		w.WriteHeader(http.StatusOK)
		logger.Sys("Preflight handled")
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		logger.Error("Got invalid http method: %v ", r.Method)
		return
	}

	if running.Load() {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, "Scenario is already running\n")
		logger.Error("Scenario is already running")
		return
	}

	body := r.Body
	defer body.Close()

	buf := new(bytes.Buffer)
	io.Copy(buf, body)
	logger.Sys("Received Scenario Run. %v", buf.String())

	// Initialize globalParams
	gp = &scenarios.GlobalParams{}
	// Keep json body
	json.Unmarshal(buf.Bytes(), &gp.Raw.ParamsFromAPI)
	// Parse Settings (like duration to run scenario)
	json.Unmarshal(buf.Bytes(), &ss)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Scenario Started\n")
	logger.Info("Starting Scenario [%s] with parameters [%s] ", ss.ScenarioName, ss.ScenarioPattern)

	go run()

}

func handleBulkRun(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")

	if r.Method == "OPTIONS" {
		// In case the request comes from a browser
		w.WriteHeader(http.StatusOK)
		logger.Sys("Preflight handled")
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		logger.Error("Got invalid http method: %v ", r.Method)
		return
	}

	if running.Load() {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, "Scenario is already running\n")
		logger.Error("Scenario is already running")
		return
	}

	body := r.Body
	defer body.Close()

	buf := new(bytes.Buffer)
	io.Copy(buf, body)
	logger.Sys("Received Bulk Scenario Run. %v", buf.String())

	var scenarioSettings []*ScenarioSettings
	if err := json.Unmarshal(buf.Bytes(), &scenarioSettings); err != nil {
		logger.Error("Failed to unmarshal scenario settings: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid scenario settings\n")
		return
	}

	var parameters []map[string]any
	json.Unmarshal(buf.Bytes(), &parameters)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Multiple Scenario Started\n")

	go func() {
		for i, settings := range scenarioSettings {
			ss = settings
			gp = &scenarios.GlobalParams{}
			gp.Raw.ParamsFromAPI = parameters[i]
			logger.Info("Starting Scenario [%s] with parameters [%s] ", ss.ScenarioName, ss.ScenarioPattern)
			run()
		}
	}()

}

func handleCancel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")

	if r.Method == "OPTIONS" {
		// In case the request comes from a browser
		w.WriteHeader(http.StatusOK)
		logger.Sys("Preflight handled")
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		logger.Error("Got invalid http method: %v ", r.Method)
		return
	}

	err := stop()
	if err != nil {
		logger.Error("Failed to stop scenario. %v", err)
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, "Scenario is not running\n")
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Scenario Stopped\n")
	logger.Info("Scenario Stopped [%s]", ss.ScenarioName)

}

func handleListReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")

	if r.Method == "OPTIONS" {
		// In case the request comes from a browser
		w.WriteHeader(http.StatusOK)
		logger.Sys("Preflight handled")
		return
	}

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		logger.Error("Got invalid http method: %v ", r.Method)
		return
	}

	files, err := os.ReadDir("/tmp")
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		logger.Error("Error reading /tmp directory: %w", err)
		return
	}

	type FileInfo struct {
		Name      string `json:"name"`
		Timestamp uint32 `json:"timestamp"`
	}
	var fileNames []*FileInfo
	for _, file := range files {
		prefix := "DIARKIS_" + report.FilePrefix
		if strings.HasPrefix(file.Name(), prefix) {
			fileInfo, _ := file.Info()
			fi := &FileInfo{
				Name:      file.Name(),
				Timestamp: uint32(fileInfo.ModTime().Unix()),
			}
			fileNames = append(fileNames, fi)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fileNames)

	logger.Info("List report handled.")
}

func handleDownloadReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")

	if r.Method == "OPTIONS" {
		// In case the request comes from a browser
		w.WriteHeader(http.StatusOK)
		logger.Sys("Preflight handled")
		return
	}

	fileName := path.Base(r.URL.Path)
	filePath := "/tmp/" + fileName
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")

	http.ServeFile(w, r, filePath)

	logger.Info("Download report handled. file: %s", fileName)

}

// WIP
func handleGetMetrics(w http.ResponseWriter, r *http.Request) {
	// body := r.Body
	// defer body.Close()

	metrics := report.GetPrometheusMetrics()
	fmt.Fprint(w, metrics)
	logger.Verbose("Get Metrics Called...")
	// w.WriteHeader(http.StatusOK)

}

func listen() error {
	address := util.GetEnv("BOT_ADDRESS")
	if address == "" {
		address = "localhost"
	}
	port := util.GetEnv("BOT_PORT")
	if port == "" {
		port = "9500"
	}
	host := strings.Join([]string{address, port}, ":")
	http.Handle("/", String("hello"))
	http.HandleFunc("/run/", handleRun)
	http.HandleFunc("/runs/", handleBulkRun)
	http.HandleFunc("/stop/", handleCancel)
	http.HandleFunc("/report/list/", handleListReport)
	http.HandleFunc("/report/download/", handleDownloadReport)
	http.HandleFunc("/metrics/", handleGetMetrics)
	logger.Info("Bot server started. listening %s ...", host)
	http.ListenAndServe(host, nil)
	return nil
}
