package eprocessor

import (
	"context"
	"sync"

	"github.com/KDarenskii/worker-service/internal/app/entity"
	ehandler "github.com/KDarenskii/worker-service/internal/app/handler/event"
	"github.com/KDarenskii/worker-service/internal/app/processor"
	"github.com/KDarenskii/worker-service/internal/pkg/broker"
	"github.com/rs/zerolog/log"
)

type orderCreatedProc struct {
	h   ehandler.OrderCreated
	bus broker.Bus[entity.EventOrderCreated]
}

func NewOrderCreatedEventsCatcher(
	h ehandler.OrderCreated,
	bus broker.Bus[entity.EventOrderCreated],
) processor.Processor {
	return &orderCreatedProc{h, bus}
}

func (p *orderCreatedProc) StartAsync(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	if err := p.bus.Subscribe(ctx, wg, p.h.CallbackOrderCreated); err != nil {
		wg.Done()
		log.Fatal().Str("topic", p.bus.QueueName()).Err(err).Msg("broker: failed to subscribe")
	}

	log.Info().Str("topic", p.bus.QueueName()).Msg("broker: subscription requested")

	wg.Done()
}
