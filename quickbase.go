// Copyright 2013 James Massara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package quickbase implements pieces of the Quickbase API.
package quickbase

// References:
//    QuickBase API: http://www.quickbase.com/api-guide/index.html.

import (
	"bytes"
	"encoding/xml"
	"net/http"
)

// QuickBase represents a QuickBase host.
type QuickBase struct {
	host     string
	ticket   string
	apptoken string
	client   *http.Client
}

// New creates a new QuickBase.
func New(host string) *QuickBase {
	return &QuickBase{
		host: host,
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		},
	}
}

// Ticket returns the ticket received after successfully authenticating.
func (q *QuickBase) Ticket() string { return q.ticket }

// SetAppToken sets the QuickBase application token used for API requests.
func (q *QuickBase) SetAppToken(apptoken string) { q.apptoken = apptoken }

// AppToken returns the QuickBase application token used for API requests.
func (q *QuickBase) AppToken() string { return q.apptoken }

// QBError represents a detailed error message returned by the QuickBase API.
type QBError struct {
	msg    string
	Code   int
	Detail string
}

func (e *QBError) Error() string { return e.msg }

// qbrequest has xml fields used in every API request.
type qbrequest struct {
	Ticket   string `xml:"ticket"`
	AppToken string `xml:"apptoken,omitempty"`
}

// qbresponse has xml fields returned in every API response.
type qbresponse struct {
	Action      string `xml:"action"`
	ErrorCode   int    `xml:"errcode"`
	ErrorText   string `xml:"errtext"`
	ErrorDetail string `xml:"errdetail"`
	Udata       string `xml:"udata"`
}

// recordField is a field when adding or editing a record.
type recordField struct {
	Id    int    `xml:"fid,attr"`
	Value string `xml:",chardata"`
}

// Helper functions

// Send an XML API request to QuickBase and set response to the XML returned
// from QuickBase.
func (qb *QuickBase) query(action, dbid string, request, response interface{}) error {
	postdata, err := xml.MarshalIndent(request, "", " ")
	// postdata, err := xml.Marshal(request)
	if err != nil {
		return err
	}

	// Convert &#39; to &apos; because xml.Marshal runs EscapeText and
	// transforms ' to &#39; which is fine for HTML but not XML.
	postdata = bytes.Replace(postdata, []byte("&#39;"), []byte("&apos;"), -1)

	var url string
	if action == "API_Authenticate" {
		url = qb.host + "/db/main"
	} else {
		url = qb.host + "/db/" + dbid
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postdata))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Quickbase-Action", action)

	resp, err := qb.client.Do(req)
	if err != nil {
		return err
	}

	err = xml.NewDecoder(resp.Body).Decode(response)
	resp.Body.Close()

	return err
}
