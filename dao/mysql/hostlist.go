package mysql

import (
	"errors"
	"fmt"
	"go-web-app/models"
)

var (
	ErrorHostExist = errors.New("主机已存在")
)

func Hostlistdataget(host *models.ParamHostDateGet) (hostgetdata []models.Hostlist, err error) {

	sqlStr := `select hostid,hostname,systemtype,hoststatus,hostip,hostlocation,hostowner,hostaddtime,hostnote from hostlist`
	if err := db.Select(&hostgetdata, sqlStr); err != nil {
		return hostgetdata, err
	}
	return
}
func Hostinfo(host *models.ParamHostDateGet) (hostgetdata interface{}, err error) {
	hostlistid := host.Hostid
	sqlStr := `select hostid,hostname,systemtype,hoststatus,hostip,hostlocation,hostowner,hostaddtime,hostnote,hostysteminfo from hostlist where hostip = ?`
	if err := db.Get(hostgetdata, sqlStr, string(hostlistid)); err != nil {
		return hostgetdata, err
	}
	return
}
func Hostedit(host *models.ParamHostDateGet) (n int64, err error) {
	fmt.Println(host)
	sqlStr := "update hostlist set hostname=?,systemtype=?,hoststatus=?,hostip=?,hostlocation=?,hostowner=?,hostnote=? where hostid=?"
	ret, err := db.Exec(sqlStr,
		host.HostName,
		host.SystemType,
		host.HostStatus,
		host.HostIP,
		host.HostLocation,
		host.HostOwner,
		host.HostNote,
		host.Hostid)
	if err != nil {
		return
	}
	n, err = ret.RowsAffected()
	if err != nil {
		return
	} else {
		fmt.Printf("更新数据为 %d 条\n", n)
	}
	return

}
func Hostadd(host *models.ParamHostDateGet) (theId int64, err error) {

	sqlStr := "insert into hostlist(hostname,systemtype,hoststatus,hostip,hostlocation,hostowner,hostnote,hostid) values (?,?,?,?,?,?,?,?)"
	ret, err := db.Exec(sqlStr,
		host.HostName,
		host.SystemType,
		host.HostStatus,
		host.HostIP,
		host.HostLocation,
		host.HostOwner,
		host.HostNote,
		host.Hostid)
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
func Hostdel(host *models.ParamHostDateGet) (n int64, err error) {

	sqlStr := "delete  from hostlist where hostid=?"
	ret, err := db.Exec(sqlStr, host.Hostid)
	if err != nil {
		return
	}
	n, err = ret.RowsAffected()
	if err != nil {
		return
	} else {
		fmt.Printf("删除数据为 %d 条\n", n)
	}
	return

}
func Hostcheck(host *models.ParamHostDateGet) (err error) {
	sqlStr := `select count(hostid)  from hostlist where hostip = ?`
	var count int
	if err := db.Get(&count, sqlStr, host.HostIP); err != nil {
		return err
	}
	if count > 0 {
		return ErrorHostExist
	}

	return
}
