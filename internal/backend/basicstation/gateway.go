package basicstation

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/brocaar/chirpstack-gateway-bridge/internal/backend/events"
	"github.com/brocaar/lorawan"
)

var (
	errGatewayDoesNotExist = errors.New("gateway does not exist")
)

type gateway struct {
	conn          *websocket.Conn
	configVersion string
}

type gateways struct {
	sync.RWMutex
	gateways map[lorawan.EUI64]gateway

	subscribeEventFunc func(events.Subscribe)
}

func (g *gateways) get(id lorawan.EUI64) (gateway, error) {
	g.RLock()
	defer g.RUnlock()

	gw, ok := g.gateways[id]
	if !ok {
		return gw, errGatewayDoesNotExist
	}
	return gw, nil
}

func (g *gateways) set(id lorawan.EUI64, gw gateway) error {
	g.Lock()
	defer g.Unlock()

	g.gateways[id] = gw

	if g.subscribeEventFunc != nil {
		g.subscribeEventFunc(events.Subscribe{Subscribe: true, GatewayID: id})
	}

	return nil
}

func (g *gateways) remove(id lorawan.EUI64) error {
	g.Lock()
	defer g.Unlock()

	if g.subscribeEventFunc != nil {
		g.subscribeEventFunc(events.Subscribe{Subscribe: false, GatewayID: id})
	}

	delete(g.gateways, id)
	return nil
}
