// beehive's web-module.
package web

import (
	"encoding/json"
	"github.com/hoisie/web"
	"github.com/muesli/beehive/modules"
	"io/ioutil"
	"log"
)

var (
	cIn chan modules.Event
)

type WebBee struct {
	Addr string
}

func (mod *WebBee) Name() string {
	return "webbee"
}

func (mod *WebBee) Description() string {
	return "A RESTful HTTP module for beehive"
}

func (mod *WebBee) Run(channelIn chan modules.Event) {
	cIn = channelIn
	go web.Run(mod.Addr)
}

func (mod *WebBee) Events() []modules.EventDescriptor {
	events := []modules.EventDescriptor{
		modules.EventDescriptor{
			Namespace:   mod.Name(),
			Name:        "post",
			Description: "A POST call was received by the HTTP server",
			Options: []modules.PlaceholderDescriptor{
				modules.PlaceholderDescriptor{
					Name:        "json",
					Description: "JSON map received from caller",
					Type:        "map",
				},
				modules.PlaceholderDescriptor{
					Name:        "ip",
					Description: "IP of the caller",
					Type:        "string",
				},
			},
		},
		modules.EventDescriptor{
			Namespace:   mod.Name(),
			Name:        "get",
			Description: "A GET call was received by the HTTP server",
			Options: []modules.PlaceholderDescriptor{
				modules.PlaceholderDescriptor{
					Name:        "query_params",
					Description: "Map of query parameters received from caller",
					Type:        "map",
				},
				modules.PlaceholderDescriptor{
					Name:        "ip",
					Description: "IP of the caller",
					Type:        "string",
				},
			},
		},
	}
	return events
}

func (mod *WebBee) Actions() []modules.ActionDescriptor {
	actions := []modules.ActionDescriptor{}
	return actions
}

func (mod *WebBee) Action(action modules.Action) []modules.Placeholder {
	outs := []modules.Placeholder{}
	return outs
}

func GetRequest(ctx *web.Context) {
	//FIXME
	ms := make(map[string]string)
	ev := modules.Event{
		Namespace: "webbee",
		Name:      "get",
		Options: []modules.Placeholder{
			modules.Placeholder{
				Name:  "query_params",
				Type:  "map",
				Value: ms,
			},
			modules.Placeholder{
				Name:  "ip",
				Type:  "string",
				Value: "tbd",
			},
		},
	}
	cIn <- ev
}

func PostRequest(ctx *web.Context) {
	b, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	var payload interface{}
	err = json.Unmarshal(b, &payload)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	ev := modules.Event{
		Namespace: "webbee",
		Name:      "post",
		Options: []modules.Placeholder{
			modules.Placeholder{
				Name:  "json",
				Type:  "map",
				Value: payload,
			},
			modules.Placeholder{
				Name:  "ip",
				Type:  "string",
				Value: "tbd",
			},
		},
	}
	cIn <- ev
}

func init() {
	w := WebBee{
		Addr: "0.0.0.0:12345",
	}
	web.Get("/event", GetRequest)
	web.Post("/event", PostRequest)

	modules.RegisterModule(&w)
}
