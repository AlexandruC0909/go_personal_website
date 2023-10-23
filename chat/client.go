package chat

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"go_api/database"
	"go_api/templates"
	"go_api/types"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	hub   *Hub
	conn  *websocket.Conn
	send  chan []byte
	store database.Methods
}

type MessageData struct {
	User        *types.User
	ChatMessage string
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
	}
}

func (c *Client) writePump(w http.ResponseWriter, r *http.Request) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			var parsedMessage map[string]interface{}
			err := json.Unmarshal(message, &parsedMessage)
			if err != nil {
				log.Println("Error parsing JSON:", err)
				return
			}
			cookie, _ := r.Cookie("email")
			db := &database.DbConnection{}
			dbname := os.Getenv("DB_NAME")
			dbPassword := os.Getenv("DB_PASSWORD")
			dbUser := os.Getenv("DB_USER")
			connString := "user=" + dbUser + " dbname=" + dbname + " password=" + dbPassword + " sslmode=disable"
			sqlDB, err := sql.Open("postgres", connString)
			if err != nil {
				log.Fatal(err)
			}
			defer sqlDB.Close()

			db.DB = sqlDB

			user, _ := db.GetUserByEmail(cookie.Value)
			chatMessage, ok := parsedMessage["chat_message"].(string)
			if !ok {
				log.Println("Error parsing chat message")
				return
			}

			tmpl, err := template.ParseFS(templates.Templates, "user/userMessage.html")
			if err != nil {
				log.Println("Error parsing template file:", err)
				return
			}

			data := MessageData{
				User:        user,
				ChatMessage: chatMessage,
			}

			var tplBuffer bytes.Buffer

			err = tmpl.Execute(&tplBuffer, data)
			if err != nil {
				log.Println("Error executing template:", err)
				return
			}

			err = c.conn.WriteMessage(websocket.TextMessage, tplBuffer.Bytes())
			if err != nil {
				log.Println(err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	go client.writePump(w, r)
	go client.readPump()
}
