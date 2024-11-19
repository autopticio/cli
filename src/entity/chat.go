package entity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type Chat struct {
	ChatContents ChatContents `json:"chat"`
}

// Chat struct with JSON tags
type ChatContents struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Models    []string  `json:"models"`
	System    string    `json:"system"`
	Options   Options   `json:"options"`
	Messages  []Message `json:"messages"`
	History   History   `json:"history"`
	Tags      []string  `json:"tags"`
	Timestamp int64     `json:"timestamp"`
}

// Options struct with JSON tags
type Options struct {
	Temperature float64 `json:"temperature"`
}

// Message struct with JSON tags
type Message struct {
	ID          string   `json:"id"`
	ParentID    *string  `json:"parentId"`
	ChildrenIDs []string `json:"childrenIds"`
	Role        string   `json:"role"`
	Content     string   `json:"content"`
	Timestamp   int64    `json:"timestamp"`
	Models      []string `json:"models"`
	Model       string   `json:"model"`
	ModelName   string   `json:"modelName"`
	UserContext *string  `json:"userContext"`
}

// History struct with JSON tags
type History struct {
	Messages  map[string]Message `json:"messages"`
	CurrentID string             `json:"currentId"`
}

type ChatRequest struct {
	Title       string  `json:"title"`
	System      string  `json:"system"`
	Question    string  `json:"question"`
	Answer      string  `json:"answer"`
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
}

func makeUiChatCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "make:chats",
		Short: "Create UI chats",
		Run: func(cmd *cobra.Command, args []string) {
			in, _ := cmd.Flags().GetString("in")
			out, _ := cmd.Flags().GetString("out")
			log.Printf("Creating UI chats from %s to %s\n", in, out)
			err := makeChats(in, out)
			if err != nil {
				log.Printf("Error creating users: %v\n", err)
			}
		},
	}
	cmd.Flags().String("in", "", "Input template path")
	cmd.Flags().String("out", "", "Output path for the chat import")

	cmd.MarkFlagRequired("in")
	cmd.MarkFlagDirname("in")
	cmd.MarkFlagRequired("out")
	cmd.MarkFlagDirname("out")
	return cmd
}

func getUiChatCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get:chats",
		Short: "Fetch UI users",
		Run: func(cmd *cobra.Command, args []string) {
			token, _ := cmd.Flags().GetString("token")
			server, _ := cmd.Flags().GetString("server")
			log.Printf("Fetching UI chats with token %s from server %s\n", token, server)
			err := getChatsFromAPI(server, token)
			if err != nil {
				log.Printf("Error fetching users: %v\n", err)
			}
		},
	}
	cmd.Flags().String("token", "", "API token")
	cmd.Flags().String("server", "", "Server URL")

	cmd.MarkFlagRequired("token")
	cmd.MarkFlagRequired("server")
	return cmd
}

func saveUiChatCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "save:chats",
		Short: "Save UI chats",
		Run: func(cmd *cobra.Command, args []string) {
			server, _ := cmd.Flags().GetString("server")
			chats, _ := cmd.Flags().GetString("chats")
			users, _ := cmd.Flags().GetString("users")
			log.Printf("Saving UI chats  to server %s from %s\n", server, chats)
			err := saveChats(server, users, chats)
			if err != nil {
				log.Printf("Error saving chats: %v\n", err)
			}
		},
	}
	cmd.Flags().String("server", "", "Server URL")
	cmd.Flags().String("users", "", "Input file path for users")
	cmd.Flags().String("chats", "", "Input file path for chats")

	// Mark flags as required
	cmd.MarkFlagRequired("server")
	cmd.MarkFlagRequired("users")
	cmd.MarkFlagFilename("users")
	cmd.MarkFlagRequired("chats")
	cmd.MarkFlagFilename("chats")

	return cmd
}

func makeChats(in, out string) error {
	err := CopyFile(in, out)
	if err != nil {
		log.Printf("Error copying directory: %v\n", err)
		return err
	}
	return nil
}

func saveChats(server, usersIn, chatsIn string) error {

	//Get a list of users from a local template file
	users, err := getUsersFromFile(usersIn)
	if err != nil {
		return err
	}

	//Signin with each user credentials and get the user token
	for _, user := range users {
		signIn, err := signIn(server, user)
		if err != nil {
			return err
		}
		//Create the chats in the UI API
		chatRequests, err := getChatsFromFile(chatsIn)
		if err != nil {
			return err
		}
		for _, chatRequest := range chatRequests {
			chat, err := createChatJSON(chatRequest)
			if err != nil {
				return err
			}
			//Post the chat to the server
			err = createChat(server, signIn, chat)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func createChat(server string, signIn SignIn, chat Chat) error {
	// Marshal the Chat struct to JSON
	jsonData, err := json.Marshal(chat)
	if err != nil {
		return fmt.Errorf("failed to marshal chat struct: %w", err)
	}

	// Define the API endpoint and authorization token
	url := server + "/api/v1/chats/new"
	auth := signIn.TokenType + " " + signIn.Token

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json")

	// Initialize HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	return nil
}

func getChatsFromAPI(server, token string) error {
	// Define the API endpoint and authorization token
	url := server + "/api/v1/chats/all/db"
	auth := "Bearer " + token

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json")

	// Initialize HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v\n", err)
		return err
	}

	// Output the response status and body
	log.Printf("Response Status: %s\n", resp.Status)
	log.Printf("Response : %s\n", string(respBody))

	return nil
}

func getChatsFromFile(chatRequestsIn string) ([]ChatRequest, error) {

	// Load JSON data from file
	filePath := chatRequestsIn
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file: %v\n", err)
		return nil, err
	}

	// Unmarshal JSON data into Config struct
	var chatRequests []ChatRequest
	if err := json.Unmarshal(fileData, &chatRequests); err != nil {
		log.Printf("Error parsing JSON: %v\n", err)
		return nil, err
	}
	return chatRequests, nil
}

// Creates a sample Chat object closely matching the specified JSON structure
func createChatJSON(req ChatRequest) (Chat, error) {

	timestamp := time.Now().Unix()
	timestampMilli := time.Now().UnixMilli()

	// Generate UUIDs for consistent message IDs
	userMessageID := uuid.New().String()
	assistantMessageID := uuid.New().String()

	// User message
	userMessage := Message{
		ID:          userMessageID,
		ParentID:    nil,
		ChildrenIDs: []string{assistantMessageID},
		Role:        "user",
		Content:     req.Question,
		Timestamp:   timestamp,
		Models:      []string{req.Model},
	}

	// Assistant message
	assistantMessage := Message{
		ID:          assistantMessageID,
		ParentID:    &userMessageID,
		ChildrenIDs: []string{},
		Role:        "assistant",
		Content:     req.Answer,
		Timestamp:   timestamp,
		Model:       req.Model,
		ModelName:   req.Model,
		UserContext: nil,
	}

	// History map
	historyMap := map[string]Message{
		userMessageID:      userMessage,
		assistantMessageID: assistantMessage,
	}

	// History object
	history := History{
		Messages:  historyMap,
		CurrentID: assistantMessageID,
	}

	// Construct the Chat object
	chat := Chat{
		ChatContents: ChatContents{
			ID:        uuid.New().String(),
			Title:     req.Title,
			Models:    []string{req.Model},
			System:    req.System,
			Options:   Options{Temperature: req.Temperature},
			Messages:  []Message{userMessage, assistantMessage},
			History:   history,
			Tags:      []string{},
			Timestamp: timestampMilli,
		},
	}

	return chat, nil
}
