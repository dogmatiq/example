package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages/events"
)

const retainEvents = 10

type AccountListSSEProjectionHandler struct {
	dogma.NoCompactBehavior
	dogma.NoTimeoutHintBehavior

	m           sync.RWMutex
	resources   map[string]string
	state       accountListResponse
	version     int
	events      []dogma.Message
	subscribers map[*http.Request]http.ResponseWriter
}

func (h *AccountListSSEProjectionHandler) Subscribe(req *http.Request, w http.ResponseWriter) {
	h.m.Lock()
	defer h.m.Unlock()

	fmt.Println("subscribing ...")

	var resumeAt int
	if s := req.Header.Get("Last-Event-ID"); s != "" {
		prev, err := strconv.Atoi(s)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resumeAt = prev + 1
		fmt.Printf("resuming at %d\n", resumeAt)
	}

	if h.subscribers == nil {
		h.subscribers = map[*http.Request]http.ResponseWriter{}
	}

	h.subscribers[req] = w

	if resumeAt == 0 {
		data, err := json.Marshal(h.state)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("sending state at %d\n", h.version)
		writeSSE(w, h.version, "domain-state", data)
		return
	}

	oldest := 1 + h.version - len(h.events)
	if resumeAt < oldest {
		http.Error(w, "Last-Event-ID is too old!", http.StatusNotFound)
		return
	}

	index := resumeAt - oldest
	for _, ev := range h.events[index:] {
		fmt.Printf("sending event at %d\n", resumeAt)
		writeSSEEvent(w, resumeAt, ev)
		resumeAt++
	}
}

func (h *AccountListSSEProjectionHandler) Unsubscribe(req *http.Request) {
	h.m.Lock()
	defer h.m.Unlock()

	fmt.Println("unsubscribed")
	delete(h.subscribers, req)
}

func (h *AccountListSSEProjectionHandler) Configure(c dogma.ProjectionConfigurer) {
	c.Identity("account-list-sse", "eb8dedb6-784c-487f-8947-6468fa960e1e")

	c.ConsumesEventType(events.AccountOpened{})
	c.ConsumesEventType(events.AccountCredited{})
	c.ConsumesEventType(events.AccountDebited{})
}

func (h *AccountListSSEProjectionHandler) HandleEvent(
	ctx context.Context,
	r, c, n []byte,
	s dogma.ProjectionEventScope,
	m dogma.Message,
) (bool, error) {
	h.m.Lock()
	defer h.m.Unlock()

	if h.resources[string(r)] != string(c) {
		return false, nil
	}

	switch m := m.(type) {
	case events.AccountOpened:
		h.state.Accounts = append(
			h.state.Accounts,
			accountListEntry{
				ID:   m.AccountID,
				Name: m.AccountName,
			},
		)

	case events.AccountCredited:
		for i, a := range h.state.Accounts {
			if a.ID == m.AccountID {
				a.Balance += m.Amount
				h.state.Accounts[i] = a
			}
		}

	case events.AccountDebited:
		for i, a := range h.state.Accounts {
			if a.ID == m.AccountID {
				a.Balance -= m.Amount
				h.state.Accounts[i] = a
			}
		}

	default:
		panic(dogma.UnexpectedMessage)
	}

	h.version++
	h.events = append(h.events, m)
	if len(h.events) > retainEvents {
		h.events = h.events[1:]
	}

	for _, w := range h.subscribers {
		writeSSEEvent(w, h.version, m)
	}

	if h.resources == nil {
		h.resources = map[string]string{}
	}

	h.resources[string(r)] = string(n)

	return true, nil
}

func (h *AccountListSSEProjectionHandler) ResourceVersion(ctx context.Context, r []byte) ([]byte, error) {
	h.m.RLock()
	defer h.m.RUnlock()

	return []byte(h.resources[string(r)]), nil
}

func (h *AccountListSSEProjectionHandler) CloseResource(ctx context.Context, r []byte) error {
	h.m.Lock()
	defer h.m.Unlock()

	delete(h.resources, string(r))
	return nil
}

// writeSSE sends a server-sent event and logs it to the console.
func writeSSE(
	w http.ResponseWriter,
	version int,
	event string,
	data []byte,
) {
	writeSSEField(w, "event", event)
	writeSSEField(w, "data", string(data))
	writeSSEField(w, "id", strconv.Itoa(version))
	fmt.Fprintf(w, "\n")
	w.(http.Flusher).Flush()
}

func writeSSEEvent(
	w http.ResponseWriter,
	version int,
	m dogma.Message,
) {
	env := []interface{}{
		reflect.TypeOf(m).String(),
		m,
	}

	data, _ := json.Marshal(env)

	writeSSE(
		w,
		version,
		"domain-event",
		data,
	)
}

// writeSSEField sends a single field within an event.
func writeSSEField(
	wr http.ResponseWriter,
	k, v string,
) {
	for _, line := range strings.Split(v, "\n") {
		fmt.Fprintf(wr, "%s: %s\n", k, line)
	}
}

type Subscribable interface {
	Subscribe(req *http.Request, w http.ResponseWriter)
	Unsubscribe(req *http.Request)
}

type accountListSSEHandler struct {
	Subscribable Subscribable
}

func (h *accountListSSEHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	h.Subscribable.Subscribe(req, w)
	defer h.Subscribable.Unsubscribe(req)
	<-req.Context().Done()
}
