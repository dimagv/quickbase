// // Copyright 2014 James Massara. All rights reserved.
// // Use of this source code is governed by a BSD-style
// // license that can be found in the LICENSE file.

package quickbase

import (
	"encoding/xml"
)

// AddRecord adds a new record to the given dbid QuickBase. It returns the
// record id of the newly created record.
func (c *Conn) AddRecord(dbid string, record *Record) (int, error) {
	req := &addRecordRequest{}
	req.Ticket = c.ticket
	req.AppToken = c.apptoken

	for id, f := range record.fields {
		req.Fields = append(req.Fields, recordField{Id: id, Value: f.value})
	}

	rsp := &addRecordResponse{}
	if err := c.do("API_AddRecord", dbid, req, rsp); err != nil {
		return -1, err
	}

	if rsp.ErrorCode != 0 {
		return -1, &QBError{msg: rsp.ErrorText, Code: rsp.ErrorCode, Detail: rsp.ErrorDetail}
	}

	return rsp.RecordId, nil
}

// addRecordRequest is the XML structure for the API_AddRecord call.
type addRecordRequest struct {
	XMLName xml.Name      `xml:"qdbapi"`
	Fields  []recordField `xml:"field"`

	qbRequest // XML fields required for every API call
}

// addRecordResponse is the XML returned from an API_AddRecord call.
type addRecordResponse struct {
	XMLName  xml.Name `xml:"qdbapi"`
	RecordId int      `xml:"rid"`
	UpdateId int64    `xml:"update_id"`

	qbResponse // XML fields returned in every API call
}
