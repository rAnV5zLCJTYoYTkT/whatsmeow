// Copyright (c) 2024 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package whatsmeow implements a WhatsApp web API client.
package whatsmeow

import (
	"context"
	"sync"
	"sync/atomic"

	"go.mau.fi/whatsme"go.mau.fi/whatsmeow/util/log"
)

// EventHandler is a function)

// Client is the main WhatsApp web client.
type Client struct {
	Store   *store.Device
	Log     log.Logger

	// Event handlers
	eventHandlersLock sync.RWMutex
	eventHandlers     []wrappedEventHandler
	lastHandlerID     uint32

	// Connection state
	connected     atomic.Bool
	connectLock   sync.Mutex

	// Context for managing goroutine lifecycle
	ctx    context.Context
	cancel context.CancelFunc
}

type wrappedEventHandler struct {
	fn EventHandler
	id uint32
}

// NewClient creates a new WhatsApp web client with the given device store and logger.
func NewClient(deviceStore *store.Device, log log.Logger) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		Store:  deviceStore,
		Log:    log,
		ctx:    ctx,
		cancel: cancel,
	}
}

// AddEventHandler adds an event handler function and returns an ID that can be
// used to remove the handler later with RemoveEventHandler.
func (cli *Client) AddEventHandler(handler EventHandler) uint32 {
	id := atomic.AddUint32(&cli.lastHandlerID, 1)
	cli.eventHandlersLock.Lock()
	cli.eventHandlers = append(cli.eventHandlers, wrappedEventHandler{fn: handler, id: id})
	cli.eventHandlersLock.Unlock()
	return id
}

// RemoveEventHandler removes the event handler with the given ID.
// Returns true if the handler was found and removed.
func (cli *Client) RemoveEventHandler(id uint32) bool {
	cli.eventHandlersLock.Lock()
	defer cli.eventHandlersLock.Unlock()
	for i, handler := range cli.eventHandlers {
		if handler.id == id {
			cli.eventHandlers = append(cli.eventHandlers[:i], cli.eventHandlers[i+1:]...)
			return true
		}
	}
	return false
}

// dispatch sends an event to all registered event handlers.
func (cli *Client) dispatch(evt interface{}) {
	cli.eventHandlersLock.RLock()
	handlers := cli.eventHandlers
	cli.eventHandlersLock.RUnlock()
	for _, handler := range handlers {
		handler.fn(evt)
	}
}

// IsConnected returns true if the client is currently connected to WhatsApp.
func (cli *Client) IsConnected() bool {
	return cli.connected.Load()
}

// IsLoggedIn returns true if the client has valid credentials stored.
func (cli *Client) IsLoggedIn() bool {
	return cli.Store != nil && cli.Store.ID != nil
}

// Disconnect disconnects the client from WhatsApp.
func (cli *Client) Disconnect() {
	if !cli.connected.Load() {
		return
	}
	cli.Log.Infof("Disconnecting from WhatsApp")
	cli.cancel()
	cli.connected.Store(false)
	cli.dispatch(&events.Disconnected{})
}

// GetJID returns the JID of the logged-in user, or an empty JID if not logged in.
func (cli *Client) GetJID() types.JID {
	if cli.Store == nil || cli.Store.ID == nil {
		return types.EmptyJID
	}
	return *cli.Store.ID
}
