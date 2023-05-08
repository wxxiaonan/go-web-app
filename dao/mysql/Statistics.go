package mysql

import (
	"fmt"
	"go-web-app/models"
	"go-web-app/pkg/snowflake"
	"go-web-app/pkg/todaytime"
	"time"
)

func AlarmAdd(host *models.ParamStatistics) (theId int64, err error) {
	sqlStr := "insert into alarmstatistics(alarmid,hostid,alarmstatus,alarmtype,alarminfo,alarmnote,alarmstarttime) values (?,?,?,?,?,?,?)"
	ret, err := db.Exec(sqlStr,
		snowflake.IdNum(),
		host.Hostid,
		1,
		host.AlarmType,
		host.AlarmInfo,
		host.AlarmNote,
		todaytime.NowTimeFull())
	if err != nil {
		return
	}
	theId, err = ret.LastInsertId()
	if err != nil {
		return theId, err
	} else {
		fmt.Printf("插入数据的id 为 %d. \n", theId)
	}
	return

}
func AlarmTotal(host *models.ParamStatistics) (total int, err error) {
	sqlStr := `select count(id)  from alarmstatistics`
	if err := db.Get(&total, sqlStr); err != nil {
		return total, err
	}
	return
}
func HostTotal(host *models.ParamStatistics) (total int, err error) {
	sqlStr := `select count(id)  from hostlist`
	if err := db.Get(&total, sqlStr); err != nil {
		return total, err
	}
	return
}

func AlarmOnline(host *models.ParamStatistics) (total int, err error) {
	host.Alarmline = 1
	sqlStr := `select count(id)  from alarmstatistics where alarmstatus= ?`
	if err := db.Get(&total, sqlStr, host.Alarmline); err != nil {
		return total, err
	}
	return
}
func HostOnline(host *models.ParamStatistics) (total int, err error) {
	host.Hostline = 1
	sqlStr := `select count(id)  from hostlist where hoststatus= ?`
	if err := db.Get(&total, sqlStr, host.Hostline); err != nil {
		return total, err

	}
	fmt.Println(total)
	return
}

func AlarmDisposeToday(host *models.ParamStatistics) (total int, err error) {
	now := time.Now()
	sqlStr := `select count(id)  from alarmstatistics where alarmstatus= ? and alarmstarttime > ? `
	if err := db.Get(&total, sqlStr, host.Hostline, now.Format("2006-01-02")+" 00:00:00"); err != nil {
		return total, err
	}
	return
}
func AlarmAddToday(host *models.ParamStatistics) (total int, err error) {
	now := time.Now()
	sqlStr := `select count(id)  from alarmstatistics where alarmstatus= ? and alarmstarttime > ? `
	if err := db.Get(&total, sqlStr, host.Hostline, now.Format("2006-01-02")+" 00:00:00"); err != nil {
		return total, err
	}
	return
}
func HostAddToday(host *models.ParamStatistics) (total int, err error) {
	host.Hostline = 1
	now := time.Now()
	sqlStr := `select count(id)  from hostlist where hoststatus= ? and hostaddtime > ? `
	if err := db.Get(&total, sqlStr, host.Hostline, now.Format("2006-01-02")+" 00:00:00"); err != nil {
		return total, err
	}
	return
}
