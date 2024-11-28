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
});

// ==================== Router ====================
const router = {
    init() {
        this.handleRoute(); // Handle the current route on page load
        window.addEventListener('popstate', () => this.handleRoute()); // Listen for back/forward navigation
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
            initializeNewChatButton();
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
        console.log('Render posts updated');
        postElement.innerHTML = `
            <a href="#" class="post-link" data-id="${post.postId}">
                <p>${post.postDescription}</p>
            </a>
            <div class="post-meta">
                <span>By: ${post.postId}</span>

                <span>By: ${post.username}</span>
                <span>Likeshbjb: ${post.like} | Dislikes: ${post.dislike}</span>
                <span>Categories: ${post.categoryName.join(', ')}</span>
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
function initializeNewChatButton() {
    const newChatBtn = document.querySelector('.new-chat-btn');
    if (newChatBtn) {
        newChatBtn.addEventListener('click', () => this.showUserSelectModal());
    }
}

// Placeholder for chat-related functions (showUserSelectModal, ChatUI, ChatManager)
// Add chat functions here if needed

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

        // Set the innerHTML directly
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
                            <i class="bx bx-like"></i> <span>${post.like}</span>
                        </button>
                    </form>

                    <form method="post" action="/dislike">
                        <input type="hidden" name="action" value="dislike" />
                        <input type="hidden" name="post_id" value="${post.postId}" />
                        <button type="submit" class="post-button">
                            <i class="bx bx-dislike"></i> <span>${post.dislike}</span>
                        </button>
                    </form>
                    <span class="post-author">by ${post.username}</span>
                </div>
                <div class="post-comments">
                    <h3>Comments</h3>
                    
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
 