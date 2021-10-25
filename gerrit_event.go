// Copyright 2013 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
    "bytes"
    "fmt"
    "encoding/json"
    "net/http"
    "io/ioutil"

    "github.com/dongwenjuan/gerritssh"
    "gopkg.in/alecthomas/kingpin.v2"

)

func main() {
    var (
        gerritAddress = kingpin.Flag("gerrit-url", "The url of gerrit.").String()
        gerritUser    = kingpin.Flag("gerrit-user", "The username to login gerrit.").String()
        gerritUserPublickey = kingpin.Flag("gerrit-user-publickey", "The Public key for user to login gerrit").String()
        webhookUrl = kingpin.Flag("webhook-url", "The address of webhook for sending event stream.").Default("").String()
    )

    kingpin.HelpFlag.Short('h')
    kingpin.Parse()

    var gerritServer = gerritssh.New(*gerritAddress, *gerritUser, *gerritUserPublickey)
    defer gerritServer.StopStreamEvents()

    if *webhookUrl == "" {
        fmt.Println("No webhook config, stop gerrit ssh conn ")
        return
    }

    // get gerrit event
    fmt.Println("Start Gerrit event")
    go gerritServer.StartStreamEvents()

    // send gerrit event to webhook
    fmt.Println("Send Gerrit event")
    var event gerritssh.StreamEvent
    for {
        select {
        case event = <-gerritServer.ResultChan:
            go func(event gerritssh.StreamEvent) {
                fmt.Println("Receive event:", event)
                bytesData, _ := json.Marshal(&event)
                resp, err := http.Post(*webhookUrl,"application/json", bytes.NewReader(bytesData))
                defer resp.Body.Close()
                if err != nil {
                    fmt.Println("response err:", err)
                    return
                }
                data, _ := ioutil.ReadAll(resp.Body)
                fmt.Println("response Status:", resp.Status, "response Body:", data)
            }(event)
        case <- gerritServer.StopChan:
            fmt.Println("Gerrit ssh stop receive event stream, return for sending webhook")
            return
        default:
            // fmt.Println("Receive no data, go to default case, do nothing")
        }
    }
}
