package neato

import (
	"encoding/json"
)

type StandardResponse struct {
	Version int                    `json:"version"`
	ReqID   string                 `json:"reqId"`
	Result  string                 `json:"result"`
	Data    map[string]interface{} `json:"data"`
}

func (r StandardResponse) String() string {
	b, _ := json.Marshal(r)
	return string(b)
	// return fmt.Sprintf("version: %d, reqId: %s, result: %s", r.Version, r.ReqID, r.Result)
}

type RobotState int

const (
	StateInvalid RobotState = 0
	StateIdle               = 1
	StateBusy               = 2
	StatePaused             = 3
	StateError              = 4
)

type RobotAction int

const (
	ActionInvalid                 = 0
	ActionHouseCleaning           = 1
	ActionSpotCleaning            = 2
	ActionManualCleaning          = 3
	ActionDocking                 = 4
	ActionUserMenuActive          = 5
	ActionSuspendedCleaning       = 6
	ActionUpdating                = 7
	ActionCopyingLogs             = 8
	ActionRecoveringLocation      = 9
	ActionIECTest                 = 10
	ActionMapCleaning             = 11
	ActionExploringMap            = 12
	ActionAcquiringPersistentMap  = 13
	ActionCreatingAndUploadingMap = 14
	ActionSuspendedExploration    = 15
)

type StateResponse struct {
	StandardResponse
	Error             string                 `json:"error"`
	Alert             string                 `json:"alert"`
	State             int                    `json:"state"`
	Action            int                    `json:"action"`
	Cleaning          map[string]interface{} `json:"cleaning"`
	Details           map[string]interface{} `json:"details"`
	AvailableCommands map[string]interface{} `json:"availableCommands"`
	AvailableServices map[string]interface{} `json:"availableServices"`
	Meta              map[string]interface{} `json:"meta"`
}

func (r StateResponse) String() string {
	b, _ := json.Marshal(r)
	return string(b)
	//return fmt.Sprintf("version: %d, reqId: %s, result: %s", r.Version, r.ReqID, r.Result)
}
