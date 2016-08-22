package main

import (
    "io"
    "log"
    "net/http"
    "strconv"
)

func userHandler(w http.ResponseWriter, r *http.Request) {
    //param1 := r.URL.Query().Get("name")
    //param2 := r.URL.Query().Get("phone")
    //io.WriteString(w,"中国公民钟某在入境越南胡志明市时发现护照被污损。总领事馆收集证据后，于27日上午向胡志明市外办提出交涉，指出污损中国护照有损中国公民国格、人格，是无耻的懦夫行为。")
    //io.WriteString(w, param1)
    //io.WriteString(w, param2)

    io.WriteString(w, "{\"feed\": [")
    for i := 0; i < 10; i++ {
        for j := 0; j < 10; j++ {
            if i + j != 0 {
                io.WriteString(w, ",")
            }
            io.WriteString(w, "{\"title\":\"")
            io.WriteString(w, strconv.Itoa(10 * i + j))
            io.WriteString(w, "指出污损中国护照有损")
            io.WriteString(w, "\"")
            io.WriteString(w, ",\"content\":\"")
            io.WriteString(w, strconv.Itoa(100 * i + j))
            io.WriteString(w, "中国公民钟某在入境越南胡志明市时发现护照被污损。总领事馆收集证据后，于27日上午向胡志明市外办提出交涉，指出污损中国护照有损中国公民国格、人格，是无耻的懦夫行为。")
            io.WriteString(w, "\",\"username\":\"forkingdog\",")
            io.WriteString(w, "\"time\":\"2015.04.10\",")
            io.WriteString(w, "\"imageName\":\"forkingdog\"}")
        }
    }

    io.WriteString(w, "]}")
}

func main() {
    http.HandleFunc("/user", userHandler)
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe:", err.Error())
    }
}
