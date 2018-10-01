package groupme

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/edjmore/edbot/db"
	"github.com/edjmore/edbot/yoda"
)

// GroupMe sends this JSON structure when a message is posted to a group.
type Message struct {
	Attachments []string `json:"attachments"`
	AvatarURL   string   `json:"avatar_url"`
	CreatedAt   int      `json:"created_at"`
	GroupID     string   `json:"group_id"`
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	SenderID    string   `json:"sender_id"`
	SenderType  string   `json:"sender_type"`
	SourceGUID  string   `json:"source_guid"`
	System      bool     `json:"system"`
	Text        string   `json:"text"`
	UserID      string   `json:"user_id"`
}

// Structure of an outgoing message from EdBot to a group.
type Response struct {
	BotID string `json:"bot_id"`
	Text  string `json:"text"`
}

func HandleMessage(m Message) {
	botID, ok := os.LookupEnv("GROUPME_BOT_ID")
	if !ok {
		panic("No GROUPME_BOT_ID found.")
	}

	// Don't send responses to other bots, or our own messages (not sure if necessary).
	if m.SenderID == botID || m.SenderType != "user" {
		fmt.Printf("Not sending response to message: %v", m)
	} else {
		r := Response{BotID: botID, Text: craftResponseText(m)}
		j, err := json.Marshal(r)
		if err != nil {
			panic(err)
		}
		http.Post("https://api.groupme.com/v3/bots/post", "application/json", bytes.NewBuffer(j))
	}

	saveToMessageLog(m)
}

func craftResponseText(m Message) string {
	if strings.HasPrefix(m.Text, "@yoda ") {
		return fmt.Sprintf(yoda.Translate(strings.TrimPrefix(m.Text, "@yoda ")))
	} else if strings.HasPrefix(m.Text, "@history ") {
		return "todo"
	} else {
		return fmt.Sprintf("%s said \"%s\"", m.Name, m.Text)
	}
}

// Save the message to persistent storage.
func saveToMessageLog(m Message) {
	conn := db.Conn()
	_, err := conn.Exec(
		"INSERT INTO messages(created_at, group_id, sender_id, text) VALUES($1, $2, $3, $4)",
		m.CreatedAt, m.GroupID, m.SenderID, m.Text)
	if err != nil {
		log.Fatalf("Unable to insert %v, %v\n", m, err)
	}
}
