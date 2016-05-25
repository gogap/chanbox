package box

import (
	"golang.org/x/net/context"
)

type Option func(*Options)

type Options struct {
	Context context.Context
}
