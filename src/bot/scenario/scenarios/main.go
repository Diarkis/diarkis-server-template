package scenarios

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"{0}/bot/scenario/lib/log"

	"github.com/Diarkis/diarkis/util"
	uuidv4 "github.com/Diarkis/diarkis/uuid/v4"
)

type UserState struct {
	sync.RWMutex
	// map[UserID]map[key]value
	state map[string]map[string]any
}

type GlobalParams struct {
	// // readonly params
	// Scenario name you want to run. Set in ScenarioFactoryList.
	ScenarioName string `json:"scenarioName"`
	// Scenario pattern you want to run. Set in the root level of config json file. (eg. OneTicket, RandomTicket and RatedTicket in ticket.json)
	ScenarioPattern string `json:"scenarioPattern"`
	// Interval to spawn clients.
	Interval int `json:"interval"`
	// Duration how long the scenario should run in Second. set 0 if it's not looping scenario
	Duration int `json:"duration"`
	// time to allow users idling in Second. If user does not do any command send or get push/response during IdleDuration, OnIdle() is triggered
	IdleDuration int `json:"idleDuration"`
	// the number of clients to be spawned
	HowMany            int            `json:"howmany"`
	Host               string         `json:"host"`
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
	// // updatable params
	sync.RWMutex
	UserState *UserState
}

// Scenario represents a scenario for a specific task.
type Scenario interface {
	// GetUserID should return the UID used for HTTP auth
	GetUserID() string

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

	// OnIdle is called when the user does not send or get push/response for a while.
	// idle time can be specified in IdleDuration
	OnIdle()
}

var logger = log.New("BOT/SCENARIO")

var ScenarioFactoryList map[string]func() Scenario = map[string]func() Scenario{
	"Connect": NewConnectScenario,
	"Ticket":  NewTicketScenario,
}

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

func NewUserState() *UserState {
	us := &UserState{}
	us.state = map[string]map[string]any{}
	return us
}

func (us *UserState) Init(userID string) {
	us.Lock()
	defer us.Unlock()
	us.state[userID] = map[string]any{}
}

func (us *UserState) Get(userID string, key string) any {
	us.RLock()
	defer us.RUnlock()
	if itf, ok := us.state[userID]; ok {
		if val, ok := itf[key]; ok {
			return val
		}
	}
	return nil
}

func (us *UserState) Set(userID string, key string, value any) error {
	us.Lock()
	defer us.Unlock()
	if _, ok := us.state[userID]; !ok {
		return errors.New("User does not exist. Call Init first.")
	}
	us.state[userID][key] = value
	return nil
}

// returns a list of userID that hit passed key and value
func (us *UserState) Search(key string, value any, limit int) []string {
	us.RLock()
	defer us.RUnlock()
	var list []string
	for userID, state := range us.state {
		if v, _ := state[key]; v == value {
			list = append(list, userID)
			if len(list) == limit {
				break
			}
		}
	}
	return list
}

func (us *UserState) Exists(userID string) bool {
	us.RLock()
	defer us.RUnlock()
	_, ok := us.state[userID]
	return ok
}

func (gp *GlobalParams) GenerateParams(index int) ([]byte, error) {
	return GenerateParams(index, gp.Raw.CommonParams, gp.Raw.ScenarioParams, gp.Raw.ParamsFromAPI)
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
