/**
 * @Description:
 * @FilePath: /bull-golang/example/example.go
 * @Author: liyibing liyibing@lixiang.com
 * @Date: 2023-07-27 18:04:13
 */
package main

import (
	"bull-go"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type MqTrigger struct {
	bull bull.BullQueueIface
}

func NewMqTrigger(opts bull.BullQueueOption) (*MqTrigger, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     opts.RedisIp,
		Password: opts.RedisPasswd,
		DB:       0,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.New("redis init failed")
	}
	return &MqTrigger{
		bull: &bull.BullQueue{
			Name:      opts.QueueName,
			KeyPrefix: opts.KeyPrefix + ":" + opts.QueueName + ":",
			Token:     uuid.New(),
			Client:    rdb,
		},
	}, nil
}

// func (m MqTrigger) OnNewMsg(data []byte){
//     m.bull.Add(nil, xxx)
// }

func main() {
	queueOp := bull.BullQueueOption{
		KeyPrefix:   "bull",
		QueueName:   "jobs",
		RedisIp:     "127.0.0.1:6379",
		RedisPasswd: "",
	}
	// q, err := bull.NewQueue(queueOp) // 暴露方法 创建队列
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	person := Person{
		Name: "ppperson",
		Age:  105,
	}
	jobdata, err := json.Marshal(person)
	if err != nil {
		fmt.Println(err)
		return
	}
	// jobOp := bull.JobOptions{}
	// _, err = q.Add(jobdata, jobOp) // 暴露方法 添加任务
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(job)

	mqtriger, err := NewMqTrigger(queueOp)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = mqtriger.bull.Add(jobdata, bull.WithPriorityOp(100), bull.WithRemoveOnCompleteOp(true), bull.WithRemoveOnFailOp(true))
	if err != nil {
		fmt.Println(err)
		return
	}
	// _, err = mqtriger.bull.Add(jobdata)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
}
