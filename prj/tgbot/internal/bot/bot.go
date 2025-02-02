package bot

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	tele "gopkg.in/telebot.v4"
)

type Config struct {
	Token string `yaml:"token"`
}

type Bot struct {
	tg *tele.Bot

	uss UserStateService
	sm  *stateManager

	menu *menu

	done chan struct{}
}

func New(
	cfg Config,
	uss UserStateService,
) (*Bot, error) {
	pref := tele.Settings{ //nolint:exhaustruct
		Token:  cfg.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second}, //nolint:exhaustruct,mnd,gomnd
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("new telebot: %w", err)
	}

	stateManager := newStateManager(uss)
	menu := newMenu(stateManager)

	bot := &Bot{
		tg:   b,
		done: make(chan struct{}),
		uss:  uss,
		sm:   stateManager,
		menu: menu,
	}

	for _, btn := range menu.Buttons() {
		b.Handle(btn, func(c tele.Context) error {
			ctx := context.Background()
			return bot.menu.ButtonCallback(btn)(ctx, c)
		})
	}

	b.Handle("/start", func(c tele.Context) error {
		ctx := context.Background()

		_ = stateManager.SwitchState(ctx, int(c.Sender().ID), string(menuStateEnter))

		return bot.menu.Handle(ctx, c)
	})

	b.Handle("/menu", func(c tele.Context) error {
		ctx := context.Background()

		_ = stateManager.SwitchState(ctx, int(c.Sender().ID), string(menuStateEnter))

		return bot.menu.Handle(ctx, c)
	})

	b.Handle(tele.OnDocument, func(c tele.Context) error {
		ctx := context.Background()
		return bot.menu.Handle(ctx, c)
	})

	return bot, nil
}

func (b *Bot) Start(ctx context.Context) error {
	go b.tg.Start()
	defer b.tg.Stop()

	slog.Info("bot started")

	select {
	case <-ctx.Done():
		return fmt.Errorf("run bot: %w", ctx.Err())

	case <-b.done:
		return nil
	}
}

func (b *Bot) Stop() error {
	b.done <- struct{}{}
	return nil
}
