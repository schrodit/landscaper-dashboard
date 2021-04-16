// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"gopkg.in/olahol/melody.v1"
)

type Request struct {
	Path string          `json:"path"`
	Data json.RawMessage `json:"data"`
}

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
}

func OkResponse(data interface{}) *Response {
	return &Response{
		Code: http.StatusOK,
		Data: data,
	}
}

type CommonRequestHandlerFunc func(data []byte) (*Response, error)

type WebSocketWriter interface {
	Write(res *Response)
}

type WebsocketRequestHandlerFunc func(writer WebSocketWriter, data []byte) error

type Router struct {
	mux            sync.RWMutex
	log            logr.Logger
	httpPaths      map[string]CommonRequestHandlerFunc
	websocketPaths map[string]WebsocketRequestHandlerFunc
	httpServer     *gin.Engine
}

// NewRouter creates a new router
func NewRouter(log logr.Logger, server *gin.Engine) *Router {
	return &Router{
		log:            log,
		httpServer:     server,
		httpPaths:      map[string]CommonRequestHandlerFunc{},
		websocketPaths: map[string]WebsocketRequestHandlerFunc{},
	}
}

func (r *Router) Register(m *melody.Melody) error {
	m.HandleMessage(r.HandleWebSocket)
	return nil
}

// HandleWebSocket handles a websocket request.
func (r *Router) HandleWebSocket(session *melody.Session, data []byte) {
	writer := NewWebSocketWriter(r.log, session)
	req := Request{}
	if err := json.Unmarshal(data, &req); err != nil {
		r.log.Error(err, "malformed request")
		writer.Write(&Response{
			Code: http.StatusMisdirectedRequest,
		})
		return
	}

	wHdlr := r.GetWebsocketHandler(req.Path)
	if wHdlr != nil {
		if err := wHdlr(writer, req.Data); err != nil {
			r.log.Error(err, "malformed request")
		}
		return
	}

	hdlr := r.GetRequestHandler(req.Path)
	if hdlr == nil {
		r.log.Info("no path found", "path", req.Path)
		writer.Write(&Response{
			Code: http.StatusMisdirectedRequest,
		})
		return
	}
	res, err := hdlr(req.Data)
	if err != nil {
		r.log.Error(err, "unable to handle request", "path", req.Path)
	}
	writer.Write(res)
}

// HandleHttpRequest handles a http request.
func (r *Router) HandleHttpRequest(c *gin.Context) {
	var data bytes.Buffer
	if _, err := io.Copy(&data, c.Request.Body); err != nil {
		r.WriteHTTPResponse(c, &Response{
			Code: http.StatusInternalServerError,
			Data: "unable to read request body",
		})
		return
	}
	req := Request{
		Path: strings.TrimPrefix(c.Request.URL.Path, "/"),
		Data: data.Bytes(),
	}

	hdlr, ok := r.httpPaths[req.Path]
	if !ok {
		r.log.Info("no path found", "path", req.Path)
		r.WriteHTTPResponse(c, &Response{
			Code: http.StatusMisdirectedRequest,
		})
		return
	}
	res, err := hdlr(req.Data)
	if err != nil {
		r.log.Error(err, "unable to handle request", "path", req.Path)
	}
	r.WriteHTTPResponse(c, res)
}

// WriteHTTPResponse writes a http response to a writer
func (r *Router) WriteHTTPResponse(ctx *gin.Context, res *Response) {
	ctx.IndentedJSON(res.Code, res.Data)
}

/*
  Paths methods
*/

// AddCommonPath adds a new handler to the router.
func (r *Router) AddCommonPath(path string, hdlr CommonRequestHandlerFunc) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.httpPaths[path] = hdlr

	// register http handler
	r.httpServer.POST("/"+path, r.HandleHttpRequest)
}

// AddWebsocketPath adds a new handler to the router.
func (r *Router) AddWebsocketPath(path string, hdlr WebsocketRequestHandlerFunc) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.websocketPaths[path] = hdlr
}

// GetRequestHandler returns the request handler for a specific response.
func (r *Router) GetRequestHandler(path string) CommonRequestHandlerFunc {
	r.mux.RLock()
	defer r.mux.RUnlock()
	return r.httpPaths[path]
}

// GetWebsocketHandler returns the request handler for a specific response.
func (r *Router) GetWebsocketHandler(path string) WebsocketRequestHandlerFunc {
	r.mux.RLock()
	defer r.mux.RUnlock()
	return r.websocketPaths[path]
}

type webSocketWriter struct {
	log     logr.Logger
	session *melody.Session
}

func NewWebSocketWriter(log logr.Logger, session *melody.Session) WebSocketWriter {
	return &webSocketWriter{
		log:     log,
		session: session,
	}
}

func (w webSocketWriter) Write(res *Response) {
	if err := WriteWebSocketResponse(w.session, res); err != nil {
		w.log.Error(err, "unable to send response")
	}
}

// WriteWebSocketResponse writes a websocket reponse to a session
func WriteWebSocketResponse(session *melody.Session, res *Response) error {
	data, err := json.Marshal(res)
	if err != nil {
		return err
	}
	return session.Write(data)
}
