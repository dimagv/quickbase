// Copyright 2013 James Massara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package quickbase

import (
	"encoding/xml"
)

// API_EditRecord request parameters.
// See http://goo.gl/3OgE49 for more details.
type EditRecordRequest struct {
	XMLName     xml.Name      `xml:"qdbapi"`
	RecordId    int           `xml:"rid,omitempty"`
	Key         string        `xml:"key,omitempty"`
	UpdateId    int           `xml:"update_id,omitempty"`
	Fields      []recordField `xml:"field"`
	DispRec     int           `xml:"disprec,omitempty"`
	FForm       int           `xml:"fform,omitempty"`
	IgnoreError int           `xml:"ignoreError,omitempty"`
	MsInUTC     int           `xml:"msInUTC,omitempty"`
	Udata       string        `xml:"udata,omitempty"`

	qbrequest // Fields required for every request
}

// Response to an API_EditRecord request.
// See http://goo.gl/3OgE49 for more details.
type EditRecordResponse struct {
	XMLName          xml.Name `xml:"qdbapi"`
	RecordId         int      `xml:"rid"`
	UpdateId         string   `xml:"update_id"`
	NumFieldsChanged int      `xml:"num_fields_changed"`

	qbresponse // Fields returned in every response
}

// UpdateField updates a record field.
func (r *EditRecordRequest) UpdateField(id int, value string) {
	r.Fields = append(r.Fields, recordField{Id: id, Value: value})
}

// EditRecord edits record to a QuickBase database.
func (qb *QuickBase) EditRecord(dbid string, req *EditRecordRequest) (*EditRecordResponse, error) {
	req.Ticket = qb.ticket
	req.AppToken = qb.apptoken

	// Only one of the these types in the request
	if req.RecordId > 0 {
		req.Key = ""
	}

	resp := new(EditRecordResponse)
	if err := qb.query("API_EditRecord", dbid, req, resp); err != nil {
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
