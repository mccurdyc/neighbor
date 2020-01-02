package search

import (
	"context"
)

type Backend interface {
	Search(context.Context, string)
	NextPage()
	ProcessResults()
}

type BackendConfig struct {
	NumberOfDesiredResults int
	Options                interface{}
}

type Factory func(context.Context, *BackendConfig) (Backend, error)
