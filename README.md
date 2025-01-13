# Real-Time Forum
This project is a single-page application that allows users to post messages and view messages posted by other users. Additionally, it enables real-time communication among users.

## Features
- Message posting
- Viewing messages
- Comments and likes on posts
- Real-time chat

## Technologies Used
### Frontend:
- HTML
- CSS
- JavaScript

### Backend:
- Golang
- SQLite for the database

### Containerization:
- Docker


## How to Run
### Option 1: Run Directly
1. Clone the repository:
   ```bash
   git clone https://github.com/Batool-Isa/real-time-forum.git
   ```
2. Navigate to the project directory:
   ```bash
    cd real-time-forum
    ```
3. Run the backend server:
   ```bash
    go run main.go
    ```
4. Open the application in your browser:
   ```bash
    http://localhost:8000
    ```
### Option 2: Run with Docker
Ensure Docker is installed on your system.

1. Navigate to the project directory:
    ```bash
    cd real-time-forum
    ```
2. Run the shell script to build and run the Docker container:
    ```bash
    ./run.sh
    ```
This script will:
Build the Docker image with the name my-forum-app.
Run the Docker container and map port 8080 to the host.
3. Access the application in your browser:
    ```bash
    http://localhost:8000
    ```

## Images
- **Home Page:** 
![Home Page](https://github.com/Batool-Isa/real-time-forum/blob/main/template/assets/images/image-1.png?raw=true)
- **Message Posting:** 
![Post Creation Page](https://github.com/Batool-Isa/real-time-forum/blob/main/template/assets/images/image-2.png?raw=true)
- **Real-Time Chat:**
![Chat](https://github.com/Batool-Isa/real-time-forum/blob/main/template/assets/images/image-3.png?raw=true)
- **Comments and Likes:** 
![Commens and likes](https://github.com/Batool-Isa/real-time-forum/blob/main/template/assets/images/image-4.png?raw=true)