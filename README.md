# GoRedis

A homemade Redis implementation built with Go, featuring a TCP server, a web interface for demonstration, and a simple CLI receiver.

## Overview

GoRedis is a custom Redis-like key-value store written in Go. It includes three main components:

- **Server**: A TCP server (`server` directory) that handles Redis-like commands and data storage.
- **Web**: A web interface (`web` directory) to demonstrate the functionality of the server, with support for Discord and GitHub authentication.
- **Receiver**: A simple command-line interface (`receiver` directory) for interacting with the server.

## Prerequisites

- **Go**: Version 1.21 or higher.
- **Air**: For live-reloading during development (`go install github.com/air@latest`).
- **Node.js**: Required for Tailwind CSS in the web component.
- **Templ**: For generating HTML templates in the web component (`go install github.com/a-h/templ/cmd/templ@latest`).

## Project Structure
GoRedis/
├── server/         # Redis-like TCP server implementation
├── web/            # Web interface for demonstration
├── receiver/       # CLI tool for interacting with the server
└── README.md

## Setup and Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/yourusername/GoRedis.git
   cd GoRedis
   ```
   
2. Install dependencies:
Ensure Go, Air, Node.js, and Templ are installed.

## Running the Project

### Server

1. **Navigate to the server directory**:
  ```bash
  cd server
  air 
  ```

### Server

1. **Navigate to the web directory**:
  ```bash
  cd web
  ```
2. Create a **.env** file in the web directory with the following content:
  ```env
    DISCORD_CLIENT_ID=your_discord_client_id
    DISCORD_CLIENT_SECRET=your_discord_client_secret
    GITHUB_CLIENT_ID=your_github_client_id
    GITHUB_CLIENT_SECRET=your_github_client_secret
  ```

3. Install Node.js dependencies (for Tailwind CSS):
   ```bash
   npm install
   air
   ```

4. In a separate terminal, generate and watch Tailwind CSS styles:
   ```bash
    make tailwindcss
   ```

5. In another terminal, generate and watch Templ templates:
   ```bash
    make templ
   ```

## Makefile (web)
```makefile
.PHONY: tailwindcss
tailwindcss:
	@npx @tailwindcss/cli -i internal/presentation/css/styles.css -o static/css/styles.css --watch

.PHONY: templ
templ:
	@templ generate view -watch
```

## License
License
This project is licensed under the MIT License.
