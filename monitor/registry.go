package monitor

import (
    "encoding/json"
    "fmt"
    "github.com/samuel/go-zookeeper/zk"
    path2 "path"
    "strings"
    "time"
)

var zkClient *zk.Conn

func Register(connectString, nodePath string, appInfo *AppInfo) {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("catch panic: %v\n", r)
        }
    }()

    zkServers := strings.Split(connectString, `,`)
    zkClient, _, _ = zk.Connect(zkServers, time.Second)

    flags := int32(0)
    path := ``
    for _, value := range strings.Split(nodePath, "/")[1:] {
        path = path2.Join(path, `/`, value)
        if strings.Compare(path, nodePath) == 0 {
            flags = zk.FlagEphemeral | zk.FlagSequence
        }

        bytes, _ := json.Marshal(appInfo)
        _, err := zkClient.Create(path, bytes, flags, zk.WorldACL(zk.PermAll))
        if err != nil && err != zk.ErrNodeExists {
            panic("register node occur error")
        }
    }
}

type AppInfo struct {
    Ip           string `json:"ip"`
    Name         string `json:"name"`
    DataCenterId int8   `json:"dataCenterId"`
    WorkerId     int8   `json:"workerId"`
    Pid          int    `json:"pid"`
    Port         int16  `json:"port"`
    Profile      string `json:"profile"`
    StartTimeMS  int64  `json:"startTimeMS"`
}
