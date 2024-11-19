package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSaveChat(t *testing.T) {
	/*err := saveChat("https://demo.autoptic.com", "/Users/pkarayan/autoptic/instance/templates/ui/users.json",
		"/Users/pkarayan/autoptic/instance/templates/ui/chats.json")
	assert.NoError(t, err)*/
}

// TestNewChat is a unit test for the NewChat function
func TestNewChat(t *testing.T) {
	req := ChatRequest{
		Title:       "test 001 chat",
		System:      "All answers must generate a single PQL program. Always answer with a single PQL (performance query language) program with 50 lines or less. When generating comments use // and NOT # . Do not summarize or explain just output the completed code. The answer is always a code snippet.",
		Question:    "What is the sum of the values in the 'value' column of the 'table' table?",
		Answer:      "SELECT SUM(value) FROM table",
		Model:       "pql:latest",
		Temperature: 0.15,
	}

	_, err := createChatJSON(req)
	assert.NoError(t, err)
}

func TestGenerateChat(t *testing.T) {
	// Prepare a sample ChatRequest
	req := ChatRequest{
		Title:       "Sample Chat",
		System:      "test-system",
		Question:    "What is the weather today?",
		Answer:      "It's sunny.",
		Model:       "test-model",
		Temperature: 0.7,
	}

	// Generate the chat
	chat, err := createChatJSON(req)

	// Ensure no error was returned
	assert.NoError(t, err)

	// Validate the ChatContents fields
	assert.Equal(t, req.Title, chat.ChatContents.Title)
	assert.Equal(t, req.System, chat.ChatContents.System)
	assert.Equal(t, req.Temperature, chat.ChatContents.Options.Temperature)
	assert.Equal(t, []string{req.Model}, chat.ChatContents.Models)
	assert.Equal(t, int64(time.Now().UnixMilli()), chat.ChatContents.Timestamp, "Timestamp should be close to the current time")

	// Validate the generated messages
	assert.Len(t, chat.ChatContents.Messages, 2)
	userMessage := chat.ChatContents.Messages[0]
	assistantMessage := chat.ChatContents.Messages[1]

	// Validate User Message
	assert.Equal(t, "user", userMessage.Role)
	assert.Equal(t, req.Question, userMessage.Content)
	assert.Equal(t, []string{req.Model}, userMessage.Models)
	assert.Nil(t, userMessage.ParentID)
	assert.NotEmpty(t, userMessage.ChildrenIDs)
	assert.Equal(t, userMessage.ChildrenIDs[0], assistantMessage.ID)

	// Validate Assistant Message
	assert.Equal(t, "assistant", assistantMessage.Role)
	assert.Equal(t, req.Answer, assistantMessage.Content)
	assert.Equal(t, req.Model, assistantMessage.Model)
	assert.Equal(t, req.Model, assistantMessage.ModelName)
	assert.Equal(t, *assistantMessage.ParentID, userMessage.ID)
	assert.Empty(t, assistantMessage.ChildrenIDs)

	// Validate the History structure
	assert.Equal(t, assistantMessage.ID, chat.ChatContents.History.CurrentID)
	assert.Contains(t, chat.ChatContents.History.Messages, userMessage.ID)
	assert.Contains(t, chat.ChatContents.History.Messages, assistantMessage.ID)
	assert.Equal(t, userMessage, chat.ChatContents.History.Messages[userMessage.ID])
	assert.Equal(t, assistantMessage, chat.ChatContents.History.Messages[assistantMessage.ID])

}

func TestGenerateChat_EmptyRequest(t *testing.T) {
	// Prepare an empty ChatRequest
	req := ChatRequest{}

	// Generate the chat
	chat, err := createChatJSON(req)

	// Ensure no error was returned
	assert.NoError(t, err)

	// Validate default ChatContents fields
	assert.NotEmpty(t, chat.ChatContents.ID)
	assert.Empty(t, chat.ChatContents.Title)
	assert.Empty(t, chat.ChatContents.System)
	assert.Equal(t, 0.0, chat.ChatContents.Options.Temperature)
	assert.Equal(t, 1, len(chat.ChatContents.Models))
	assert.Len(t, chat.ChatContents.Messages, 2)

	// Validate empty request still generates unique message IDs and history
	userMessage := chat.ChatContents.Messages[0]
	assistantMessage := chat.ChatContents.Messages[1]

	assert.Equal(t, assistantMessage.ID, chat.ChatContents.History.CurrentID)
	assert.Contains(t, chat.ChatContents.History.Messages, userMessage.ID)
	assert.Contains(t, chat.ChatContents.History.Messages, assistantMessage.ID)
}
