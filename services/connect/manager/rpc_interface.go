package manager

import (
	"XtTalkServer/pb"
	"context"
)

var RpcApi = new(_RpcApi)

type _RpcApi struct {
}

func (_RpcApi) LogicData(ctx context.Context, req *pb.LogicDataReq, res *pb.LogicDataRes) error {
	return LogicRpcClient.Call(ctx, "LogicData", req, res)
}
