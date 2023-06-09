package config

type Config struct {
	ActivationCode string `mapstructure:"activation_code"` //激活码授权
	Debug          bool   //调试部分日志
	Mysql          struct {
		Host        string //服务器
		Port        int    //端口
		Name        string //数据库名称
		User        string //用户名
		Password    string //密码
		Local       string //时区
		Prefix      string //表前缀
		Charset     string //编码集
		MaxIdleConn int    `mapstructure:"max_idle_conn"` //设置空闲连接池中连接的最大数量
		MaxOpenConn int    `mapstructure:"max_open_conn"` //开数据库连接的最大数量
		MaxLifeTime int    `mapstructure:"max_life_time"` //设置了连接可复用的最大时间(单位分钟)
	}
	//redis配置
	Redis struct {
		Host     string //服务器
		Port     int    //端口
		Password string //密码
		Db       int    //选择的库
	}
	Mongo struct {
		Host     string //服务器
		Port     int    //端口
		Name     string //数据库名称
		User     string //用户名
		Password string //密码
	}

	//jwt配置
	Jwt struct {
		Key    string //jwt秘钥
		Expire int64  //有效时间,单位秒
	}

	//Zookeeper服务发现
	Zookeeper struct {
		Root    string   //根目录
		Servers []string //服务地址列表
	}
	RabbitMq struct {
		User     string //用户名
		Password string //密码
		Host     string //服务器地址
		Port     int    //服务器端口
	} `mapstructure:"rabbitmq"`

	Services struct {
		//网关层,用于动态分配目标服务器
		Gateway struct {
			HttpPort int `mapstructure:"http_port"` //http服务端口
		}
		//接入层(需要单独配置唯一ID)
		Connect struct {
			Id      uint32 //服务ID
			WsPort  int    `mapstructure:"ws_port"`  //ws服务端口
			TcpPort int    `mapstructure:"tcp_port"` //tcp服务端口
			RpcPort int    `mapstructure:"rpc_port"` //rpc服务端口
		}
		//业务层(支持无状态多服务器)
		Logic struct {
			Id      uint32 //这里暂时不确定是否需要设定ID,因为需要雪花算法
			RpcPort int    `mapstructure:"rpc_port"` //rpc服务端口
		}
		//消息中心(支持无状态多服务器)
		Message struct {
		}
	}
}
