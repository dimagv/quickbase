// Copyright 2013 James Massara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package quickbase_test

import (
	"encoding/xml"
	"fmt"
	"github.com/jmassara/quickbase"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var (
	mux      *http.ServeMux
	server   *httptest.Server
	qbdomain *url.URL
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewTLSServer(mux)
	qbdomain, _ = url.Parse(server.URL)
}

func teardown() {
	server.Close()
}

func testAction(t *testing.T, h http.Header, expect string) {
	if got, ok := h["Quickbase-Action"]; ok {
		if got[0] != expect {
			t.Errorf("expected: %q, got: %q", expect, got[0])
		}
	} else {
		t.Errorf("Missing header: Quickbase-Action")
	}
}

var authtests = []struct {
	username string
	password string
	path     string
	xmlresp  string
	qbresp   *quickbase.AuthResponse
	err      quickbase.QBError
}{
	{
		username: "PTBarnum",
		password: "TopSecret",
		path:     "/success",
		xmlresp:  authSuccess,
		qbresp: &quickbase.AuthResponse{
			XMLName:   xml.Name{Local: "qdbapi"},
			Action:    "api_authenticate",
			ErrorCode: 0,
			ErrorText: "No error",
			Ticket:    "2_beeinrxmv_dpvx_b_crf8ttndjwyf9bui94rhciirqcs",
			UserId:    "112245.efy7",
		},
		err: quickbase.QBError{},
	},
	{
		username: "PTBarnum",
		password: "WrongPassword",
		path:     "/failure",
		xmlresp:  authFailure,
		qbresp:   nil,
		err: quickbase.QBError{
			Code:   20,
			Detail: "Sorry! You entered the wrong E-Mail or Screen Name or Password. Try again.",
		},
	},
}

func TestAuthentication(t *testing.T) {
	setup()
	defer teardown()

	for _, tt := range authtests {
		qburl := qbdomain
		qburl.Path = tt.path

		mux.HandleFunc(tt.path+"/db/main",
			func(w http.ResponseWriter, r *http.Request) {
				testAction(t, r.Header, "API_Authenticate")
				fmt.Fprint(w, tt.xmlresp)
			},
		)

		qb := quickbase.New(qburl)
		resp, err := qb.Authenticate(&quickbase.AuthRequest{
			Username: tt.username,
			Password: tt.password,
		})

		if err != nil {
			if err.Detail != tt.err.Detail {
				t.Errorf("expected: %s, got: %s", tt.err.Detail, err.Detail)
			}
		}

		if !reflect.DeepEqual(resp, tt.qbresp) {
			t.Errorf("expected: %+v, got: %+v", tt.qbresp, resp)
		}
	}
}

// QuickBase API responses

// API_Authenticate
var authSuccess = `<?xml version="1.0" ?>
<qdbapi>
	<action>api_authenticate</action>
	<errcode>0</errcode>
	<errtext>No error</errtext>
	<ticket>2_beeinrxmv_dpvx_b_crf8ttndjwyf9bui94rhciirqcs</ticket>
	<userid>112245.efy7</userid>
</qdbapi>
`

var authFailure = `<?xml version="1.0" ?>
<qdbapi>
	<action>API_Authenticate</action>
	<errcode>20</errcode>
	<errtext>Unknown username/password</errtext>
	<errdetail>Sorry! You entered the wrong E-Mail or Screen Name or Password. Try again.</errdetail>
</qdbapi>
`
