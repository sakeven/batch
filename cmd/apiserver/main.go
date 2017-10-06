package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sakeven/batch/pkg/api"

	"github.com/coreos/etcd/clientv3"
	log "github.com/sirupsen/logrus"
)

// EtcdClient holds an etcd connetion
type EtcdClient struct {
	client *clientv3.Client
}

// NewEtcdClient create a EtcdClient instance
func NewEtcdClient() *EtcdClient {
	conf := clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	}
	client, err := clientv3.New(conf)
	if err != nil {
		panic("can't connect etcd")
	}
	return &EtcdClient{client: client}
}

var etcdStore *EtcdClient

func createJob(w http.ResponseWriter, r *http.Request) {
	job := &api.Job{}
	err := json.NewDecoder(r.Body).Decode(job)
	if err != nil {
		log.Errorf("unmarshal job error %s", err)
		return
	}

	b, err := json.Marshal(job)
	if err != nil {
		log.Errorf("marshal job error %s", err)
		return
	}
	etcdStore.client.Put(context.TODO(), fmt.Sprintf("/batch/job/%s", job.Name), string(b))
}

func listJob(w http.ResponseWriter, r *http.Request) {
	resp, err := etcdStore.client.Get(context.TODO(), "/batch/job", clientv3.WithPrefix())
	if err != nil {
		log.Errorf("can't list jobs")
		return
	}

	jobs := make([]*api.Job, 0, resp.Count)
	for _, kv := range resp.Kvs {
		job := &api.Job{}
		err := json.Unmarshal(kv.Value, job)
		if err != nil {
			log.Errorf("marshal job error %s", err)
		}
		jobs = append(jobs, job)
	}
	json.NewEncoder(w).Encode(jobs)
}

func bind(w http.ResponseWriter, r *http.Request) {

}

func main() {
	etcdStore = NewEtcdClient()
	router := newRouter()
	log.Infof("APIserver runs on 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
