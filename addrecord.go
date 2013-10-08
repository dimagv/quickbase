// Copyright 2013 James Massara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package quickbase

import (
	"encoding/xml"
	"fmt"
)

// API_AddRecord request parameters.
// See http://goo.gl/SELQDd for more details.
type AddRecordRequest struct {
	XMLName     xml.Name         `xml:"qdbapi"`
	Fields      []addRecordField `xml:"field"`
	DispRec     int              `xml:"disprec,omitempty"`
	FForm       int              `xml:"fform,omitempty"`
	IgnoreError int              `xml:"ignoreError,omitempty"`
	Ticket      string           `xml:"ticket"`
	AppToken    string           `xml:"apptoken,omitempty"`
	Udata       string           `xml:"udata,omitempty"`
	MsInUTC     int              `xml:"msInUTC,omitempty"`
}

// Response to an API_AddRecord request.
// See http://goo.gl/SELQDd for more details.
type AddRecordResponse struct {
	XMLName     xml.Name `xml:"qdbapi"`
	Action      string   `xml:"action"`
	ErrorCode   int      `xml:"errcode"`
	ErrorText   string   `xml:"errtext"`
	ErrorDetail string   `xml:"errdetail"`
	Udata       string   `xml:"udata"`
	RecordId    string   `xml:"rid"`
	UpdateId    string   `xml:"update_id"`
}

// addRecordField represents a field when adding a record in QuickBase.
// The XML structure for fields (and may other things) are inconsistent across
// the Quickbase API.
type addRecordField struct {
	Id    int    `xml:"fid,attr"`
	Value string `xml:",chardata"`
}

// AddField adds a field to the record.
func (r *AddRecordRequest) AddField(id int, value string) {
	r.Fields = append(r.Fields, addRecordField{Id: id, Value: value})
}

// AddRecord adds a new record to a QuickBase database (dbid).
func (qb *QuickBase) AddRecord(dbid string, rec *AddRecordRequest) (*AddRecordResponse, *QBError) {
	params := makeParams("API_AddRecord")
	params["url"] = fmt.Sprintf("https://%s/db/%s", qb.domain, dbid)

	resp := new(AddRecordResponse)
	if err := qb.query(params, rec, resp); err != nil {
		return nil, &QBError{msg: err.Error()}
	}

	if resp.ErrorCode != 0 {
		return nil, &QBError{msg: resp.ErrorText, Detail: resp.ErrorDetail}
	}

	return resp, nil
}
