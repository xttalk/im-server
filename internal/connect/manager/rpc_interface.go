package manager

import (
	"XtTalkServer/pb"
	"context"
)

var RpcApi = new(rpcApi)

type rpcApi struct {
}

func (rpcApi) LogicData(ctx context.Context, req *pb.LogicDataReq, res *pb.LogicDataRes) error {
	client, err := RpcClient.GetLogicClient(ctx)
	if err != nil {
		return err
	}
	return client.Call(ctx, "LogicData", req, res)
}
