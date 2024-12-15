// ==================== DOM Content Loaded ====================
document.addEventListener('DOMContentLoaded', () => {
    // Initialize the router to handle the current route
    router.init();
    fetchPosts();

    // Initialize navigation links for SPA navigation
    document.querySelectorAll('.nav__link[data-path]').forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            const path = link.getAttribute('data-path');

            // Update active state for navigation links
            document.querySelectorAll('.nav__link').forEach(l => l.classList.remove('active'));
            link.classList.add('active');

            // Navigate to the selected route
            router.navigate(path);
        });
    });

    // Initialize chat manager for WebSocket functionality
    const chatManager = new ChatManager();
    chatManager.chatUI.loadAllChats();
    //chatManager.chatUI.loadAllUsers(); // Explicitly call to ensure users are loaded


    document.querySelectorAll('.toggle-auth').forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            const targetPath = `/${e.target.dataset.target}`;
            router.navigate(targetPath);
        });
    });

    const loginForm = document.getElementById('auth-form');

    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            clearErrorMessages('login');

            const formData = new FormData(loginForm);
            const data = new URLSearchParams(formData);

            try {
                const response = await fetch('/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded',
                    },
                    body: data.toString(),
                });

                if (!response.ok) {
                    const result = await response.json();
                    if (result.errors) {
                        showErrorMessages(result.errors, 'login');
                    }
                } else {
                    window.location.href = '/';
                }
            } catch (error) {
                console.error('Error during login:', error);
            }
        });
    }

    function clearErrorMessages(context) {
        document.querySelectorAll(`#${context}-form .error-message`).forEach(element => {
            element.style.display = 'none';
            element.textContent = '';
        });
    }

    function showErrorMessages(errors, context) {
        for (const [field, message] of Object.entries(errors)) {
            const errorElement = document.querySelector(`#error-${field}`);
            if (errorElement) {
                errorElement.textContent = message;
                errorElement.style.display = 'block';
            }
        }
    }

});

// ==================== Router ====================
const router = {
    init() {
        // this.handleRoute(); // Handle the current route on page load
        // window.addEventListener('popstate', () => this.handleRoute()); // Listen for back/forward navigation
        // Check session status first
        fetch('/api/session-status')
            .then(response => {
                if (!response.ok) {
                    document.body.classList.add('no-session');
                    this.navigate('/login');
                }
            });
        this.handleRoute();
    },

    routes: {
        '/': () => {
            hideAllSections();
            document.querySelector('.home').style.display = 'block'; // Display home section
        },
        '/create': () => {
            hideAllSections();
            document.querySelector('.create-post').style.display = 'block'; // Display create post section
        },
        '/chat': () => {
            hideAllSections();
            document.querySelector('.chat-container').style.display = 'flex'; // Show chat container
            //initializeNewChatButton();
        },
        '/post': () => {
            hideAllSections();
            document.querySelector('.post-details').style.display = 'block'; // Show post details section
            const postId = new URLSearchParams(window.location.search).get('id');
            if (postId) {
                fetchPostDetails(postId); // Fetch the details of the specific post
            } else {
                document.querySelector('.post-details').innerHTML = '<p>Invalid post ID.</p>';
            }
        },
        '/logout': () => {
            fetch('/logout', { method: 'POST', credentials: 'include' })
                .then(() => { window.location.href = '/'; });
        },
        '/login': () => {
            hideAllSections();
            document.querySelector('.login-section').style.display = 'block'; // Display login section
        },
        '/register': () => {
            hideAllSections();
            document.querySelector('.register-section').style.display = 'block'; // Display register section
        }
    },

    navigate(path) {
        history.pushState(null, '', path); // Update the browser's URL without reloading the page
        this.handleRoute(); // Handle the new route
    },

    handleRoute() {
        const path = window.location.pathname; // Get the current path
        const handler = this.routes[path] || this.routes['/']; // Default to the home route if path not found
        handler();
    }
};

// ==================== Section Management ====================
function hideAllSections() {
    document.querySelector('.post-details').style.display = 'none';
    document.querySelector('.home').style.display = 'none';
    document.querySelector('.create-post').style.display = 'none';
    document.querySelector('.chat-container').style.display = 'none';
    document.querySelector('.login-section').style.display = 'none';
    document.querySelector('.register-section').style.display = 'none';
}

// ==================== Posts Management ====================
const postsContainer = document.getElementById('posts-container');
let allPosts = [];
let currentPage = 1;
const postsPerPage = 10;

async function fetchPosts(category = 'all') {
    let apiUrl = '/api/posts';
    if (category !== 'all') apiUrl += `?filter-category=${category}`;

    try {
        const response = await fetch(apiUrl, { headers: { Accept: 'application/json' } });
        if (!response.ok) throw new Error('Failed to fetch posts');
        const { posts } = await response.json();
        allPosts = posts;
        renderPosts(); // Render posts for the current page
        updatePaginationControls(); // Update pagination controls
    } catch (error) {
        console.error('Error fetching posts:', error);
        postsContainer.innerHTML = '<p>Failed to load posts. Please try again.</p>';
    }
}



function renderPosts() {
    postsContainer.innerHTML = ''; // Clear existing posts
    const start = (currentPage - 1) * postsPerPage;
    const end = start + postsPerPage;
    const postsToShow = allPosts.slice(start, end);
    postsToShow.forEach((post) => {
        const postElement = document.createElement('article');
        postElement.classList.add('post');
        postElement.setAttribute('data-id', post.postId); // Add a unique identifier
        postElement.innerHTML = `
            <a href="#" class="post-link" data-id="${post.postId}">
                <p>${post.postDescription}</p>
            </a>
            <div class="post-meta">
                <span>By: ${post.username}</span>
                <span>👍 <span class="like-count">${post.like}</span> | 👎 <span class="dislike-count">${post.dislike}</span></span>
                <span>Categories: ${post.categoryName.join(', ')}</span>
            </div>
            <div class="post-actions">
            <button class="action-btn like-btn" data-id="${post.postId}">
                <i class='bx bx-like'></i>
                <span class="like-count">${post.like}</span>
            </button>
            <button class="action-btn dislike-btn" data-id="${post.postId}">
                <i class='bx bx-dislike'></i>
                <span class="dislike-count">${post.dislike}</span>
            </button>
        </div>
        `;
        postsContainer.appendChild(postElement);
    });
    // Add event listeners for post links
    document.querySelectorAll('.post-link').forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            const postId = link.getAttribute('data-id');
            router.navigate(`/post?id=${postId}`);
        });
    });
    // Add event listeners for like and dislike buttons
    document.querySelectorAll('.like-btn').forEach((button) => {
        button.addEventListener('click', () => {
            console.log(`Like clicked for post ID: ${button.dataset.id}`);
            handleLike(button.dataset.id);
        });
    });
    document.querySelectorAll('.dislike-btn').forEach((button) => {
        button.addEventListener('click', () => {
            console.log(`Dislike clicked for post ID: ${button.dataset.id}`);
            handleDislike(button.dataset.id);

        });
    });

}

function updatePaginationControls() {
    const totalPages = Math.ceil(allPosts.length / postsPerPage);
    document.querySelector('.pagination__info').textContent = `Page ${currentPage} of ${totalPages}`;
    document.querySelector('.pagination__prev').disabled = currentPage === 1;
    document.querySelector('.pagination__next').disabled = currentPage === totalPages;
}

// Pagination controls
document.querySelector('.pagination__prev').addEventListener('click', () => {
    if (currentPage > 1) {
        currentPage--;
        renderPosts();
        updatePaginationControls();
    }
});

document.querySelector('.pagination__next').addEventListener('click', () => {
    if (currentPage < Math.ceil(allPosts.length / postsPerPage)) {
        currentPage++;
        renderPosts();
        updatePaginationControls();
    }
});

// ==================== Chat Management ====================

/*==================== CHAT FUNCTIONALITY ====================*/

// Chat UI Class
class ChatUI {
    constructor() {
        this.ws = new WebSocket('ws://localhost:8000/ws'); // WebSocket setup
        // this.setupWebSocket();
        if (!this.ws || this.ws.readyState === WebSocket.CLOSED) {

            this.setupWebSocket();
        } else {
            console.warn('WebSocket is already open or connecting. Skipping setup.');
        }
        this.currentUserId = this.getCurrentUser();
        this.chatContainer = document.querySelector('.chat-container');
        this.usersList = document.querySelector('.users-list');
        this.messagesContainer = document.querySelector('.messages-container');
        this.currentChatHeader = document.querySelector('.current-chat-user');
        this.messageInput = document.getElementById('message-text');
        this.sendButton = document.querySelector('.send-message');

        this.activeChats = new Map();
        this.currentRecipient = null;
        this.users = [];

        // Event listeners
        this.sendButton.onclick = () => this.sendMessage();
        this.messageInput.onkeypress = (e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                this.sendMessage();
            }
        };

        // Load available users and initialize chat features
        this.loadAllUsers();
        this.initializeNewChatButton();
        this.onlineUsersList = document.querySelector('.online-use-list');
        this.loadAllUserswithStatus();

        // Refresh online users every 5 seconds
        // setInterval(() => this.loadAllUserswithStatus(), 5000);
        // setInterval(() => this.loadAllChats(), 5000);

        this.messageOffset = 0;
        this.isLoading = false;
        this.hasMoreMessages = true;

        // Throttle the scroll event
        this.messagesContainer.addEventListener('scroll', this.throttle(this.handleScroll.bind(this), 1000)); // 1-second throttle
    }

    throttle(func, delay) {
        let lastCall = 0;
        return function (...args) {
            const now = new Date().getTime();
            if (now - lastCall >= delay) {
                lastCall = now;
                func(...args);
            }
        };
    }

    handleScroll() {
        // If loading is in progress or no more messages, do nothing
        if (this.isLoading || !this.hasMoreMessages) return;

        // Check if we're near the top (trigger loading more messages)
        if (this.messagesContainer.scrollTop <= 100) {
            this.loadMoreMessages();
        }
    }

    async loadMoreMessages() {
        if (!this.currentRecipient || this.isLoading) return;

        this.isLoading = true;

        try {
            const response = await fetch(
                `/api/chat/history?receiverId=${this.currentRecipient.UserID}&offset=${this.messageOffset}`,
                { credentials: 'include' }
            );
            const messages = await response.json();

            // If there are fewer than 10 messages, stop loading more
            if (messages.length < 10) {
                this.hasMoreMessages = false;
            }

            // Preserve scroll position
            const scrollPos = this.messagesContainer.scrollHeight - this.messagesContainer.scrollTop;

            // Add messages to top
            messages.reverse().forEach(msg => {
                const messageElement = document.createElement('div');
                messageElement.className = `message ${msg.senderId === this.currentUserId ? 'sent' : 'received'}`;
                messageElement.innerHTML = `
                    <div class="message-content">${msg.content}</div>
                    <div class="message-header">
                        <span class="message-time">${new Date(msg.createdAt).toLocaleTimeString()}</span>
                    </div>
                `;
                this.messagesContainer.prepend(messageElement);
            });

            // Maintain scroll position
            this.messagesContainer.scrollTop = this.messagesContainer.scrollHeight - scrollPos;

            // Update the message offset
            this.messageOffset += messages.length;
        } catch (error) {
            console.error('Error loading more messages:', error);
        } finally {
            this.isLoading = false;
        }
    }


    setupWebSocket() {
        this.ws = new WebSocket('ws://localhost:8000/ws');

        this.ws.onopen = () => {
            console.log('WebSocket connection established');
        };


        // this.ws.onmessage = (event) => {
        //     try {
        //         const data = JSON.parse(event.data);

        //         // // Check the type of message received
        //         // switch (data.type) {
        //         //     case 'chat': // For chat messages
        //         //     const message = data.payload;
        //         //     console.log('Received message via WebSocket:', message);
        //         //     this.displayMessage({
        //         //         content: message.content,
        //         //         timestamp: message.timestamp,
        //         //         sent: message.senderId === this.currentUserId,
        //         //     });
        //         //     this.scrollToBottom();
        //         //         break;

        //         //     // case 'presence': // For user online/offline status
        //         //     // console.log('Received presence info via WebSocket:', data.payload);
        //         //     // const presenceInfo = data.payload;
        //         //     // this.updatePresence(presenceInfo.userId, presenceInfo.online);
        //         //     // break;


        //         //     // case 'online_users': // For the list of all online users
        //         //     //     console.log('Online users list received:', data.users);
        //         //     //     this.renderOnlineUsers(data.payload.users);
        //         //     //    // this.renderOnlineUsers(data.users); // Update the UI for online users
        //         //     //     break;

        //         //     default:
        //         //         console.warn('Unknown WebSocket message type:', data.type);
        //         }
        //     } catch (error) {
        //         console.error('Error handling WebSocket message:', error);
        //     }
        // };




        this.ws.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        this.ws.onclose = () => {
            console.warn('WebSocket connection closed. Attempting to reconnect...');
            setTimeout(() => this.setupWebSocket(), 3000); // Retry connection after 3 seconds
        };
    }


    updatePresence(userId, online) {
        const userElement = this.onlineUsersList.querySelector(`.user-item[data-userid="${userId}"]`);
        if (userElement) {
            const statusDot = userElement.querySelector('.user-status');
            statusDot.style.backgroundColor = online ? 'green' : 'gray';
        }

    }



    sendMessage() {
        const content = this.messageInput.value.trim();
        if (content && this.currentRecipient) {
            const message = {
                content: content,
                receiverId: this.currentRecipient.UserID,
                timestamp: new Date().toISOString()
            };

            // Remove the placeholder message if it exists
            const placeholder = this.messagesContainer.querySelector('.placeholder-message');
            if (placeholder) {
                placeholder.remove();
            }

            // Only send via WebSocket
            if (this.ws.readyState === WebSocket.OPEN) {
                this.ws.send(JSON.stringify(message));
            } else {
                console.error('WebSocket is not open. Current state:', this.ws.readyState);
                alert('WebSocket connection is not open. Please refresh the page.');
            }

            console.log('Sending message:', message);

            // Display sent message immediately
            // this.displayMessage({
            //     content: content,
            //     timestamp: message.timestamp,
            //     sent: true
            // });
            //Clear input
            this.messageInput.value = '';
            this.scrollToBottom();
        }
    }
    startChat(user) {
        this.currentRecipient = user;
        this.currentChatHeader.textContent = `Chat with ${user.FirstName} ${user.LastName}`;
        document.querySelector('.message-input').style.display = 'flex';

        // Reset messages and start with initial fetch
        this.messageOffset = 0;
        this.hasMoreMessages = true;
        this.messagesContainer.innerHTML = ''; // Clear existing messages
        this.lastMessageID = 0; // Initially load the latest messages
        this.loading = false; // Flag to prevent multiple simultaneous fetches

        this.loadMoreMessages();
        // Load the first batch of messages (last 10)
        //this.loadMessages();

        // Listen for scroll events to load older messages when scrolling up
        this.messagesContainer.addEventListener('scroll', () => {
            if (!this.loading && this.messagesContainer.scrollTop === 0) {
                this.loadMoreMessages();
            }
        });
    }

    loadMessages() {
        this.loading = true; // Set loading flag to true

        const url = `/api/chat/history?receiverId=${this.currentRecipient.UserID}&before=${this.lastMessageID}`;

        fetch(url, { credentials: 'include' })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Failed to fetch chat history');
                }
                return response.json();
            })
            .then(messages => {
                // If there are no messages or empty messages, stop the loading flag
                if (!messages || messages.length === 0) {
                    this.loading = false;
                    return;
                }

                // Prepend new messages to the messages container
                messages.forEach(msg => this.displayMessage({
                    content: msg.content,
                    timestamp: msg.createdAt,
                    sent: msg.senderId === this.currentUserId,
                }));

                // Update the `lastMessageID` to the ID of the last loaded message
                this.lastMessageID = messages[messages.length - 1].messageID;

                this.scrollToBottom(); // Ensure the latest messages are at the bottom

                // Stop the loading flag after messages are loaded
                this.loading = false;
            })
            .catch(error => {
                console.error('Error loading chat history:', error);
                this.messagesContainer.innerHTML = '<p>Failed to load chat history. Please try again later.</p>';
                this.loading = false; // Stop loading flag on error
            });
    }

    // displayMessage({ content, timestamp, sent }) {
    //     const messageElement = document.createElement('div');
    //     messageElement.className = sent ? 'sent-message' : 'received-message';
    //     messageElement.innerHTML = `
    //         <div class="message-content">${content}</div>
    //         <div class="message-timestamp">${new Date(timestamp).toLocaleTimeString()}</div>
    //     `;
    //     this.messagesContainer.appendChild(messageElement);
    // }

    scrollToBottom() {
        this.messagesContainer.scrollTop = this.messagesContainer.scrollHeight;
    }


    addUserToList(user, prepend = false) {
        // Check if the user already exists in the list
        const existingUser = this.usersList.querySelector(`.user-item[data-userid="${user.UserID}"]`);
        if (existingUser) {


            if (prepend) {
                this.usersList.removeChild(existingUser);
                this.usersList.prepend(existingUser);
            }
            console.warn(`User with ID ${user.UserID} is already in the list.`);
            return;
        }

        // Create the user element
        const userElement = document.createElement('div');
        userElement.className = 'user-item';
        userElement.dataset.userid = user.UserID;
        userElement.innerHTML = `
            <div class="user-info">
                <span class="user-name">${user.FirstName} ${user.LastName}</span>
                <span class="user-username">@${user.Username}</span>
                <span class="user-status" style="width: 10px; height: 10px; border-radius: 50%; background-color: gray; display: inline-block; margin-left: 10px;"></span>
            </div>
        `;
        userElement.addEventListener('click', () => this.startChat(user));

        // Add the user to the top or bottom of the list
        if (prepend) {
            this.usersList.prepend(userElement); // Place at the top
        } else {
            this.usersList.appendChild(userElement); // Place at the bottom
        }
    }



    loadAllUsers() {
        fetch('/api/users')
            .then(response => response.json())
            .then(users => {
                this.users = users;
                users.forEach(user => this.addUserToList(user, true)); // Prepend to the list
            });

    }

    displayMessage(message) {
        console.log('Displaying message:sads', message);
        // Remove the placeholder message if it exists
        const placeholder = this.messagesContainer.querySelector('.placeholder-message');
        if (placeholder) {
            placeholder.remove();
        }

        const messageElement = document.createElement('div');
        messageElement.className = `message ${message.sent ? 'sent' : 'received'}`;
        console.log(message.sent);

        messageElement.innerHTML = `
            <div class="message-content">${message.content}</div>
            <div class="message-header">
<span class="message-time">${new Date(message.timestamp).toLocaleString()}</span>
            </div>
        `;

        this.messagesContainer.appendChild(messageElement);

        // Move the user to the top of the user list
        if (this.currentRecipient) {
            this.addUserToList(this.currentRecipient, true); // Move user to top
        }


    }


    getCurrentUser() {
        fetch('/api/session-status', {
            credentials: 'include'
        })
            .then(response => response.json())
            .then(data => {
                this.currentUserId = data.userId;
            })
            .catch(error => console.error('Error getting current user:', error));
    }
    scrollToBottom() {
        this.messagesContainer.scrollTop = this.messagesContainer.scrollHeight;
    }

    formatTime(timestamp) {
        return new Date(timestamp).toLocaleTimeString([], {
            hour: '2-digit',
            minute: '2-digit'
        });
    }

    initializeNewChatButton() {
        const newChatBtn = document.querySelector('.new-chat-btn');
        if (newChatBtn) {
            newChatBtn.addEventListener('click', () => this.showUserSelectModal());
        }
    }
    loadAllChats() {
        fetch('/api/user-chat', { credentials: 'include' })
            .then(response => response.json())
            .then(users => {
                this.usersList.innerHTML = ''; // Clear the existing user list
                users.forEach(user => this.addUserToList(user, false));

                // Initial fetch of online statuses
                // fetchOnlineUsers();
            })
            .catch(error => console.error('Error loading chats:', error));
    }


    loadAllUserswithStatus() {
        fetch('/api/all-online-users', { credentials: 'include' })
            .then(response => response.json())
            .then(users => {
                console.log('Raw API Response:', users); // Log the exact response
                this.onlineUsersList.innerHTML = ''; // Clear the online users list

                users.forEach(user => {
                    console.log('User Object:', user); // Log each user object to inspect properties
                    const statusColor = user.online ? 'green' : 'gray';
                    const userElement = document.createElement('div');
                    userElement.className = 'user-item';
                    userElement.dataset.userid = user.UserID;

                    // Handle missing properties
                    const fullName = `${user.first_name || ''} ${user.last_name || ''}`.trim();
                    const username = user.username;

                    userElement.innerHTML = `
                   <div class="user-info">
                    <span class="user-name">${fullName}</span>
                    <span class="user-username">@${username}</span>
                    <span class="user-status" 
                        style="width: 10px; height: 10px; border-radius: 50%; background-color: ${statusColor}; display: inline-block; margin-left: 10px;">
                    </span>
                   </div>
                      `;
                    this.onlineUsersList.appendChild(userElement);
                });
            })
            .catch(error => console.error('Error fetching online users:', error));

    }

    showUserSelectModal() {
        fetch('/api/users')
            .then(response => response.json())
            .then(users => {
                const modal = this.createModal(users);
                document.body.appendChild(modal.overlay);
                document.body.appendChild(modal.modal);
            });
    }

    createModal(users) {

        // Sort users alphabetically by their first and last name
        users.sort((a, b) => {
            const fullNameA = `${a.FirstName} ${a.LastName}`.toLowerCase();
            const fullNameB = `${b.FirstName} ${b.LastName}`.toLowerCase();
            return fullNameA.localeCompare(fullNameB);
        });

        const modal = document.createElement('div');
        modal.className = 'user-select-modal active';

        const overlay = document.createElement('div');
        overlay.className = 'modal-overlay active';

        modal.innerHTML = `
            <h3>Select User</h3>
            <div class="user-select-list">
                ${users.map(user => `
                    <div class="user-item" data-userid="${user.UserID}">
                        <div class="user-info">
                            <span class="user-name">${user.FirstName} ${user.LastName}</span>
                            <span class="user-username">@${user.Username}</span>
                        </div>
                    </div>
                `).join('')}
            </div>
        `;
        modal.querySelectorAll('.user-item').forEach(item => {
            item.onclick = () => {
                const userId = item.dataset.userid;
                const selectedUser = this.users.find(u => u.UserID === parseInt(userId));
                this.startChat(selectedUser);
                console.log('Selected User:', selectedUser);
                this.addUserToList(selectedUser, true);
                console.log('Selected User:', selectedUser);
                modal.remove();
                overlay.remove();
            };
        });
        overlay.onclick = () => {
            modal.remove();
            overlay.remove();
        };
        return { modal, overlay };
    }
}
function fetchOnlineUsers() {
    fetch('/api/online-users') // Backend endpoint that returns a list of online user IDs
        .then(response => response.json())
        .then(onlineUsers => {
            document.querySelectorAll('.user-item').forEach(userElement => {
                const userId = parseInt(userElement.dataset.userid);
                const statusDot = userElement.querySelector('.user-status');

                if (onlineUsers.includes(userId)) {
                    statusDot.style.backgroundColor = 'green'; // Online
                } else {
                    statusDot.style.backgroundColor = 'gray'; // Offline
                }
            });
        })
        .catch(error => console.error('Error fetching online users:', error));
}


// Chat Manager Class
class ChatManager {
    constructor() {
        this.ws = new WebSocket('ws://localhost:8000/ws');
        this.chatUI = new ChatUI();

        this.ws.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);

                // Check the type of message received
                switch (data.type) {
                    case 'chat': // For chat messages
                        const message = data.payload;
                        console.log('Received message via WebSocket:', message);

                        // Only display the message if it's for the current chat
                        if (
                            this.chatUI.currentRecipient &&
                            (message.senderId === this.chatUI.currentRecipient.UserID ||
                                message.receiverId === this.chatUI.currentRecipient.UserID)
                        ) {
                            this.chatUI.displayMessage({
                                content: message.content,
                                timestamp: message.timestamp,
                                sent: message.senderId === this.chatUI.currentUserId,
                            });
                            this.chatUI.scrollToBottom();
                          
                        } else {
                            // Optionally show a notification for other incoming chats
                            notifyNewMessage(message.senderId, message.content);
                        }
                        break;
                    case 'presence': // For user online/offline status
                        const { userId, online } = data.payload;
                        this.chatUI.updatePresence(userId, online);
                        //  updateUserPresence(userId, online);

                        break;


                    // case 'online_users': // For the list of all online users
                    // this.renderOnlineUsers(data.payload.users);

                    // console.log('Online users list received:', data.users);
                    //this.chatUI.renderOnlineUsers(data.users); // Update the UI for online users
                    //  break;

                    default:
                        console.warn('Unknown WebSocket message type:', data.type);
                }
            } catch (error) {
                console.error('Error handling WebSocket message:', error);
            }
        };

    }
}

function notifyNewMessage(senderId, content) {
    const sender = this.chatUI.users.find((user) => user.UserID === senderId);
    const senderName = sender ? `${sender.FirstName} ${sender.LastName}` : "Unknown";

    // Show a notification banner
    const notification = document.createElement("div");
    notification.className = "notification";
    notification.textContent = `New message from ${senderName}: ${content}`;
    document.body.appendChild(notification);

    // Automatically remove the notification after 5 seconds
    setTimeout(() => notification.remove(), 5000);
}


function updateUserPresence(userId, isOnline) {
    const userElement = document.querySelector(`.user-item[data-userid="${userId}"]`);
    if (userElement) {
        const statusDot = userElement.querySelector('.user-status');
        if (statusDot) {
            statusDot.style.backgroundColor = isOnline ? 'red' : 'gray';
        }
    }
}
// ==================== Navigation Bar Management ====================
const showMenu = (headerToggle, navbarId) => {
    const toggleBtn = document.getElementById(headerToggle);
    const nav = document.getElementById(navbarId);

    if (headerToggle && navbarId) {
        toggleBtn.addEventListener('click', () => {
            nav.classList.toggle('show-menu');
            toggleBtn.classList.toggle('bx-x');
        });
    }
};
showMenu('header-toggle', 'navbar');

const linkColor = document.querySelectorAll('.nav__link');

function colorLink() {
    linkColor.forEach(l => l.classList.remove('active'));
    this.classList.add('active');
}

linkColor.forEach(l => l.addEventListener('click', colorLink));

// ==================== Helper Functions ====================
function getCurrentSelectedUser() {
    const selectedUser = document.querySelector('.user-item.selected');
    return selectedUser ? selectedUser.querySelector('.user-name').textContent : null;
}


async function fetchPostDetails(postId) {
    try {
        // Fetch the post details from the server
        const response = await fetch(`/post?id=${postId}`);
        if (!response.ok) throw new Error('Failed to fetch post details');

        // Parse the response JSON
        const post = await response.json();

        // Get the container for post details
        const postContainer = document.querySelector('.post-details');

        // Check if there are comments, and generate the comments section accordingly
        const commentsHTML = post.comments && post.comments.length > 0
            ? post.comments.map(comment => `
                <li class="comment" data-comment-id="${comment.commentId}">
                    <div class="comment-content">
                        <span class="comment-author">${comment.username}</span>
                        <p class="comment-text">${comment.commentText}</p>
                    </div>
                    <div class="comment-actions">
                        <button class="action-btn comment-like-btn" data-id="${comment.commentId}">
                            👍 <span class="comment-like-count">${comment.likes}</span>
                        </button>
                        <button class="action-btn comment-dislike-btn" data-id="${comment.commentId}">
                            👎 <span class="comment-dislike-count">${comment.dislikes}</span>
                        </button>
                    </div>
                </li>
            `).join('')
            : '<p>No comments yet. Be the first to comment!</p>';

        // Set the innerHTML for the post details, including the comments section
        postContainer.innerHTML = `
            <article class="post">
                <pre><p class="post-description">${post.postDescription}</p></pre>
                <div class="post-category">
                    <span>Categories: ${post.categoryName.join(', ')}</span>
                </div>
                <div class="post-info">
                    <span class="post-author">by ${post.username}</span>
                </div>
                <div class="post-comments">
                    <h3>Comments</h3>
                    <ul class="comments-list">${commentsHTML}</ul>
                    <form onsubmit="event.preventDefault(); submitComment(${post.postId});">
                        <textarea id="comment-input" placeholder="Write your comment..." required></textarea>
                        <button type="submit" class="submit-comment-button">Submit</button>
                    </form>
                </div>
            </article>
        `;

        // Add event listeners for comment like/dislike buttons
        document.querySelectorAll('.comment-like-btn').forEach(button => {
            button.addEventListener('click', () => {
                const commentId = button.getAttribute('data-id');
                console.log("like clicked for comment ID:", commentId);

                // Ensure the comment ID is passed to the function
                handleCommentLike(commentId);
            });
        });

        document.querySelectorAll('.comment-dislike-btn').forEach((button) => {
            button.addEventListener('click', (e) => {
                e.preventDefault(); // Prevent default browser behavior
                const commentId = button.getAttribute('data-id');
                console.log(`Dislike clicked for comment ID: ${commentId}`);
                handleCommentDislike(commentId);
            });
        });



    } catch (error) {
        console.error('Error fetching post details:', error);
        const postContainer = document.querySelector('.post-details');
        postContainer.innerHTML = '<p>Failed to load post details.</p>';
    }
}


async function submitComment(postId) {
    const commentInput = document.querySelector('#comment-input');
    const commentText = commentInput.value.trim();

    if (!commentText) {
        alert('Comment cannot be empty.');
        return;
    }

    try {
        const response = await fetch('/api/comment', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ post_id: postId, comment: commentText }),
            credentials: 'include', // Include cookies for session handling
        });

        if (!response.ok) {
            const error = await response.json();
            alert(`Error: ${error.error}`);
            return;
        }

        const result = await response.json();
        alert(result.message);

        // Add the new comment to the comment list dynamically
        const commentsContainer = document.querySelector('.post-comments ul');
        const newComment = document.createElement('li');
        newComment.classList.add('comment');
        newComment.innerHTML = `
            <div class="comment-content">
                <span class="comment-author">You</span>
                <p class="comment-text">${commentText}</p>
            </div>
            <div class="comment-actions">
                <button class="action-btn comment-like-btn" data-id="${result.commentId}">
                    👍 <span class="comment-like-count">0</span>
                </button>
                
                <button class="action-btn comment-dislike-btn" data-id="${result.commentId}">
                    👎 <span class="comment-dislike-count">0</span>
                </button>

                 
            </div>
        `;
        commentsContainer.appendChild(newComment);

        // Reapply event listeners
        newComment.querySelector('.comment-like-btn').addEventListener('click', () => {
            handleCommentLike(result.commentId);
        });
        newComment.querySelector('.comment-dislike-btn').addEventListener('click', () => {
            handleCommentDislike(result.commentId);
        });

        // Clear input field
        commentInput.value = '';
    } catch (error) {
        console.error('Error submitting comment:', error);
        alert('An error occurred while submitting your comment.');
    }
}


async function handleCommentLike(commentId) {
    try {
        const response = await fetch(`/api/comments/like`, { // Updated endpoint
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ comment_id: parseInt(commentId) }), // Send comment_id in the request body
            credentials: 'include', // Include cookies for session handling
        });

        if (!response.ok) {
            const errorText = await response.text();
            console.error('Error response for comment like:', errorText);
            alert(`Error: ${errorText}`);
            return;
        }

        const data = await response.json();
        console.log('Updated comment after like:', data);

        // Update the UI for the comment
        updateCommentUI(data);
    } catch (error) {
        console.error('Error handling comment like:', error);
        alert('An error occurred while liking the comment.');
    }
}

async function handleCommentDislike(commentId) {
    console.log(`Sending dislike request for comment ID: ${commentId}`);
    try {
        const response = await fetch('/api/comments/dislike', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ comment_id: parseInt(commentId) }), // Ensure comment_id is a number
            credentials: 'include',
        });

        if (!response.ok) {
            const errorText = await response.text();
            console.error('Error response for comment dislike:', errorText);
            alert(`Error: ${errorText}`);
            return;
        }

        const data = await response.json();
        console.log('Comment disliked successfully:', data);
        updateCommentUI(data);
    } catch (error) {
        console.error('Error during dislike operation:', error);
    }
}



function updateCommentUI(commentData) {
    const commentElement = document.querySelector(`.comment[data-comment-id="${commentData.commentId}"]`);
    if (commentElement) {
        commentElement.querySelector('.comment-like-count').textContent = commentData.likes;
        commentElement.querySelector('.comment-dislike-count').textContent = commentData.dislikes;
    }
}


async function handleLike(postId) {
    console.log(`Sending like request for post ID: ${postId}`);
    try {
        const response = await fetch(`/api/like`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ post_id: parseInt(postId) }),
            credentials: 'include',
        });

        if (!response.ok) {
            const errorText = await response.text();
            console.error('Error response for like:', errorText);
            alert(`Error: ${errorText}`);
            return;
        }

        const data = await response.json();
        console.log('Updated post after like:', data);

        // Update the UI with the new like count
        updatePostUI(data);

        // If the post was previously disliked, update the dislike count
        const dislikeBtn = document.querySelector(`.dislike-btn[data-id="${postId}"]`);
        if (dislikeBtn) {
            const dislikeCountElem = dislikeBtn.querySelector('.dislike-count');
            if (parseInt(dislikeCountElem.textContent) > 0) {
                dislikeCountElem.textContent = parseInt(dislikeCountElem.textContent) - 1;
            }
        }
    } catch (error) {
        console.error('Error handling like:', error);
        alert('An error occurred while liking the post.');
    }
}

async function handleDislike(postId) {
    console.log(`Sending dislike request for post ID: ${postId}`);
    try {
        const response = await fetch('/api/dislike-post', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ post_id: parseInt(postId) }),
            credentials: 'include',
        });

        if (!response.ok) {
            const errorText = await response.text();
            console.error('Error response for dislike:', errorText);
            alert(`Error: ${errorText}`);
            return;
        }

        const data = await response.json();
        console.log('Updated post after dislike:', data);

        // Update the UI with the new dislike count
        updatePostUI(data);

        // If the post was previously liked, update the like count
        const likeBtn = document.querySelector(`.like-btn[data-id="${postId}"]`);
        if (likeBtn) {
            const likeCountElem = likeBtn.querySelector('.like-count');
            if (parseInt(likeCountElem.textContent) > 0) {
                likeCountElem.textContent = parseInt(likeCountElem.textContent) - 1;
            }
        }
    } catch (error) {
        console.error('Error handling dislike:', error);
        alert('An error occurred while disliking the post.');
    }
}


function updatePostUI(postData) {
    const postElement = document.querySelector(`.post[data-id="${postData.postId}"]`);
    if (postElement) {
        postElement.querySelector('.like-count').textContent = postData.like;
        postElement.querySelector('.dislike-count').textContent = postData.dislike;
    }
}

