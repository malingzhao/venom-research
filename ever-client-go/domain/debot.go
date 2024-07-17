package domain

import "math/big"

var DebotErrorCode map[string]int

type (
	// DebotHandle - Handle of registered in SDK debot.
	DebotHandle = int

	// DebotAction - Describes a debot action in a Debot Context.
	DebotAction struct {
		Description string `json:"description"`
		Name        string `json:"name"`
		ActionType  int    `json:"action_type"`
		To          int    `json:"to"`
		Attributes  string `json:"attributes"`
		Misc        string `json:"misc"`
	}

	// DebotInfo - Describes DeBot metadata.
	DebotInfo struct {
		Name        string   `json:"name,omitempty"`
		Version     string   `json:"version,omitempty"`
		Publisher   string   `json:"publisher,omitempty"`
		Caption     string   `json:"caption,omitempty"`
		Author      string   `json:"author,omitempty"`
		Support     string   `json:"support,omitempty"`
		Hello       string   `json:"hello,omitempty"`
		Language    string   `json:"language,omitempty"`
		Dabi        string   `json:"dabi,omitempty"`
		Icon        string   `json:"icon,omitempty"`
		Interfaces  []string `json:"interfaces"`
		DabiVersion string   `json:"dabiVersion"`
	}

	DebotActivity struct {
		ValueEnumType interface{}
	}

	DebotActivityTransaction struct {
		Msg              string     `json:"msg"`
		Dst              string     `json:"dst"`
		Out              []Spending `json:"out"`
		Fee              *big.Int   `json:"fee"`
		Setcode          bool       `json:"setcode"`
		Signkey          string     `json:"signkey"`
		SigningBoxHandle int        `json:"signing_box_handle"`
	}

	// Spending - Describes how much funds will be debited from the target contract balance as a result of the transaction.
	Spending struct {
		Amount *big.Int
		Dst    string
	}

	// ParamsOfInit - Parameters to init DeBot.
	ParamsOfInit struct {
		Address string `json:"address"`
	}

	// RegisteredDebot - Structure for storing debot handle returned from init function.
	RegisteredDebot struct {
		DebotHandle DebotHandle `json:"debot_handle"`
		DebotAbi    string      `json:"debot_abi"`
		Info        *DebotInfo  `json:"info"`
	}

	// ParamsOfStart - Parameters to start DeBot. DeBot must be already initialized with init() function.
	ParamsOfStart struct {
		DebotHandle DebotHandle `json:"debot_handle"`
	}

	// ParamsOfAppDebotBrowser - Debot Browser callbacks.
	ParamsOfAppDebotBrowser struct {
		ValueEnumType interface{}
	}

	ParamsOfAppDebotBrowserLog struct {
		Msg string `json:"msg"`
	}

	ParamsOfAppDebotBrowserSwitch struct {
		ContextID int `json:"context_id"`
	}

	ParamsOfAppDebotBrowserSwitchCompleted struct{}

	ParamsOfAppDebotBrowserShowAction struct {
		Action *DebotAction `json:"action"`
	}

	ParamsOfAppDebotBrowserInput struct {
		Prompt string `json:"prompt"`
	}

	ParamsOfAppDebotBrowserGetSigningBox struct{}

	ParamsOfAppDebotBrowserInvokeDebot struct {
		DebotAddr string       `json:"debot_addr"`
		Action    *DebotAction `json:"action"`
	}

	ParamsOfAppDebotBrowserSend struct {
		Message string `json:"message"`
	}

	ParamsOfAppDebotBrowserApprove struct {
		Activity *DebotActivity `json:"activity"`
	}

	// ResultOfAppDebotBrowser - Returning values from Debot Browser callbacks.
	ResultOfAppDebotBrowser struct {
		ValueEnumType interface{}
	}

	ResultOfAppDebotBrowserInput struct {
		Value string `json:"value"`
	}

	ResultOfAppDebotBrowserGetSigningBox struct {
		SigningBox SigningBoxHandle `json:"signing_box"`
	}

	ResultOfAppDebotBrowserInvokeDebot struct{}

	ResultOfAppDebotBrowserApprove struct {
		Approved bool `json:"approved"`
	}

	// ParamsOfFetch - Parameters to fetch DeBot metadata.
	ParamsOfFetch struct {
		Address string `json:"address"`
	}

	ResultOfFetch struct {
		Info *DebotInfo `json:"info"`
	}

	// ParamsOfExecute - Parameters for executing debot action.
	ParamsOfExecute struct {
		DebotHandle DebotHandle  `json:"debot_handle"`
		Action      *DebotAction `json:"action"`
	}

	// ParamsOfSend - Parameters of send function.
	ParamsOfSend struct {
		DebotHandle DebotHandle `json:"debot_handle"`
		Message     string      `json:"message"`
	}

	ParamsOfRemove struct {
		DebotHandle DebotHandle `json:"debot_handle"`
	}

	DebotUseCase interface {
		Init(*ParamsOfInit, AppDebotBrowser) (*RegisteredDebot, error)
		Start(*ParamsOfStart) error
		Fetch(*ParamsOfFetch) (*ResultOfFetch, error)
		Execute(*ParamsOfExecute) error
		Send(*ParamsOfSend) error
		Remove(*ParamsOfRemove) error
	}
)

func init() {
	DebotErrorCode = map[string]int{
		"DebotStartFailed":           801,
		"DebotFetchFailed":           802,
		"DebotExecutionFailed":       803,
		"DebotInvalidHandle":         804,
		"DebotInvalidJsonParams":     805,
		"DebotInvalidFunctionId":     806,
		"DebotInvalidAbi":            807,
		"DebotGetMethodFailed":       808,
		"DebotInvalidMsg":            809,
		"DebotExternalCallFailed":    810,
		"DebotBrowserCallbackFailed": 811,
		"DebotOperationRejected":     812,
		"DebotNoCode":                813,
	}
}
