// // Copyright 2014 James Massara. All rights reserved.
// // Use of this source code is governed by a BSD-style
// // license that can be found in the LICENSE file.

package quickbase

import (
	"encoding/xml"
)

// EditRecord edits record in the given dbid QuickBase. It returns the number of
// fields changed.
func (c *Conn) EditRecord(dbid string, record *Record) (int, error) {
	req := &editRecordRequest{RecordId: record.Id}
	req.Ticket = c.ticket
	req.AppToken = c.apptoken

	for id, f := range record.fields {
		req.Fields = append(req.Fields, recordField{Id: id, Value: f.value})
	}

	rsp := &editRecordResponse{}
	if err := c.do("API_EditRecord", dbid, req, rsp); err != nil {
		return -1, err
	}

	if rsp.ErrorCode != 0 {
		return -1, &QBError{msg: rsp.ErrorText, Code: rsp.ErrorCode, Detail: rsp.ErrorDetail}
	}

	return rsp.NumFieldsChanged, nil
}

// editRecordRequest is the XML structure for the API_EditRecord call.
type editRecordRequest struct {
	XMLName  xml.Name      `xml:"qdbapi"`
	RecordId int           `xml:"rid,omitempty"`
	Fields   []recordField `xml:"field"`

	qbRequest // XML fields required for every API call
}

// editRecordResponse is the XML returned from an API_EditRecord call.
type editRecordResponse struct {
	XMLName          xml.Name `xml:"qdbapi"`
	RecordId         int      `xml:"rid"`
	UpdateId         int64    `xml:"update_id"`
	NumFieldsChanged int      `xml:"num_fields_changed"`

	qbResponse // XML fields returned in every API call
}
