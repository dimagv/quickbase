quickbase
=========
Partial implementation of the QuickBase API (http://www.quickbase.com/api-guide/index.html)

Example
=======
```go
package main

import (
    "fmt"
    "github.com/jmassara/quickbase"
    "log"
)

const targetdomain = "somecorp.quickbase.com"

func main() {
    qb := quickbase.New(targetdomain)

    auth, err := qb.Authenticate(&quickbase.Authenticate{
        Username: "PTBarnum",
        Password: "TopSecret",
        Hours:    1,
    })

    if err != nil {
        log.Fatalf("Failed to authenticate to QuickBase (%s): %s\n", targetdomain, err)
    }

    query, err := qb.DoQuery("bddfa5nbx", &quickbase.DoQuery{
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
        log.Fatalf("Failed to query QuickBase (%s): %s\n", targetdomain, err)
    }

    for _, record := range query.GetRecords() {
        fmt.Printf("Business Phone Number: %s\n", record["Business Phone Number"])
        fmt.Printf("                Email: %s\n", record["Email"])
    }
}
```
