package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
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

var dataCenterId int64 = 1
var workerId int64 = 1

var sequence int64 = 0
var lastTimestamp int64 = -1

var lck sync.Mutex

func main() {
    preCheck()

    http.HandleFunc("/id/next", func(writer http.ResponseWriter, request *http.Request) {
        io.WriteString(writer, strconv.FormatInt(nextId(), 10))
    })

    addr := fmt.Sprintf("%s%d", ":", PORT)
    log.Printf("IDGenerator startup, port%s\n", addr)
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
