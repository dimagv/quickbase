// Copyright 2013 James Massara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package quickbase

import (
	"encoding/xml"
)

// API_AddRecord request parameters.
// See http://goo.gl/SELQDd for more details.
type AddRecordRequest struct {
	XMLName     xml.Name      `xml:"qdbapi"`
	Fields      []recordField `xml:"field"`
	DispRec     int           `xml:"disprec,omitempty"`
	FForm       int           `xml:"fform,omitempty"`
	IgnoreError int           `xml:"ignoreError,omitempty"`
	MsInUTC     int           `xml:"msInUTC,omitempty"`
	Udata       string        `xml:"udata,omitempty"`

	qbrequest // Fields required for every request
}

// Response to an API_AddRecord request.
// See http://goo.gl/SELQDd for more details.
type AddRecordResponse struct {
	XMLName  xml.Name `xml:"qdbapi"`
	RecordId string   `xml:"rid"`
	UpdateId string   `xml:"update_id"`

	qbresponse // Fields returned in every response
}

// AddField adds a field to the record.
func (r *AddRecordRequest) AddField(id int, value string) {
	r.Fields = append(r.Fields, recordField{Id: id, Value: value})
}

// AddRecord adds a new record to a QuickBase database.
func (qb *QuickBase) AddRecord(dbid string, req *AddRecordRequest) (*AddRecordResponse, error) {
	req.Ticket = qb.ticket
	req.AppToken = qb.apptoken

	resp := new(AddRecordResponse)
	if err := qb.query("API_AddRecord", dbid, req, resp); err != nil {
		return nil, err
	}

	if resp.ErrorCode != 0 {
		return nil, &QBError{
			msg:    resp.ErrorText,
			Code:   resp.ErrorCode,
			Detail: resp.ErrorDetail,
		}
	}

	return resp, nil
}
