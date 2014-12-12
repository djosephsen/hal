package haltest

import (
	"github.com/djosephsen/hal"
)

type Robot struct {
	hal.Robot
}

type ResponseRecorder struct {
}

func NewRecorder() *ResponseRecorder {
	return &ResponseRecorder{}
}
