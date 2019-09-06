package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "github.com/juju/errors"
    "github.com/xuezhongde/id-generator/id"
    "log"
    "net/http"
    "os"
)

const (
    DefaultAppName = "id-generator"
    DefaultProfile = "dev"
    DefaultPort    = 8000
    DefaultRouter  = "/id"
    ProbeRouter    = "/probe"

    DefaultStartTimestamp int64  = 1563764872049
    DefaultDataCenterBits uint16 = 5
    DefaultWorkerBits     uint16 = 5
    DefaultSequenceBits   uint16 = 12
)

var (
    h    bool
    v    bool
    c    string
    d, w int64
    p    int
    r    string
)

func init() {
    flag.BoolVar(&h, "h", false, "this help")
    flag.BoolVar(&v, "v", false, "show version and exit")
    flag.StringVar(&c, "c", "./etc/id.toml", "id generator config file")
    flag.Int64Var(&d, "d", -1, "data center id")
    flag.Int64Var(&w, "w", -1, "worker id")
    flag.IntVar(&p, "p", -1, "listen on port")
    flag.StringVar(&r, "r", "", "router path")

    flag.Usage = usage
}

func main() {
    flag.Parse()

    if h {
        flag.Usage()
        return
    }

    if v {
        fmt.Println("1.0.0")
        return
    }

    cfg, err := id.LoadConfig(c)
    if err != nil {
        println(errors.ErrorStack(err))
        return
    }

    if p >= 0 {
        cfg.Port = p
    }

    if len(r) > 0 {
        cfg.Router = r
    }

    if d >= 0 {
        cfg.DateCenterId = d
    }

    if w >= 0 {
        cfg.WorkerId = w
    }

    if cfg.Port <= 0 {
        cfg.Port = DefaultPort
    }

    if len(cfg.Router) <= 0 {
        cfg.Router = DefaultRouter
    }

    if len(cfg.AppName) <= 0 {
        cfg.AppName = DefaultAppName
    }

    if len(cfg.Profile) <= 0 {
        cfg.Profile = DefaultProfile
    }

    port := cfg.Port
    router := cfg.Router
    gen, _ := id.NewGenerator(DefaultStartTimestamp, DefaultDataCenterBits, DefaultWorkerBits, DefaultSequenceBits, cfg.DateCenterId, cfg.WorkerId)

    http.HandleFunc(router, func(writer http.ResponseWriter, request *http.Request) {
        _id, _ := gen.NextId()
        apiRsp := &ApiResponse{0, "success", _id}
        writeJsonResponse(&writer, apiRsp)
    })

    http.HandleFunc(ProbeRouter, func(writer http.ResponseWriter, request *http.Request) {
        probeResult := &ProbeResult{0, "success", cfg.AppName, cfg.Profile}
        writeJsonResponse(&writer, probeResult)
    })

    addr := fmt.Sprintf("%s%d", ":", port)
    log.Printf("IDGenerator startup, port%s, router: %s, dataCenterId: %d, workerId: %d\n", addr, router, gen.DataCenterId, gen.WorkerId)
    log.Fatal(http.ListenAndServe(addr, nil))
}

func writeJsonResponse(writer *http.ResponseWriter, v interface{}) {
    jsonBytes, _ := json.Marshal(v)
    (*writer).Header().Set("Content-Type", "application/json")
    _, _ = (*writer).Write(jsonBytes)
}

func usage() {
    _, _ = fmt.Fprintf(os.Stdout, "Usage: id-gen [-hv] [-c config file] [-d data center id] [-w worker id] [-p port] [-r router]\nOptions:\n")
    flag.PrintDefaults()
}

type ApiResponse struct {
    Code int16  `json:"code"`
    Msg  string `json:"msg"`
    Data int64  `json:"data"`
}

type ProbeResult struct {
    Code    int16  `json:"code"`
    Msg     string `json:"msg"`
    AppName string `json:"appName"`
    Profile string `json:"profile"`
}
