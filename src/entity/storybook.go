package entity

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Storybooks Entity Commands
func StorybookCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "storybooks",
		Short: "Commands related to storybooks content",
	}

	cmd.AddCommand(makeStorybooksCommand())
	cmd.AddCommand(saveStorybooksCommand())
	return cmd
}

func makeStorybooksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "make",
		Short: "Create Storybooks data",
		Run: func(cmd *cobra.Command, args []string) {
			in, _ := cmd.Flags().GetString("in")
			out, _ := cmd.Flags().GetString("out")
			log.Printf("Creating Storybooks from %s to %s\n", in, out)
			err := makeStorybooks(in, out)
			if err != nil {
				log.Println(err)
			}

		},
	}
	cmd.Flags().String("in", "", "Input template path")
	cmd.Flags().String("out", "", "Output path for the Storybooks data")

	cmd.MarkFlagRequired("in")
	cmd.MarkFlagDirname("in")
	cmd.MarkFlagRequired("out")
	cmd.MarkFlagDirname("out")
	return cmd
}

func saveStorybooksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "save",
		Short: "Save Storybooks data",
		Run: func(cmd *cobra.Command, args []string) {
			// Retrieve flag values
			server, _ := cmd.Flags().GetString("server")
			in, _ := cmd.Flags().GetString("in")
			ep, _ := cmd.Flags().GetString("ep")
			token, _ := cmd.Flags().GetString("token")

			// Log the action
			log.Printf("Saving Storybooks data to server %s from %s\n", server, in)

			// Call the save function with the collected flags
			err := saveStorybooksContent(in, server, ep, token)
			if err != nil {
				log.Println(err)
			}
		},
	}

	// Define the flags
	cmd.Flags().String("server", "", "Server URL")
	cmd.Flags().String("in", "", "Input path for Storybooks data")
	cmd.Flags().String("ep", "", "Endpoint ID for the instance")
	cmd.Flags().String("token", "", "API Access token")

	// Mark flags as required
	cmd.MarkFlagRequired("server")
	cmd.MarkFlagRequired("in")
	cmd.MarkFlagDirname("in") // Ensures "in" is a valid directory
	cmd.MarkFlagRequired("ep")

	return cmd
}

func makeStorybooks(in, out string) error {
	log.Printf("Creating Storybooks from %s to %s\n", in, out)

	err := CopyDir(in, out)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func saveStorybooksContent(in, server, endpoint, token string) error {

	dataPath := in
	apiEndpointURL := server + "/story/ep/" + endpoint

	// Check if required flags are provided
	if apiEndpointURL == "" || dataPath == "" {
		log.Fatal("Both endpointID and dataDir arguments are required")
	}

	//check if apiEndpoint is a valid URL
	if _, err := http.Get(apiEndpointURL); err != nil {
		log.Fatalf("Invalid API endpoint: %v", err)
	}

	//check if dataPath is a valid directory name
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		log.Fatalf("Invalid data directory: %v", err)
	}

	// Define paths to /pql, /env, and /brief folders within the provided data directory
	pqlDir := filepath.Join(dataPath, "pqls")
	envDir := filepath.Join(dataPath, "environments")
	briefDir := filepath.Join(dataPath, "briefs")

	// Read files from /pql, /env, and /brief directories
	processDirectory(pqlDir, apiEndpointURL, token, ".pql", submitPQL)
	processDirectory(envDir, apiEndpointURL, token, ".json", submitEnv)
	processDirectory(briefDir, apiEndpointURL, token, ".json", submitBrief)

	return nil
}

// Helper function to process all files in a directory and submit them
func processDirectory(dir, apiEndpoint, token, fileExtension string, handlerFunc func(string, string, string, string) error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed to read directory %s: %v", dir, err)
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == fileExtension {
			filePath := filepath.Join(dir, file.Name())
			log.Printf("Processing file: %s", filePath)

			// Now pass the full file name (including extension) to the handler function
			err := handlerFunc(dir, apiEndpoint, file.Name(), token)
			if err != nil {
				log.Printf("Failed to process file %s: %v", file.Name(), err)
			}
		}
	}
}

// Submit PQL data to the API
func submitPQL(baseDataPath, apiEndpoint, pqlFileName, token string) error {
	fileContent, err := readFileContent(baseDataPath, pqlFileName)
	if err != nil {
		return err
	}
	// Extract the PQL ID by removing the extension from the filename
	pqlID := extractResourceID(pqlFileName, ".pql")
	apiURL := fmt.Sprintf("%s/pql/%s", apiEndpoint, pqlID)
	return postToAPI(apiURL, fileContent, token)
}

// Submit Environment data to the API
func submitEnv(baseDataPath, apiEndpoint, envFileName, token string) error {
	fileContent, err := readFileContent(baseDataPath, envFileName)
	if err != nil {
		return err
	}
	// Extract the Env ID by removing the extension from the filename
	envID := extractResourceID(envFileName, ".json")
	apiURL := fmt.Sprintf("%s/env/%s", apiEndpoint, envID)
	return postToAPI(apiURL, fileContent, token)
}

// Submit Brief data to the API
func submitBrief(baseDataPath, endpointID, briefFileName, token string) error {
	fileContent, err := readFileContent(baseDataPath, briefFileName)
	if err != nil {
		return err
	}
	// Extract the Brief ID by removing the extension from the filename
	briefID := extractResourceID(briefFileName, ".json")
	apiURL := fmt.Sprintf("%s/brief/%s", endpointID, briefID)
	return postToAPI(apiURL, fileContent, token)
}

// Read and base64 encode file content
func readFileContent(folder, fileName string) (string, error) {
	// Use the full file name (including extension)
	filePath := filepath.Join(folder, fileName)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}
	encodedContent := base64.StdEncoding.EncodeToString(content)
	return encodedContent, nil
}

// Post data to the API
func postToAPI(url, data, token string) error {
	client := &http.Client{}

	reqBody := []byte(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("x-api-token", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}

	log.Printf("Successfully submitted data to %s", url)
	return nil
}

// Extract resourceID from the file name by removing the file extension
// Assuming the file name format is {resource_id}.pql or {resource_id}.json
func extractResourceID(fileName, fileExtension string) string {
	baseName := filepath.Base(fileName)
	nameWithoutExt := baseName[:len(baseName)-len(fileExtension)]
	return nameWithoutExt
}
