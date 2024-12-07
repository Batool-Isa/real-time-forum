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
            handleLikeDislike(button.dataset.id, 'dislike');
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
    this.getCurrentUser();
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
        this.loadOnlineUsers(); // Load online users on initialization

        // Refresh online users every 5 seconds
        setInterval(() => this.loadOnlineUsers(), 5000);
    }

    setupWebSocket() {
        this.ws = new WebSocket('ws://localhost:8000/ws');
    
        this.ws.onopen = () => {
            console.log('WebSocket connection established');
        };
    
        this.ws.onmessage = (event) => {
            const message = JSON.parse(event.data);
            console.log("sdgds"+event.data);

            console.log('Received message via WebSocket:', message);
            this.displayMessage({
                content: message.content,
                timestamp: message.timestamp,
                sent: false
            });
            this.scrollToBottom();
        };
    
        this.ws.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
    
        this.ws.onclose = () => {
            console.warn('WebSocket connection closed. Attempting to reconnect...');
            setTimeout(() => this.setupWebSocket(), 3000); // Retry connection after 3 seconds
        };
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

            // Display sent message immediately
            this.displayMessage({
                content: content,
                timestamp: message.timestamp,
                sent: true
            });
            // Clear input
            this.messageInput.value = '';
            this.scrollToBottom();
        }
    }
    
    startChat(user) {
        this.currentRecipient = user;
        this.currentChatHeader.textContent = `Chat with ${user.FirstName} ${user.LastName}`;
        document.querySelector('.message-input').style.display = 'flex';
    
        fetch(`/api/chat/history?receiverId=${user.UserID}`, { credentials: 'include' })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Failed to fetch chat history');
                }
                return response.json();
            })
            .then(messages => {
                this.messagesContainer.innerHTML = ''; // Clear existing messages
    
                // Check if the response contains messages
                if (!messages || messages.length === 0) {
                    const placeholderMessage = document.createElement('p');
                    placeholderMessage.className = 'placeholder-message';
                    placeholderMessage.textContent = 'This is the start of your conversation. Say hi!';
                    placeholderMessage.style.textAlign = 'center';
                    placeholderMessage.style.color = 'gray';
                    this.messagesContainer.appendChild(placeholderMessage);
                } else {
                    // Display existing messages
                    messages.forEach(msg => this.displayMessage({
                        content: msg.content,
                        timestamp: msg.createdAt,
                        sent: msg.senderId === this.currentUserId,
                    }));
                }
    
                this.scrollToBottom();
            })
            .catch(error => {
                console.error('Error loading chat history:', error);
                this.messagesContainer.innerHTML = '<p>Failed to load chat history. Please try again later.</p>';
            });
    }
    
    
    addUserToList(user) {

         // Check if the user already exists in the list
    const existingUser = this.usersList.querySelector(`.user-item[data-userid="${user.UserID}"]`);
    if (existingUser) {
        console.warn(`User with ID ${user.UserID} is already in the list.`);
        return; // Do not add the user again
    }

    // If the user doesn't exist, add them to the list
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

     //   this.usersList.appendChild(userElement);
          // Prepend the new user element to the top of the list
    this.usersList.prepend(userElement);
    }
    
    

    loadAllUsers() {
        fetch('/api/users')
            .then(response => response.json())
            .then(users => {
                this.users = users;
                users.forEach(user => this.addUserToList(user));
            });

    }
    
    displayMessage(message) {
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
                <span class="message-time">${new Date(message.timestamp).toLocaleTimeString()}</span>
            </div>
        `;
        this.messagesContainer.appendChild(messageElement);
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
                users.forEach(user => this.addUserToList(user));
    
                // Initial fetch of online statuses
                fetchOnlineUsers();
            })
            .catch(error => console.error('Error loading chats:', error));
    }
    
    
    loadOnlineUsers() {
        fetch('/api/all-online-users', { credentials: 'include' })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Failed to fetch online users');
                }
                return response.json();
            })
            .then(onlineUsers => {
                console.log('Online users fetched:', onlineUsers); // Log the response
                if (!onlineUsers || !Array.isArray(onlineUsers)) {
                    console.error('Invalid response format:', onlineUsers);
                    return;
                }
                this.renderOnlineUsers(onlineUsers);
            })
            .catch(error => {
                console.error('Error fetching online users:', error);
            });
    }
    
    
    renderOnlineUsers(users) {
        this.onlineUsersList.innerHTML = ''; // Clear the current list
    
        if (!users || users.length === 0) {
            console.warn('No users online to display.');
            const emptyMessage = document.createElement('p');
            emptyMessage.textContent = 'No users are currently online.';
            emptyMessage.style.textAlign = 'center';
            this.onlineUsersList.appendChild(emptyMessage);
            return;
        }
    
        console.log('Rendering users:', users); // Log the valid users array
        users.forEach(user => {
            const userElement = document.createElement('div');
            userElement.className = 'user-item';
            userElement.dataset.userid = user.UserID;
    
            userElement.innerHTML = `
                <div class="user-info">
                    <span class="user-name">${user.FirstName} ${user.LastName}</span>
                    <span class="user-username">@${user.Username}</span>
                    <span class="user-status" style="width: 10px; height: 10px; border-radius: 50%; background-color: green; display: inline-block; margin-left: 10px;"></span>
                </div>
            `;
            userElement.addEventListener('click', () => this.startChat(user));
            this.onlineUsersList.appendChild(userElement);
        });
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
                this.addUserToList(selectedUser);
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

// Fetch online users every 5 seconds
setInterval(fetchOnlineUsers, 5000);

// Chat Manager Class
class ChatManager {
    constructor() {
        this.ws = new WebSocket('ws://localhost:8000/ws');
        this.chatUI = new ChatUI();

        this.ws.onmessage = (event) => {
            const message = JSON.parse(event.data);
            console.log("sdgds11111111111"+event.data);
            const placeholder = this.messagesContainer.querySelector('.placeholder-message');
            if (placeholder) {
                placeholder.remove();
            }

            this.chatUI.displayMessage({
                content: message.content,
                timestamp: message.createdAt,
                sent: false
            });
        };
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
                <li>
                    <p>${comment.commentText}</p>
                    <span>By: ${comment.username}</span>
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
                <form method="post" action="/like">
                    <input type="hidden" name="action" value="like" />
                    <input type="hidden" name="post_id" value="${post.postId}" />
                    <button type="submit" class="post-button">
                        👍 <span>${post.like}</span>
                    </button>
                </form>
    
                <form method="post" action="/dislike">
                    <input type="hidden" name="action" value="dislike" />
                    <input type="hidden" name="post_id" value="${post.postId}" />
                    <button type="submit" class="post-button">
                        👎 <span>${post.dislike}</span>
                    </button>
                </form>
                <span class="post-author">by ${post.username}</span>
            </div>
            <div class="post-comments">
                <h3>Comments</h3>
                <ul>${commentsHTML}</ul>
                <form onsubmit="event.preventDefault(); submitComment(${post.postId});">
                    <textarea id="comment-input" placeholder="Write your comment..." required></textarea>
                    <button type="submit">Submit</button>
                </form>
            </div>
        </article>
    `;
    
    } catch (error) {
        // Log the error and display a fallback error message
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

        // Dynamically add the comment to the comment section
        const commentsContainer = document.querySelector('.post-comments ul');
        const newComment = document.createElement('li');
        newComment.innerHTML = `
            <p>${commentText}</p>
            <span>By: ${result.username}</span>
        `;
        commentsContainer.appendChild(newComment);

        // Clear the input field
        commentInput.value = '';
    } catch (error) {
        console.error('Error submitting comment:', error);
        alert('An error occurred while submitting your comment.');
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

        console.log('Raw response for like:', response);

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
    } catch (error) {
        console.error('Error handling like:', error);
        alert('An error occurred while liking the post.');
    }
}

function updatePostUI(postData) {
    const postElement = document.querySelector(`.post[data-id="${postData.postId}"]`);
    if (postElement) {
        postElement.querySelector('.like-count').textContent = postData.like;
        postElement.querySelector('.dislike-count').textContent = postData.dislike;
    }
}


async function handleLikeDislike(postId, action) {
    console.log("Sending request for:", action, "with post ID:", postId); // Debug log
    try {
        const response = await fetch(`/api/${action}`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ post_d: parseInt(postId) }),
            credentials: 'include',
        });

        console.log('Raw response:', response);

        if (!response.ok) {
            const errorText = await response.text(); // Read text response in case of errors
            console.error('Error response:', errorText);
            alert(`Error: ${errorText}`);
            return;
        }

        // Handle empty or non-JSON responses
        const contentType = response.headers.get('Content-Type');
        if (contentType && contentType.includes('application/json')) {
            const data = await response.json();
            console.log('Updated post:', data);
            updatePostUI(data); // Update the UI with the new post data
        } else {
            console.log('No JSON response body');
        }
    } catch (error) {
        console.error('Error handling like/dislike:', error);
        alert('An error occurred while updating the post.');
    }
}
