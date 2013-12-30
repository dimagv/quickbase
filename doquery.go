// Copyright 2013 James Massara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package quickbase

import (
	"encoding/xml"
)

// API_DoQuery request parameters.
// See http://goo.gl/vHzW5K for more details.
type DoQueryRequest struct {
	XMLName          xml.Name `xml:"qdbapi"`
	Query            string   `xml:"query,omitempty"`
	Qid              int      `xml:"qid,omitempty"`
	Qname            string   `xml:"qname,omitempty"`
	Clist            string   `xml:"clist,omitempty"`
	Slist            string   `xml:"slist,omitempty"`
	ReturnPercentage int      `xml:"returnpercentage,omitempty"`
	Options          string   `xml:"options,omitempty"`
	Udata            string   `xml:"udata,omitempty"`

	qbrequest      // Fields required for every request
	doQueryRequest // Fields defined in doQueryRequest have hardcoded values
}

type doQueryRequest struct {
	Fmt         string `xml:"fmt"`         // defaults to "structured"
	IncludeRids int    `xml:"includeRids"` // defaults to 1
}

// Response to an API_DoQuery request.
// See http://goo.gl/vHzW5K for more details.
type DoQueryResponse struct {
	XMLName xml.Name `xml:"qdbapi"`
	Udata   string   `xml:"udata"`

	Records []struct {
		Rid      int    `xml:"rid,attr"`
		UpdateId string `xml:"update_id"`
		Fields   []struct {
			Id    int    `xml:"id,attr"`
			Value string `xml:",chardata"`
		} `xml:"f"`
	} `xml:"table>records>record"`

	FieldLabels []struct {
		Id    int    `xml:"id,attr"`
		Label string `xml:"label"`
	} `xml:"table>fields>field"`

	// Field label keyed by its id
	labels map[int]string `xml:"-"`

	qbresponse // Fields returned in every response
}

// DoQuery queries a QuickBase database.
func (qb *QuickBase) DoQuery(dbid string, req *DoQueryRequest) (*DoQueryResponse, error) {
	req.Ticket = qb.ticket
	req.AppToken = qb.apptoken
	req.Fmt = "structured"
	req.IncludeRids = 1

	// Only pass one of the query types in the request
	if req.Query != "" {
		req.Qid = 0
		req.Qname = ""
	} else if req.Qid > 0 {
		req.Qname = ""
	}

	resp := new(DoQueryResponse)
	if err := qb.query("API_DoQuery", dbid, req, resp); err != nil {
		return nil, err
	}

	if resp.ErrorCode != 0 {
		return nil, &QBError{
			msg:    resp.ErrorText,
			Code:   resp.ErrorCode,
			Detail: resp.ErrorDetail,
		}
	}

	// Map of record field id to its label name
	resp.labels = make(map[int]string)
	for _, field := range resp.FieldLabels {
		resp.labels[field.Id] = field.Label
	}

	return resp, nil
}

// Record is a Quickbase record result
type Record struct {
	Id       int
	UpdateId string
	// Each field in the record with the map key set to the field id.
	Fields map[int]struct{ Label, Value string }
}

func (r *DoQueryResponse) GetRecords() []Record {
	records := make([]Record, len(r.Records))
	for i, record := range r.Records {
		fields := make(map[int]struct{ Label, Value string })
		for _, field := range record.Fields {
			fields[field.Id] = struct{ Label, Value string }{
				Label: r.labels[field.Id],
				Value: field.Value,
			}
		}
		records[i].Id = record.Rid
		records[i].UpdateId = record.UpdateId
		records[i].Fields = fields
	}
	return records
}
