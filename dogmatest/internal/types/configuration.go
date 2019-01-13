package types

import (
	"fmt"
	"reflect"
)

type Configuration struct {
	controllers map[string]Controller
	commands    map[reflect.Type]Controller
	events      map[reflect.Type][]Controller
}

func (cfg *Configuration) RegisterController(c Controller) {
	if x, ok := cfg.controllers[c.Name()]; ok {
		panic(fmt.Sprintf(
			"can not use name '%s' for %#v because the name is already used by %#v",
			c.Name(),
			c.Handler(),
			x.Handler(),
		))
	}

	if cfg.controllers == nil {
		cfg.controllers = map[string]Controller{}
	}

	cfg.controllers[c.Name()] = c
}

func (cfg *Configuration) RouteCommand(t reflect.Type, c Controller) {
	if x, ok := cfg.commands[t]; ok {
		routeConflict(t, c, "commands", x, "commands")
	}

	if x, ok := cfg.events[t]; ok {
		routeConflict(t, c, "commands", x[0], "events")
	}

	if cfg.commands == nil {
		cfg.commands = map[reflect.Type]Controller{}
	}
	cfg.commands[t] = c
}

func (cfg *Configuration) RegisterEvent(t reflect.Type, c Controller) {
	if x, ok := cfg.commands[t]; ok {
		routeConflict(t, c, "events", x, "commands")
	}

	if cfg.events == nil {
		cfg.events = map[reflect.Type][]Controller{}
	}

	cfg.events[t] = append(cfg.events[t], c)
}

func (cfg *Configuration) Controllers() []Controller {
	var controllers []Controller

	for _, c := range cfg.controllers {
		controllers = append(controllers, c)
	}

	return controllers
}

func (cfg *Configuration) Classes() map[reflect.Type]MessageClass {
	m := map[reflect.Type]MessageClass{}

	for t := range cfg.commands {
		m[t] = Command
	}

	for t := range cfg.events {
		m[t] = Event
	}

	return m
}

func (cfg *Configuration) Routes() map[reflect.Type][]Controller {
	m := map[reflect.Type][]Controller{}

	for t, c := range cfg.commands {
		m[t] = append(m[t], c)
	}

	for t, c := range cfg.events {
		m[t] = append(m[t], c...)
	}

	return m
}

func routeConflict(
	t reflect.Type,
	c Controller, ct string,
	x Controller, xt string,
) {
	panic(fmt.Sprintf(
		"can not route %s of type %s to '%s' because they are already routed to '%s' as %s",
		t,
		ct,
		c.Name(),
		x.Name(),
		xt,
	))
}
