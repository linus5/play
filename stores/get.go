package stores

import (
	"github.com/dave/flux"
	"github.com/dave/jsgo/server/messages"
	"github.com/dave/play/actions"
)

func NewGetStore(app *App) *GetStore {
	s := &GetStore{
		app: app,
	}
	return s
}

type GetStore struct {
	app *App
}

func (s *GetStore) Handle(payload *flux.Payload) bool {
	switch action := payload.Action.(type) {
	case *actions.GetStart:
		s.app.Log("downloading")
		s.app.Dispatch(&actions.Dial{
			Url:  defaultUrl(),
			Open: func() flux.ActionInterface { return &actions.GetOpen{Path: action.Path} },
			Message: func(m interface{}) flux.ActionInterface {
				return &actions.GetMessage{Path: action.Path, Message: m, Save: action.Save}
			},
			Close: func() flux.ActionInterface { return &actions.GetClose{} },
		})
		payload.Notify()
	case *actions.GetOpen:
		message := messages.Get{
			Path: action.Path,
		}
		s.app.Dispatch(&actions.Send{
			Message: message,
		})
	case *actions.GetMessage:
		switch message := action.Message.(type) {
		case messages.Downloading:
			if len(message.Message) > 0 {
				s.app.Log(message.Message)
			}
		case messages.GetComplete:
			s.app.Dispatch(&actions.LoadSource{Source: message.Source, Save: action.Save})
			var count int
			for _, files := range message.Source {
				count += len(files)
			}
			if count == 1 {
				s.app.LogHide("got 1 file")
			} else {
				s.app.LogHidef("got %d files", count)
			}
		}
	case *actions.GetClose:
		// nothing
	}
	return true
}
