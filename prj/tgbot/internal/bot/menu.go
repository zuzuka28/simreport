package bot

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"maps"
	"slices"
	"time"

	"github.com/zuzuka28/simreport/prj/tgbot/internal/model"
	tele "gopkg.in/telebot.v4"
	"gopkg.in/yaml.v2"
)

type menuButton struct {
	event Event
	btn   *tele.Btn
}

type menu struct {
	btns        map[string]*menuButton
	markup      *tele.ReplyMarkup
	transitions []Transition
}

func newMenu(
	ds DocumentService,
	ss SimilarityService,
) *menu {
	buttons := map[string]*menuButton{
		"uploadFileBtn": {
			event: eventAddFile,
			btn: &tele.Btn{ //nolint:exhaustruct
				Unique: "uploadFileBtn",
				Text:   "Add new file",
			},
		},
		"searchSimilarBtn": {
			event: eventSearchFile,
			btn: &tele.Btn{ //nolint:exhaustruct
				Unique: "searchSimilarBtn",
				Text:   "Search similar file",
			},
		},
	}

	markup := &tele.ReplyMarkup{ResizeKeyboard: true} //nolint:exhaustruct

	rows := make([]tele.Row, 0, len(buttons))
	for _, btn := range buttons {
		rows = append(rows, markup.Row(*btn.btn))
	}

	markup.Inline(rows...)

	return &menu{
		btns:   buttons,
		markup: markup,
		transitions: []Transition{
			{
				From:   menuStateEnter,
				Event:  eventAddFile,
				To:     menuStateAddFile,
				Action: sendMessage("send file to upload"),
			},
			{
				From:   menuStateAddFile,
				Event:  eventFileRecieved,
				To:     menuStateEnter,
				Action: handleFileUpload(ds),
			},
			{
				From:   menuStateAddFile,
				Event:  eventEnterMenu,
				To:     menuStateEnter,
				Action: sendMenuChoice(markup),
			},
			{
				From:   menuStateEnter,
				Event:  eventSearchFile,
				To:     menuStateSearchFile,
				Action: sendMessage("send the file to search for similar ones"),
			},
			{
				From:   menuStateSearchFile,
				Event:  eventTextRecieved,
				To:     menuStateEnter,
				Action: handleSimilaritySearch(ss),
			},
			{
				From:   menuStateSearchFile,
				Event:  eventEnterMenu,
				To:     menuStateEnter,
				Action: sendMenuChoice(markup),
			},
			{
				From:   menuStateEnter,
				Event:  eventEnterMenu,
				To:     menuStateEnter,
				Action: sendMenuChoice(markup),
			},
		},
	}
}

func (m *menu) Buttons() []*menuButton {
	return slices.Collect(maps.Values(m.btns))
}

func (m *menu) Markup() *tele.ReplyMarkup {
	return m.markup
}

func (m *menu) Transitions() []Transition {
	return m.transitions
}

func handleFileUpload(ds DocumentService) func(context.Context, tele.Context) error {
	return func(ctx context.Context, c tele.Context) error {
		file := c.Message().Document

		r, err := c.Bot().File(file.MediaFile())
		if err != nil {
			return c.Send("Failed to retrieve the file.")
		}

		hasher := sha256.New()
		data := &bytes.Buffer{}

		if _, err := io.Copy(io.MultiWriter(hasher, data), r); err != nil {
			return fmt.Errorf("load file: %w", err)
		}

		var desc docDetails

		if text := c.Message().Caption; text != "" {
			if err := yaml.Unmarshal([]byte(text), &desc); err != nil {
				return c.Send("failed to parse document details. Example:\n" + docDetailsExample)
			}
		}

		res, err := ds.Save(ctx, model.DocumentSaveCommand{
			Item: model.Document{
				ParentID: desc.ParentID,
				Name:     desc.Name,
				Version:  desc.Version,
				GroupID:  desc.GroupID,
				Source: model.File{
					Name:        file.FileName,
					Content:     data.Bytes(),
					Sha256:      hex.EncodeToString(hasher.Sum(nil)),
					LastUpdated: time.Now(),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("upload file: %w", err)
		}

		return c.Send("document uploaded: " + res.ID)
	}
}

func handleSimilaritySearch(ss SimilarityService) func(context.Context, tele.Context) error {
	return func(ctx context.Context, c tele.Context) error {
		id := c.Message().Text

		res, err := ss.Search(ctx, model.SimilarityQuery{
			ID: id,
		})
		if err != nil {
			return fmt.Errorf("search similar: %w", err)
		}

		buf := &bytes.Buffer{}

		err = yaml.NewEncoder(buf).Encode(mapSimilarityMatchesToResponse(res))
		if err != nil {
			return fmt.Errorf("marshal response: %w", err)
		}

		doc := &tele.Document{
			File:                 tele.FromReader(buf),
			Thumbnail:            nil,
			Caption:              "search results for document" + id,
			MIME:                 "",
			FileName:             "similar_to_" + id + ".yaml",
			DisableTypeDetection: false,
		}

		return c.Reply(doc)
	}
}

func sendMessage(msg string) func(context.Context, tele.Context) error {
	return func(_ context.Context, c tele.Context) error {
		return c.Send(msg)
	}
}

func sendMenuChoice(markup *tele.ReplyMarkup) func(context.Context, tele.Context) error {
	return func(_ context.Context, c tele.Context) error {
		return c.Send("Choose an action:", markup)
	}
}
