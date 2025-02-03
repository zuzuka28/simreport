package bot

import (
	"context"
	"errors"
	"fmt"

	tele "gopkg.in/telebot.v4"
)

type Event string

type State string

const (
	botStateStart                       State = "start"
	menuStateEnter                      State = "menu.Enter"
	menuStateAddFile                    State = "menu.AddFile"
	menuStateAddFileAwaitingDocument    State = "menu.AddFile.AwaitingDocument"
	menuStateSearchFile                 State = "menu.SearchFile"
	menuStateSearchFileAwaitingDocument State = "menu.SearchFile.AwaitingDocument"
	menuStateExit                       State = "menu.Exit"
)

const (
	eventAddFile      Event = "AddFile"
	eventSearchFile   Event = "SearchFile"
	eventFileUploaded Event = "FileUploaded"
	eventEnterMenu    Event = "EnterMenu"
	eventExitMenu     Event = "ExitMenu"
	eventFileSearched Event = "FileSearched"
)

var errInvalidTransition = errors.New("invalid transition")

type Transition struct {
	From   State
	Event  Event
	To     State
	Action func(ctx context.Context, c tele.Context) error
}

type StateMachine struct {
	transitions map[State]map[Event]Transition
	usm         *stateManager
}

func NewStateMachine(sm *stateManager) *StateMachine {
	return &StateMachine{
		transitions: make(map[State]map[Event]Transition),
		usm:         sm,
	}
}

func (sm *StateMachine) AddTransitions(transitions []Transition) {
	for _, t := range transitions {
		if sm.transitions[t.From] == nil {
			sm.transitions[t.From] = make(map[Event]Transition)
		}

		sm.transitions[t.From][t.Event] = t
	}
}

func (sm *StateMachine) CurrentState(ctx context.Context, userID int) (State, error) {
	state, err := sm.usm.CurrentState(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("get state: %w", err)
	}

	return State(state), nil
}

func (sm *StateMachine) Trigger(ctx context.Context, c tele.Context, event string) error {
	currentState, err := sm.CurrentState(ctx, int(c.Sender().ID))
	if err != nil {
		return fmt.Errorf("get current state: %w", err)
	}

	transition, ok := sm.findTransition(currentState, Event(event))
	if !ok {
		return fmt.Errorf("%w: from %s on event %s", errInvalidTransition, currentState, event)
	}

	return sm.executeTransition(ctx, c, transition)
}

func (sm *StateMachine) findTransition(state State, event Event) (Transition, bool) {
	transitions, ok := sm.transitions[state]
	if !ok {
		return Transition{}, false //nolint:exhaustruct
	}

	transition, ok := transitions[event]

	return transition, ok
}

func (sm *StateMachine) executeTransition(ctx context.Context, c tele.Context, t Transition) error {
	if err := t.Action(ctx, c); err != nil {
		return fmt.Errorf("execute transition action: %w", err)
	}

	if err := sm.usm.SwitchState(ctx, int(c.Sender().ID), string(t.To)); err != nil {
		return fmt.Errorf("switch state: %w", err)
	}

	return nil
}
