package container

import (
	"context"
)

type Config struct {
	Image string
	Cmd   []string
}

type Options struct {
	ShowLogs bool
	Force    bool
}

type Manager interface {
	PullImage(ctx context.Context, imageURL string) error
	CreateContainer(ctx context.Context, config Config) (string, error)
	StartContainer(ctx context.Context, id string) error
	WaitContainer(ctx context.Context, id string) error
	GetLogs(ctx context.Context, id string, opts Options) (string, error)
	RemoveContainer(ctx context.Context, id string, opts Options) error
}
