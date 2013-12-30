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
    // QuickBase field ids
    BusinessPhoneNumber = 6
    Email               = 7
)

func main() {
    qb := quickbase.New("https://somecorp.quickbase.com")
    if err := qb.Authenticate("PTBarnum", "TopSecret", 1); err != nil {
        log.Fatalf("Failed to authenticate to QuickBase: %s", err)
    }

    qb.SetAppToken("dtmd897bfsw85bb6bneceb6wnze3")

    query, err := qb.DoQuery("bddfa5nbx", &quickbase.DoQueryRequest{
        Udata:   "mydata",
        Query:   "{'5'.CT.'Ragnar Lodbrok'}AND{'5'.CT.'Acquisitions'}",
        Clist:   "5.6.7.22.3",
        Slist:   "3",
        Options: "num-4.sortorder-A.skp-10.onlynew",
    })

    if err != nil {
        log.Fatalf("Failed to query QuickBase (%s): %s\n", qbhost, err)
    }

    for _, r := range query.GetRecords() {
        log.Printf("Business Phone Number: %s\n", r.Fields[BusinessPhoneNumber].Value)
        log.Printf("                Email: %s\n", r.Fields[Email].Value)
    }
}
```
