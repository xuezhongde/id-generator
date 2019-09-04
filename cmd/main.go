package main

import (
    "flag"
    "fmt"
    "github.com/juju/errors"
    conf "github.com/xuezhongde/id-generator"
    "io"
    "log"
    "net/http"
    "os"
    "strconv"
    "sync"
    "time"
)

const PORT = 8000

const startTimestamp int64 = 1563764872049

const dataCenterBits = 5
const workerBits = 5
const sequenceBits = 12

const maxDataCenterNum = -1 ^ (-1 << dataCenterBits)
const maxWorkerNum = -1 ^ (-1 << workerBits)
const maxSequence = -1 ^ (-1 << sequenceBits)

const workerShift = sequenceBits
const dataCenterShift = sequenceBits + workerBits
const timestampShift = dataCenterShift + dataCenterBits

var port int
var router string
var dataCenterId int64
var workerId int64

var sequence int64 = 0
var lastTimestamp int64 = -1

var lck sync.Mutex

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

    cfg, err := conf.LoadConfig(c)
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
        cfg.Port = PORT
    }

    if len(cfg.Router) <= 0 {
        cfg.Router = "/id"
    }

    port = cfg.Port
    router = cfg.Router
    dataCenterId = cfg.DateCenterId
    workerId = cfg.WorkerId

    preCheck()

    http.HandleFunc(router, func(writer http.ResponseWriter, request *http.Request) {
        io.WriteString(writer, strconv.FormatInt(nextId(), 10))
    })

    addr := fmt.Sprintf("%s%d", ":", port)
    log.Printf("IDGenerator startup, port%s, router: %s, dataCenterId: %d, workerId: %d\n", addr, router, dataCenterId, workerId)
    log.Fatal(http.ListenAndServe(addr, nil))
}

func preCheck() {
    if dataCenterId > maxDataCenterNum || dataCenterId < 0 {
        panic(fmt.Sprintf("dataCenterId can't be greater than %d or less than 0, actual: %d", maxDataCenterNum, dataCenterId))
    }

    if workerId > maxWorkerNum || workerId < 0 {
        panic(fmt.Sprintf("workerId can't be greater than %d or less than 0, actual: %d", maxWorkerNum, workerId))
    }
}

func nextId() int64 {
    currentTimestamp := currentTimeMillis();
    if currentTimestamp < lastTimestamp {
        //TODO error
        log.Fatal("Clock moved backwards.  Refusing to generate id")
    }

    lck.Lock()
    if currentTimestamp == lastTimestamp {
        sequence = (sequence + 1) & maxSequence
        //同一毫秒的序列数已经达到最大
        if sequence == 0 {
            currentTimestamp = getNextTimestamp()
        }
    } else {
        sequence = 0
    }

    lastTimestamp = currentTimestamp
    lck.Unlock()

    return (currentTimestamp-startTimestamp)<<timestampShift | dataCenterId<<dataCenterShift | workerId<<workerShift | sequence
}

func currentTimeMillis() int64 {
    return time.Now().UnixNano() / 1e6
}

func getNextTimestamp() int64 {
    currentTimestamp := currentTimeMillis()
    for currentTimestamp <= lastTimestamp {
        currentTimestamp = currentTimeMillis()
    }
    return currentTimestamp
}

func usage() {
    fmt.Fprintf(os.Stdout, "Usage: id-gen [-hv] [-c config file] [-d data center id] [-w worker id] [-p port] [-r router]\nOptions:\n")
    flag.PrintDefaults()
}
