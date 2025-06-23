package consumer

import (
	"5DOJ/judger/global"
	"5DOJ/judger/service"
	"5DOJ/pkg/constant/contestMode"
	"5DOJ/pkg/constant/language"
	"5DOJ/pkg/constant/topic"
	"5DOJ/pkg/kafka"
	"context"

	"github.com/IBM/sarama"
	"github.com/to404hanga/pkg404/saramax"
)

type JudgerSubmitConsumer struct {
	svc service.IJudgerService
}

var _ kafka.Consumer = (*JudgerSubmitConsumer)(nil)

func NewJudgerSubmitConsumer(svc service.IJudgerService) *JudgerSubmitConsumer {
	return &JudgerSubmitConsumer{
		svc: svc,
	}
}

func (j *JudgerSubmitConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient(topic.TopicSubmitEvent, global.Kafka)
	if err != nil {
		return err
	}

	go func() {
		err = cg.Consume(context.Background(), []string{topic.TopicSubmitEvent}, saramax.NewHandler[topic.SubmitEvent](global.L, j.Consume))
		if err != nil {
			// TODO: 记录日志
			panic(err)
		}
	}()

	return nil
}

func (j *JudgerSubmitConsumer) Consume(msg *sarama.ConsumerMessage, evt topic.SubmitEvent) (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, _, _, err = j.svc.Judge(ctx, evt.RecordId, evt.ProblemId, language.LanguageType(evt.Language), evt.FilenameWithoutExt, evt.Code, contestMode.ContestModeType(evt.Mode))
	return
}
