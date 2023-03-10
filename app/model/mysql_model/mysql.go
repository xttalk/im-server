package mysql_model

const (
	IDX_UserFriendRequestOnlyRequest = "idx_only_request"
)

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

// UserFriend 好友表
type UserFriend struct {
	ID         uint64 //好友关系ID
	Uid        uint64 //用户ID
	Fid        uint64 //好友ID
	AddTime    int64  //成为好友时间
	RemarkName string //好友备注
	Friend     User   `gorm:"foreignKey:id;references:fid"` //好友的模型
}

type UserFriendRequestStatus int

const (
	UserFriendRequestStatusWait   UserFriendRequestStatus = 1 //等待处理
	UserFriendRequestStatusAccept UserFriendRequestStatus = 2 //同意
	UserFriendRequestStatusReject UserFriendRequestStatus = 3 //拒绝
)

// UserFriendRequest 好友验证消息
type UserFriendRequest struct {
	Id          int                     //验证ID
	FromUid     uint64                  //发起方
	ToUid       uint64                  //接受方
	RequestTime int64                   //请求时间
	Reason      string                  //验证消息
	Status      UserFriendRequestStatus //验证状态
	LastTime    int64                   //最后更新时间,指的是每次新增更新和首次新增

	RejectReason string //拒绝原因
}
