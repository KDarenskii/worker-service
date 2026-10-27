package eorder

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/KDarenskii/worker-service/internal/app/entity"
	ehandler "github.com/KDarenskii/worker-service/internal/app/handler/event"
	"github.com/KDarenskii/worker-service/internal/app/service"
	"github.com/KDarenskii/worker-service/internal/pkg/broker"
	"github.com/rs/zerolog/log"
)

type handler struct {
	delivery service.Delivery
	bus      broker.Bus[entity.EventOrderDeliveryCalculated]
}

func NewHandler(
	delivery service.Delivery,
	bus broker.Bus[entity.EventOrderDeliveryCalculated],
) ehandler.OrderCreated {
	return &handler{
		delivery: delivery,
		bus:      bus,
	}
}

func (h *handler) CallbackOrderCreated(
	ctx context.Context, ev *entity.EventOrderCreated, _ []broker.Header,
) error {
	log.Info().Ctx(ctx).EmbedObject(ev).Msg("broker: received message")

	if ev.OrderGUID == "" {
		return broker.NotCriticalError(fmt.Errorf("broker: invalid order_guid: %s", ev.OrderGUID))
	}

	deliveryPrice, err := h.delivery.Calculate(ctx, ev.Currency)
	if err != nil {
		if errors.Is(err, entity.ErrFixerCurrencyNotFound) {
			return broker.NotCriticalError(err)
		}
		return fmt.Errorf("calculate delivery price %s: %w", ev.Currency, err)
	}

	calculatedEvent := entity.EventOrderDeliveryCalculated{
		OrderGUID:     ev.OrderGUID,
		DeliveryPrice: deliveryPrice,
		Currency:      ev.Currency,
		CalculatedAt:  time.Now().UTC().Format(time.RFC3339),
	}

	if err := h.bus.Send(
		ctx, &calculatedEvent, entity.BrokerHeaderOrderDeliveryCalculatedEventID(),
		entity.BrokerHeaderOrderDeliveryCalculatedType(),
	); err != nil {
		return fmt.Errorf("send EventOrderDeliveryCalculated: %w", err)
	}

	log.Info().
		Ctx(ctx).
		EmbedObject(calculatedEvent).
		Msg("Опубликовано событие order.delivery.calculated")

	return nil
}
