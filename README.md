# Chat Server in Go

## Overview
This project implements a comprehensive multi-client chat server developed in Go. It facilitates smooth communication with support for user registration, authentication, private messaging, and room-based chats. Additionally, it offers robust utilities for password generation, validation, and color-coded system messages for clarity and ease of use. Designed with scalability in mind, this server provides both public and private room management and a seamless user experience.

---

## Features
- Multi-client support with concurrent handling using Go routines.
- Room-based messaging, supporting both public and password-protected private rooms.
- User registration with secure authentication.
- Real-time password strength validation with suggestions for improvement.
- Direct private messaging between users.
- Dynamic creation and switching of chat rooms.
- Auto-generated passwords meeting strong security criteria.
- Color-coded system messages to enhance readability and usability.
- A lightweight, in-memory database for user and room management.
- Robust error handling and logging for smooth operation.

---

## Requirements
- **Go** (1.18 or newer)
- A terminal client such as `telnet` for connecting to the server.

---

## Installation
1. Clone the repository:
   ```bash
   git clone github.com/vovakirdan/gochat
   cd gochat
   ```

2. Build the application using Go:
   ```bash
   go build -o chat-server
   ```

3. Run the server executable:
   ```bash
   ./chat-server
   ```

By default, the server starts on `127.0.0.1:7878`. Customize the IP and port with command-line flags:
```bash
./chat-server -ip=192.168.1.100 -port=9000
```

---

## Usage

### Starting the Server
To start the server with default configurations:
```bash
./chat-server
```
To specify a custom IP and port:
```bash
./chat-server -ip=0.0.0.0 -port=8080
```

### Connecting to the Server
Utilize any TCP client to establish a connection. For example, using `telnet`:
```bash
telnet 127.0.0.1 7878
```

Upon connecting, the system will guide you through logging in or registering.

---

## Commands
### General Commands
- `/help` - Display all available commands.
- `/help <command>` - Detailed information about a specific command.

### Room Management
- `/list rooms` - View all existing chat rooms.
- `/create room <room_name> <password>` - Create a new chat room with an optional password.
- `/switch room <room_name> <password>` - Switch to another room. A password is needed for private rooms.

### Messaging
- `@<username> <message>` - Send a private message to a specific user.

### User Management
- `/list users` - Display all registered users and their current online/offline statuses.
- `/quit` - Log out and disconnect from the server.

---

## Utilities

### Password Strength Validation
The server evaluates password strength during user registration based on these criteria:
- Minimum length of 6 characters.
- Presence of at least one numeric character.
- Use of both uppercase and lowercase letters.
- Inclusion of at least one special symbol for added security.

For weak passwords, the server provides constructive feedback or offers an option to auto-generate a strong password that meets all criteria.

### Auto-Generated Passwords
The server includes a password generator capable of producing secure, random passwords. These passwords are immediately displayed to the user with recommendations for saving them securely.

---

## File Structure
- **`main.go`**: Entry point of the server, handling initialization and network binding.
- **`server.go`**: Core server logic, managing client connections and commands.
- **`client.go`**: Defines the client context and its attributes.
- **`database.go`**: In-memory database implementation for storing user credentials and room data.
- **`auth.go`**: Functions for password validation and random password generation.
- **`utils.go`**: Utility functions, including text coloring for terminal outputs.

---

## Contribution Guidelines
We welcome contributions to enhance this project. To contribute:
1. Fork the repository.
2. Make your changes in a feature branch.
3. Submit a pull request with a detailed description of your changes.

---

## License
This project is open-source and licensed under the MIT License. Refer to the `LICENSE` file for further details.

---

## Contact
For issues, feature requests, or questions, please open an issue on GitHub or reach out to the project maintainer through the provided contact information. Feel free to contribute to the discussion and help improve this project for everyone.

