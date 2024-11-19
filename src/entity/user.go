package entity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type SignIn struct {
	ID              string `json:"id"`
	Email           string `json:"email"`
	Name            string `json:"name"`
	Role            string `json:"role"`
	ProfileImageURL string `json:"profile_image_url"`
	Token           string `json:"token"`
	TokenType       string `json:"token_type"`
}

func makeUiUserCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "make:users",
		Short: "Create UI users",
		Run: func(cmd *cobra.Command, args []string) {
			in, _ := cmd.Flags().GetString("in")
			out, _ := cmd.Flags().GetString("out")
			log.Printf("Creating UI users from %s to %s\n", in, out)
			err := makeUsers(in, out)
			if err != nil {
				log.Printf("Error creating users: %v\n", err)
			}
		},
	}
	cmd.Flags().String("in", "", "Input template path")
	cmd.Flags().String("out", "", "Output path for the user import")

	cmd.MarkFlagRequired("in")
	cmd.MarkFlagDirname("in")
	cmd.MarkFlagRequired("out")
	cmd.MarkFlagDirname("out")
	return cmd
}

func getUiUserCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get:users",
		Short: "Fetch UI users",
		Run: func(cmd *cobra.Command, args []string) {
			token, _ := cmd.Flags().GetString("token")
			server, _ := cmd.Flags().GetString("server")
			log.Printf("Fetching UI users with token %s from server %s\n", token, server)
			err := getUsersFromAPI(server, token)
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

func saveUiUserCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "save:users",
		Short: "Save UI users",
		Run: func(cmd *cobra.Command, args []string) {
			token, _ := cmd.Flags().GetString("token")
			server, _ := cmd.Flags().GetString("server")
			in, _ := cmd.Flags().GetString("in")
			log.Printf("Saving UI users with token %s to server %s from %s\n", token, server, in)
			err := saveUsers(server, token, in)
			if err != nil {
				log.Printf("Error saving users: %v\n", err)
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

func makeUsers(in, out string) error {
	// Read JSON data from input file
	data, err := os.ReadFile(in)
	if err != nil {
		log.Printf("Failed to read input file: %v\n", err)
		return err
	}

	// Unmarshal JSON data into a slice of User structs
	var users []User
	err = json.Unmarshal(data, &users)
	if err != nil {
		log.Printf("Failed to parse JSON data: %v\n", err)
		return err
	}

	// Generate GUID for each user's password
	for i := range users {
		users[i].Password = uuid.New().String()
	}

	// Marshal modified users slice back to JSON
	modifiedData, err := json.MarshalIndent(users, "", "    ")
	if err != nil {
		log.Printf("Failed to marshal modified data: %v\n", err)
		return err
	}

	// Write the modified JSON data to the output file
	err = os.WriteFile(out, modifiedData, 0644)
	if err != nil {
		log.Printf("Failed to write to output file: %v\n", err)
		return err
	}

	log.Println("Output file generated successfully with GUIDs for passwords.")
	return nil
}

func saveUsers(server, token, in string) error {

	// Load users from file
	users, err := getUsersFromFile(in)
	if err != nil {
		log.Printf("Error loading users: %v\n", err)
		return err
	}

	//Iterate over the users and save them
	for _, user := range users {
		// Marshal User struct back to JSON for the request payload
		jsonData, err := json.Marshal(user)
		if err != nil {
			log.Printf("Error marshalling JSON: %v\n", err)
			return err
		}

		// Prepare the HTTP POST request
		url := server + `/api/v1/auths/add`
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
	}

	return nil
}

func getUsersFromAPI(server, token string) error {
	// Prepare the HTTP GET request
	url := server + `/api/v1/users/`
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
	log.Printf("Response : %s\n", string(respBody))

	return nil
}

func signIn(server string, user User) (SignIn, error) {

	//Signin with the user credentials and get the user token
	log.Printf("Signing in to server %s\n with user %s\n", server, user.Email)
	// Create the payload
	var signIn SignIn
	payload := map[string]string{
		"email":    user.Email,
		"password": user.Password,
	}

	// Marshal the payload into JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return signIn, fmt.Errorf("failed to marshal JSON payload: %v", err)
	}

	// Set up the request
	url := server + "/api/v1/auths/signin"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return signIn, fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return signIn, fmt.Errorf("failed to execute HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return signIn, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the JSON response to extract the token
	if err := json.Unmarshal(body, &signIn); err != nil {
		return signIn, fmt.Errorf("failed to unmarshal JSON response: %v", err)
	}

	// Retrieve the token
	token := signIn.Token

	if token == "" {
		return signIn, fmt.Errorf("token is empty")
	}

	return signIn, nil
}

func getUsersFromFile(in string) ([]User, error) {

	// Load JSON data from file
	filePath := in
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file: %v\n", err)
		return nil, err
	}

	// Unmarshal JSON data into Config struct
	var users []User
	if err := json.Unmarshal(fileData, &users); err != nil {
		log.Printf("Error parsing JSON: %v\n", err)
		return nil, err
	}
	return users, nil
}
