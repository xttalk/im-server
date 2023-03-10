package core

import (
	"XtTalkServer/internal"
	"context"
	"math/rand"
	"runtime"
	"time"
)

var Entry = new(_Entry)

type _Entry struct {
}

func (_Entry) ExecStart(ctx context.Context, modules []internal.InitCtx) {
	rand.Seed(time.Now().UTC().UnixNano()) //随机数种子
	runtime.GOMAXPROCS(runtime.NumCPU())   //设置cpu运行核心数量

	for _, m := range modules {
		m.Initialize(ctx)
	}
}
