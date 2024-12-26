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
	Type    ActionType  `json:"type"`
	Data    interface{} `json:"data"`
	Inputs  []Input     `json:"inputs"`
	Title   string      `json:"title"`
	Message string      `json:"message"`
}

// NewActionMessage creates a new action message
func NewActionMessage(actionType ActionType, title, message string, inputs []Input) ActionMessage {
	return ActionMessage{
		Type:    actionType,
		Title:   title,
		Message: message,
		Inputs:  inputs,
	}
}

type Input struct {
	Name        string    `json:"name"`
	Type        InputType `json:"type"`
	Placeholder string    `json:"placeholder"`
	Value       *string   `json:"value"`
	Options     []string  `json:"options"`
	Required    bool      `json:"required"`
	Disabled    bool      `json:"disabled"`
	ReadOnly    bool      `json:"readOnly"`
	Order       int       `json:"order"`
}

// NewSelectInput creates a new select input
func NewSelectInput(name string, options []string) Input {
	return Input{
		Name:     name,
		Type:     InputTypeSelect,
		Options:  options,
		Required: true,
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
	ActionTypeCritical = "critical"
	ActionTypeError    = "error"
	ActionTypeRequest  = "request"
	ActionTypeResponse = "response"
	ActionTypeQuery    = "query"
)

const (
	InputTypeText     = "text"
	InputTypePassword = "password"
	InputTypeSelect   = "select"
	InputTypeCheckbox = "checkbox"
	InputTypeRadio    = "radio"
	InputTypeDropzone = "dropzone"
	InputTypeHidden   = "hidden"
	InputTypeFile     = "file"
	InputTypeTel      = "tel"
	InputTypeEmail    = "email"
	InputTypeUrl      = "url"
	InputTypeNumber   = "number"
	InputTypeRange    = "range"
	InputTypeDate     = "date"
	InputTypeTime     = "time"
	InputTypeColor    = "color"
	InputTypeCombo    = "combo"
	InputTypeTextarea = "textarea"
)
