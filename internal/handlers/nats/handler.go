package natshandler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Dorrrke/rent-analytcs-service/internal/domain"
	"github.com/Dorrrke/rent-analytcs-service/internal/service"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

const handlerTimeout = 10 * time.Second

type Handler struct {
	svc    *service.AanalyticsService
	logger *zap.Logger
	subs   []*nats.Subscription
}

func New(svc *service.AanalyticsService, logger *zap.Logger) *Handler {
	return &Handler{svc: svc, logger: logger}
}

func (h *Handler) Subscribe(nc *nats.Conn) error {
	routes := []struct {
		subject string
		handler nats.MsgHandler
	}{
		{domain.SubjectUserRegistered, h.handlerUserRegistered},
		{domain.SubjectUserLoggedIn, h.handlerUserLoggedIn},
		{domain.SubjectTokenRefreshed, h.handlerTokenRefreshed},
		{domain.SubjectCarAdded, h.handlerCarAdded},
		{domain.SubjectRentStarted, h.handlerRentStarted},
		{domain.SubjectRentEnded, h.handlerRentEnded},
	}

	for _, r := range routes {
		sub, err := nc.Subscribe(r.subject, r.handler)
		if err != nil {
			return err
		}
		h.subs = append(h.subs, sub)
		h.logger.Debug("subscribe to NATs subject", zap.String("subject", r.subject))
	}

	return nil
}

func (h *Handler) Unsubscribe() {
	for _, sub := range h.subs {
		if err := sub.Drain(); err != nil {
			h.logger.Warn("failed to drain subscription",
				zap.String("subject", sub.Subject),
				zap.Error(err),
			)
		}
	}
}

// --- Handlers --------------------------------

func (h *Handler) handlerUserRegistered(msg *nats.Msg) {
	var e domain.UserRegisteredEvent
	if err := decode(msg.Data, &e); err != nil {
		h.logDecodeError(msg.Subject, err)
		return
	}

	h.dispatch(msg.Subject, func(ctx context.Context) error {
		return h.svc.HandleUserRegistered(ctx, e)
	})
}

func (h *Handler) handlerUserLoggedIn(msg *nats.Msg) {
	var e domain.UserLoggedInEvent
	if err := decode(msg.Data, &e); err != nil {
		h.logDecodeError(msg.Subject, err)
		return
	}

	h.dispatch(msg.Subject, func(ctx context.Context) error {
		return h.svc.HandleUserLoggedIn(ctx, e)
	})
}

func (h *Handler) handlerTokenRefreshed(msg *nats.Msg) {
	var e domain.TokenRefreshedEvent
	if err := decode(msg.Data, &e); err != nil {
		h.logDecodeError(msg.Subject, err)
		return
	}

	h.dispatch(msg.Subject, func(ctx context.Context) error {
		return h.svc.HandleTokenRefreshed(ctx, e)
	})
}

func (h *Handler) handlerCarAdded(msg *nats.Msg) {
	var e domain.CarAddedEvent
	if err := decode(msg.Data, &e); err != nil {
		h.logDecodeError(msg.Subject, err)
		return
	}

	h.dispatch(msg.Subject, func(ctx context.Context) error {
		return h.svc.HandleCarAdded(ctx, e)
	})
}

func (h *Handler) handlerRentStarted(msg *nats.Msg) {
	var e domain.RentStartedEvent
	if err := decode(msg.Data, &e); err != nil {
		h.logDecodeError(msg.Subject, err)
		return
	}

	h.dispatch(msg.Subject, func(ctx context.Context) error {
		return h.svc.HandleRentStarted(ctx, e)
	})
}

func (h *Handler) handlerRentEnded(msg *nats.Msg) {
	var e domain.RentEndedEvent
	if err := decode(msg.Data, &e); err != nil {
		h.logDecodeError(msg.Subject, err)
		return
	}

	h.dispatch(msg.Subject, func(ctx context.Context) error {
		return h.svc.HandleRentEnded(ctx, e)
	})
}

// --- Helpers --------------------------------

func (h *Handler) dispatch(subject string, fn func(ctx context.Context) error) {
	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeout)
	defer cancel()

	start := time.Now()
	if err := fn(ctx); err != nil {
		h.logger.Error("failed to process event",
			zap.String("subject", subject),
			zap.Duration("elapsed", time.Since(start)),
			zap.Error(err),
		)
		return
	}

	h.logger.Debug("event processed",
		zap.String("subject", subject),
		zap.Duration("elapsed", time.Since(start)),
	)
}

func (h *Handler) logDecodeError(subject string, err error) {
	h.logger.Error("failed to decode event payload",
		zap.String("subject", subject),
		zap.Error(err),
	)
}

func decode(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
