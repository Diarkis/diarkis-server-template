package report

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"{0}/bot/scenario/lib/log"
)

const interval = 2

var logger = log.New("BOT/REPORT")

type Report map[string]int
type metrics struct {
	total atomic.Uint32
}

func NewMetrics() *metrics {
	m := metrics{}
	go m.start()
	return &m
}

func (m *metrics) start() {
	for {
		// mean
		time.Sleep(interval * time.Second)
	}
}
func (m *metrics) Add(delta uint32) {
	logger.Notice("%#v", m)
	m.total.Add(delta)
}
func (m *metrics) GetTotal() uint32 {
	return m.total.Load()
}

type CustomMetrics struct {
	counter map[string]*metrics
}

var customMetrics = &CustomMetrics{
	counter: map[string]*metrics{},
}

func NewCustomMetrics() *CustomMetrics {
	cm := CustomMetrics{
		counter: map[string]*metrics{},
	}
	return &cm
}

func (cm *CustomMetrics) Increment(key string) {
	if _, ok := cm.counter[key]; !ok {
		cm.counter[key] = NewMetrics()
	}
	cm.counter[key].Add(1)

	if _, ok := customMetrics.counter[key]; !ok {
		customMetrics.counter[key] = NewMetrics()
	}
	customMetrics.counter[key].Add(1)
}

func (cm *CustomMetrics) Print() {
	for key, c := range cm.counter {
		logger.Info("key:%s call count:%d", key, c.total.Load())
	}
}
func PrintCustomMetrics() {
	if customMetrics == nil {
		logger.Notice("Nothing called")
		return
	}
	for key, m := range customMetrics.counter {
		logger.Notice("key:%s call count:%d", key, m.total.Load())

	}
}

type CommandMetrics struct {
	sync.RWMutex
	cType   string
	counter map[uint8]map[uint16]*metrics
	// metrics map[uint8]map[uint16]struct {
	// 	counter  atomic.Uint32
	// 	duration time.Duration
	// }
}

var callCommandMetrics = &CommandMetrics{cType: "send"}
var pushCommandMetrics = &CommandMetrics{cType: "push"}
var responseCommandMetrics = &CommandMetrics{cType: "response"}

func (cm *CommandMetrics) Increment(ver uint8, cmd uint16) {
	cm.Lock()
	defer cm.Unlock()
	if cm.counter == nil {
		cm.counter = map[uint8]map[uint16]*metrics{}
	}
	if _, ok := cm.counter[ver]; !ok {
		cm.counter[ver] = map[uint16]*metrics{}
	}
	if _, ok := cm.counter[ver][cmd]; !ok {
		cm.counter[ver][cmd] = NewMetrics()
	}
	cm.counter[ver][cmd].total.Add(1)
}

func (cm *CommandMetrics) Print() {
	for ver, cmds := range cm.counter {
		for cmd, m := range cmds {
			logger.Notice("[%s]	ver:%d cmd:%d count:%d", cm.cType, ver, cmd, m.total.Load())
		}
	}
}

func IncrementCallCommandMetrics(ver uint8, cmd uint16) {
	callCommandMetrics.Increment(ver, cmd)
}
func IncrementPushMetrics(ver uint8, cmd uint16) {
	pushCommandMetrics.Increment(ver, cmd)
}
func IncrementResponseMetrics(ver uint8, cmd uint16) {
	responseCommandMetrics.Increment(ver, cmd)
}

func GetCallCommandMetrics() map[uint8]map[uint16]*metrics {
	return callCommandMetrics.counter
}
func GetPushMetrics() map[uint8]map[uint16]*metrics {
	return pushCommandMetrics.counter
}
func GetResponseMetrics() map[uint8]map[uint16]*metrics {
	return responseCommandMetrics.counter
}

func PrintAllCommandMetrics() {
	callCommandMetrics.Print()
	pushCommandMetrics.Print()
	responseCommandMetrics.Print()
}

type APIMetrics struct {
}

type statistic struct {
	userId         string
	startTime      time.Time
	duration       time.Duration
	scenarioReport any
}

type report struct {
	Name       string
	startTime  time.Time
	statistics []*statistic
}

var uidMap struct {
	sync.RWMutex
	IDs map[string]string
}

func NewReport(name string) *report {
	rpt := new(report)

	logger.Info("generating report.... %v", name)
	rpt.Name = name
	rpt.startTime = time.Now()
	return rpt
}

func (rpt *report) Start(apiID uint16, userId string, startTime_ time.Time) *statistic {
	sts := new(statistic)
	if startTime_.IsZero() {
		sts.startTime = time.Now()
	} else {
		sts.startTime = startTime_
	}
	sts.userId = userId
	logger.Info("starting statistic.... %v", sts)
	return sts

}

func (rpt *report) Write(sts *statistic, scenarioReport any) {
	sts.duration = time.Since(sts.startTime)
	sts.scenarioReport = scenarioReport
	rpt.statistics = append(rpt.statistics, sts)
	logger.Info("adding statistic.... %v", sts)
}

func (rpt *report) AppendUID(uid string, cosmosID string) {
	uidMap.Lock()
	if uidMap.IDs == nil {
		uidMap.IDs = map[string]string{}
	}
	uidMap.IDs[uid] = cosmosID
	uidMap.Unlock()
}

func (rpt *report) WriteUserIDMap() {
	file, err := os.Create("/tmp/usermap.csv")
	if err != nil {
		logger.Error(err)
	}
	defer file.Close()
	logger.Notice(uidMap.IDs)
	for uid, cid := range uidMap.IDs {
		file.WriteString(fmt.Sprintf("%v,%v\n", uid, cid))
	}
}

// func (rpt *report) PrintResult(totalCnt int, errCnt int) {
// 	logger.Info("Report Result: %v", rpt)
// 	var totalMatching int
// 	var totalTicketError uint
// 	var latestStartTime time.Time
// 	var matchingRate float32
// 	var avrDuration time.Duration
// 	var minDuration time.Duration
// 	var maxDuration time.Duration
// 	var durations time.Duration

// 	for _, sts := range rpt.statistics {
// 		logger.Info("userId: %v", sts.userId)
// 		logger.Info(" started:	%v", sts.startTime)
// 		logger.Info(" took:		%v", sts.duration)
// 		if flags, ok := sts.scenarioReport.(load.MatchingFlags); ok {
// 			values := reflect.ValueOf(flags)
// 			keys := values.Type()
// 			for i := 0; i < values.NumField(); i++ {
// 				key := keys.Field(i).Name
// 				value := values.Field(i).Interface()
// 				logger.Info(" %v:		%v", key, value)
// 			}
// 			// individual report statistics
// 			for _, duration := range flags.MatchingDurations {
// 				durations = durations + duration
// 				if minDuration == 0 || duration < minDuration {
// 					minDuration = duration
// 				}
// 				if maxDuration < duration {
// 					maxDuration = duration
// 				}
// 			}
// 			totalMatching = totalMatching + len(flags.MatchingDurations)
// 			totalTicketError = totalTicketError + uint(flags.IsTicketError.Load())
// 			if sts.startTime.After(latestStartTime) {
// 				latestStartTime = sts.startTime
// 			}
// 		}

// 	}

// 	clientCount := len(rpt.statistics)
// 	if totalMatching != 0 {
// 		avrDuration = time.Duration(int(durations) / totalMatching)
// 	}
// 	if clientCount != 0 {
// 		matchingRate = float32(totalMatching) / float32(clientCount)
// 	}

// 	errorRate := float32(errCnt) / float32(totalCnt)
// 	fmt.Printf("~~~~~~~~~~~~~~~~~~~~~~~~ report summary ~~~~~~~~~~~~~~~~~~~~~~~~ \n")
// 	fmt.Printf("time: %v \n", time.Since(rpt.startTime))
// 	fmt.Printf("client count: %v  \n", clientCount)
// 	fmt.Printf("all clients spawned in: %v  \n", latestStartTime.Sub(rpt.startTime))
// 	fmt.Printf("error rate for target API: %.2f％ \n", errorRate*100)
// 	fmt.Println()
// 	fmt.Printf("total matching: %v \n", totalMatching/2)
// 	fmt.Printf("total ticket error: %v \n", totalTicketError)
// 	fmt.Printf("matching rate: %.3f times / client \n", matchingRate)
// 	fmt.Printf("average matching duration: %v \n", avrDuration)
// 	fmt.Printf("minimum matching duration: %v \n", minDuration)
// 	fmt.Printf("maximum matching duration: %v \n", maxDuration)
// 	fmt.Printf("~~~~~~~~~~~~~~~~~~~~~~~~ report summary ~~~~~~~~~~~~~~~~~~~~~~~~ \n")

// }

// func (rpt *report) PrintSessionResult(totalCnt int, errCnt int, createRate float64) {
// 	logger.Info("Report Result: %v", rpt)
// 	var totalSession uint32
// 	var latestStartTime time.Time
// 	var requestCount int
// 	totalRtt := map[uint16]time.Duration{}
// 	apiCall := map[uint16]int{}
// 	errorCounts := map[uint16]uint32{}

// 	for _, sts := range rpt.statistics {
// 		logger.Info("userId: %v", sts.userId)
// 		logger.Info(" started:	%v", sts.startTime)
// 		logger.Info(" took:		%v", sts.duration)
// 		if flags, ok := sts.scenarioReport.(load.SessionFlags); ok {
// 			values := reflect.ValueOf(flags)
// 			keys := values.Type()
// 			for i := 0; i < values.NumField(); i++ {
// 				key := keys.Field(i).Name
// 				value := values.Field(i).Interface()
// 				logger.Info(" %v:		%v", key, value)
// 			}
// 			// individual report statistics
// 			totalSession = totalSession + flags.IsCreatedSession.Load()

// 			for apiID, rtts := range flags.RTT {
// 				for _, rtt := range rtts {
// 					totalRtt[apiID] = totalRtt[apiID] + rtt
// 					apiCall[apiID]++
// 					requestCount++
// 				}
// 			}
// 			for apiID, cnt := range flags.ErrorCount {
// 				errorCounts[apiID] = errorCounts[apiID] + cnt
// 			}
// 			if sts.startTime.After(latestStartTime) {
// 				latestStartTime = sts.startTime
// 			}
// 		}
// 	}

// 	clientCount := len(rpt.statistics)
// 	duration := time.Since(rpt.startTime)
// 	rps := float32(requestCount) / (float32(duration) / float32(time.Second))
// 	fmt.Printf("~~~~~~~~~~~~~~~~~~~~~~~~ report summary ~~~~~~~~~~~~~~~~~~~~~~~~ \n")
// 	fmt.Printf("time: %v \n", duration)
// 	fmt.Printf("client count: %v  \n", clientCount)
// 	fmt.Printf("session create rate: %.1f%%  \n", createRate*100)
// 	fmt.Printf("all clients spawned in: %v  \n", latestStartTime.Sub(rpt.startTime))
// 	fmt.Println()
// 	fmt.Printf("total session created: %v \n", totalSession)
// 	fmt.Printf("RPS (Cosmos): %.2f times / second\n", rps)
// 	fmt.Printf("Average RTT:\n")
// 	for apiID, rtt := range totalRtt {
// 		apiIDstr := fmt.Sprintf("%d", apiID)
// 		fmt.Printf("   [%d]%s: %v\n", apiID, strings.Repeat(" ", 5-len(apiIDstr)), rtt/time.Duration(apiCall[apiID]))
// 	}
// 	fmt.Printf("Error Rate:\n")
// 	for apiID, callCount := range apiCall {
// 		apiIDstr := fmt.Sprintf("%d", apiID)
// 		errorCount := errorCounts[apiID]
// 		fmt.Printf("   [%d]%s: %.3f%% (%d times called)\n", apiID, strings.Repeat(" ", 5-len(apiIDstr)), float32(errorCount)/float32(callCount), callCount)
// 	}
// 	fmt.Printf("~~~~~~~~~~~~~~~~~~~~~~~~ report summary ~~~~~~~~~~~~~~~~~~~~~~~~ \n")

// }
