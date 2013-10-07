quickbase (alpha)
=================
Partial implementation of the QuickBase API (http://www.quickbase.com/api-guide/index.html).
API is in flux and not considered stable.

Example
=======
```go
package main

import (
    "github.com/jmassara/quickbase"
    "log"
)

const (
    const qbdomain = "somecorp.quickbase.com"
    // QuickBase field ids
    const BusinessPhoneNumber = 6
    const Email               = 7
)

func main() {
    qb := quickbase.New(qbdomain)
    auth, err := qb.Authenticate(&quickbase.AuthRequest{
        Username: "PTBarnum",
        Password: "TopSecret",
        Hours:    1,
    })
    if err != nil {
        log.Fatalf("Failed to authenticate to QuickBase (%s): %s\n", qbdomain, err)
    }

    query, err := qb.DoQuery("bddfa5nbx", &quickbase.DoQueryRequest{
        Ticket:      auth.Ticket,
        AppToken:    "dtmd897bfsw85bb6bneceb6wnze3",
        Udata:       "mydata",
        IncludeRids: 1,
        Query:       "{'5'.CT.'Ragnar Lodbrok'}AND{'5'.CT.'Acquisitions'}",
        Clist:       "5.6.7.22.3",
        Slist:       "3",
        Options:     "num-4.sortorder-A.skp-10.onlynew",
    })
    if err != nil {
        log.Fatalf("Failed to query QuickBase (%s): %s\n", qbdomain, err)
    }

    for _, r := range query.GetRecords() {
        log.Printf("Business Phone Number: %s\n", r.Fields[BusinessPhoneNumber].Value)
        log.Printf("                Email: %s\n", r.Fields[Email].Value)
    }
}
```
