package handlers

import (
	"github.com/taimast/gotgproto/ext"
)

// CallbackResponse is the function which will be called on a handler's execution.
type CallbackResponse func(*ext.Context, *ext.Update) error
