package manager

import (
	"XtTalkServer/app/conts"
	"XtTalkServer/global"
	"XtTalkServer/pb"
	"XtTalkServer/services/connect/types"
	"XtTalkServer/services/logic/logic_model"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
)

//核心服务集合

// CoreServer 服务器事件
var CoreServer = new(_CoreServer)

type _CoreServer struct {
	ServerId uint32 //服务器ID
}

// OnStart 服务已经启动
func (c *_CoreServer) OnStart() {
	c.ServerId = global.Config.Services.Connect.Id
}

// OnConnect 客户端连接
func (c *_CoreServer) OnConnect(client *Client) {
	client.Infof("客户端连接")
}

// OnClose 客户端断开
func (c *_CoreServer) OnClose(client *Client) {
	client.Infof("客户端断开")
	rk := fmt.Sprintf(conts.RK_ClientAuth, client.SessionId)

	if bytes, err := global.Redis.Get(client.Context, rk).Bytes(); err == nil {
		if err := global.Redis.Del(client.Context, rk).Err(); err != nil {
			client.Warningf("解析用户连接数据失败: %s", err.Error())
		}
		var userInfo logic_model.UserClient
		if err := gjson.Unmarshal(bytes, &userInfo); err != nil {
			client.Warningf("解析用户数据失败: %s", err.Error())
			return
		}
		//删除设备登录
		urk := fmt.Sprintf(conts.RK_UserDevice, userInfo.Uid)
		if err := global.Redis.HDel(client.Context, urk, client.SessionId).Err(); err != nil {
			client.Warningf("移除多设备登录记录失败: %s", err.Error())
		}
	}

	//client.SendBytes([]byte("客户端断开"))
}

// OnMessage 消息接收
func (c *_CoreServer) OnMessage(client *Client, head *types.ImHeadDataPack, data []byte) {
	client.Infof("收到客户端数据: %d -> %s", head.DataFormat, data)
	//数据解包
	//直接向业务服务层发起登录
	client.LastDataFormat = head.DataFormat
	//掉一个rpc
	req := pb.LogicDataReq{
		DataFormat: uint32(head.DataFormat),
		Command:    pb.Packet(head.Command),
		Data:       data,
		SessionId:  client.SessionId,
		ServerId:   c.ServerId,
	}
	res := pb.LogicDataRes{}
	if err := RpcApi.LogicData(client.Context, &req, &res); err != nil {
		client.Warningf("向逻辑层发送数据失败: %s", err.Error())
		client.Close()
		return
	}
	if !res.GetIsSend() {
		return //不发送
	}
	dataFormat := types.DataFormat(res.GetDataFormat())
	switch dataFormat {
	case types.DataFormatJson:
	case types.DataFormatPb:
	default:
		client.Warningf("未知的数据协议: %v", res.GetDataFormat())
		client.Close()
		return
	}
	client.Debugf("发送客户端数据: %s", res.Data)
	if err := client.SendClientPacket(head.Sequence, res.Command, types.DataFormat(res.GetDataFormat()), res.Data); err != nil {
		client.Warningf("发送客户端数据失败: %s", err.Error())
		client.Close()
		return
	}
}