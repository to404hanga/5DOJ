package producer

import (
	"5DOJ/pkg/constant/topic"
	"5DOJ/submitter/global"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/IBM/sarama"
)

type SubmitterProducer struct {
	producer sarama.SyncProducer
}

var _ Producer = (*SubmitterProducer)(nil)

func NewSubmitterProducer() *SubmitterProducer {
	p, err := sarama.NewSyncProducerFromClient(global.Kafka)
	if err != nil {
		panic(fmt.Errorf("初始化消费者失败: %s", err))
	}
	return &SubmitterProducer{
		producer: p,
	}
}

func (p *SubmitterProducer) Produce(ctx context.Context, evt topic.SubmitEvent) error {
	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Key:   sarama.StringEncoder(strconv.FormatUint(evt.RecordId, 10)),
		Topic: topic.TopicSubmitEvent,
		Value: sarama.ByteEncoder(data),
	})
	return err
}
