package bot

import "errors"

type botState string

const (
	botStateStart botState = "start"
	botStateMenu  botState = "menu"
)

type menuState string

const (
	menuStateEnter                      menuState = "menu.Enter"
	menuStateAddFile                    menuState = "menu.AddFile"
	menuStateAddFileAwaitingDocument    menuState = "menu.AddFile.AwaitingDocument"
	menuStateSearchFile                 menuState = "menu.SearchFile"
	menuStateSearchFileAwaitingDocument menuState = "menu.SearchFile.AwaitingDocument"
	menuStateExit                       menuState = "menu.Exit"
)

var errInvalidState = errors.New("unknown state")
