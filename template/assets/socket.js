class WebSocketManager {
    constructor() {
        this.url = 'ws://localhost:8888/ws';
        this.initializeSocket();
    }

    initializeSocket() {
        this.socket = new WebSocket(this.url);

        this.socket.onopen = () => {
            console.log('WebSocket Connected');
        };

        this.socket.onerror = (error) => {
            console.error('WebSocket Error:', error);
        };

        this.socket.onclose = (event) => {
            console.warn('WebSocket Disconnected:', event.reason);
            setTimeout(() => this.initializeSocket(), 3000); // Reconnect after 3 seconds
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

        // Display the message in chat UI
        const chatMessages = document.getElementById('chat-messages');
        if (chatMessages) {
            const messageElement = document.createElement('p');
            messageElement.textContent = `${message.username || 'User'}: ${message.content}`;
            chatMessages.appendChild(messageElement);
        }
    }
}

document.addEventListener('DOMContentLoaded', () => {
    window.webSocketManager = new WebSocketManager();
});
