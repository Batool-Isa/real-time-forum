// Constants and State
const SELECTORS = {
    MENU_ITEMS: '.menu-item',
    CREATE_POST: 'label[for="create-post"]',
    FEEDS: '.feeds',
    CREATE_POST_SECTION: '.create-post',
    CREATE_POST_FORM: '#create-post-form',
    MESSAGES: '.messages',
    MESSAGE_ITEMS: '.messages .message',
    MIDDLE_SECTION: '.middle'
};

let feedsContent = '';

// UI Controller
class UIController {
    constructor() {
        this.elements = {
            menuItems: document.querySelectorAll(SELECTORS.MENU_ITEMS),
            feedsSection: document.querySelector(SELECTORS.FEEDS),
            createPostSection: document.querySelector(SELECTORS.CREATE_POST_SECTION),
            messageItems: document.querySelectorAll(SELECTORS.MESSAGE_ITEMS),
            middleSection: document.querySelector(SELECTORS.MIDDLE_SECTION)
        };
        this.wsManager = new WebSocketManager();
        this.chatManager = new ChatManager(this.wsManager);
        this.initializeElements();
        this.bindEvents();
        this.storeFeedsContent();
    }

    bindEvents() {
        // Menu items (Home and Create Post)
        this.elements.menuItems.forEach(item => {
            item.addEventListener('click', () => {
                const menuText = item.querySelector('h3').textContent;
                
                // Clear all active states
                this.clearAllActiveStates();
                
                // Set new active state and show content
                item.classList.add('active');
                
                if (menuText === 'Home') {
                    this.showHome();
                } else if (menuText === 'Create Post') {
                    this.showCreatePost();
                }
            });
        });
    }

    showHome() {
        this.elements.middleSection.innerHTML = feedsContent;
        this.elements.createPostSection.classList.add('hidden');
        this.elements.feedsSection.classList.remove('hidden');
    }

    showCreatePost() {
        this.elements.feedsSection.classList.add('hidden');
        this.elements.createPostSection.classList.remove('hidden');
        this.elements.middleSection.appendChild(this.elements.createPostSection);
    }

    createChatInterface(username, profilePic) {
        const container = document.createElement('div');
        container.className = 'chat-container';

        // Create header
        const header = document.createElement('div');
        header.className = 'chat-header';
        
        const profilePicDiv = document.createElement('div');
        profilePicDiv.className = 'profile-pic';
        const img = document.createElement('img');
        img.src = profilePic;
        img.alt = `${username}'s profile`;
        profilePicDiv.appendChild(img);

        const userInfo = document.createElement('div');
        userInfo.className = 'chat-user-info';
        const nameHeading = document.createElement('h4');
        nameHeading.textContent = username;
        const status = document.createElement('p');
        status.className = 'status';
        status.textContent = 'online';
        userInfo.append(nameHeading, status);

        header.append(profilePicDiv, userInfo);

        // Create messages area
        const messagesArea = document.createElement('div');
        messagesArea.className = 'chat-messages';
        messagesArea.id = 'chat-messages';

        // Create input form
        const form = document.createElement('form');
        form.className = 'chat-input';
        form.id = 'message-form';

        const input = document.createElement('input');
        input.type = 'text';
        input.id = 'message-input';
        input.placeholder = 'Type a message...';
        input.required = true;

        const button = document.createElement('button');
        button.type = 'submit';
        button.className = 'btn btn-primary';
        button.textContent = 'Send';

        form.append(input, button);
        form.addEventListener('submit', this.handleMessageSubmit.bind(this));

        // Assemble all parts
        container.append(header, messagesArea, form);
        return container;
    }

    handleMessageSubmit(event) {
        event.preventDefault();
        const input = event.target.querySelector('#message-input');
        const message = input.value.trim();
        
        if (message) {
            // Send message through WebSocket
            this.wsManager.sendMessage({
                type: 'chat_message',
                content: message,
                receiverId: this.currentChatUserId
            });

            // Add message to UI
            const messagesContainer = document.getElementById('chat-messages');
            const messageDiv = document.createElement('div');
            messageDiv.className = 'message sent';
            messageDiv.textContent = message;
            messagesContainer.appendChild(messageDiv);
            
            input.value = '';
            messagesContainer.scrollTop = messagesContainer.scrollHeight;
        }
    }

    clearAllActiveStates() {
        // Remove active class from menu items
        this.elements.menuItems.forEach(item => item.classList.remove('active'));
        
        // Remove active class from messages
        this.elements.messageItems.forEach(msg => msg.classList.remove('active'));
    }

    showHome() {
        this.elements.middleSection.innerHTML = feedsContent;
        this.elements.feedsSection.classList.remove('hidden');
        this.elements.createPostSection.classList.add('hidden');
    }

    showCreatePost() {
        this.elements.feedsSection.classList.add('hidden');
        this.elements.createPostSection.classList.remove('hidden');
    }

    storeFeedsContent() {
        feedsContent = this.elements.middleSection.innerHTML;
    }
}
// Initialize when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    new UIController();
});

class WebSocketManager {
    constructor() {
        this.socket = new WebSocket('ws://your-server:port');
        this.messageHandlers = new Map();
        this.setupSocketListeners();
    }

    setupSocketListeners() {
        this.socket.onopen = () => {
            console.log('WebSocket Connected');
        };

        this.socket.onmessage = (event) => {
            const message = JSON.parse(event.data);
            this.handleIncomingMessage(message);
        };
    }

    handleIncomingMessage(message) {
        const handler = this.messageHandlers.get(message.type);
        if (handler) {
            handler(message);
        }
    }

    sendMessage(messageData) {
        if (this.socket.readyState === WebSocket.OPEN) {
            this.socket.send(JSON.stringify(messageData));
        }
    }
}

class ChatManager {
    constructor(wsManager) {
        this.wsManager = wsManager;
        this.activeChats = new Map();
        this.setupMessageHandlers();
    }

    setupMessageHandlers() {
        this.wsManager.messageHandlers.set('chat_message', (message) => {
            this.handleIncomingChatMessage(message);
        });
    }

    handleIncomingChatMessage(message) {
        const chatId = `chat_${message.senderId}`;
        if (this.activeChats.has(chatId)) {
            const chat = this.activeChats.get(chatId);
            chat.messages.push({
                content: message.content,
                type: 'received',
                timestamp: new Date()
            });
            
            // Trigger UI update
            document.dispatchEvent(new CustomEvent('newMessage', {
                detail: { chatId, message }
            }));
        }
    }
}