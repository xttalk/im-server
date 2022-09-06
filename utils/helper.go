package utils

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"math"
	"time"
)

// JsonToMap Json转map[string]interface{}
func JsonToMap(json string) map[string]interface{} {
	result, err := gjson.DecodeToJson(json)
	if err != nil {
		return map[string]interface{}{}
	}
	return result.Map()
}
func DataToJson(v interface{}) string {
	return gjson.MustEncodeString(v)
}
func DataToJsonBytes(v interface{}) []byte {
	return gjson.MustEncode(v)
}

// CalcTime 计算时间
func CalcTime(s int64) string {
	_s := float64(s)
	var _i float64 = 60           //1分钟
	var _h float64 = 60 * 60      //一小时
	var _d float64 = 60 * 60 * 24 //一天
	switch {
	case _i > _s:
		return fmt.Sprintf("%v秒", _s)
	case _s >= _i && _h > _s:
		return fmt.Sprintf("%v分钟", math.Ceil(_s/_i))
	case _s >= _h && _d > _s:
		return fmt.Sprintf("%v小时", math.Ceil(_s/_h))
	default:
		return fmt.Sprintf("%v天", math.Ceil(_s/_d))
	}
}
func TimeFormat(time time.Time, format string) string {
	if 0 > time.Unix() {
		return ""
	}
	return gtime.New(time).Format(format)
}

// ComputeFileSize 计算文件大小,字节大小
func ComputeFileSize(byteSize int64) string {
	size := byteSize / 1024
	if 1024 > size {
		return fmt.Sprintf("%dKb", size)
	} else if size >= 1024 && 1024*1024 > size {
		return fmt.Sprintf("%dMb", size/1024)
	} else {
		return fmt.Sprintf("%dGb", size/1024/1024)
	}
}
