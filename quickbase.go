package quickbase

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"net/http"
)

type QuickBase struct {
	Domain string
}

func New(domain string) *QuickBase {
	return &QuickBase{Domain: domain}
}

func (qb *QuickBase) query(params map[string]string, request, response interface{}) error {
	postdata, err := xml.MarshalIndent(request, "", " ")
	if err != nil {
		return err
	}

	// Convert `&#39;' to `&apos;' because xml.Marshal runs EscapeText which
	// and transforms `'' to `&#39;', fine for HTML but not XML.
	postdata = bytes.Replace(postdata, []byte("&#39;"), []byte("&apos;"), -1)

	req, err := http.NewRequest("POST", params["url"], bytes.NewBuffer(postdata))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Quickbase-Action", params["action"])

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	return xml.NewDecoder(res.Body).Decode(response)
}

type QuickBaseError struct {
	msg    string
	Detail string
}

func (e *QuickBaseError) Error() string { return e.msg }

// Helper functions

// makeParams returns the initial map with the action set.
func makeParams(action string) map[string]string {
	params := make(map[string]string)
	params["action"] = action
	return params
}
