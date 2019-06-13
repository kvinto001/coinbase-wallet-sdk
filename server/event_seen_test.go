// Copyright (c) 2019 Coinbase, Inc. See LICENSE

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CoinbaseWallet/walletlinkd/store/models"
	"github.com/stretchr/testify/require"
)

func TestMarkEventSeenNoSession(t *testing.T) {
	srv := NewServer(nil)

	req, err := http.NewRequest("POST", "/events/123/seen", nil)
	require.Nil(t, err)

	sessionID := "456"
	sessionKey := "789"
	req.SetBasicAuth(sessionID, sessionKey)

	rr := httptest.NewRecorder()
	srv.router.ServeHTTP(rr, req)

	resp := rr.Result()
	require.Equal(t, 401, resp.StatusCode)

	body := markEventSeenResponse{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	require.Nil(t, err)
	require.Equal(t, responseErrorInvalidSessionCredentials, body.Error)
	require.False(t, body.Success)
}

func TestMarkEventSeenInvalidSessionKey(t *testing.T) {
	srv := NewServer(nil)
	sessionID := "456"

	s := models.Session{ID: sessionID, Key: "correctKey"}
	s.Save(srv.store)

	req, err := http.NewRequest("POST", "/events/123/seen", nil)
	require.Nil(t, err)

	sessionKey := "incorrectKey"
	req.SetBasicAuth(sessionID, sessionKey)

	rr := httptest.NewRecorder()
	srv.router.ServeHTTP(rr, req)

	resp := rr.Result()
	require.Equal(t, 401, resp.StatusCode)

	body := markEventSeenResponse{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	require.Nil(t, err)
	require.Equal(t, responseErrorInvalidSessionCredentials, body.Error)
	require.False(t, body.Success)
}

func TestMarkEventSeenNoEvent(t *testing.T) {
	srv := NewServer(nil)
	sessionID := "123"
	sessionKey := "456"

	s := models.Session{ID: sessionID, Key: sessionKey}
	err := s.Save(srv.store)
	require.Nil(t, err)

	req, err := http.NewRequest("POST", "/events/789/seen", nil)
	require.Nil(t, err)

	req.SetBasicAuth(sessionID, sessionKey)

	rr := httptest.NewRecorder()
	srv.router.ServeHTTP(rr, req)

	resp := rr.Result()
	require.Equal(t, 200, resp.StatusCode)

	body := markEventSeenResponse{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	require.Nil(t, err)
	require.Empty(t, body.Error)
	require.True(t, body.Success)
}

func TestMarkEventSeen(t *testing.T) {
	srv := NewServer(nil)
	sessionID := "123"
	sessionKey := "456"
	eventID := "789"

	s := models.Session{ID: sessionID, Key: sessionKey}
	err := s.Save(srv.store)
	require.Nil(t, err)

	name := "name"
	data := "data"
	e := models.Event{ID: eventID, Event: name, Data: data}
	err = e.Save(srv.store, sessionID)
	require.Nil(t, err)

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("/events/%s/seen", eventID),
		nil,
	)
	require.Nil(t, err)

	req.SetBasicAuth(sessionID, sessionKey)

	rr := httptest.NewRecorder()
	srv.router.ServeHTTP(rr, req)

	resp := rr.Result()
	require.Equal(t, 200, resp.StatusCode)

	body := markEventSeenResponse{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	require.Nil(t, err)
	require.Empty(t, body.Error)
	require.True(t, body.Success)
}
