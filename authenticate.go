package quickbase

import (
	"encoding/xml"
	"fmt"
)

// API_Authenticate request parameters.
// See http://goo.gl/eQSiZy for more details.
type Authenticate struct {
	XMLName  xml.Name `xml:"qdbapi"`
	Username string   `xml:"username"`
	Password string   `xml:"password"`
	Hours    int      `xml:"hours,omitempty"`
	Udata    string   `xml:"udata,omitempty"`
}

// Response to an API_Authenticate request.
// See http://goo.gl/eQSiZy for more details.
type AuthenticateResponse struct {
	XMLName     xml.Name `xml:"qdbapi"`
	Action      string   `xml:"action"`
	ErrorCode   int      `xml:"errcode"`
	ErrorText   string   `xml:"errtext"`
	ErrorDetail string   `xml:"errdetail"`
	Ticket      string   `xml:"ticket"`
	UserId      string   `xml:"userid"`
	Udata       string   `xml:"udata"`
}

func (qb *QuickBase) Authenticate(auth *Authenticate) (*AuthenticateResponse, *QuickBaseError) {
	params := makeParams("API_Authenticate")
	params["url"] = fmt.Sprintf("https://%s/db/main", qb.Domain)

	resp := &AuthenticateResponse{}
	if err := qb.query(params, auth, resp); err != nil {
		return nil, &QuickBaseError{msg: err.Error()}
	}

	if resp.ErrorCode != 0 {
		return nil, &QuickBaseError{msg: resp.ErrorText, Detail: resp.ErrorDetail}
	}

	return resp, nil
}
