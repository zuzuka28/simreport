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

	menu *menu

	done chan struct{}
	fsm  *StateMachine
}

func New(
	cfg Config,
	uss UserStateService,
	ds DocumentService,
	ss SimilarityService,
	m Metrics,
) (*Bot, error) {
	pref := tele.Settings{ //nolint:exhaustruct
		Token:  cfg.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second}, //nolint:exhaustruct,mnd,gomnd
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("new telebot: %w", err)
	}

	b.Use(newMetricsMiddleware(m))

	sm := newStateManager(uss)

	menu := newMenu(ds, ss)

	bot := &Bot{
		tg:   b,
		done: make(chan struct{}),
		uss:  uss,
		menu: menu,
		fsm:  NewStateMachine(sm),
	}

	bot.fsm.AddTransitions([]Transition{
		{
			From:   botStateStart,
			Event:  eventEnterMenu,
			To:     menuStateEnter,
			Action: sendMenuChoice(menu.markup),
		},
	})

	bot.fsm.AddTransitions(menu.Transitions())

	b.Handle("/start", func(c tele.Context) error {
		ctx := context.Background()
		return bot.fsm.Trigger(ctx, c, string(eventEnterMenu))
	})

	b.Handle(tele.OnDocument, func(c tele.Context) error {
		ctx := context.Background()
		return bot.fsm.Trigger(ctx, c, string(eventFileRecieved))
	})

	b.Handle(tele.OnText, func(c tele.Context) error {
		ctx := context.Background()
		return bot.fsm.Trigger(ctx, c, string(eventTextRecieved))
	})

	for _, v := range menu.Buttons() {
		b.Handle(v.btn, func(c tele.Context) error {
			ctx := context.Background()
			return bot.fsm.Trigger(ctx, c, string(v.event))
		})
	}

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
