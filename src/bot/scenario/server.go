// Used to wait for the scenario run online like to use in k8s.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"{0}/bot/scenario/lib/report"

	"github.com/Diarkis/diarkis/util"
)

type String string

func (s String) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, s)
}

func handleRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		logger.Error("Got invalid http method: %v ", r.Method)
	}
	body := r.Body
	defer body.Close()

	buf := new(bytes.Buffer)
	io.Copy(buf, body)

	json.Unmarshal(buf.Bytes(), &scenarioSettings)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Scenario Started\n")
	logger.Info("Starting Scenario [%s] with parameters [%s] ", scenarioSettings.ScenarioName, scenarioSettings.ScenarioPattern)

	go run()

}

// WIP
func handleGetMetrics(w http.ResponseWriter, r *http.Request) {
	// body := r.Body
	// defer body.Close()

	callCnt := report.GetCallCommandMetrics()
	for ver, cmds := range callCnt {
		for cmd, m := range cmds {
			// TODO: use {} rather than name all
			fmt.Fprint(w, "# HELP COMMAND_CALL_COUNT_VER_%d_CMD_%d Call count for each commands\n", ver, cmd)
			fmt.Fprint(w, "# TYPE COMMAND_CALL_COUNT_VER_%d_CMD_%d counter\n", ver, cmd)
			fmt.Fprint(w, "COMMAND_CALL_COUNT_VER_%d_CMD_%d %d\n", ver, cmd, m.GetTotal())

		}
	}

	// // todo
	// pushCnt := report.GetPushMetrics()
	// // todo
	// resCnt := report.GetResponseMetrics()

	w.WriteHeader(http.StatusOK)

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
	http.HandleFunc("/metrics/", handleGetMetrics)
	// todo: dynamic port
	logger.Info("Bot server started. listening %s ...", host)
	http.ListenAndServe(host, nil)
	return nil
}
