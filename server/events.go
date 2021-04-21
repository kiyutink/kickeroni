package server

import (
	"github.com/kiyutink/kickeroni/lib"
	"github.com/labstack/echo/v4"
	"github.com/slack-go/slack/slackevents"
)

func Events(c echo.Context) error {
	eventsAPIEvent := new(slackevents.EventsAPIEvent)
	err := (&lib.MultiBinder{}).Bind(eventsAPIEvent, c)

	if err != nil {
		return err
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		challengeResponse := new(slackevents.ChallengeResponse)
		if err := c.Bind(challengeResponse); err != nil {
			return err
		}
		return c.JSON(200, challengeResponse)

	}
	return echo.NewHTTPError(400)
}
