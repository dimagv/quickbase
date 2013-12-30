// Copyright 2013 James Massara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package quickbase

import (
	"encoding/xml"
)

// API_Authenticate request parameters.
// See http://goo.gl/eQSiZy for more details.
type authRequest struct {
	XMLName  xml.Name `xml:"qdbapi"`
	Username string   `xml:"username"`
	Password string   `xml:"password"`
	Hours    int      `xml:"hours,omitempty"`
	Udata    string   `xml:"udata,omitempty"`
}

// Response to an API_Authenticate request.
// See http://goo.gl/eQSiZy for more details.
type authResponse struct {
	XMLName xml.Name `xml:"qdbapi"`
	Ticket  string   `xml:"ticket"`
	UserId  string   `xml:"userid"`

	qbresponse // Fields returned in every response
}

// Autenticate gets a time-based ticket from QuickBase to use with other API
// requests.
func (qb *QuickBase) Authenticate(username, password string, hours int) error {
	request := authRequest{Username: username, Password: password, Hours: hours}
	response := new(authResponse)

	if err := qb.query("API_Authenticate", "", request, response); err != nil {
		return err
	}

	if response.ErrorCode != 0 {
		return &QBError{
			msg:    response.ErrorText,
			Code:   response.ErrorCode,
			Detail: response.ErrorDetail,
		}
	} else {
		qb.ticket = response.Ticket
	}

	return nil
}
