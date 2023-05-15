package logic

import (
	"fmt"
	"go-web-app/dao/mysql"
	"go-web-app/models"
	"go-web-app/pkg/codeconversion"
	"go-web-app/pkg/jwt"
	"go-web-app/pkg/snowflake"
	"math"
	"os/exec"
	"strconv"
	"strings"
)

// 存放业务逻辑代码

func SignUp(p *models.ParamSignUp) (err error) {

	//1.判断用户可用性
	err = mysql.CheckUserByUsername(p.Username)
	if err != nil {
		//数据库查询出错
		return err
	}

	//2.生成UID
	userId := snowflake.GenID()
	//构造一个user示例
	user := &models.User{
		UserId:   userId,
		Username: p.Username,
		Password: p.Password,
	}
	//3.用户数据入库
	return mysql.InsertUser(user)

}

func Login(p *models.ParamLogin) (token string, err error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	//传递的是指针，可以获取到user.UserId
	if err := mysql.Login(user); err != nil {
		return "", err

	}
	fmt.Println(err)

	//生成JWTtoken
	return jwt.GenToken(user.UserId, user.Username)

}

func NetworkSentSpeed(p *models.ParamSystemGet) (s string, err error) {
	switch {
	case p.ParameterType == "cpu":
		//cpu使用率
		command := " typeperf -si 1 -sc 1 \"\\Processor(_Total)\\% Processor Time\" |findstr /V \"Processor\" |findstr /V \"?\" "
		commaninput := exec.Command("powershell.exe", command)
		output, _ := commaninput.CombinedOutput()
		outstring := codeconversion.ConvertByte2String(output, "GB18030")
		lastoutstring := strings.Split(outstring, ",\"")
		lastoutstringdel := strings.Split(lastoutstring[1], "\"")
		lastoupspeed, _ := strconv.ParseFloat(lastoutstringdel[0], 8)
		s = fmt.Sprintf("%.0f", math.Floor(lastoupspeed))
		return s, err
	case p.ParameterType == "uns":
		//网络上传速率
		command := "typeperf -si 1 -sc 1 \"\\Network Interface(*)\\Bytes Sent/sec\"  |findstr \",\" | findstr /V \"Interface\""
		commaninput := exec.Command("powershell.exe", command)
		output, _ := commaninput.CombinedOutput()
		outstring := codeconversion.ConvertByte2String(output, "GB18030")
		lastoutstring := strings.Split(outstring, "\",\"")
		lastoupspeed, _ := strconv.ParseFloat(lastoutstring[1], 8)
		lastoupspeedend := lastoupspeed / 1000
		s = fmt.Sprintf("%.2f", lastoupspeedend)
		return s, err
	case p.ParameterType == "dns":
		//网络下载速率
		command := "typeperf -si 1 -sc 1 \"\\Network Interface(*)\\Bytes Received/sec\"  |findstr \",\" | findstr /V \"Interface\""
		commaninput := exec.Command("powershell.exe", command)
		output, _ := commaninput.CombinedOutput()
		outstring := codeconversion.ConvertByte2String(output, "GB18030")
		lastoutstring := strings.Split(outstring, "\",\"")
		lastoupspeed, _ := strconv.ParseFloat(lastoutstring[1], 8)
		lastoupspeedend := lastoupspeed / 1000
		s = fmt.Sprintf("%.2f", lastoupspeedend)
		return s, err
	case p.ParameterType == "mp":
		//内存使用率
		command := "typeperf -si 1 -sc 1 \"\\Memory\\% Committed Bytes In Use\" |findstr /V \"Memory\" |findstr /V \"?\""
		commaninput := exec.Command("powershell.exe", command)
		output, _ := commaninput.CombinedOutput()
		outstring := codeconversion.ConvertByte2String(output, "GB18030")
		lastoutstring := strings.Split(outstring, ",\"")
		lastoutstringdel := strings.Split(lastoutstring[1], "\"")
		lastoupspeed, err := strconv.ParseFloat(lastoutstringdel[0], 8)
		s = fmt.Sprintf("%.0f", math.Floor(lastoupspeed))
		return s, err
	case p.ParameterType == "dt":
		//系统磁盘容量
		command := "wmic LogicalDisk where \"Caption='C:'\" get  Size /value | findstr \"Size\""
		commaninput := exec.Command("powershell.exe", command)
		output, _ := commaninput.CombinedOutput()
		outstring := codeconversion.ConvertByte2String(output, "GB18030")
		lastoutstring := strings.Split(outstring, "=")
		b := strings.Replace(lastoutstring[1], "\r\n", "", -1)
		c, err := strconv.ParseInt(b, 10, 64)
		d := fmt.Sprintf("%.0f", math.Floor(float64(c/1073741824)))
		return d, err
	case p.ParameterType == "fdp":
		//系统磁盘剩余空间占总比的
		command := "typeperf -si 1 -sc 1 \"\\LogicalDisk(C:)\\% Free Space\" |findstr /V \"Space\"| findstr /V \"?\""
		commaninput := exec.Command("powershell.exe", command)
		output, _ := commaninput.CombinedOutput()
		outstring := codeconversion.ConvertByte2String(output, "GB18030")
		lastoutstring := strings.Split(outstring, ",\"")
		lastoutstringdel := strings.Split(lastoutstring[1], "\"")
		lastoupspeed, _ := strconv.ParseFloat(lastoutstringdel[0], 8)
		s = fmt.Sprintf("%.0f", math.Floor(lastoupspeed))
		return s, err
	case p.ParameterType == "mt":
		//系统总内存
		command := " wmic ComputerSystem get TotalPhysicalMemory | findstr /V \"Total\""
		commaninput := exec.Command("powershell.exe", command)
		output, _ := commaninput.CombinedOutput()
		outstring := codeconversion.ConvertByte2String(output, "GB18030")
		a := strings.Replace(outstring, "\r\n", "", -1)
		a = strings.Replace(a, " ", "", -1)
		b, _ := strconv.ParseInt(a, 10, 64)
		s := ((b / 1073741824) + 1)
		d := strconv.FormatInt(s, 10)
		return d, err
	case p.ParameterType == "sup":
		//系统运行时间
		Time := []string{"Days", "Hours", "Minutes", "Seconds"}
		var TimeString []string
		for _, v := range Time {
			command := "(get-date) - (gcim Win32_OperatingSystem).LastBootUpTime | findstr /V  Total |findstr " + v
			commaninput := exec.Command("powershell.exe", command)
			output, _ := commaninput.CombinedOutput()
			outstring := codeconversion.ConvertByte2String(output, "GB18030")
			lastoutstring := strings.Split(outstring, ":")
			a := strings.Replace(lastoutstring[1], "\r\n", "", -1)
			TimeString = append(TimeString, a+v)
		}
		s := strings.Join(TimeString, "")
		s = strings.Replace(s, " ", "", -1)
		s = strings.Replace(s, "Days", "天", -1)
		s = strings.Replace(s, "Hours", "小时", -1)
		s = strings.Replace(s, "Minutes", "分钟", -1)
		s = strings.Replace(s, "Seconds", "秒", -1)
		return s, err

	}
	return s, err

}

func Hostdataget(p *models.ParamHostDateGet) (s interface{}, err error) {
	switch {
	case p.TypeOperation == "hostinit":
		s, err := mysql.Hostlistalarm(p)
		if err != nil {
			return s, err
		}
		return s, err
	case p.TypeOperation == "init":
		s, err := mysql.Hostlistdataget(p)
		if err != nil {
			return s, err
		}
		return s, err

	case p.TypeOperation == "add":

		err = mysql.Hostcheck(p)
		if err != nil {
			s = 0
			return s, err
		}
		s, err := mysql.Hostadd(p)
		if err != nil {
			return s, err
		}
		return s, err
	case p.TypeOperation == "del":

		s, err := mysql.Hostdel(p)
		if err != nil {
			return s, err
		}
		return s, err
	case p.TypeOperation == "edit":
		s, err := mysql.Hostedit(p)
		if err != nil {
			return s, err
		}
		return s, err
	}

	return
}

func Statistics(p *models.ParamStatistics) (s interface{}, err error) {
	switch {
	case p.StatisticsType == "alarmedit":
		s, err := mysql.AlarmEdit(p)
		if err != nil {
			return s, err
		}
		return s, err
	case p.StatisticsType == "alarmtotal":
		s, err := mysql.AlarmTotal(p)
		if err != nil {
			return s, err
		}
		return s, err
	case p.StatisticsType == "hosttotal":
		s, err := mysql.HostTotal(p)
		if err != nil {
			return s, err
		}
		return s, err
	case p.StatisticsType == "alarmonline":
		s, err := mysql.AlarmOnline(p)
		if err != nil {
			return s, err
		}
		return s, err
	case p.StatisticsType == "hostonline":
		s, err := mysql.HostOnline(p)
		if err != nil {
			return s, err
		}
		return s, err
	case p.StatisticsType == "hostaddtoday":
		s, err := mysql.HostAddToday(p)
		if err != nil {
			return s, err
		}
		return s, err

	case p.StatisticsType == "alarmdispose":
		s, err := mysql.AlarmDisposeToday(p)
		if err != nil {
			return s, err
		}
		return s, err
	case p.StatisticsType == "alarmaddtoday":
		s, err := mysql.AlarmAddToday(p)
		if err != nil {
			return s, err
		}
		return s, err
	case p.StatisticsType == "alarmtodaytotal":
		s, err := mysql.AlarmTodayTotal(p)
		if err != nil {
			return s, err
		}
		return s, err
	case p.StatisticsType == "alarmadd":
		s, err := mysql.AlarmAdd(p)
		if err != nil {
			return s, err
		}
		return s, err
	case p.StatisticsType == "alarminit":
		s, err := mysql.AlarmInit(p)
		if err != nil {
			return s, err
		}
		return s, err
	case p.StatisticsType == "alarmonlineinit":
		s, err := mysql.AlarmOnlineInit(p)
		if err != nil {
			return s, err
		}
		return s, err

	}

	return
}

func ClientData(p *models.ParamSystemGet) (Reply interface{}, err error) {
	switch {
	case p.ParameterType == "uptime":
		Reply, err := mysql.ClientUptime(p)
		if err != nil {
			return Reply, err
		}
		return Reply, err
	case p.ParameterType == "Confirm":
		Reply, err := mysql.ClientConfirm(p)
		if err != nil {
			return Reply, err
		}
		return Reply, err
	case p.ParameterType == "basemonitoring":
		Reply, err := mysql.BaseMonitoring(p)
		if err != nil {
			return Reply, err
		}
		return Reply, err
	case p.ParameterType == "systeminfo":
		Reply, err := mysql.ClientSystemInfo(p)
		if err != nil {
			return Reply, err
		}
		return Reply, err
	case p.ParameterType == "systeminfoget":
		Reply, err := mysql.ClientSystemInfoGet(p)
		if err != nil {
			return Reply, err
		}
		return Reply, err
	}
	return
}
