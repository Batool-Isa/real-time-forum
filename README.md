# Real-Time Forum
This project is a single-page application (SPA) that enables users to post messages, view posts, comment, and engage in private real-time messaging. It builds upon the previous forum project, introducing premium features like private messages, a dynamic online/offline user list, and a modernized user experience.


## Features

### Core Features:
- **Registration and Login:** Secure user registration and login system with validations.
- **Message Posting:** Users can create categorized posts and view them in a feed.
- **Comments on Posts:** Users can comment on posts, visible only when the post is opened.
- **Private Messages:** Real-time private messaging system with:
  - Online/offline status indicators.
  - Message history and dynamic scrolling.
  - Notifications for new messages.

### Real-Time Functionality:
- WebSocket-powered updates for real-time interactions.
- Instant notifications for private messages.
- Live updates for posts, comments, and user status.

### User Interface:
- **Single Page Application:** All page changes are handled via JavaScript to ensure seamless navigation.
- Dynamic rendering of online users and chat history.
- Organized chat interface (similar to Discord).

## Technologies Used

### Frontend:
- **HTML:** Structuring the elements of the page.
- **CSS:** Styling the interface for a responsive and user-friendly design.
- **JavaScript:** Managing events, WebSocket communication, and dynamic content updates.

### Backend:
- **Golang:** Handling data operations, WebSocket communication, and session management.
- **SQLite:** Storing user data, posts, comments, and message history.

### Containerization:
- **Docker:** Simplified setup and deployment.

## Installation and Setup

### Prerequisites
- **Golang** and **SQLite** installed on your system.
- **Docker** (optional, for containerized deployment).

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
    ./run_docker.sh
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