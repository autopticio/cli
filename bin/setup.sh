#!/bin/bash

# Define variables
BASE_PATH="/path/to/local/resources"
UI_SERVER_URL="https://localhost:8888"
UI_TOKEN="sk-ui-bearer-token"

API_SERVER_URL="http://localhost:8080"
API_ENDPOINT_ID="autoptic-api-server-epid"
API_TOKEN="autoptic-api-server-auth-token"

# Paths to template files and output files
USER_TEMPLATE="${BASE_PATH}/templates/ui/users.json"
CHAT_TEMPLATE="${BASE_PATH}/templates/ui/chats.json"
PROMPT_TEMPLATE="${BASE_PATH}/templates/ui/prompts.json"
SUGGESTION_TEMPLATE="${BASE_PATH}/templates/ui/suggestions.json"
STORYBOOK_TEMPLATE="${BASE_PATH}/templates/storybooks"

# Paths to generated output files
USER_OUTPUT="${BASE_PATH}/data/ui/users.json"
CHAT_OUTPUT="${BASE_PATH}/data/ui/chats.json"
PROMPT_OUTPUT="${BASE_PATH}/data/ui/prompts.json"
SUGGESTION_OUTPUT="${BASE_PATH}/data/ui/suggestions.json"
STORYBOOK_OUTPUT="${BASE_PATH}/data/storybooks"
INVENTORY_OUTPUT="${BASE_PATH}/data/inventory/service_inventory.json"

# Determine OS and architecture
OS="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Path to the autopticli executable
AUTOPTI_CLI="${BASE_PATH}/bin/exe/${OS}-${ARCH}/autopticli"

if [[ ! -x "$AUTOPTI_CLI" ]]; then
    echo "autopticli executable not found or not executable at: $AUTOPTI_CLI"
    exit 1
fi

echo "Using autopticli executable: $AUTOPTI_CLI"

export PATH="${PATH}:${BASE_PATH}/bin"
echo "Creating resources..."

# Create users from a template
"$AUTOPTI_CLI" ui make:users --in "$USER_TEMPLATE" --out "$USER_OUTPUT"
echo "Users created: $USER_OUTPUT"

# Create chats from a template
"$AUTOPTI_CLI" ui make:chats --in "$CHAT_TEMPLATE" --out "$CHAT_OUTPUT"
echo "Chats created: $CHAT_OUTPUT"

# Create prompts from a template
"$AUTOPTI_CLI" ui make:prompts --in "$PROMPT_TEMPLATE" --out "$PROMPT_OUTPUT"
echo "Prompts created: $PROMPT_OUTPUT"

# Create suggestions from a template
"$AUTOPTI_CLI" ui make:suggestions --in "$SUGGESTION_TEMPLATE" --out "$SUGGESTION_OUTPUT"
echo "Suggestions created: $SUGGESTION_OUTPUT"

# Create Storybooks data from a template
"$AUTOPTI_CLI" storybooks make --in "$STORYBOOK_TEMPLATE" --out "$STORYBOOK_OUTPUT"
echo "Storybooks data created: $STORYBOOK_OUTPUT"

# Function to confirm before saving
confirm_save() {
    local prompt_message="$1"
    local command="$2"
    
    read -p "$prompt_message (y/n): " confirm
    if [[ $confirm == [yY] ]]; then
        eval "$command"
        echo "Save command executed."
    else
        echo "Skipped saving."
    fi
}

echo "Proceeding to save resources to the server..."

# Prompt user to save each resource to the server
confirm_save "Do you want to save the users to the server?" \
    "\"$AUTOPTI_CLI\" ui save:users --in \"$USER_OUTPUT\" --server \"$UI_SERVER_URL\" --token \"$UI_TOKEN\""

confirm_save "Do you want to save the chats to the server?" \
    "\"$AUTOPTI_CLI\" ui save:chats --chats \"$CHAT_OUTPUT\" --users \"$USER_OUTPUT\" --server \"$UI_SERVER_URL\""

confirm_save "Do you want to save the prompts to the server?" \
    "\"$AUTOPTI_CLI\" ui save:prompts --in \"$PROMPT_OUTPUT\" --server \"$UI_SERVER_URL\" --token \"$UI_TOKEN\""

confirm_save "Do you want to save the suggestions to the server?" \
    "\"$AUTOPTI_CLI\" ui save:suggestions --in \"$SUGGESTION_OUTPUT\" --server \"$UI_SERVER_URL\" --token \"$UI_TOKEN\""

confirm_save "Do you want to save the Storybooks data to the server?" \
    "\"$AUTOPTI_CLI\" storybooks save --in \"$STORYBOOK_OUTPUT\" --server \"$API_SERVER_URL\" --token \"$API_TOKEN\" --ep \"$API_ENDPOINT_ID\""

