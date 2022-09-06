package types

import (
	"bytes"
	"encoding/binary"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/panjf2000/gnet/v2"
)

type ImHeadDataPack struct {
	ProtocolVersion uint16     //协议版本号
	DataFormat      DataFormat //数据交互协议,Json、Pb
	Command         uint16     //命令码
	Sequence        uint32     //数据序列号
	Length          uint32     //数据包长度
}

const (
	DataPackHeaderLength        = 2 + 2 + 2 + 4 + 4 //数据包长度
	ProtocolVersion      uint16 = 0x01              //协议版本号
)

type DataFormat uint16

const (
	DataFormatJson DataFormat = 0x01 //Json数据交互
	DataFormatPb   DataFormat = 0x02 //Pb数据交互
)

type XtTalkTcpCodec struct {
}

// DecodeBytes 拆包完整数据为包头,数据流,错误
func (XtTalkTcpCodec) DecodeBytes(data []byte) (*ImHeadDataPack, []byte, error) {
	if DataPackHeaderLength > len(data) {
		return nil, nil, gerror.Newf("数据包长度不足,无法解析")
	}
	headBuffer := bytes.NewBuffer(data[0:DataPackHeaderLength])
	var dataHead ImHeadDataPack
	if err := binary.Read(headBuffer, binary.LittleEndian, &dataHead.ProtocolVersion); err != nil {
		return nil, nil, gerror.Newf("读取数据协议版本失败: %s", err.Error())
	}
	if err := binary.Read(headBuffer, binary.LittleEndian, &dataHead.DataFormat); err != nil {
		return nil, nil, gerror.Newf("读取数据协议失败: %s", err.Error())
	}
	if err := binary.Read(headBuffer, binary.LittleEndian, &dataHead.Command); err != nil {
		return nil, nil, gerror.Newf("读取数据命令失败: %s", err.Error())
	}
	if err := binary.Read(headBuffer, binary.LittleEndian, &dataHead.Sequence); err != nil {
		return nil, nil, gerror.Newf("读取数据序列号失败: %s", err.Error())
	}
	if err := binary.Read(headBuffer, binary.LittleEndian, &dataHead.Length); err != nil {
		return nil, nil, gerror.Newf("读取数据长度失败: %s", err.Error())
	}
	dataBytes := data[DataPackHeaderLength:]
	return &dataHead, dataBytes, nil
}
func (XtTalkTcpCodec) Decode(conn gnet.Conn) (*ImHeadDataPack, []byte, error) {
	//拆包
	headBytes, err := conn.Next(DataPackHeaderLength)
	if err != nil {
		return nil, nil, gerror.Newf("读取包头失败: %s", err.Error())
	}
	headBuffer := bytes.NewBuffer(headBytes)
	var dataHead ImHeadDataPack
	if err := binary.Read(headBuffer, binary.LittleEndian, &dataHead.ProtocolVersion); err != nil {
		return nil, nil, gerror.Newf("读取数据协议版本失败: %s", err.Error())
	}
	if err := binary.Read(headBuffer, binary.LittleEndian, &dataHead.DataFormat); err != nil {
		return nil, nil, gerror.Newf("读取数据协议失败: %s", err.Error())
	}
	if err := binary.Read(headBuffer, binary.LittleEndian, &dataHead.Command); err != nil {
		return nil, nil, gerror.Newf("读取数据命令失败: %s", err.Error())
	}
	if err := binary.Read(headBuffer, binary.LittleEndian, &dataHead.Sequence); err != nil {
		return nil, nil, gerror.Newf("读取数据序列号失败: %s", err.Error())
	}
	if err := binary.Read(headBuffer, binary.LittleEndian, &dataHead.Length); err != nil {
		return nil, nil, gerror.Newf("读取数据长度失败: %s", err.Error())
	}
	//判断当前数据协议版本号
	if dataHead.ProtocolVersion != ProtocolVersion {
		return nil, nil, gerror.Newf("客户端数据协议版本不匹配")
	}
	//读取数据包长度
	dataBytes, err := conn.Next(int(dataHead.Length))
	if err != nil {
		return nil, nil, gerror.Newf("读取详细数据失败: %s", err.Error())
	}
	return &dataHead, dataBytes, nil
}
