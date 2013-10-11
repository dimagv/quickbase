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
	"net/url"
)

// QuickBase represents a QuickBase domain.
type QuickBase struct {
	url    *url.URL
	client *http.Client
}

// New creates a new QuickBase.
func New(url *url.URL) *QuickBase {
	if url.Scheme == "" {
		url.Scheme = "https"
	}
	return &QuickBase{
		url: url,
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		},
	}
}

// QBError represents a detailed error message returned by the QuickBase API
type QBError struct {
	msg    string
	Code   int
	Detail string
}

func (e *QBError) Error() string { return e.msg }

// Helper functions

// Send an XML API request to QuickBase and populate `response' with the XML
// response returned from QuickBase.
func (qb *QuickBase) query(params map[string]string, request, response interface{}) error {
	postdata, err := xml.MarshalIndent(request, "", " ")
	if err != nil {
		return err
	}

	// Convert `&#39;' to `&apos;' because xml.Marshal runs EscapeText and
	// transforms `'' to `&#39;', which is fine for HTML but not XML.
	postdata = bytes.Replace(postdata, []byte("&#39;"), []byte("&apos;"), -1)

	req, err := http.NewRequest("POST", params["url"], bytes.NewBuffer(postdata))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Quickbase-Action", params["action"])

	res, err := qb.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return xml.NewDecoder(res.Body).Decode(response)
}

// makeParams returns an initial map with the action set.
func makeParams(action string) map[string]string {
	params := make(map[string]string)
	params["action"] = action
	return params
}
