quickbase (beta)
=================
Partial implementation of the QuickBase API (http://www.quickbase.com/api-guide/index.html).
API is in flux and not considered stable. Pull Requests are welcome!

Example
=======
```go
package main

import (
    "log"

    "github.com/jmassara/quickbase"
)

const (
    // QuickBase field ids
    BusinessPhoneNumber = 6
    Email               = 7
)

func main() {
    conn, err := quickbase.Login("https://somecorp.quickbase.com", "PTBarnum", "TopSecret")
    if err != nil {
        log.Fatalf("Failed to authenticate to QuickBase: %s", err)
    }

    conn.SetAppToken("dtmd897bfsw85bb6bneceb6wnze3")

    opts := &quickbase.DoQueryOptions{
        Clist:   "5.6.7.22.3",
        Slist:   "3",
        Options: "num-4.sortorder-A.skp-10.onlynew",
    }

    records, err := conn.DoQuery("bddfa5nbx", "{'5'.CT.'Ragnar Lodbrok'}AND{'5'.CT.'Acquisitions'}", opts)
    if err != nil {
        log.Fatalf("Failed to query QuickBase: %s\n", err)
    }

    for _, record := range records {
        fields := record.GetFieldsById()
        log.Printf("Business Phone Number: %s\n", fields[BusinessPhoneNumber])
        log.Printf("                Email: %s\n", fields[Email])
    }
}
```
