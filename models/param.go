package models

//定义请求的参数结构体

// 用户注册参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
	Email      string `json:"email"`
}

// 用户登陆参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ParamSystemGet struct {
	ParameterType string `json:"parametertype" binding:"required"`
}
type ParamHostDateGet struct {
	TypeOperation string `json:"typeoperation" binding:"required"`
	Hostlist
}
type Hostlist struct {
	Hostid         int    `json:"hostid" bindding:"required"`
	HostName       string `json:"hostname" bindding:"required"`
	SystemType     string `json:"systemtype" bindding:"required"`
	HostStatus     int    `json:"hoststatus" bindding:"required"`
	HostIP         string `json:"hostip" bindding:"required"`
	HostLocation   string `json:"hostlocation" bindding:"required"`
	HostOwner      string `json:"hostowner" bindding:"required"`
	HostAddTime    string `json:"hostaddtime" bindding:"required"`
	HostNote       string `json:"hostnote" bindding:"required"`
	HostSystemInfo string `json:"hostsysteminfo"`
}
