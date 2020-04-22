package main

import (
    "flag"
    "fmt"
    logger "github.com/sirupsen/logrus"
    "net/http"
)
var (
    path = flag.String("d", "./", "path")
    port = flag.Int("p", 8080, "port")
)
func main()  {
    flag.Parse()
    fmt.Printf("run %v serve: %s", *port, *path)
    http.Handle("/", http.FileServer(http.Dir(*path)))
    err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
    if err != nil {
        logger.Error(err)
    }

}
