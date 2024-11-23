// Add this to your existing main.js
document.addEventListener('DOMContentLoaded', () => {
    const navLinks = document.querySelectorAll('.nav__link');
    
    navLinks.forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            
            // Remove active class from all links
            navLinks.forEach(l => l.classList.remove('active'));
            
            // Add active class to clicked link
            link.classList.add('active');
            
            // Get the path and handle navigation
            const path = link.getAttribute('data-path');
            router.navigate(path);
        });
    });
});

const router = {
    init() {
        this.handleRoute();
        window.addEventListener('popstate', () => this.handleRoute());
    },

    routes: {
        '/': () => {
            hideAllSections();
            document.querySelector('.home').style.display = 'block';
        },
        '/create': () => {
            hideAllSections();
            document.querySelector('.create-post').style.display = 'block';
        },
        '/chat': () => {
            hideAllSections();
            document.querySelector('.chat-container').style.display = 'flex';
        },
        '/logout': () => {
            fetch('/logout', {
                method: 'POST',
                credentials: 'include'
            }).then(() => {
                window.location.href = '/';
            });
        }
    },

    navigate(path) {
        history.pushState(null, '', path);
        this.handleRoute();
    },

    handleRoute() {
        const path = window.location.pathname;
        const handler = this.routes[path] || this.routes['/'];
        handler();
    }
};

function hideAllSections() {
    document.querySelector('.home').style.display = 'none';
    document.querySelector('.create-post').style.display = 'none';
    document.querySelector('.chat-container').style.display = 'none';
}

document.addEventListener('DOMContentLoaded', () => {
    router.init();

    document.querySelectorAll('.nav__link[data-path]').forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            const path = link.getAttribute('data-path');
            
            // Update active state
            document.querySelectorAll('.nav__link').forEach(l => l.classList.remove('active'));
            link.classList.add('active');
            
            router.navigate(path);
        });
    });
});/*==================== SHOW NAVBAR ====================*/
const showMenu = (headerToggle, navbarId) =>{
  const toggleBtn = document.getElementById(headerToggle),
  nav = document.getElementById(navbarId)
  
  // Validate that variables exist
  if(headerToggle && navbarId){
      toggleBtn.addEventListener('click', ()=>{
          // We add the show-menu class to the div tag with the nav__menu class
          nav.classList.toggle('show-menu')
          // change icon
          toggleBtn.classList.toggle('bx-x')
      })
  }
}
showMenu('header-toggle','navbar')

/*==================== LINK ACTIVE ====================*/
const linkColor = document.querySelectorAll('.nav__link')

function colorLink(){
  linkColor.forEach(l => l.classList.remove('active'))
  this.classList.add('active')
}

linkColor.forEach(l => l.addEventListener('click', colorLink))


/*==================== Pagination ====================*/

const minPostsPerPage = 6; // Number of posts per page when there are fewer than 20 posts
const defaultPostsPerPage = 10; // Default posts per page if there are 20 or more posts
let postsPerPage = defaultPostsPerPage; // Default to 10 posts per page
let currentPage = 1;

// Get elements
const postsContainer = document.querySelector('.posts');
const prevButton = document.querySelector('.pagination__prev');
const nextButton = document.querySelector('.pagination__next');
const infoSpan = document.querySelector('.pagination__info');

function showPage(page) {
const totalPages = Math.ceil(posts.length / postsPerPage);

// Validate page number
if (page < 1 || page > totalPages) return;

// Hide all posts
posts.forEach((post, index) => {
  post.style.display = 'none';
  if (index >= (page - 1) * postsPerPage && index < page * postsPerPage) {
    post.style.display = 'block';
  }
});

// Update pagination controls
prevButton.disabled = (page === 1);
nextButton.disabled = (page === totalPages);
infoSpan.textContent = `Page ${page} of ${totalPages}`;
}

// Chat UI Toggle
document.getElementById('chat-link').addEventListener('click', () => {
    const chatContainer = document.querySelector('.chat-container');
    const mainContent = document.querySelector('main');
    
    if (chatContainer.style.display === 'none') {
        chatContainer.style.display = 'flex';
        mainContent.style.display = 'none';
    } else {
        chatContainer.style.display = 'none';
        mainContent.style.display = 'block';
    }
});


document.querySelectorAll('.user-item').forEach(item => {
    item.addEventListener('click', () => {
        // Remove selected class from all users
        document.querySelectorAll('.user-item').forEach(user => {
            user.classList.remove('selected');
        });
        // Add selected class to clicked user
        item.classList.add('selected');
    });
});

class ChatUI {
    constructor() {
        this.ws = new WebSocket('ws://localhost:8080/ws');
        this.setupWebSocket();
        this.chatContainer = document.querySelector('.chat-container');
        this.usersList = document.querySelector('.users-list');
        this.messagesContainer = document.querySelector('.messages-container');
        this.currentChatHeader = document.querySelector('.current-chat-user');
        this.messageInput = document.getElementById('message-text');
        this.sendButton = document.querySelector('.send-message');
        
        this.activeChats = new Map();
        this.currentRecipient = null;
        this.users = [];

        this.loadAvailableUsers();
        this.initializeNewChatButton();
    }

    setupWebSocket() {
        this.ws.onmessage = this.handleMessage.bind(this);
    }

    bindEvents() {
        document.querySelector('.new-chat-btn').addEventListener('click', () => {
            this.showUserSelectModal();
        });
    }

    showUserSelectModal() {
        fetch('/api/users')
            .then(response => response.json())
            .then(users => {
                const modal = document.createElement('div');
                modal.className = 'user-select-modal active';
                modal.innerHTML = `
                    <h3>Select User</h3>
                    <div class="user-select-list">
                        ${users.map(user => `
                            <div class="user-item" data-userid="${user.UserID}">
                                <img src="assets/img/perfil.jpg" alt="${user.Username}" class="user-avatar">
                                <div class="user-info">
                                    <span class="user-name">${user.FirstName} ${user.LastName}</span>
                                    <span class="user-username">@${user.Username}</span>
                                </div>
                            </div>
                        `).join('')}
                    </div>
                `;

                const overlay = document.createElement('div');
                overlay.className = 'modal-overlay active';

                document.body.appendChild(overlay);
                document.body.appendChild(modal);

                modal.querySelectorAll('.user-item').forEach(item => {
                    item.addEventListener('click', () => {
                        const userId = item.dataset.userid;
                        this.startChat(users.find(u => u.UserID === parseInt(userId)));
                        modal.remove();
                        overlay.remove();
                    });
                });

                overlay.addEventListener('click', () => {
                    modal.remove();
                    overlay.remove();
                });
            });
    }

    startChat(user) {
        document.querySelector('.message-input').style.display = 'flex';
        document.querySelector('.current-chat-user').textContent = `Chat with ${user.FirstName} ${user.LastName}`;
        // Add user to chat list and load chat history
    }
}
    class ChatManager {
      constructor() {
          this.ws = new WebSocket('ws://localhost:8080/ws');
          this.chatUI = new ChatUI();
        
          this.ws.onmessage = (event) => {
              const message = JSON.parse(event.data);
              this.chatUI.displayMessage({
                  content: message.content,
                  timestamp: message.timestamp,
                  senderName: message.senderName,
                  sent: false
              });
          };

          // Update the event listener to use this.chatUI
          document.querySelector('.send-message').addEventListener('click', () => {
              this.chatUI.sendMessage();
          });
      }
  }  

  // Initialize chat when DOM loads
  document.addEventListener('DOMContentLoaded', () => {
      router.init();

      // Navigation links
      document.querySelectorAll('.nav__link').forEach(link => {
          link.addEventListener('click', (e) => {
              e.preventDefault();
              const path = link.getAttribute('data-path');
              router.navigate(path);
          });
      });

      const chatManager = new ChatManager();
  });

  // Helper functions
  function getCurrentSelectedUser() {
      const selectedUser = document.querySelector('.user-item.selected');
      return selectedUser ? selectedUser.querySelector('.user-name').textContent : null;
  }

  document.getElementById('message-text').addEventListener('keypress', (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault();
        const chatManager = new ChatManager();
        chatManager.chatUI.sendMessage();
    }
});
