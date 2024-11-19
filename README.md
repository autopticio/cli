# setup.sh
The script creates default content (users,pql,env,prompts,chats) for a newly installed Autoptic instance. This Bash script automates the process of generating and saving various UI and storybook data resources using the `autopticli`. Here's a short summary:

1. **Setup:** The script defines several variables, including paths to local resources and server URLs, authentication tokens, and template file paths.

2. **Resource Creation:** It generates different resource files—users, chats, prompts, suggestions, and storybooks—using templates by calling the `autopticli` command-line tool.

3. **Output Paths:** The generated files are saved to specific output directories.

4. **User Confirmation:** The script includes a function to prompt the user to confirm before saving each generated resource to the server.

5. **Save to Server:** Based on user confirmation, the script saves the generated resources to the specified servers by executing save commands with the appropriate server URL and authentication tokens.

6. **Completion:** It outputs a message indicating the completion of the save operations.



# Autopticli CLI

`autopticli` is a command-line tool for managing various resources, including API services, UI elements, Storybooks content, and inventory files. This tool provides a standardized interface to create, fetch, and save data for different categories.

## Installation

Download or clone this repository and ensure `autopticli` is added to your PATH.

## Usage

The basic syntax for using `autopticli` is:

```sh
autopticli [command]
```

For detailed information on any specific command, use the `--help` flag:

```sh
autopticli [command] --help
```

## Commands and Complete Examples

### API Commands

Commands related to API services can be managed with `autopticli api`, although specific subcommands are not provided in the documentation.

#### Example

```sh
autopticli api [subcommand]
```

### Inventory Commands

The `inventory` command allows you to create inventory files.

#### Example: Create an Inventory File

This command creates an inventory file from the provided data and saves it to the specified output path.

```sh
autopticli inventory make --out /path/to/output/inventory.json
```

### Storybooks Commands

Manage Storybooks data using the `storybooks` command, which includes options to create or save data.

#### Example: Create Storybooks Data

This command creates Storybooks data based on a template file and outputs it to a specified location.

```sh
autopticli storybooks make --in /path/to/template.json --out /path/to/output/storybooks.json
```

#### Example: Save Storybooks Data to a Server

This command saves Storybooks data to a specified server.

```sh
autopticli storybooks save --in /path/to/storybooks.json --server https://example.com --token YOUR_API_TOKEN --ep ENDPOINT_ID
```

### UI Commands

Manage UI resources like users, chats, prompts, and suggestions using `autopticli ui`. Each of these resources has subcommands to fetch, create, and save data.

#### Fetch Commands

Fetch data from the server using `get` commands.

- **Fetch Users**: Retrieve a list of users from the server.

  ```sh
  autopticli ui get:users --server https://example.com --token YOUR_API_TOKEN
  ```

- **Fetch Chats**: Retrieve chat data from the server.

  ```sh
  autopticli ui get:chats --server https://example.com --token YOUR_API_TOKEN
  ```

- **Fetch Prompts**: Retrieve prompts from the server.

  ```sh
  autopticli ui get:prompts --server https://example.com --token YOUR_API_TOKEN
  ```

- **Fetch Suggestions**: Retrieve suggestions from the server.

  ```sh
  autopticli ui get:suggestions --server https://example.com --token YOUR_API_TOKEN
  ```

#### Create Commands

Create new UI resources based on an input template file.

- **Create Users**: Create users based on a template.

  ```sh
  autopticli ui make:users --in /path/to/user_template.json --out /path/to/output/users.json
  ```

- **Create Chats**: Generate chat data from a template.

  ```sh
  autopticli ui make:chats --in /path/to/chat_template.json --out /path/to/output/chats.json
  ```

- **Create Prompts**: Generate prompt data from a template.

  ```sh
  autopticli ui make:prompts --in /path/to/prompt_template.json --out /path/to/output/prompts.json
  ```

- **Create Suggestions**: Generate suggestions data from a template.

  ```sh
  autopticli ui make:suggestions --in /path/to/suggestion_template.json --out /path/to/output/suggestions.json
  ```

#### Save Commands

Save data back to a server or file system.

- **Save Users**: Save user data to the server.

  ```sh
  autopticli ui save:users --in /path/to/users.json --server https://example.com --token YOUR_API_TOKEN
  ```

- **Save Chats**: Save chat data, including associated user data, to the server.

  ```sh
  autopticli ui save:chats --chats /path/to/chats.json --users /path/to/users.json --server https://example.com
  ```

- **Save Prompts**: Save prompt data to the server.

  ```sh
  autopticli ui save:prompts --in /path/to/prompts.json --server https://example.com --token YOUR_API_TOKEN
  ```

- **Save Suggestions**: Save suggestion data to the server.

  ```sh
  autopticli ui save:suggestions --in /path/to/suggestions.json --server https://example.com --token YOUR_API_TOKEN
  ```

### Completion Command

Generate autocompletion scripts for a specified shell:

```sh
autopticli completion [shell]
```

### Help Command

Get help for any specific command or subcommand:

```sh
autopticli help [command]
```

## Flags

Common flags available for various commands:

- `--server` - Specify the server URL for data operations.
- `--token` - API token for authentication.
- `--in` - Path to an input file.
- `--out` - Path to an output file.

## Additional Topics

Some advanced commands include:

- **Model Commands** - Manage model-related services.
- **Scheduler Commands** - Manage scheduling services.

To access additional information about these topics:

```sh
autopticli model --help
autopticli scheduler --help
```

## Examples

Here are a few additional examples demonstrating specific use cases.

#### Fetching Prompts Example

```sh
autopticli ui get:prompts --server https://example.com --token YOUR_API_TOKEN
```

#### Creating a Suggestions File Example

```sh
autopticli ui make:suggestions --in /path/to/suggestion_template.json --out /path/to/output/suggestions.json
```

#### Saving User Data to a Server Example

```sh
autopticli ui save:users --in /path/to/users.json --server https://example.com --token YOUR_API_TOKEN
```

## Contributing

To contribute to `autopticli`, please submit a pull request with a description of your changes.


This updated README includes complete examples for each command related to different resources, providing a clear guide for users to execute the commands effectively.