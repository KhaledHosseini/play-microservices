package kafka

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/config"
	models "github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/internal/models"
	"github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/pkg/logger"
	"github.com/segmentio/kafka-go"
)

const (
	minBytes               = 10e3 // 10KB
	maxBytes               = 10e6 // 10MB
	queueCapacity          = 100
	heartbeatInterval      = 3 * time.Second
	commitInterval         = 0
	partitionWatchInterval = 5 * time.Second
	maxAttempts            = 3
	dialTimeout            = 3 * time.Minute
)

type JobsConsumerGroup struct {
	log   logger.Logger
	jobDB models.JobDB
}

func NewJobsConsumerGroup(
	log logger.Logger,
	jobDB models.JobDB,
) *JobsConsumerGroup {
	return &JobsConsumerGroup{
		log:   log,
		jobDB: jobDB,
	}
}

func (jcg *JobsConsumerGroup) getNewKafkaReader(brokers []string, topic string, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:                brokers,
		GroupID:                groupID,
		Topic:                  topic,
		MinBytes:               minBytes,
		MaxBytes:               maxBytes,
		QueueCapacity:          queueCapacity,
		HeartbeatInterval:      heartbeatInterval,
		CommitInterval:         commitInterval,
		PartitionWatchInterval: partitionWatchInterval,
		Logger:                 kafka.LoggerFunc(jcg.log.Debugf),
		ErrorLogger:            kafka.LoggerFunc(jcg.log.Errorf),
		MaxAttempts:            maxAttempts,
		Dialer: &kafka.Dialer{
			Timeout: dialTimeout,
		},
	})
}

// Run init producers writers
func (jcg *JobsConsumerGroup) Run(ctx context.Context, cancel context.CancelFunc, conn *kafka.Conn, cfg *config.Config) {

	go func() {
		r := jcg.getNewKafkaReader(cfg.KafkaBrokers, cfg.TopicJobRunResult, cfg.TopicJobRunResultConsumerGroupID)
		defer cancel()
		defer func() {
			if err := r.Close(); err != nil {
				jcg.log.Errorf("r.Close", err)
				cancel()
			}
		}()
		jcg.log.Infof("Starting consumer group: %v", r.Config().GroupID)

		wg := &sync.WaitGroup{}
		for i := 0; i <= cfg.TopicJobRunResultConsumerWorkerCount; i++ {
			wg.Add(1)
			go jcg.createJobWorker(ctx, cancel, r, wg, i)
		}
		wg.Wait()
	}()
}

func (jcg *JobsConsumerGroup) createJobWorker(
	ctx context.Context,
	cancel context.CancelFunc,
	r *kafka.Reader,
	wg *sync.WaitGroup,
	workerID int,
) {
	defer wg.Done()
	defer cancel()

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			jcg.log.Errorf("FetchMessage", err)
			return
		}

		jcg.log.Infof(
			"WORKER: %v, message at topic/partition/offset %v/%v/%v: %s = %s\n",
			workerID,
			m.Topic,
			m.Partition,
			m.Offset,
			string(m.Key),
			string(m.Value),
		)

		var job models.Job
		if err := json.Unmarshal(m.Value, &job); err != nil {
			jcg.log.Errorf("json.Unmarshal", err)
			continue
		}

		jobID := job.Id.Hex()
		job2, err2 := jcg.jobDB.GetByID(ctx, jobID)
		if err2 != nil {
			jcg.log.Errorf("JobsConsumerGroup.CreateJobWorker.GetByID: %v", err)
			continue
		}

		job2.Status = job.Status

		_, err3 := jcg.jobDB.Update(ctx, job2)
		if err3 != nil {
			jcg.log.Errorf("JobsConsumerGroup.CreateJobWorker.Update: %v", err)
			continue
		}

		if err := r.CommitMessages(ctx, m); err != nil {
			jcg.log.Errorf("CommitMessages", err)
			continue
		}
	}
}
