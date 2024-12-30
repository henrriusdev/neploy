package websocket

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type (
	ActionType string
	InputType  string
)

type ProgressMessage struct {
	Type     string  `json:"type"`
	Progress float64 `json:"progress"`
	Message  string  `json:"message"`
}

// NewProgressMessage creates a new progress message
func NewProgressMessage(progress float64, message string) ProgressMessage {
	return ProgressMessage{
		Type:     "progress",
		Progress: progress,
		Message:  message,
	}
}

type ActionMessage struct {
	Type    string      `json:"type"`
	Action  string      `json:"action"`
	Data    interface{} `json:"data,omitempty"`
	Inputs  []Input     `json:"inputs"`
	Title   string      `json:"title"`
	Message string      `json:"message"`
}

// NewActionMessage creates a new action message
func NewActionMessage(actionType ActionType, title, message string, inputs []Input) ActionMessage {
	return ActionMessage{
		Type:    string(actionType),
		Title:   title,
		Message: message,
		Inputs:  inputs,
	}
}

type ActionResponse struct {
	Type   string                 `json:"type"`
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data"`
}

type Input struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Placeholder string   `json:"placeholder"`
	Value       string   `json:"value,omitempty"`
	Options     []string `json:"options,omitempty"`
	Required    bool     `json:"required"`
	Disabled    bool     `json:"disabled"`
	ReadOnly    bool     `json:"readOnly"`
	Order       int      `json:"order"`
}

// NewSelectInput creates a new select input
func NewSelectInput(name string, options []string) Input {
	return Input{
		Name:     name,
		Type:     "select",
		Options:  options,
		Required: true,
	}
}

// NewTextInput creates a new text input
func NewTextInput(name, placeholder string) Input {
	return Input{
		Name:        name,
		Type:        "text",
		Placeholder: placeholder,
		Required:    true,
	}
}

type Client struct {
	Conn *websocket.Conn
	Mu   sync.Mutex
}

// NewClient creates a new client
func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		Conn: conn,
	}
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

const (
	// Action types
	ActionTypeInfo     ActionType = "info"
	ActionTypeWarning  ActionType = "warning"
	ActionTypeCritical ActionType = "critical"
	ActionTypeError    ActionType = "error"
	ActionTypeRequest  ActionType = "request"
	ActionTypeResponse ActionType = "response"
	ActionTypeQuery    ActionType = "query"
)

const (
	// Input types
	InputTypeText     InputType = "text"
	InputTypePassword InputType = "password"
	InputTypeSelect   InputType = "select"
	InputTypeCheckbox InputType = "checkbox"
	InputTypeRadio    InputType = "radio"
	InputTypeDropzone InputType = "dropzone"
	InputTypeHidden   InputType = "hidden"
	InputTypeFile     InputType = "file"
	InputTypeTel      InputType = "tel"
	InputTypeEmail    InputType = "email"
	InputTypeUrl      InputType = "url"
	InputTypeNumber   InputType = "number"
	InputTypeRange    InputType = "range"
	InputTypeDate     InputType = "date"
	InputTypeTime     InputType = "time"
	InputTypeColor    InputType = "color"
	InputTypeCombo    InputType = "combo"
	InputTypeTextarea InputType = "textarea"
)
