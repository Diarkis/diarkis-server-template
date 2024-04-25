package report

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Diarkis/diarkis-server-template/bot/scenario/lib/log"

	"github.com/Diarkis/diarkis/util"
)

// 15 is default scraping interval for prometheus
var Interval = 15

var logger = log.New("BOT/REPORT")

type Report map[string]int
type metrics struct {
	name    string
	running bool
	sync.RWMutex
	// a number of metrics elements while scenario is being executed
	counter atomic.Uint32
	// a number of metrics elements in the current `interval`
	gauge atomic.Uint32
	// an array for elements
	values []float64
	// total values while scenario is being executed
	total float64
	// total values in the current `interval`
	subTotal float64
}

func NewMetrics(name string) *metrics {
	m := metrics{name: name}
	go m.start()
	return &m
}

func (m *metrics) start() {
	m.running = true
	var prevCount uint32
	var prevTotal float64
	for {
		time.Sleep(time.Duration(Interval) * time.Second)
		if !m.running {
			break
		}
		var total float64
		m.RLock()
		for _, v := range m.values {
			total += v
		}

		count := uint32(len(m.values))
		m.RUnlock()

		m.Lock()
		m.total = total
		m.subTotal = total - prevTotal
		m.Unlock()

		m.gauge.Store(count - prevCount)

		m.counter.Store(count)

		prevCount = count
		prevTotal = total
	}
}

func (m *metrics) Stop() {
	m.stop()
}

func (m *metrics) stop() {
	m.running = false

	var total float64
	m.RLock()
	for _, v := range m.values {
		total += v
	}

	count := uint32(len(m.values))
	m.RUnlock()

	m.Lock()
	m.total = total
	m.Unlock()

	m.counter.Store(count)

}
func (m *metrics) Add(delta float64) {
	m.Lock()
	defer m.Unlock()

	m.values = append(m.values, delta)
}
func (m *metrics) GetTotalCount() uint32 {
	return m.counter.Load()
}
func (m *metrics) GetSubTotalCount() uint32 {
	return m.gauge.Load()
}

func (m *metrics) GetTotal() float64 {
	m.RLock()
	defer m.RUnlock()
	return m.total
}

func (m *metrics) GetSubTotal() float64 {
	m.RLock()
	defer m.RUnlock()
	return m.subTotal
}

func (m *metrics) GetAverage() float64 {
	m.RLock()
	defer m.RUnlock()

	total := m.GetTotal()
	totalCount := m.GetTotalCount()
	if totalCount == 0 {
		return 0.0
	}
	return total / float64(totalCount)
}

func (m *metrics) GetIntervalAverage() float64 {
	m.RLock()
	defer m.RUnlock()
	subTotal := m.GetSubTotalCount()
	if subTotal == 0 {
		return 0.0
	}
	return m.subTotal / float64(subTotal)
}
func (m *metrics) Reset() {
	m.counter.Store(0)
	m.gauge.Store(0)
	m.values = []float64{}
	m.total = 0
	m.subTotal = 0
}

// // // // // // CCU // // // // // //
type ActiveUsers struct {
	sync.RWMutex
	list  map[string]time.Time
	gauge atomic.Int32
}

var activeUsers = &ActiveUsers{list: map[string]time.Time{}}

func IsActive(userID string) bool {
	activeUsers.RLock()
	lastTouched := activeUsers.list[userID]
	activeUsers.RUnlock()

	duration := time.Since(lastTouched)
	return duration < time.Second*time.Duration(Interval)
}
func TouchAsActiveUser(userID string) {
	activeUsers.RLock()
	now := time.Now()

	if lastTouched, ok := activeUsers.list[userID]; !ok || lastTouched.Add(time.Second*time.Duration(Interval)).Before(now) {
		activeUsers.gauge.Add(1)
	}
	activeUsers.RUnlock()

	activeUsers.Lock()
	activeUsers.list[userID] = now
	activeUsers.Unlock()

	decrement := func() {
		time.Sleep(time.Second * time.Duration(Interval))
		if lastTouched := activeUsers.list[userID]; lastTouched.Equal(now) {
			activeUsers.Lock()
			activeUsers.gauge.Add(-1)
			activeUsers.Unlock()
		}
	}
	go decrement()
}

func (au *ActiveUsers) GetMetrics() string {
	label := fmt.Sprintf("Bot_Active_Users", Interval)
	var metrics string
	metrics += fmt.Sprintf("# HELP %s number of users that issued a command in %d seconds\n", label, Interval)
	metrics += fmt.Sprintf("# TYPE %s gauge\n", label)
	metrics += fmt.Sprintf("%s %d\n", label, au.gauge.Load())

	return metrics
}
func (au *ActiveUsers) Stop() {
	// do nothing as this does not have a loop to stop
}

func (au *ActiveUsers) GetAsKV(kv *KeyValue) {
	(*kv)["active-users"] = au.gauge.Load()
}

func (au *ActiveUsers) Print() {
	logger.Notice("active users %d", au.gauge.Load())
}

// // // // // // ScenarioError // // // // // //
type ScenarioError struct {
	sync.RWMutex
	m *metrics
}

var scenarioError = &ScenarioError{}

// func (se *ScenarioError) Increment() {
// 	if se.m == nil {
// 		se.m = NewMetrics()
// 	}
// 	se.m.counter.Add(1)
// }

func IncrementScenarioError() {
	if scenarioError.m == nil {
		scenarioError.m = NewMetrics("error")
	}
	scenarioError.m.counter.Add(1)
}

// func ResetScenarioError() {
// 	scenarioError.m = &metrics{}
// }

func (se *ScenarioError) Reset() {
	se.m = &metrics{}
}

// // // // // // Custom Metrics // // // // // //
type CustomMetrics struct {
	sync.RWMutex
	m map[string]map[string]*metrics
}

var customMetrics = &CustomMetrics{
	m: map[string]map[string]*metrics{},
}

func NewCustomMetrics() *CustomMetrics {
	cm := CustomMetrics{
		m: map[string]map[string]*metrics{},
	}
	return &cm
}

// Increment adds 1 to a metrics value to the specified name and key.
// It will be collected and calculated for total, average and subtotal for a specific period.
func (cm_ *CustomMetrics) Increment(name, key string) {
	increment := func(cm *CustomMetrics) {
		cm.Lock()
		if _, ok := cm.m[name]; !ok {
			cm.m[name] = map[string]*metrics{}
		}
		if _, ok := cm.m[name][key]; !ok {
			cm.m[name][key] = NewMetrics(fmt.Sprintf("custom-%s-%s", name, key))
		}
		// cm.m[name][key].counter.Add(1)
		cm.m[name][key].Add(1.0)
		cm.Unlock()
	}

	// individual
	increment(cm_)
	// total
	increment(customMetrics)
}

// Add adds a metrics value to the specified name and key.
// It will be collected and calculated for total, average and subtotal for a specific period.
func (cm *CustomMetrics) Add(name, key string, value float64) {
	increment := func(cm *CustomMetrics) {
		cm.Lock()
		if _, ok := cm.m[name]; !ok {
			cm.m[name] = map[string]*metrics{}
		}
		if _, ok := cm.m[name][key]; !ok {
			cm.m[name][key] = NewMetrics(fmt.Sprintf("custom-%s-%s", name, key))
		}
		cm.m[name][key].Add(value)
		cm.Unlock()
	}

	// individual
	increment(cm)
	// total
	increment(customMetrics)
}

// Print outputs metrics to the log
func (cm *CustomMetrics) Print() {
	for name, keys := range cm.m {
		logger.Notice("metrics %s", name)
		for key, m := range keys {
			total := m.GetTotal()
			average := m.GetAverage()
			hasZeroDecimal := func(f float64) bool {
				return f == float64(int(f))
			}
			if (average == 1.0 || average == 0.0) && hasZeroDecimal(total) {
				logger.Notice(" key:%s total:%d", key, int(total))
			} else {
				logger.Notice(" key:%s total:%v  average:%v", key, total, average)
			}
		}
	}
}

func (cm *CustomMetrics) GetAsKV(kv *KeyValue) {
	for name, keys := range cm.m {
		for k, m := range keys {
			key := name
			if k != "" {
				key = strings.Join([]string{name, k}, "-")
			}

			var value any
			average := m.GetAverage()
			total := m.GetTotal()
			hasZeroDecimal := func(f float64) bool {
				return f == float64(int(f))
			}
			if (average == 1.0 || average == 0.0) && hasZeroDecimal(total) {
				value = total
			} else {
				key += "-average"
				value = average
			}

			(*kv)[key] = value
			logger.Debug("Getting custom metrics as key value key:%s value:%v", key, value)
		}
	}
	logger.Debug("Returning merged key value for custom metrics :%+v", *kv)
}

// GetMetrics returns a string for Prometheus metrics
func (cm *CustomMetrics) GetMetrics() string {
	var c strings.Builder
	var g strings.Builder
	var a strings.Builder
	if cm != nil && len(cm.m) > 0 {
		for name, keys := range cm.m {
			counterLabel := fmt.Sprintf("Bot_Custom_Metrics_%s_total", name)
			gaugeLabel := fmt.Sprintf("Bot_Custom_Metrics_%s", name, Interval)
			averageLabel := fmt.Sprintf("Bot_Custom_Metrics_%s_average", name)
			fmt.Fprintf(&c, "# HELP %s total value for custom metrics\n", counterLabel)
			fmt.Fprintf(&c, "# TYPE %s counter\n", counterLabel)
			fmt.Fprintf(&g, "# HELP %s total value for custom metrics %s in %d seconds\n", gaugeLabel, name, Interval)
			fmt.Fprintf(&g, "# TYPE %s gauge\n", gaugeLabel)
			fmt.Fprintf(&a, "# HELP %s average value for custom metrics %s in %d seconds\n", averageLabel, name, Interval)
			fmt.Fprintf(&a, "# TYPE %s gauge\n", averageLabel)
			for key, m := range keys {

				var labelMatcher string
				if key != "" {
					labelMatcher = fmt.Sprintf("{key=\"%s\"}", key)
				}
				m.RLock()
				fmt.Fprintf(&c, "%s%s %f\n", counterLabel, labelMatcher, m.GetTotal())
				fmt.Fprintf(&g, "%s%s %f\n", gaugeLabel, labelMatcher, m.GetSubTotal())
				fmt.Fprintf(&a, "%s%s %f\n", averageLabel, labelMatcher, m.GetIntervalAverage())
				m.RUnlock()
			}
		}
	}
	return util.StrConcat(g.String(), c.String(), a.String())
}

// Stop stops metrics loop that is set as Interval
func (cm *CustomMetrics) Stop() {

	if cm != nil && len(cm.m) > 0 {
		for _, keys := range cm.m {
			for _, m := range keys {
				m.Stop()
			}
		}
	}
}

func (cm *CustomMetrics) WriteCSV(name string) {
	var keys []string
	var values []string

	var kv = KeyValue{}
	cm.GetAsKV(&kv)

	for key, value := range kv {
		keys = append(keys, key)
		values = append(values, fmt.Sprintf("%v", value))
	}
	header := strings.Join(keys, ",")
	data := strings.Join(values, ",")
	output := strings.Join([]string{header, data}, "\n")

	filename := "Bot_Report_Custom_"
	filename += name
	filename += "_"
	filename += time.Now().Format("20060102150405")
	filename += ".csv"
	logger.Info("Writing report to csv file. file name:%s, data: \n%s", filename, output)
	util.WriteToTmp(filename, output)
}

func PrintCustomMetrics() {
	if customMetrics == nil {
		logger.Warn("Nothing called")
		return
	}
	customMetrics.Print()
}

// // // // // // Command Metrics // // // // // //
type CommandMetrics struct {
	sync.RWMutex
	cType string
	m     map[uint8]map[uint16]*metrics
	// metrics map[uint8]map[uint16]struct {
	// 	counter  atomic.Uint32
	// 	duration time.Duration
	// }
}

var callCommandMetrics = &CommandMetrics{cType: "Send"}
var pushCommandMetrics = &CommandMetrics{cType: "Push"}
var responseCommandMetrics = &CommandMetrics{cType: "Response"}

func (cm *CommandMetrics) Increment(ver uint8, cmd uint16) {
	cm.Lock()
	defer cm.Unlock()

	if cm.m == nil {
		cm.m = map[uint8]map[uint16]*metrics{}
	}
	if _, ok := cm.m[ver]; !ok {
		cm.m[ver] = map[uint16]*metrics{}
	}
	if _, ok := cm.m[ver][cmd]; !ok {
		cm.m[ver][cmd] = NewMetrics(fmt.Sprintf("command-%d-%d", ver, cmd))
	}
	cm.m[ver][cmd].Add(1.0)
}

func (cm *CommandMetrics) Print() {
	for ver, cmds := range cm.m {
		for cmd, m := range cmds {
			// logger.Notice(m.values)
			logger.Notice("[%s]	ver:%d cmd:%d count:%d", cm.cType, ver, cmd, m.GetTotalCount())
		}
	}
}

func (cm *CommandMetrics) GetAsKV(kv *KeyValue) {
	for ver, cmds := range cm.m {
		for cmd, m := range cmds {
			key := fmt.Sprintf("ver%d-cmd%d-%s", ver, cmd, cm.cType)
			value := m.GetTotal()
			(*kv)[key] = int(value)
			logger.Debug("Getting command metrics as key value key:%s value:%v", key, value)
		}
	}
	logger.Debug("Returning merged key value :%+v", *kv)
}

func (cm *CommandMetrics) GetMetrics() string {
	var g strings.Builder
	var c strings.Builder
	if cm != nil && len(cm.m) > 0 {
		gaugeLabel := fmt.Sprintf("Bot_Command_%s", cm.cType, Interval)
		counterLabel := fmt.Sprintf("Bot_Command_%s_total", cm.cType)
		fmt.Fprintf(&g, "# HELP %s sub total %s count for each commands in %d seconds\n", gaugeLabel, cm.cType, Interval)
		fmt.Fprintf(&g, "# TYPE %s gauge\n", gaugeLabel)
		fmt.Fprintf(&c, "# HELP %s total %s count for each commands\n", counterLabel, cm.cType)
		fmt.Fprintf(&c, "# TYPE %s counter\n", counterLabel)
		for ver, cmds := range cm.m {
			for cmd, m := range cmds {
				fmt.Fprintf(&g, "%s{ver=\"%d\",cmd=\"%d\"} %d\n", gaugeLabel, ver, cmd, m.GetSubTotalCount())
				fmt.Fprintf(&c, "%s{ver=\"%d\",cmd=\"%d\"} %d\n", counterLabel, ver, cmd, m.GetTotalCount())
			}
		}
	}
	return util.StrConcat(g.String(), c.String())
}

func (cm *CommandMetrics) Stop() {
	if cm != nil && len(cm.m) > 0 {
		for _, cmds := range cm.m {
			for _, m := range cmds {
				m.Stop()
			}
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
	return callCommandMetrics.m
}
func GetPushMetrics() map[uint8]map[uint16]*metrics {
	return pushCommandMetrics.m
}
func GetResponseMetrics() map[uint8]map[uint16]*metrics {
	return responseCommandMetrics.m
}

type KeyValue map[string]any
type MetricsCollector interface {
	Print()
	GetAsKV(kv *KeyValue)
	GetMetrics() string
	Stop()
	// Reset()
}

// metrics' that implement collector
var allMetrics = []MetricsCollector{
	activeUsers,
	callCommandMetrics,
	pushCommandMetrics,
	responseCommandMetrics,
	customMetrics,
}

func PrintAllMetrics() {
	for _, metrics := range allMetrics {
		metrics.Print()
	}
}

func WriteCSV(name string, inputs ...map[string]string) {
	var keys []string
	var values []string
	for _, input := range inputs {
		for k, v := range input {
			keys = append(keys, k)
			values = append(values, v)
		}
	}

	var kv = KeyValue{}
	for _, metrics := range allMetrics {
		metrics.GetAsKV(&kv)
	}
	for key, value := range kv {
		keys = append(keys, key)
		values = append(values, fmt.Sprintf("%v", value))
	}
	header := strings.Join(keys, ",")
	data := strings.Join(values, ",")
	output := strings.Join([]string{header, data}, "\n")

	filename := "Bot_Report_"
	filename += name
	filename += "_"
	filename += time.Now().Format("20060102150405")
	filename += ".csv"
	logger.Info("Writing report to csv file. file name:%s, data: \n%s", filename, output)
	util.WriteToTmp(filename, output)
}

func ResetAllMetrics() {
	callCommandMetrics.m = map[uint8]map[uint16]*metrics{}
	pushCommandMetrics.m = map[uint8]map[uint16]*metrics{}
	responseCommandMetrics.m = map[uint8]map[uint16]*metrics{}
	customMetrics.m = map[string]map[string]*metrics{}
	scenarioError.Reset()
}

func StopAllMetrics() {
	for _, metrics := range allMetrics {
		metrics.Stop()
	}
}

func GetPrometheusMetrics() string {
	var output string
	for _, metrics := range allMetrics {
		output += metrics.GetMetrics()
	}

	return output
}
