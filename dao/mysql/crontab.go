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

var (
	clinet   *clientv3.Client
	kv       clientv3.KV
	lease    clientv3.Lease
	G_JobMgr *models.JobMgr
	oldjob   *models.Job
)

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
	G_JobMgr = &models.JobMgr{
		Clinet: clinet,
		Kv:     kv,
		Lease:  lease,
	}

	return
}
func SaveJob(jobmgr *models.JobMgr, job models.Job) (oldJob *models.Job, err error) {

	var (
		jobKey   string
		jobValue []byte
		putResp  *clientv3.PutResponse
	)
	jobKey = "/cron/jobs/" + job.Name
	if jobValue, err = json.Marshal(job); err != nil {
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
func CrontabAdd(client *models.ParameCrontab) (Reply int64, err error) {
	if oldjob, err = SaveJob(G_JobMgr, client.Job); err != nil {
		fmt.Println(err)
	}
	sqlStr := "insert into joblist(jobid,jobname,jobshell,jobstarttime,jobstatus,jobcronexpr) values (?,?,?,?,?,?)"
	ret, err := db.Exec(sqlStr,
		snowflake.IdNum(),
		client.Job.Name,
		client.Job.Command,
		todaytime.NowTimeFull(),
		client.JobStatus,
		client.Job.CronExpr,
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
