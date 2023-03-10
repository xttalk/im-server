package internal

import "context"

type InitCtx interface {
	Initialize(context.Context)
}
