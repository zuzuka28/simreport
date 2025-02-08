package bot

import (
	"strconv"

	"gopkg.in/telebot.v4"
)

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
