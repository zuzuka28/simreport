package bot

import (
	"context"
	"fmt"

	tele "gopkg.in/telebot.v4"
)

type menuButton struct {
	nextState menuState
	btn       *tele.Btn
}

type menu struct {
	btns   []*menuButton
	markup *tele.ReplyMarkup

	sm       *stateManager
	handlers map[menuState]func(context.Context, tele.Context) error
}

func newMenu(stateManager *stateManager) *menu {
	btnAddFile := menuButton{
		nextState: menuStateAddFile,
		btn: &tele.Btn{ //nolint:exhaustruct
			Unique: "uploadFileBtn",
			Text:   "Add new file",
		},
	}
	btnSearchFile := menuButton{
		nextState: menuStateSearchFile,
		btn: &tele.Btn{ //nolint:exhaustruct
			Unique: "searchSimilarBtn",
			Text:   "Search similar file",
		},
	}

	markup := &tele.ReplyMarkup{ResizeKeyboard: true} //nolint:exhaustruct

	markup.Inline(
		markup.Row(*btnAddFile.btn),
		markup.Row(*btnSearchFile.btn),
	)

	menu := &menu{
		btns: []*menuButton{
			&btnAddFile,
			&btnSearchFile,
		},
		markup:   markup,
		sm:       stateManager,
		handlers: make(map[menuState]func(context.Context, tele.Context) error),
	}

	menu.handlers[menuStateEnter] = func(_ context.Context, c tele.Context) error {
		return c.Send("Choose an action:", menu.markup)
	}

	menu.handlers[menuStateAddFile] = func(ctx context.Context, c tele.Context) error {
		_ = stateManager.SwitchState(ctx, int(c.Sender().ID), string(menuStateAddFileAwaitingDocument))
		return c.Send("send file to upload")
	}

	menu.handlers[menuStateAddFileAwaitingDocument] = func(ctx context.Context, c tele.Context) error {
		userID := int(c.Sender().ID)
		file := c.Message().Document

		fmt.Println(file.FileName)

		_ = stateManager.SwitchState(ctx, userID, string(menuStateEnter))

		return c.Send("document uploaded")
	}

	menu.handlers[menuStateSearchFile] = func(ctx context.Context, c tele.Context) error {
		_ = stateManager.SwitchState(ctx, int(c.Sender().ID), string(menuStateSearchFileAwaitingDocument))
		return c.Send("Please send the file to search for similar ones.")
	}

	menu.handlers[menuStateSearchFileAwaitingDocument] = func(ctx context.Context, c tele.Context) error {
		userID := int(c.Sender().ID)
		file := c.Message().Document

		fmt.Println(file.FileName)

		_ = stateManager.SwitchState(ctx, userID, string(menuStateExit))

		return c.Send("searching by sample...")
	}

	menu.handlers[menuStateExit] = func(ctx context.Context, c tele.Context) error {
		_ = stateManager.SwitchState(ctx, int(c.Sender().ID), string(botStateStart))
		return c.Send("Exiting menu.")
	}

	return menu
}

func (m *menu) Buttons() []*tele.Btn {
	res := make([]*tele.Btn, 0, len(m.btns))
	for _, v := range m.btns {
		res = append(res, v.btn)
	}

	return res
}

func (m *menu) ButtonCallback(btn *tele.Btn) func(ctx context.Context, c tele.Context) error {
	btnuniq := make(map[string]*menuButton)

	for _, v := range m.btns {
		btnuniq[v.btn.Unique] = v
	}

	return func(ctx context.Context, c tele.Context) error {
		_ = m.sm.SwitchState(ctx, int(c.Sender().ID), string(btnuniq[btn.Unique].nextState))
		return m.Handle(ctx, c)
	}
}

func (m *menu) Handle(ctx context.Context, c tele.Context) error {
	currentState, err := m.sm.CurrentState(ctx, int(c.Sender().ID))
	if err != nil {
		return fmt.Errorf("retrieve state: %w", err)
	}

	if action, exists := m.handlers[menuState(currentState)]; exists {
		return action(ctx, c)
	}

	_ = c.Send("Unknown menu option.")

	return nil
}

func (m *menu) Markup() *tele.ReplyMarkup {
	return m.markup
}
