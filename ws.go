package main

import (
    "log"
    "net/http"

    "github.com/gorilla/websocket"
)

func InitWebSocket() {
    http.HandleFunc("/ws", handleWebSocketConnections)
}

func handleWebSocketConnections(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Error upgrading to WebSocket:", err)
        return
    }
    defer ws.Close()

    clients[ws] = true

    for {
        _, _, err := ws.ReadMessage()
        if err != nil {
            delete(clients, ws)
            break
        }
    }
}
