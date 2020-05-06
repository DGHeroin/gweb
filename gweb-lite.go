package main

import (
    "crypto/md5"
    "flag"
    "fmt"
    logger "github.com/sirupsen/logrus"
    "net/http"

    "strings"
    "time"
)
var (
    servePath = flag.String("d", "./", "path")
    servePort = flag.Int("p", 8080, "port")
)

func checkEtag(w http.ResponseWriter, r *http.Request) (done bool) {
    etag := w.Header().Get("Etag")

    if inm := r.Header.Get("If-None-Match"); inm != "" {
        if etag == "" {
            return false
        }
        if r.Method != "GET" && r.Method != "HEAD" {
            return false
        }
        if inm == etag || inm == "*" {
            h := w.Header()
            delete(h, "Content-Type")
            delete(h, "Content-Length")
            w.WriteHeader(http.StatusNotModified)
            return true
        }
    }
    return false
}

func checkLastModified(w http.ResponseWriter, r *http.Request, modtime time.Time) bool {
    if modtime.IsZero() {
        return false
    }

    if t, err := time.Parse(http.TimeFormat, r.Header.Get("If-Modified-Since")); err == nil && modtime.Before(t.Add(1*time.Second)) {
        h := w.Header()
        delete(h, "Content-Type")
        delete(h, "Content-Length")
        w.WriteHeader(http.StatusNotModified)
        return true
    }
    w.Header().Set("Last-Modified", modtime.UTC().Format(http.TimeFormat))
    return false
}

type FileServer struct {
    dir string
    fileSys http.FileSystem
    handler http.Handler
}

func (fs *FileServer) init()  {
    fs.fileSys = http.Dir(fs.dir)
    fs.handler = http.FileServer(fs.fileSys)
}
func (fs FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // check etag
    upath := r.URL.Path
    if !strings.HasPrefix(upath, "/") {
        upath = "/" + upath
        r.URL.Path = upath
    }

    if f, err := fs.fileSys.Open(upath); err != nil {
        w.WriteHeader(http.StatusNotFound)
        return
    } else {
        defer f.Close()
        if d, err := f.Stat(); err != nil {
            w.WriteHeader(http.StatusNotFound)
            return
        } else {
            // 根据文件名和上次修改时间 生成 etag
            m := md5.New()
            m.Write([]byte(fmt.Sprintf("%v%v", d.Name(), d.ModTime())))
            etag := fmt.Sprintf("%x", m.Sum(nil))
            w.Header().Set("Etag", etag)
            if checkEtag(w, r) {
                return
            }
        }
    }

    fs.handler.ServeHTTP(w, r)
}

func main()  {
    flag.Parse()
    fmt.Printf("run %v serve: %s", *servePort, *servePath)
    //http.Handle("/", http.FileServer(http.Dir(*path)))
    server := &FileServer{dir:*servePath}
    server.init()
    err := http.ListenAndServe(fmt.Sprintf(":%d", *servePort), server)
    if err != nil {
        logger.Error(err)
    }

}
