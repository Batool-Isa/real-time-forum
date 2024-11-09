class WebSocketManager {
    constructor() {
        this.socket = new WebSocket('ws://localhost:8888/ws');
        this.messageHandlers = new Map();
        this.setupSocketListeners();
    }

    setupSocketListeners() {
        this.socket.onopen = () => {
            console.log('WebSocket Connected');
        };

        this.socket.onerror = (error) => {
            console.log('WebSocket Error:', error);
        };

        this.socket.onmessage = (event) => {
            const message = JSON.parse(event.data);
            this.handleIncomingMessage(message);
        };
    }

    sendMessage(message) {
        if (this.socket.readyState === WebSocket.OPEN) {
            this.socket.send(JSON.stringify(message));
        } else {
            console.error('WebSocket is not open. Unable to send message:', message);
        }
    }

    handleIncomingMessage(message) {
        console.log("Received message:", message);

        // Display the message in chat UI (update as needed)
        const chatMessages = document.getElementById('chat-messages');
        if (chatMessages) {
            const messageElement = document.createElement('p');
            messageElement.textContent = `${message.username || 'User'}: ${message.content}`;
            chatMessages.appendChild(messageElement);
        }
    }
}

// Initialize WebSocketManager globally
document.addEventListener('DOMContentLoaded', () => {
    window.webSocketManager = new WebSocketManager();
});
