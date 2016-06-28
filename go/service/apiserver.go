// Copyright 2016 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

package service

import (
	"github.com/keybase/client/go/libkb"
	keybase1 "github.com/keybase/client/go/protocol"
	rpc "github.com/keybase/go-framed-msgpack-rpc"
	jsonw "github.com/keybase/go-jsonw"
	"golang.org/x/net/context"
)

type APIServerHandler struct {
	*BaseHandler
	libkb.Contextified
}

func NewAPIServerHandler(xp rpc.Transporter, g *libkb.GlobalContext) *APIServerHandler {
	return &APIServerHandler{
		BaseHandler:  NewBaseHandler(xp),
		Contextified: libkb.NewContextified(g),
	}
}

func (a *APIServerHandler) Get(_ context.Context, arg keybase1.GetArg) (keybase1.APIRes, error) {
	return a.doGet(arg)
}

func (a *APIServerHandler) Post(_ context.Context, arg keybase1.PostArg) (keybase1.APIRes, error) {
	return a.doPost(arg)
}

func (a *APIServerHandler) PostJSON(_ context.Context, arg keybase1.PostJSONArg) (keybase1.APIRes, error) {
	return a.doPostJSON(arg)
}

type GenericArg interface {
	GetEndpoint() string
	GetHTTPArgs() []keybase1.StringKVPair
	GetHttpStatuses() []int
	GetAppStatusCodes() []int
}

func (a *APIServerHandler) setupArg(arg GenericArg) libkb.APIArg {
	// Form http arg dict
	kbargs := make(libkb.HTTPArgs)
	for _, harg := range arg.GetHTTPArgs() {
		kbargs[harg.Key] = libkb.S{Val: harg.Value}
	}

	// Acceptable http status list
	var httpStatuses []int
	for _, hstat := range arg.GetHttpStatuses() {
		httpStatuses = append(httpStatuses, hstat)
	}

	// Acceptable app status code list
	var appStatusCodes []int
	for _, ac := range arg.GetAppStatusCodes() {
		appStatusCodes = append(appStatusCodes, ac)
	}

	// Do the API call
	kbarg := libkb.APIArg{
		Endpoint:       arg.GetEndpoint(),
		NeedSession:    true,
		Args:           kbargs,
		HTTPStatus:     httpStatuses,
		AppStatusCodes: appStatusCodes,
		Contextified:   libkb.NewContextified(a.G()),
	}

	return kbarg
}

func (a *APIServerHandler) doGet(arg keybase1.GetArg) (keybase1.APIRes, error) {
	res, err := a.G().API.Get(a.setupArg(arg))
	if err != nil {
		return keybase1.APIRes{}, err
	}
	return a.convertRes(res), nil
}

func (a *APIServerHandler) doPost(arg keybase1.PostArg) (keybase1.APIRes, error) {
	res, err := a.G().API.Post(a.setupArg(arg))
	if err != nil {
		return keybase1.APIRes{}, err
	}
	return a.convertRes(res), nil
}

func (a *APIServerHandler) doPostJSON(rawarg keybase1.PostJSONArg) (keybase1.APIRes, error) {

	arg := a.setupArg(rawarg)
	jsonPayload := make(libkb.JSONPayload)
	for _, kvpair := range rawarg.JSONPayload {
		jsonPayload[kvpair.Key] = kvpair.Value
	}
	arg.JSONPayload = jsonPayload

	res, err := a.G().API.PostJSON(arg)
	if err != nil {
		return keybase1.APIRes{}, err
	}

	return a.convertRes(res), nil
}

func (a *APIServerHandler) convertRes(res *libkb.APIRes) keybase1.APIRes {
	// Translate the result
	var ares keybase1.APIRes
	mstatus, err := res.Status.Marshal()
	if err == nil {
		ares.Status = string(mstatus[:])
	}
	mbody, err := res.Body.Marshal()
	if err == nil {
		ares.Body = string(mbody[:])
	}
	ares.HttpStatus = res.HTTPStatus

	appStatus := jsonw.NewWrapper(res.AppStatus)
	mappstatus, err := appStatus.Marshal()
	if err == nil {
		ares.AppStatus = string(mappstatus[:])
	}

	return ares
}
