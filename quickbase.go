// Copyright 2014 James Massara. All rights reserved.
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

// QBError represents a detailed error message returned by the QuickBase API.
type QBError struct {
	msg    string
	Code   int
	Detail string
}

func (e *QBError) Error() string { return e.msg }

// A Conn represents an authenticated connection to a QuickBase domain.
type Conn struct {
	url      string
	ticket   string
	apptoken string
	client   *http.Client
}

// Login logs into a QuickBase domain with the given username and password.
//
// See http://www.quickbase.com/api-guide/index.html#authenticate.html for more
// information.
func Login(url, username, password string) (*Conn, error) {
	conn := &Conn{
		url: url,
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		},
	}
	req := &authRequest{Username: username, Password: password}
	rsp := &authResponse{}

	if err := conn.do("API_Authenticate", "", req, rsp); err != nil {
		return nil, err
	}

	if rsp.ErrorCode != 0 {
		return nil, &QBError{msg: rsp.ErrorText, Code: rsp.ErrorCode, Detail: rsp.ErrorDetail}
	} else {
		conn.ticket = rsp.Ticket
	}

	return conn, nil
}

// Ticket returns the authentication ticket received after successfully logging
// into a QuickBase domain.
func (c *Conn) Ticket() string { return c.ticket }

// SetAppToken sets the QuickBase application token used in subsequent API
// requests. This only needs to be called if the database is configured to
// require app tokens.
func (c *Conn) SetAppToken(apptoken string) { c.apptoken = apptoken }

// AppToken returns the QuickBase application token used for API requests.
func (c *Conn) AppToken() string { return c.apptoken }

// do handles making QuickBase API calls for a given action.
func (c *Conn) do(action, dbid string, request, response interface{}) error {
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
		url = c.url + "/db/main"
	} else {
		url = c.url + "/db/" + dbid
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postdata))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Quickbase-Action", action)

	rsp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	err = xml.NewDecoder(rsp.Body).Decode(response)
	rsp.Body.Close()

	return err
}

// field is a field id's label and value.
type field struct {
	label, value string
}

// Record is a QuickBase record.
type Record struct {
	Id       int
	UpdateId int64

	// fields is a map of field ids to its label and value.
	fields map[int]field
}

// NewRecord creates an empty record.
func NewRecord() *Record {
	return &Record{fields: make(map[int]field)}
}

// GetFieldsById returns a map if fields keyed by the field id.
func (r *Record) GetFieldsById() map[int]string {
	fields := make(map[int]string)
	for id, f := range r.fields {
		fields[id] = f.value
	}
	return fields
}

// GetFieldsByName returns a map if fields keyed by the field label.
func (r *Record) GetFieldsByName() map[string]string {
	fields := make(map[string]string)
	for _, f := range r.fields {
		fields[f.label] = f.value
	}
	return fields
}

// AddFields adds fields to a record. It only supports adding fields by the
// field id. Field ids in QuickBase will never change while field labels can.
func (r *Record) AddFields(fields map[int]string) {
	for id, value := range fields {
		r.fields[id] = field{value: value}
	}
}

// UpdateFields updates the fields in a record. If the field doesn't exist in
// the record then it will be added to it.
//
// Like AddFields, it only supports adding fields by the field id.
func (r *Record) UpdateFields(fields map[int]string) {
	for id, value := range fields {
		if _, ok := r.fields[id]; !ok {
			r.fields[id] = field{value: value}
		} else {
			r.fields[id] = field{value: value, label: r.fields[id].label}
		}
	}
}

// recordField is the XML structure for adding or editing a record.
type recordField struct {
	Id    int    `xml:"fid,attr"`
	Value string `xml:",chardata"`
}

// qbRequest has XML fields used in every API call.
type qbRequest struct {
	Ticket   string `xml:"ticket"`
	AppToken string `xml:"apptoken,omitempty"`
}

// qbResponse has XML fields returned from every API call.
type qbResponse struct {
	Action      string `xml:"action"`
	ErrorCode   int    `xml:"errcode"`
	ErrorText   string `xml:"errtext"`
	ErrorDetail string `xml:"errdetail"`
}

type authRequest struct {
	XMLName  xml.Name `xml:"qdbapi"`
	Username string   `xml:"username"`
	Password string   `xml:"password"`
}

type authResponse struct {
	XMLName xml.Name `xml:"qdbapi"`
	Ticket  string   `xml:"ticket"`
	UserId  string   `xml:"userid"`

	qbResponse // XML fields returned in every API call
}
