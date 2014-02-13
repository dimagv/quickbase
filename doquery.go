// Copyright 2014 James Massara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package quickbase

import "encoding/xml"

// DoQueryOptions represents the options that can be passed to API_DoQuery.
//
// See http://www.quickbase.com/api-guide/do_query.html#Request_Parameters for
// more information.
type DoQueryOptions struct {
	// Clist is a period delimited list of field ids to be returned.
	Clist string

	// Slist is a period delimited list of field ids used to determine sorting.
	Slist string

	// Options is a period delimited list of return options for the query.
	// (e.g. "num-4.sortorder-A.skp-10.onlynew")
	Options string

	// ReturnPercentage specifies whether numeric percent values in the query
	// result will be percentage format (e.g. 10% is shown as 10) or decimal
	// format (e.g. 10% is shown as .1). Set this to true for percentage format.
	ReturnPercentage bool

	// These options are mutually exclusive and are set by their associated
	// DoQuery* functions.
	query string // DoQuery
	qName string // DoQueryByName
	qId   int    // DoQueryById
}

// DoQuery queries the given dbid using the QuickBase raw query syntax. If the
// query parameter is an empty string, all records from the dbid will be
// returned.
//
// See http://www.quickbase.com/api-guide/do_query.html#queryOperators for more
// information on the query syntax.
func (c *Conn) DoQuery(dbid, query string, opts *DoQueryOptions) ([]Record, error) {
	req := &doQueryRequest{Query: query}
	setQueryOptions(req, opts)
	return c.doQuery(dbid, req)
}

// DoQueryByName queries the given dbid using a stored query in the database
// with the given query name.
func (c *Conn) DoQueryByName(dbid, qname string, opts *DoQueryOptions) ([]Record, error) {
	req := &doQueryRequest{Qname: qname}
	setQueryOptions(req, opts)
	return c.doQuery(dbid, req)
}

// DoQueryById queries the given dbid using a stored query in the database with
// the given query id.
func (c *Conn) DoQueryById(dbid string, qid int, opts *DoQueryOptions) ([]Record, error) {
	req := &doQueryRequest{Qid: qid}
	setQueryOptions(req, opts)
	return c.doQuery(dbid, req)
}

func (c *Conn) doQuery(dbid string, req *doQueryRequest) ([]Record, error) {
	req.Ticket = c.ticket
	req.AppToken = c.apptoken
	req.Fmt = "structured"
	req.IncludeRids = 1

	rsp := &doQueryResponse{}
	if err := c.do("API_DoQuery", dbid, req, rsp); err != nil {
		return nil, err
	}

	if rsp.ErrorCode != 0 {
		return nil, &QBError{msg: rsp.ErrorText, Code: rsp.ErrorCode, Detail: rsp.ErrorDetail}
	}

	// Map record field ids to its label name
	labels := make(map[int]string)
	for _, field := range rsp.FieldLabels {
		labels[field.Id] = field.Label
	}

	records := make([]Record, len(rsp.Records))
	for i, r := range rsp.Records {
		fields := make(map[int]field)
		for _, f := range r.Fields {
			fields[f.Id] = field{label: labels[f.Id], value: f.Value}
		}
		records[i].Id = r.Rid
		records[i].UpdateId = r.UpdateId
		records[i].fields = fields
	}

	return records, nil
}

func setQueryOptions(req *doQueryRequest, opts *DoQueryOptions) {
	if opts == nil {
		return
	}

	if opts.Clist != "" {
		req.Clist = opts.Clist
	}

	if opts.Slist != "" {
		req.Slist = opts.Slist
	}

	if opts.ReturnPercentage {
		req.ReturnPercentage = 1
	}

	if opts.Options != "" {
		req.Options = opts.Options
	}
}

// doQueryRequest is the XML structure for the API_DoQuery call.
type doQueryRequest struct {
	XMLName          xml.Name `xml:"qdbapi"`
	Query            string   `xml:"query,omitempty"`
	Qid              int      `xml:"qid,omitempty"`
	Qname            string   `xml:"qname,omitempty"`
	Clist            string   `xml:"clist,omitempty"`
	Slist            string   `xml:"slist,omitempty"`
	ReturnPercentage int      `xml:"returnpercentage,omitempty"`
	Options          string   `xml:"options,omitempty"`
	Fmt              string   `xml:"fmt"`
	IncludeRids      int      `xml:"includeRids"`

	qbRequest // XML fields required for every API call
}

// doQueryResponse is the XML returned from an API_DoQuery call.
type doQueryResponse struct {
	XMLName xml.Name `xml:"qdbapi"`

	FieldLabels []struct {
		Id    int    `xml:"id,attr"`
		Label string `xml:"label"`
	} `xml:"table>fields>field"`

	Records []struct {
		Rid      int   `xml:"rid,attr"`
		UpdateId int64 `xml:"update_id"`
		Fields   []struct {
			Id    int    `xml:"id,attr"`
			Value string `xml:",chardata"`
		} `xml:"f"`
	} `xml:"table>records>record"`

	qbResponse // XML fields returned in every API call
}
