package mysql

import (
	"fmt"
	"go-web-app/models"
	"go-web-app/pkg/snowflake"
	"go-web-app/pkg/todaytime"
	"time"
)

func AlarmOnlineInit(host *models.ParamStatistics) (hostgetdata []models.Ararmlist, err error) {
	sqlStr := ` select a.hostid as alarmid,
        a.alarmtype,
        hostlist.hostowner as alarmhostonwer,
        hostlist.hostname as alarmhostname,
        hostlist.hostip as alarmhostip,
        a.alarminfo,
        a.alarmstarttime  
 from alarmstatistics as a join 
     hostlist on a.hostid=hostlist.hostid
 where a.alarmstatus=1;`
	if err := db.Select(&hostgetdata, sqlStr); err != nil {
		return hostgetdata, err
	}
	return
}
func AlarmInit(host *models.ParamStatistics) (hostgetdata []models.Ararmlist, err error) {
	sqlStr := ` select a.hostid as alarmid,a.alarmtype,a.alarmstatus,a.alarmowner,hostlist.hostname as alarmhostname  from alarmsetting as a join hostlist on a.hostid=hostlist.hostid;`
	if err := db.Select(&hostgetdata, sqlStr); err != nil {
		return hostgetdata, err
	}
	return
}
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

func AlarmEdit(host *models.ParamStatistics) (n int64, err error) {
	sqlStr := `update alarmsetting set alarmtype=?,alarmstatus=?,alarmowner=? where hostid=? `
	ret, err := db.Exec(sqlStr,
		host.AlarmType,
		host.AlarmStatus,
		host.AlarmHostOnwer,
		host.Alarmid,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	n, err = ret.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Printf("更新数据为 %d 条\n", n)
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
func AlarmTodayTotal(host *models.ParamStatistics) (total int, err error) {
	now := time.Now()

	sqlStr := `select count(id)  from alarmstatistics where alarmstarttime > ? `
	if err := db.Get(&total, sqlStr, now.Format("2006-01-02")+" 00:00:00"); err != nil {
		return total, err
	}
	return
}
func AlarmAddToday(host *models.ParamStatistics) (total int, err error) {
	now := time.Now()
	host.Hostline = 1
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
