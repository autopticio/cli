package entity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// UI Entity Commands
func UICommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ui",
		Short: "Commands related to UI resources like users,chats,suggestions and prompts",
	}

	//Suggested prompts
	cmd.AddCommand(makeUiSuggestionsCommand())
	cmd.AddCommand(getUiSuggestionsCommand())
	cmd.AddCommand(saveUiSuggestionsCommand())

	//Saved prompts
	cmd.AddCommand(makeUiPromptsCommand())
	cmd.AddCommand(getUiPromptsCommand())
	cmd.AddCommand(saveUiPromptsCommand())

	//Users
	cmd.AddCommand(makeUiUserCommand())
	cmd.AddCommand(getUiUserCommand())
	cmd.AddCommand(saveUiUserCommand())

	//Users
	cmd.AddCommand(makeUiChatCommand())
	cmd.AddCommand(getUiChatCommand())
	cmd.AddCommand(saveUiChatCommand())

	return cmd
}

func makeUiSuggestionsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "make:suggestions",
		Short: "Create UI suggestions",
		Run: func(cmd *cobra.Command, args []string) {
			in, _ := cmd.Flags().GetString("in")
			out, _ := cmd.Flags().GetString("out")
			log.Printf("Creating UI suggestions from %s to %s\n", in, out)
			err := makeSuggestions(in, out)
			if err != nil {
				log.Printf("Error creating suggestions: %v\n", err)
			}
		},
	}
	cmd.Flags().String("in", "", "Input template path")
	cmd.Flags().String("out", "", "Output path for the suggestions")

	cmd.MarkFlagRequired("in")
	cmd.MarkFlagDirname("in")
	cmd.MarkFlagRequired("out")
	cmd.MarkFlagDirname("out")

	return cmd
}

func getUiSuggestionsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get:suggestions",
		Short: "Fetch UI suggestions",
		Run: func(cmd *cobra.Command, args []string) {
			token, _ := cmd.Flags().GetString("token")
			server, _ := cmd.Flags().GetString("server")
			log.Printf("Fetching UI suggestions with token %s from server %s\n", token, server)
			err := getSuggestions(server, token)
			if err != nil {
				log.Printf("Error fetching suggestions: %v\n", err)
			}
		},
	}
	cmd.Flags().String("token", "", "API token")
	cmd.Flags().String("server", "", "Server URL")

	cmd.MarkFlagRequired("token")
	cmd.MarkFlagRequired("server")

	return cmd
}

func saveUiSuggestionsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "save:suggestions",
		Short: "Save UI suggestions",
		Run: func(cmd *cobra.Command, args []string) {
			token, _ := cmd.Flags().GetString("token")
			server, _ := cmd.Flags().GetString("server")
			in, _ := cmd.Flags().GetString("in")
			log.Printf("Saving UI suggestions with token %s to server %s from %s\n", token, server, in)
			err := saveSuggestions(server, token, in)
			if err != nil {
				log.Printf("Error saving suggestions: %v\n", err)
			}
		},
	}
	cmd.Flags().String("token", "", "API token")
	cmd.Flags().String("server", "", "Server URL")
	cmd.Flags().String("in", "", "Input file path")

	cmd.MarkFlagRequired("token")
	cmd.MarkFlagRequired("server")
	cmd.MarkFlagRequired("in")
	cmd.MarkFlagDirname("in")

	return cmd
}

// Prompts (UI) Command Handlers
func makeUiPromptsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "make:prompts",
		Short: "Create UI prompts",
		Run: func(cmd *cobra.Command, args []string) {

			in, _ := cmd.Flags().GetString("in")
			out, _ := cmd.Flags().GetString("out")
			log.Printf("Creating UI prompts from %s to %s\n", in, out)
			err := makePrompts(in, out)
			if err != nil {
				log.Printf("Error creating prompts: %v\n", err)
			}
		},
	}

	cmd.Flags().String("in", "", "Input template path")
	cmd.Flags().String("out", "", "Output path for the prompts")

	cmd.MarkFlagRequired("in")
	cmd.MarkFlagDirname("in")
	cmd.MarkFlagRequired("out")
	cmd.MarkFlagDirname("out")

	return cmd
}

func getUiPromptsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get:prompts",
		Short: "Fetch UI prompts",
		Run: func(cmd *cobra.Command, args []string) {
			token, _ := cmd.Flags().GetString("token")
			server, _ := cmd.Flags().GetString("server")
			log.Printf("Fetching UI prompts with token %s from server %s\n", token, server)
			err := getPrompts(server, token)
			if err != nil {
				log.Printf("Error fetching prompts: %v\n", err)
			}
		},
	}
	cmd.Flags().String("token", "", "API token")
	cmd.Flags().String("server", "", "Server URL")

	cmd.MarkFlagRequired("token")
	cmd.MarkFlagRequired("server")

	return cmd
}

func saveUiPromptsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "save:prompts",
		Short: "Save UI prompts",
		Run: func(cmd *cobra.Command, args []string) {
			token, _ := cmd.Flags().GetString("token")
			server, _ := cmd.Flags().GetString("server")
			in, _ := cmd.Flags().GetString("in")
			log.Printf("Saving UI prompts with token %s to server %s from %s\n", token, server, in)
			err := savePrompts(server, token, in)
			if err != nil {
				log.Printf("Error saving prompts: %v\n", err)
			}
		},
	}
	cmd.Flags().String("token", "", "API token")
	cmd.Flags().String("server", "", "Server URL")
	cmd.Flags().String("in", "", "Input file path")

	cmd.MarkFlagRequired("token")
	cmd.MarkFlagRequired("server")
	cmd.MarkFlagRequired("in")
	cmd.MarkFlagDirname("in")
	return cmd
}

// Config represents the structure of the suggestions data read from the file
type Config struct {
	Suggestions []struct {
		Title   []string `json:"title"`
		Content string   `json:"content"`
	} `json:"suggestions"`
}

type Prompt struct {
	Command string `json:"command"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func makeSuggestions(in, out string) error {
	err := CopyFile(in, out)
	if err != nil {
		log.Printf("Error copying directory: %v\n", err)
		return err
	}
	return nil
}

func saveSuggestions(server, token, in string) error {
	// Load JSON data from file
	filePath := in
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file: %v\n", err)
		return err
	}

	// Unmarshal JSON data into Config struct
	var config Config
	if err := json.Unmarshal(fileData, &config); err != nil {
		log.Printf("Error parsing JSON: %v\n", err)
		return err
	}

	// Marshal Config struct back to JSON for the request payload
	jsonData, err := json.Marshal(config)
	if err != nil {
		log.Printf("Error encoding JSON: %v\n", err)
		return err
	}

	// Prepare the HTTP POST request
	url := server + `/api/v1/configs/default/suggestions`
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v\n", err)
		return err
	}

	// Output the response status and body
	log.Printf("Response Status: %s\n", resp.Status)
	log.Printf("Response Body Size: %v\n", len(respBody))

	return nil
}

func getSuggestions(server, token string) error {
	// Prepare the HTTP GET request
	url := server + `/api/v1/configs/default/suggestions`
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return err
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v\n", err)
		return err
	}

	// Output the response status and body
	log.Printf("Response Status: %s\n", resp.Status)
	log.Printf("Response Body Size: %v\n", len(respBody))

	return nil
}

func getPrompts(server, token string) error {
	// Implement fetching prompts from server
	// Prepare the HTTP GET request
	url := server + `/api/v1/prompts`
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return err
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v\n", err)
		return err
	}

	// Output the response status and body
	log.Printf("Response Status: %s\n", resp.Status)
	log.Println(string(respBody))

	return nil
}

func makePrompts(in, out string) error {
	err := CopyFile(in, out)
	if err != nil {
		log.Printf("Error copying directory: %v\n", err)
		return err
	}
	return nil
}

func savePrompts(server, token, in string) error {
	// Replace with your API endpoint URL
	apiCreateEndpoint := server + "/api/v1/prompts/create"
	// Path to the JSON file containing prompts
	promptsFile := in

	// Read the JSON file
	fileData, err := os.ReadFile(promptsFile)
	if err != nil {
		log.Fatalf("Failed to read prompts file: %v", err)
	}

	// Parse JSON data into a slice of Prompt structs
	var prompts []Prompt
	if err := json.Unmarshal(fileData, &prompts); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Iterate through each prompt and post to API endpoint
	for _, prompt := range prompts {
		// Convert the prompt to JSON
		promptData, err := json.Marshal(prompt)
		if err != nil {
			log.Printf("Failed to encode prompt %v: %v", prompt, err)
			return err
		}

		// Create HTTP POST request
		req, err := http.NewRequest("POST", apiCreateEndpoint, bytes.NewBuffer(promptData))
		if err != nil {
			log.Printf("Failed to create request: %v", err)
			return err
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		// Execute the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error sending request: %v\n", err)
			return err
		}
		defer resp.Body.Close()

		// Check response status
		if resp.StatusCode != http.StatusOK {

			if resp.StatusCode == http.StatusBadRequest {
				apiUpdateEndpoint := server + "/api/v1/prompts/command" + prompt.Command + "/update"
				updateReq, err := http.NewRequest("POST", apiUpdateEndpoint, bytes.NewBuffer(promptData))
				if err != nil {
					log.Printf("Failed to create update request: %v", err)
					return err
				}
				updateReq.Header.Set("Content-Type", "application/json")
				updateReq.Header.Set("Accept", "application/json")
				updateReq.Header.Set("Authorization", "Bearer "+token)

				updateResp, err := client.Do(updateReq)
				if err != nil {
					log.Printf("Error sending update request: %v\n", err)
					return err
				}
				defer updateResp.Body.Close()

				if updateResp.StatusCode != http.StatusOK {
					log.Printf("Received non-OK response %v for prompt %v", updateResp.StatusCode, prompt)
				} else {
					log.Printf("Successfully updated prompt: %s\n", prompt.Title)
				}
			} else {
				log.Printf("Received non-OK response %v for prompt %v", resp.StatusCode, prompt)
			}
		} else {
			log.Printf("Successfully sent prompt: %s\n", prompt.Title)
		}

		// Close response body to avoid resource leaks
		resp.Body.Close()
	}

	return nil
}
