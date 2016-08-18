package main

import (
    "io"
    "log"
    "net/http"
)
func userHandler(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w,"中国公民钟某在入境越南胡志明市时发现护照被污损。总领事馆收集证据后，于27日上午向胡志明市外办提出交涉，指出污损中国护照有损中国公民国格
、人格，是无耻的懦夫行为。")
}

func main() {
    http.HandleFunc("/user", userHandler)
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe:", err.Error())
    }
}