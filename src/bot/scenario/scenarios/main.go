package scenarios

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"{0}/bot/scenario/lib/log"
	"{0}/bot/scenario/lib/report"

	"github.com/Diarkis/diarkis/util"
	uuidv4 "github.com/Diarkis/diarkis/uuid/v4"
)

type GlobalParams struct {
	ScenarioName    string `json:"scenarioName"`
	ScenarioPattern string `json:"scenarioPattern"`
	Host            string `json:"host"`
	Interval        int    `json:"interval"`
	Duration        int    `json:"duration"`
	HowMany         int    `json:"howmany"`
	LogLevel        int    `json:"logLevel"`
	ReceiveByteSize int    `json:"receiveByteSize"`
	UDPSendInterval int64  `json:"udpSendInterval"`
	RawConfigs      struct {
		CommonParams   map[string]any
		ScenarioParams map[string]any
	}
}

type Scenario interface {
	// GetName() string
	Run(globalParams *GlobalParams) error
	// todo: @params: error:
	ParseParam(index int, params []byte) error
	// todo: doc
	OnScenarioEnd() error
	// should return global report and individual report.
	// todo: readme
	WriteReport() (*report.Report, map[string]*report.Report)
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
	DefaultValue any    `json:"defaultValue"`
	Options      []any  `json:"options"`
	OptionRates  []any  `json:"optionRates"`
	IsRandom     bool   `json:"isRandom"`
	IsSequential bool   `json:"isSequential"`
	Type         string `json:"type"`
}

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
						ret[key] = index + int(baseValue)
						continue
					}

					// set sequential value for string value
					baseValueStr, ok := apiParam.DefaultValue.(string)
					if ok {
						baseValue, err := strconv.Atoi(baseValueStr)
						if err == nil {
							ret[key] = strconv.Itoa(index + baseValue)
							continue
						}
					}

					ret[key] = index
				} else {

					// default value
					ret[key] = apiParam.DefaultValue

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

func DynamicCast(sType string, target any) any {
	// todo: cover other types
	switch sType {
	case "uint8":
		if iValue, ok := target.(float64); ok {
			return uint8(iValue)
		}
	case "uint32":
		if iValue, ok := target.(float64); ok {
			return uint32(iValue)
		}
	}
	return target

}
