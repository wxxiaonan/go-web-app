package mysql

import (
	"context"
	"encoding/json"
	"fmt"
	"go-web-app/models"
	"go-web-app/pkg/snowflake"
	"go-web-app/pkg/todaytime"
	"go-web-app/settings"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

const (
	JobDir  = "/cron/jobs/"
	JobKill = "/cron/kill/"
)

var (
	clinet  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	GJobmgr *models.JobMgr
	oldjob  *models.Job
)

type JobMgr struct {
	Kv     clientv3.KV
	Lease  clientv3.Lease
	Clinet *clientv3.Client
}

func InitCrontab(cfg *settings.EtcdConfig) (err error) {
	config := clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: time.Duration(cfg.DialTimeout) * time.Millisecond,
		Username:    cfg.Username,
		Password:    cfg.Password,
	}
	if clinet, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
	//获取KV和Lease的API子集
	kv = clientv3.NewKV(clinet)
	lease = clientv3.NewLease(clinet)
	//赋值单例
	GJobmgr = &models.JobMgr{
		Clinet: clinet,
		Kv:     kv,
		Lease:  lease,
	}

	return
}
func SaveJob(jobmgr *models.JobMgr, job models.CrontabJob) (oldJob *models.Job, err error) {

	var (
		jobKey   string
		jobValue []byte
		putResp  *clientv3.PutResponse
	)
	jobetcd := models.Job{
		Name:     job.JobName,
		Command:  job.JobShell,
		CronExpr: job.JobCronExpr,
	}
	jobKey = JobDir + job.JobName
	if jobValue, err = json.Marshal(jobetcd); err != nil {
		return
	}
	if putResp, err = jobmgr.Kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return
	}
	if putResp.PrevKv != nil {
		if err = json.Unmarshal(putResp.PrevKv.Value, &oldJob); err != nil {
			err = nil
			return
		}

	}
	return
}
func DeleteJob(jobmgr *models.JobMgr, job models.Job) (oldJob *models.Job, err error) {

	var (
		jobKey  string
		DelResp *clientv3.DeleteResponse
	)
	jobKey = JobDir + job.Name
	if DelResp, err = jobmgr.Kv.Delete(context.TODO(), jobKey); err != nil {
		return
	}
	if len(DelResp.PrevKvs) != 0 {
		if err = json.Unmarshal(DelResp.PrevKvs[0].Value, &oldjob); err != nil {
			err = nil
			return
		}
	}
	return
}
func KillJob(client *models.ParameCrontab) (r int, err error) {
	var (
		killKey        string
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseid        clientv3.LeaseID
	)
	killKey = JobKill + client.JobName

	if leaseGrantResp, err = GJobmgr.Lease.Grant(context.TODO(), 1); err != nil {
		fmt.Println(err)
		return
	}
	leaseid = leaseGrantResp.ID

	if _, err = GJobmgr.Kv.Put(context.TODO(), killKey, "", clientv3.WithLease(leaseid)); err != nil {
		return
	}
	r = 1
	return
}
func CheckJob(client *models.ParameCrontab) (err error) {
	sqlStr := `select count(jobname)  from joblist where jobname=?`
	var count int
	if err := db.Get(&count, sqlStr, client.JobName); err != nil {
		return err
	}
	if count > 0 {
		return ErrorHostExist
	}

	return
}
func CrontabAdd(client *models.ParameCrontab) (Reply int64, err error) {
	if oldjob, err = SaveJob(GJobmgr, client.CrontabJob); err != nil {
		fmt.Println(err)
	}
	sqlStr := "insert into joblist(jobid,jobname,jobshell,jobstarttime,jobstatus,jobcronexpr) values (?,?,?,?,?,?)"
	ret, err := db.Exec(sqlStr,
		snowflake.IdNum(),
		client.JobName,
		client.JobShell,
		todaytime.NowTimeFull(),
		client.JobStatus,
		client.JobCronExpr,
	)
	if err != nil {
		fmt.Println(err)
		if oldjob, err = DeleteJob(GJobmgr, client.Job); err != nil {
			fmt.Println(err)
		}
		return
	}
	Reply, err = ret.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Printf("更新数据为 %d 条\n", Reply)
	}
	return

}
func CrontabDel(client *models.ParameCrontab) (Reply int64, err error) {
	if oldjob, err = DeleteJob(GJobmgr, client.Job); err != nil {
		fmt.Println(err)
	}
	sqlStr := "delete  from joblist where jobname=?"
	ret, err := db.Exec(sqlStr, client.JobName)
	if err != nil {
		return
	}
	Reply, err = ret.RowsAffected()
	if err != nil {
		return
	} else {
		fmt.Printf("主机表删除数据为 %d 条\n", Reply)
	}
	return

}
func CrontabEdit(client *models.ParameCrontab) (Reply int64, err error) {
	if oldjob, err = SaveJob(GJobmgr, client.CrontabJob); err != nil {
		fmt.Println(err)
	}
	sqlStr := `update joblist set jobstatus=?,jobshell=?,jobname=?,jobcronexpr=? where jobid=? `
	ret, err := db.Exec(sqlStr,
		client.JobStatus,
		client.JobShell,
		client.JobName,
		client.JobCronExpr,
		client.JobId,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	Reply, err = ret.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Printf("更新数据为 %d 条\n", Reply)
	}
	return

}
func CrontabSelect(client *models.ParameCrontab) (Reply []models.CrontabJob, err error) {
	sqlStr := "select jobid,jobname,jobshell,jobstarttime,jobstatus,jobcronexpr from joblist "
	if err := db.Select(&Reply, sqlStr); err != nil {
		return Reply, err
	}
	return

}
func CrontabTotal(client *models.ParameCrontab) (total int, err error) {
	sqlStr := `select count(jobid)  from joblist`
	if err := db.Get(&total, sqlStr); err != nil {
		return total, err
	}
	return
}
func CrontabOnline(client *models.ParameCrontab) (total int, err error) {
	client.JobStatus = 1
	sqlStr := `select count(jobid)  from joblist where jobstatus= ?`
	if err := db.Get(&total, sqlStr, client.JobStatus); err != nil {
		return total, err
	}
	return
}
func CrontabTodayTotal(client *models.ParameCrontab) (total int, err error) {
	now := time.Now()

	sqlStr := `select count(jobid)  from joblist where jobstarttime > ? `
	if err := db.Get(&total, sqlStr, now.Format("2006-01-02")+" 00:00:00"); err != nil {
		return total, err
	}
	return
}
func CrontabAddToday(client *models.ParameCrontab) (total int, err error) {
	now := time.Now()
	client.JobStatus = 1
	sqlStr := `select count(jobid)  from joblist where jobstatus= ? and jobstarttime > ? `
	if err := db.Get(&total, sqlStr, client.JobStatus, now.Format("2006-01-02")+" 00:00:00"); err != nil {
		return total, err
	}
	return
}
