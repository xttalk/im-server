package model

import (
	"XtTalkServer/pb"
	"github.com/gogf/gf/v2/container/gvar"
	"gorm.io/gorm"
	"math"
)

// NavPageReq 分页请求
type NavPageReq struct {
	Page int64 `d:"1" json:"page"`  //页数
	Size int64 `d:"10" json:"size"` //数量
}

// UseNav 使用数据分页
func (c *NavPageReq) UseNav(db *gorm.DB) *gorm.DB {
	return db.Offset(c.GetOffset()).Limit(c.GetLimit())
}

func (c *NavPageReq) GetPage() int {
	if 0 >= c.Page {
		return 1
	}
	return int(c.Page)
}
func (c *NavPageReq) GetOffset() int {
	page := c.GetPage()
	offset := (page - 1) * c.GetLimit()
	return offset
}
func (c *NavPageReq) GetLimit() int {
	return int(c.Size)
}

type NavPageRes struct {
	Page    int64 `json:"page"`    //当前页数
	Size    int64 `json:"size"`    //请求显示数量
	Count   int64 `json:"count"`   //最大数量
	MaxPage int64 `json:"maxPage"` //最大页数
}

func BuildNav(req NavPageReq, count int64) NavPageRes {
	maxPage := gvar.New(math.Ceil(float64(count) / (float64(req.Size)))).Int64()
	return NavPageRes{
		Page:    int64(req.GetPage()),
		Size:    int64(req.GetLimit()),
		MaxPage: maxPage,
		Count:   count,
	}
}
func BuildNavPb(req NavPageReq, count int64) *pb.Nav {
	maxPage := gvar.New(math.Ceil(float64(count) / (float64(req.Size)))).Int64()
	return &pb.Nav{
		CurrentPage: int64(req.GetPage()),
		MaxPage:     maxPage,
		Total:       count,
	}
}
