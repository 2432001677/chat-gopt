package gpt

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/2432001677/chat-gopt/db"
	"go.mongodb.org/mongo-driver/bson"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type AskReq struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}
type Choice struct {
	Message Message `json:"message"`
}
type AskRes struct {
	Choices []Choice `json:"choices"`
}

func AskMe(ip, question string) (string, error) {
	keys := strings.Split(os.Getenv("OPENAI_API_KEY"), ",")
	orgs := strings.Split(os.Getenv("OPENAI_ORGANIZATION"), ",")
	idx := rand.Intn(len(keys))

	key := keys[idx]
	org := orgs[idx]
	if org == "null" {
		org = ""
	}

	database := db.GetMongo()
	cursor, err := database.Collection("qa").Find(context.Background(), bson.D{{Key: "ip", Value: ip}})
	if err != nil {
		return "", err
	}
	defer cursor.Close(context.Background())

	messages := []Message{}
	for cursor.Next(context.Background()) {
		var qa db.Qa
		if err = cursor.Decode(&qa); err != nil {
			continue
		}
		messages = append(messages, Message{
			Role:    "user",
			Content: qa.Question,
		})
		messages = append(messages, Message{
			Role:    "assistant",
			Content: qa.Answer,
		})
	}
	messages = append(messages, Message{
		Role:    "user",
		Content: question,
	})

	askReq := AskReq{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
	}
	jsonBytes, err := json.Marshal(askReq)
	if err != nil {
		return "", err
	}
	client := http.Client{Timeout: 2 * time.Minute}

	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}
	defer req.Body.Close()

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+key)
	req.Header.Add("OpenAI-Organization", org)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var askRes AskRes
	if err := json.Unmarshal(respBytes, &askRes); err != nil {
		return "", err
	}

	return askRes.Choices[0].Message.Content, nil
}
