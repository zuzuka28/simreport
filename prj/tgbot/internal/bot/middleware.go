package bot

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/zuzuka28/simreport/prj/tgbot/internal/model"
	"gopkg.in/telebot.v4"
)

const (
	contextKey = "context"
)

func newInjectContextMiddleware() telebot.MiddlewareFunc {
	return func(hf telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			ctx := context.Background()

			c.Set(contextKey, ctx)

			return hf(c)
		}
	}
}

func newInjectRequestIDMiddleware() telebot.MiddlewareFunc {
	return func(hf telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			ctx := c.Get(contextKey).(context.Context) //nolint:forcetypeassert

			rid := uuid.NewString()

			ctx = context.WithValue(ctx, model.RequestIDKey, rid)

			c.Set(contextKey, ctx)

			return hf(c)
		}
	}
}

func newMetricsMiddleware(m Metrics) telebot.MiddlewareFunc {
	return func(hf telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			m.IncBotRequestsByUser(c.Sender().Username, strconv.Itoa(int(c.Sender().ID)))

			err := hf(c)
			if err != nil {
				m.IncBotErrors(err.Error())
			}

			return err
		}
	}
}

func newLoggingMiddleware() telebot.MiddlewareFunc {
	return func(hf telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			ctx := c.Get(contextKey).(context.Context) //nolint:forcetypeassert

			t := time.Now()

			err := hf(c)

			slog.Info(
				"request processed",
				"request_id", ctx.Value(model.RequestIDKey),
				"elapsed_time", time.Since(t),
			)

			return err
		}
	}
}
