package producer

import (
	"5DOJ/pkg/constant/topic"
	"context"
)

type Producer interface {
	Produce(ctx context.Context, evt topic.SubmitEvent) error
}
