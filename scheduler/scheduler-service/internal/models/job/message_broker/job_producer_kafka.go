package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/config"
	models "github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/internal/models"
	"github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/pkg/logger"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
)

const (
	writerReadTimeout  = 10 * time.Second
	writerWriteTimeout = 10 * time.Second
	writerRequiredAcks = -1
	writerMaxAttempts  = 3
)

type jobsProducer struct {
	log          logger.Logger
	createWriter *kafka.Writer
	updateWriter *kafka.Writer
	runWriter    *kafka.Writer
}

// NewJobsProducer constructor
func NewJobsProducer(log logger.Logger) *jobsProducer {
	return &jobsProducer{log: log}
}

// Run init producers writers
func (p *jobsProducer) Run(conn *kafka.Conn, cfg *config.Config) {

	// producers are responsible for creation of topics
	topic1 := kafka.TopicConfig{Topic: cfg.TopicJobCreate, NumPartitions: cfg.TopicJobCreatePartitions, ReplicationFactor: cfg.TopicJobCreateReplicas}
	topic2 := kafka.TopicConfig{Topic: cfg.TopicJobUpdate, NumPartitions: cfg.TopicJobUpdatePartitions, ReplicationFactor: cfg.TopicJobUpdateReplicas}
	topic3 := kafka.TopicConfig{Topic: cfg.TopicJobRun, NumPartitions: cfg.TopicJobRunPartitions, ReplicationFactor: cfg.TopicJobRunReplicas}
	conn.CreateTopics(topic1, topic2, topic3)

	p.createWriter = p.getNewKafkaWriter(cfg.TopicJobCreate, cfg.KafkaBrokers)
	p.updateWriter = p.getNewKafkaWriter(cfg.TopicJobUpdate, cfg.KafkaBrokers)
	p.runWriter = p.getNewKafkaWriter(cfg.TopicJobRun, cfg.KafkaBrokers)
}

// GetNewKafkaWriter Create new kafka writer
func (jp *jobsProducer) getNewKafkaWriter(topic string, brokers []string) *kafka.Writer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: writerRequiredAcks,
		MaxAttempts:  writerMaxAttempts,
		Logger:       kafka.LoggerFunc(jp.log.Debugf),
		ErrorLogger:  kafka.LoggerFunc(jp.log.Errorf),
		Compression:  compress.Snappy,
		ReadTimeout:  writerReadTimeout,
		WriteTimeout: writerWriteTimeout,
	}
	return w
}

// Close close writers
func (p jobsProducer) Close() {
	p.createWriter.Close()
	p.updateWriter.Close()
	p.runWriter.Close()
}

//implement models.job.JobsProducer for jobsProducer

func (p *jobsProducer) PublishCreate(ctx context.Context, job *models.Job) error {

	jobBytes, err := json.Marshal(&job)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}

	return p.createWriter.WriteMessages(ctx, kafka.Message{
		Value: jobBytes,
		Time:  time.Now().UTC(),
	})
}

func (p *jobsProducer) PublishUpdate(ctx context.Context, job *models.Job) error {

	jobBytes, err := json.Marshal(&job)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}

	return p.updateWriter.WriteMessages(ctx, kafka.Message{
		Value: jobBytes,
		Time:  time.Now().UTC(),
	})
}

func (p *jobsProducer) PublishRun(ctx context.Context, job *models.Job) error {

	jobBytes, err := json.Marshal(&job)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}

	return p.runWriter.WriteMessages(ctx, kafka.Message{
		Value: jobBytes,
		Time:  time.Now().UTC(),
	})
}
