// Copyright 2013 James Massara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package quickbase

import (
	"encoding/xml"
	"fmt"
)

// API_DoQuery request parameters.
// See http://goo.gl/vHzW5K for more details.
// Note: The following fields have hardcoded values for XML parsing
//  * `Fmt` is set to `structured`
//  * `IncludeRids` is set to `1`
type DoQueryRequest struct {
	XMLName          xml.Name `xml:"qdbapi"`
	Query            string   `xml:"query,omitempty"`
	Qid              string   `xml:"qid,omitempty"`
	Qname            string   `xml:"qname,omitempty"`
	Clist            string   `xml:"clist,omitempty"`
	Slist            string   `xml:"slist,omitempty"`
	ReturnPercentage int      `xml:"returnpercentage,omitempty"`
	Options          string   `xml:"options,omitempty"`
	IncludeRids      int      `xml:"includeRids,omitempty"`
	Ticket           string   `xml:"ticket"`
	AppToken         string   `xml:"apptoken,omitempty"`
	Udata            string   `xml:"udata,omitempty"`
	Fmt              string   `xml:"fmt"`
}

// Response to an API_DoQuery request.
// See http://goo.gl/vHzW5K for more details.
type DoQueryResponse struct {
	XMLName     xml.Name `xml:"qdbapi"`
	Action      string   `xml:"action"`
	ErrorCode   int      `xml:"errcode"`
	ErrorText   string   `xml:"errtext"`
	ErrorDetail string   `xml:"errdetail"`
	Ticket      string   `xml:"ticket"`
	UserId      string   `xml:"userid"`
	Udata       string   `xml:"udata"`

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

	// Private fields
	labels map[int]string `xml:"-"`
}

// DoQuery queries a QuickBase database (dbid).
func (qb *QuickBase) DoQuery(dbid string, q *DoQueryRequest) (*DoQueryResponse, *QBError) {
	params := makeParams("API_DoQuery")
	params["url"] = fmt.Sprintf("https://%s/db/%s", qb.domain, dbid)

	// Set defaults
	q.Fmt = "structured"
	q.IncludeRids = 1

	// Only pass one of the query types in the request
	if q.Query != "" {
		q.Qid = ""
		q.Qname = ""
	} else if q.Qid != "" {
		q.Qname = ""
	}

	resp := new(DoQueryResponse)
	if err := qb.query(params, q, resp); err != nil {
		return nil, &QBError{msg: err.Error()}
	}

	if resp.ErrorCode != 0 {
		return nil, &QBError{msg: resp.ErrorText, Code: resp.ErrorCode, Detail: resp.ErrorDetail}
	}

	// Map of record field id to its label name
	resp.labels = make(map[int]string)
	for _, field := range resp.FieldLabels {
		resp.labels[field.Id] = field.Label
	}

	return resp, nil
}

// Record represents a Quickbase record result
type Record struct {
	Id       int
	UpdateId string
	// Each field of the record where the map key is the field id.
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
