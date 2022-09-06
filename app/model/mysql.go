package model

// User 用户表
type User struct {
	ID           uint64 //用户ID
	Username     string //用户账号(唯一)
	Nickname     string //用户昵称
	Email        string //用户邮箱(唯一)
	Password     string //密码
	Salt         string //密码盐
	RegisterTime int64  //注册时间
	RegisterIp   string //注册IP
	LoginTime    int64  //上次登录时间
	LoginIp      string //上次登录IP
	Status       int    //账号状态
}

type UserFriend struct {
	ID         uint64 //好友关系ID
	Uid        uint64 //用户ID
	Fid        uint64 //好友ID
	AddTime    int64  //成为好友时间
	RemarkName string //好友备注
	Friend     User   `gorm:"foreignKey:id;references:fid"` //好友的模型
}
