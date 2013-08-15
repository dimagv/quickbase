// Copyright 2013 James Massara. All rights reserved.

package quickbase

import (
	"encoding/xml"
	"fmt"
)

// API_DoQuery request parameters.
// See http://goo.gl/vHzW5K for more details.
// Note: The `Fmt' field is hardcoded to `structured'.
type DoQuery struct {
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

	Field []struct {
		Id    int    `xml:"id,attr"`
		Label string `xml:"label"`
	} `xml:"table>fields>field"`

	Record []struct {
		UpdateId string `xml:"update_id"`
		Data     []struct {
			Id    int    `xml:"id,attr"`
			Value string `xml:",chardata"`
		} `xml:"f"`
	} `xml:"table>records>record"`
}

// DoQuery queries a QuickBase database (dbid).
func (qb *QuickBase) DoQuery(dbid string, q *DoQuery) (*DoQueryResponse, *QuickBaseError) {
	params := makeParams("API_DoQuery")
	params["url"] = fmt.Sprintf("https://%s/db/%s", qb.Domain, dbid)

	// Hardcoded to `structured' for XML decoding
	q.Fmt = "structured"

	// Only pass one of the query types in the request
	if q.Query != "" {
		q.Qid = ""
		q.Qname = ""
	} else if q.Qid != "" {
		q.Qname = ""
	}

	resp := &DoQueryResponse{}
	if err := qb.query(params, q, resp); err != nil {
		return nil, &QuickBaseError{msg: err.Error()}
	}

	if resp.ErrorCode != 0 {
		return nil, &QuickBaseError{msg: resp.ErrorText, Detail: resp.ErrorDetail}
	}

	return resp, nil
}

// GetRecords returns a map of records from the DoQuery request.
func (q *DoQueryResponse) GetRecords() map[int]map[string]string {
	fieldMap := make(map[int]string)
	records := make(map[int]map[string]string)

	for _, field := range q.Field {
		fieldMap[field.Id] = field.Label
	}

	for num, record := range q.Record {
		recordData := make(map[string]string)
		for _, data := range record.Data {
			recordData[fieldMap[data.Id]] = data.Value
		}
		records[num] = recordData
	}

	return records
}
