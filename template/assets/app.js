class UIController {
    constructor() {
        this.elements = {
            menuItems: document.querySelectorAll('.sidebar .menu-item'),
            feedsSection: document.querySelector('.feeds'),
            createPostSection: document.querySelector('.create-post'),
            messageItems: document.querySelectorAll('.messages .message'),
            middleSection: document.querySelector('.middle')
        };
        
        // Store initial feeds content
        this.feedsContent = this.elements.middleSection.innerHTML;
        
        // Bind event handlers
        this.bindMenuClicks();
        this.bindMessageClicks();
    }

    bindMenuClicks() {
        this.elements.menuItems.forEach(item => {
            item.addEventListener('click', () => {
                const menuText = item.querySelector('h3').textContent.trim();
                
                // Update active menu item
                this.elements.menuItems.forEach(mi => mi.classList.remove('active'));
                item.classList.add('active');
                
                switch(menuText) {
                    case 'Home':
                        this.showFeeds();
                        break;
                    case 'Create Post':
                        this.showCreatePost();
                        break;
                }
            });
        });
    }

    bindMessageClicks() {
        this.elements.messageItems.forEach(item => {
            item.addEventListener('click', () => {
                const username = item.querySelector('.message-body h5').textContent;
                const profilePic = item.querySelector('.profile-pic img').src;
                this.showChat(username, profilePic);
            });
        });
    }

    showFeeds() {
        this.elements.middleSection.innerHTML = this.feedsContent;
        this.elements.feedsSection.style.display = 'block';
        this.elements.createPostSection.style.display = 'none';
    }

    showCreatePost() {
        this.elements.middleSection.innerHTML = '';
        this.elements.feedsSection.style.display = 'none';
        this.elements.createPostSection.style.display = 'block';
        this.elements.middleSection.appendChild(this.elements.createPostSection);
    }

    showChat(username, profilePic) {
        const chatInterface = this.createChatInterface(username, profilePic);
        this.elements.middleSection.innerHTML = '';
        this.elements.feedsSection.style.display = 'none';
        this.elements.createPostSection.style.display = 'none';
        this.elements.middleSection.appendChild(chatInterface);
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
}

// Initialize when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    new UIController();
});

class WebSocketManager {
    constructor() {
        // Use localhost and correct port for development
        this.socket = new WebSocket('ws://localhost:8080/ws');
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
}


