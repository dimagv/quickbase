// Copyright 2013 James Massara. All rights reserved.

package quickbase

import (
	"encoding/xml"
	"fmt"
)

// API_AddRecord request parameters.
// See http://goo.gl/SELQDd for more details.
type AddRecord struct {
	XMLName     xml.Name `xml:"qdbapi"`
	Fields      []Field  `xml:"field"`
	DispRec     int      `xml:"disprec,omitempty"`
	FForm       int      `xml:"fform,omitempty"`
	IgnoreError int      `xml:"ignoreError,omitempty"`
	Ticket      string   `xml:"ticket"`
	AppToken    string   `xml:"apptoken,omitempty"`
	Udata       string   `xml:"udata,omitempty"`
	MsInUTC     int      `xml:"msInUTC,omitempty"`
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

// AddField adds a field to the record.
func (r *AddRecord) AddField(id int, value string) {
	r.Fields = append(r.Fields, Field{Id: id, Value: value})
}

// AddRecord adds a new record to a QuickBase database (dbid).
func (qb *QuickBase) AddRecord(dbid string, rec *AddRecord) (*AddRecordResponse, *QuickBaseError) {
	params := makeParams("API_AddRecord")
	params["url"] = fmt.Sprintf("https://%s/db/%s", qb.Domain, dbid)

	resp := new(AddRecordResponse)
	if err := qb.query(params, rec, resp); err != nil {
		return nil, &QuickBaseError{msg: err.Error()}
	}

	if resp.ErrorCode != 0 {
		return nil, &QuickBaseError{msg: resp.ErrorText, Detail: resp.ErrorDetail}
	}

	return resp, nil
}
