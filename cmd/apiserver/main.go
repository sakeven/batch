package main

import (
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
	}

	b, err := json.Marshal(job)
	if err != nil {
		log.Errorf("marshal job error %s", err)
	}
	etcdStore.client.Put(nil, fmt.Sprintf("/batch/job/%s", job.Name), string(b))
}

func main() {
	etcdStore = NewEtcdClient()
	router := newRouter()
	http.ListenAndServe(":80", router)
}
