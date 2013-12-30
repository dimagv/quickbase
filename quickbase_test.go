// Copyright 2013 James Massara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package quickbase_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmassara/quickbase"
)

var (
	mux      *http.ServeMux
	server   *httptest.Server
	qbdomain string
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	qbdomain = server.URL
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

var authTests = []struct {
	path     string
	username string
	password string
	xmlresp  string
	err      string
}{
	{
		path:     "/success",
		username: "PTBarnum",
		password: "TopSecret",
		xmlresp:  authSuccess,
	},
	{
		path:     "/failure",
		username: "PTBarnum",
		password: "WrongPassword",
		xmlresp:  authFailure,
		err:      "Unknown username/password",
	},
}

func TestAuthentication(t *testing.T) {
	setup()
	defer teardown()
	for _, tt := range authTests {
		mux.HandleFunc(tt.path+"/db/main",
			func(w http.ResponseWriter, r *http.Request) {
				testAction(t, r.Header, "API_Authenticate")
				fmt.Fprint(w, tt.xmlresp)
			},
		)

		qb := quickbase.New(qbdomain + tt.path)
		err := qb.Authenticate(tt.username, tt.password, 1)

		if err != nil && err.Error() != tt.err {
			t.Fatal(err)
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
