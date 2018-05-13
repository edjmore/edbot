package groupme

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

func Respond(m Message) {
	botID, ok := os.LookupEnv("GROUPME_BOT_ID")
	if !ok {
		panic("No GROUPME_BOT_ID found.")
	}

	// Don't send responses to other bots, or our own messages (not sure if necessary).
	if m.SenderID == botID || m.SenderType != "user" {
		fmt.Printf("Not sending response to message: %v", m)
		return
	}
	r := Response{BotID: botID, Text: fmt.Sprintf("%s said \"%s\"", m.Name, m.Text)}
	j, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	http.Post("https://api.groupme.com/v3/bots/post", "application/json", bytes.NewBuffer(j))
}
