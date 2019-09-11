package id

import (
    "errors"
    "fmt"
    "sync"
    "time"
)

type Generator struct {
    StartTimestamp int64

    DataCenterBits uint16
    WorkerBits     uint16
    SequenceBits   uint16
    DataCenterId   int64
    WorkerId       int64

    maxSequence      int64
    maxWorkerNum     int64
    maxDataCenterNum int64
    workerShift      uint16
    dataCenterShift  uint16
    timestampShift   uint16

    sequence      int64
    lastTimestamp int64

    lck sync.Mutex
}

func NewGenerator(startTimestamp int64, dataCenterBits uint16, workerBits uint16, sequenceBits uint16, dataCenterId int64, workerId int64) (*Generator, error) {
    gen := &Generator{
        StartTimestamp: startTimestamp,
        DataCenterBits: dataCenterBits,
        WorkerBits:     workerBits,
        SequenceBits:   sequenceBits,
        DataCenterId:   dataCenterId,
        WorkerId:       workerId,
    }

    gen.maxSequence = -1 ^ (-1 << gen.SequenceBits)
    gen.maxWorkerNum = -1 ^ (-1 << gen.WorkerBits)
    gen.maxDataCenterNum = -1 ^ (-1 << gen.DataCenterBits)
    gen.workerShift = gen.SequenceBits
    gen.dataCenterShift = gen.WorkerBits + gen.SequenceBits
    gen.timestampShift = gen.DataCenterBits + gen.WorkerBits + gen.SequenceBits

    if gen.DataCenterId > gen.maxDataCenterNum || gen.DataCenterId < 0 {
        panic(fmt.Sprintf("DataCenterId can't be greater than %d or less than 0, actual: %d", gen.maxDataCenterNum, gen.DataCenterId))
    }

    if gen.WorkerId > gen.maxWorkerNum || gen.WorkerId < 0 {
        panic(fmt.Sprintf("WorkerId can't be greater than %d or less than 0, actual: %d", gen.maxWorkerNum, gen.WorkerId))
    }

    return gen, nil
}

func (gen *Generator) NextId() (int64, error) {
    currentTimestamp := currentTimeMillis();
    if currentTimestamp < gen.lastTimestamp {
        return -1, errors.New("Clock moved backwards.  Refusing to generate id")
    }

    gen.lck.Lock()
    if currentTimestamp == gen.lastTimestamp {
        gen.sequence = (gen.sequence + 1) & gen.maxSequence
        //同一毫秒的序列数已经达到最大
        if gen.sequence == 0 {
            currentTimestamp = gen.getNextTimestamp()
        }
    } else {
        gen.sequence = 0
    }

    gen.lastTimestamp = currentTimestamp
    gen.lck.Unlock()

    return (currentTimestamp-gen.StartTimestamp)<<gen.timestampShift | gen.DataCenterId<<gen.dataCenterShift | gen.WorkerId<<gen.workerShift | gen.sequence, nil
}

func (gen *Generator) getNextTimestamp() int64 {
    currentTimestamp := currentTimeMillis()
    for currentTimestamp <= gen.lastTimestamp {
        currentTimestamp = currentTimeMillis()
    }
    return currentTimestamp
}

func currentTimeMillis() int64 {
    return time.Now().UnixNano() / 1e6
}
