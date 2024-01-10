package scenarios

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	bot_client "{0}/bot/scenario/lib/client"
	"{0}/bot/scenario/lib/log"
	"{0}/bot/scenario/lib/report"

	"github.com/Diarkis/diarkis/util"
	uuidv4 "github.com/Diarkis/diarkis/uuid/v4"
)

type Common struct {
	globalParams *GlobalParams
	client       *bot_client.UDPClient
	metrics      *report.CustomMetrics
	UID          string
}

type GlobalParams struct {
	ScenarioName       string         `json:"scenarioName"`
	ScenarioPattern    string         `json:"scenarioPattern"`
	Host               string         `json:"host"`
	Interval           int            `json:"interval"`
	Duration           int            `json:"duration"`
	HowMany            int            `json:"howmany"`
	LogLevel           int            `json:"logLevel"`
	ReceiveByteSize    int            `json:"receiveByteSize"`
	UDPSendInterval    int64          `json:"udpSendInterval"`
	MetricsInterval    int            `json:"metricsInterval"`
	Configs            map[string]int `json:"configs"`
	InputKeysForReport []string       `json:"keysForReport"`
	Raw                struct {
		CommonParams   map[string]any
		ScenarioParams map[string]any
		ParamsFromAPI  map[string]any
	}
}

// Scenario represents a scenario for a specific task.
type Scenario interface {
	// Run should perform the scenario's main logic.
	// It gives globalParams as a struct that is used for the running scenario.
	// It is shared with all scenarios
	Run(globalParams *GlobalParams) error

	// ParseParam should parse and sets the scenario-specific parameters from given params.
	// It gives the index of which describes the number of clients and the raw parameter data as input.
	ParseParam(index int, params []byte) error

	// OnScenarioEnd is called when the scenario execution has completed.
	// Implementations of this method should handle any necessary cleanup or finalization.
	OnScenarioEnd() error
}

var logger = log.New("BOT/SCENARIO")

var ScenarioFactoryList map[string]func() Scenario = map[string]func() Scenario{
	"Connect": NewConnectScenario,
	"Ticket":  NewTicketScenario,
}

// func StructTypeByName(name string) (reflect.Type, bool) {
// 	for _, t := range []reflect.Type{reflect.TypeOf(MyStruct1{}), reflect.TypeOf(MyStruct2{}), reflect.TypeOf(MyStruct3{})} {
// 		if t.Name() == name {
// 			return t, true
// 		}
// 	}
// 	return nil, false
// }

type ApiParamAttributes struct {
	DefaultValue any   `json:"defaultValue"`
	Options      []any `json:"options"`
	OptionRates  []any `json:"optionRates"`
	Range        struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"range"`
	IsRandom     bool   `json:"isRandom"`
	IsSequential bool   `json:"isSequential"`
	Type         string `json:"type"`
}

func (ap *ApiParamAttributes) isRangeSet() bool {
	return ap.Range.Min != 0 && ap.Range.Max != 0
}

// set index if you intend to generate "isSequential" value
func GenerateParams(index int, params ...map[string]any) ([]byte, error) {
	ret := map[string]any{}
	uuid, _ := uuidv4.New()
	seed := time.Now().UnixNano()
	rand := rand.New(rand.NewSource(seed))
	for _, paramsMap := range params {
		for key, paramIf := range paramsMap {
			if param, ok := paramIf.(map[string]any); !ok {
				ret[key] = paramIf
			} else {
				apiParam := ApiParamAttributes{}
				jsonParams, err := json.Marshal(param)
				if err != nil {
					continue
				}
				err = json.Unmarshal(jsonParams, &apiParam)
				if err != nil {
					continue
				}

				options := apiParam.Options
				optionRates := apiParam.OptionRates
				if len(options) > 0 && len(optionRates) > 0 {
					// random value with specified rate
					randomValue := rand.Float64()

					c := 0.0
					for i, option := range options {
						var rate float64
						if i < len(optionRates) {
							if rate_, ok := optionRates[i].(float64); ok {
								rate = rate_
							}
						}
						c += rate
						if randomValue < c || i == len(options)-1 {
							ret[key] = option
							break
						}

					}
				} else if apiParam.IsRandom {
					if len(options) > 0 {
						// random value from option
						ret[key] = options[util.RandomInt(0, len(options)-1)]
					} else {
						// uuid
						ret[key] = uuid.String
					}
				} else if apiParam.IsSequential {
					// set sequential value for int value
					baseValue, ok := apiParam.DefaultValue.(float64)
					if ok {
						v := int(baseValue)
						if apiParam.isRangeSet() {
							delta := apiParam.Range.Max - apiParam.Range.Min + 1
							ret[key] = int(math.Mod(float64(v), float64(delta))) + apiParam.Range.Min
							continue
						}
						ret[key] = index + v
						continue
					}

					// set sequential value for string value
					baseValueStr, ok := apiParam.DefaultValue.(string)
					if ok {
						baseValue, err := strconv.Atoi(baseValueStr)
						if err == nil {
							if apiParam.isRangeSet() {
								delta := apiParam.Range.Max - apiParam.Range.Min + 1
								ret[key] = int(math.Mod(float64(baseValue), float64(delta))) + apiParam.Range.Min
								continue
							}
							ret[key] = strconv.Itoa(index + baseValue)
							continue
						}
					}
					ret[key] = index
				} else if apiParam.isRangeSet() {
					// generate a value within the range
					ret[key] = util.RandomInt(apiParam.Range.Min, apiParam.Range.Max)
				} else {
					// default value
					if apiParam.DefaultValue != nil {
						ret[key] = apiParam.DefaultValue
					} else {
						ret[key] = param
					}

				}
				if apiParam.Type != "" {
					ret[key] = DynamicCast(apiParam.Type, ret[key])
				}
			}
		}
	}
	bytes, err := json.Marshal(ret)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func CollectInputParameters(keys []string, params ...map[string]any) map[string]string {
	ret := map[string]string{}
	for _, input := range params {
		for _, key := range keys {
			if paramIf, ok := input[key]; ok {
				if param, ok := paramIf.(map[string]any); ok {

					apiParam := ApiParamAttributes{}
					jsonParams, err := json.Marshal(param)
					if err != nil {
						continue
					}
					err = json.Unmarshal(jsonParams, &apiParam)
					if err != nil {
						continue
					}
					if apiParam.IsRandom {
						ret[key] = "Random"
						continue
					}
					if apiParam.IsSequential {
						ret[key] = "Sequential"
						continue
					}
					if apiParam.isRangeSet() {
						ret[key] = fmt.Sprintf("Range [%d-%d]", apiParam.Range.Min, apiParam.Range.Max)
					}
					ret[key] = fmt.Sprintf("%v", paramIf)
				} else {
					ret[key] = fmt.Sprintf("%v", paramIf)
				}

			}
		}
	}
	return ret
}

func DynamicCast(sType string, target any) any {
	// todo: cover other types
	// switch sType {
	// case "uint8":
	// 	if iValue, ok := target.(float64); ok {
	// 		return uint8(iValue)
	// 	}
	// case "uint32":
	// 	if iValue, ok := target.(float64); ok {
	// 		return uint32(iValue)
	// 	}
	// }
	return target

}
