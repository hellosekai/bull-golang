/**
 * @Description: data struct and operations with queue
 * @FilePath: /bull-golang/queue.go
 * @Author: liyibing liyibing@lixiang.com
 * @Date: 2023-07-19 15:55:49
 */
package bull

import (
	"github.com/hellosekai/bull-golang/internal/luaScripts"
	"github.com/hellosekai/bull-golang/internal/redisAction"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type BullQueueIface interface {
	Add(jobData JobData, options ...withOption) (Job, error)
	Ping() error
}

var _ BullQueueIface = (*BullQueue)(nil)

const (
	SingleNode = 0
	Cluster    = 1
)

type BullQueue struct {
	Name      string
	Token     uuid.UUID
	KeyPrefix string
	Client    redis.Cmdable
}

type BullQueueOption struct {
	Mode        int
	KeyPrefix   string
	QueueName   string
	RedisIp     string
	RedisPasswd string
}

func NewBullQueue(opts BullQueueOption) (*BullQueue, error) {
	q := &BullQueue{
		Name:      opts.QueueName,
		Token:     uuid.New(),
		KeyPrefix: opts.KeyPrefix + ":" + opts.QueueName + ":",
	}

	redisIp := opts.RedisIp
	redisPasswd := opts.RedisPasswd
	redisMode := opts.Mode
	var err error
	q.Client, err = redisAction.Init(redisIp, redisPasswd, redisMode)
	if err != nil {
		return nil, err
	}

	return q, nil
}

func (q *BullQueue) Init(opts BullQueueOption) error {
	q.Name = opts.QueueName
	q.Token = uuid.New()
	q.KeyPrefix = opts.KeyPrefix + ":" + opts.QueueName + ":"

	redisIp := opts.RedisIp
	redisPasswd := opts.RedisPasswd
	redisMode := opts.Mode
	var err error
	q.Client, err = redisAction.Init(redisIp, redisPasswd, redisMode)
	if err != nil {
		return err
	}

	return nil
}

/**
 * @description:add a job into queue
 * @param {JobData} jobData
 * @param {...withOption} options
 * @return {*}
 */
func (q *BullQueue) Add(jobData JobData, options ...withOption) (Job, error) {
	distOption := &JobOptions{}

	for _, withOptionFunc := range options {
		withOptionFunc(distOption)
	}

	name := _DEFAULT_JOB_NAME
	job, err := newJob(name, jobData, *distOption)
	if err != nil {
		return job, err
	}
	err = q.addJob(job)
	return job, err
}

func (q *BullQueue) addJob(job Job) error {
	rdb := q.Client
	keys := q.getKeys()
	args := q.getArgs(job)
	err := redisAction.ExecLua(luaScripts.AddJobLua, rdb, keys, args)
	if err != nil {
		return err
	}
	return nil
}

func (q *BullQueue) getKeys() []string {
	keys := make([]string, 0, 6)
	keys = append(keys, q.KeyPrefix+"wait")
	keys = append(keys, q.KeyPrefix+"paused")
	keys = append(keys, q.KeyPrefix+"meta-paused")
	keys = append(keys, q.KeyPrefix+"id")
	keys = append(keys, q.KeyPrefix+"delayed")
	keys = append(keys, q.KeyPrefix+"priority")

	return keys
}

func (q *BullQueue) getArgs(job Job) []interface{} {
	args := make([]interface{}, 0, 11)
	args = append(args, q.KeyPrefix)
	args = append(args, job.Id)
	args = append(args, job.Name)
	args = append(args, job.Data)
	args = append(args, job.OptsByJson)
	args = append(args, job.TimeStamp)
	args = append(args, job.Delay)
	args = append(args, job.DelayTimeStamp)
	args = append(args, job.Opts.Priority)
	if job.Opts.Lifo == "RPUSH" {
		args = append(args, "RPUSH")
	} else {
		args = append(args, "LPUSH")
	}
	args = append(args, q.Token)

	return args
}

func (q *BullQueue) Ping() error {
	return redisAction.Ping(q.Client)
}
