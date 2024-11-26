// ==================== ROUTER SETUP ====================
// Router for handling client-side navigation
const router = {
    init() {
        console.log("Router initialized"); // Debug log
        this.handleRoute();
        // Handle browser navigation (back/forward buttons)
        window.addEventListener("popstate", () => this.handleRoute());
    },

    routes: {
        // Root route: Show login or redirect to posts if logged in
        "/": () => {
            const isLoggedIn = document.body.dataset.isLoggedIn === "true";
            console.log("User logged in:", isLoggedIn); // Debug log
            if (isLoggedIn) {
                router.navigate("/api/posts");
            } else {
                hideAllSections();
                document.querySelector(".login-section").style.display = "block";
                console.log("Displaying login section"); // Debug log
            }
        },

        // Posts route: Fetch and render posts
        "/api/posts": () => {
            console.log("Route '/api/posts' matched"); // Debug log
            hideAllSections();
            const homeSection = document.querySelector(".home");
            if (!homeSection) {
                console.error("Home section not found in DOM");
                return;
            }
            homeSection.style.display = "block";
            console.log("Home section displayed:", homeSection.style.display); // Debug log
            //fetchPosts(); // Fetch and render posts dynamically
        },
        // Route to create a post
        "/create": () => {
            hideAllSections();
            document.querySelector(".create-post").style.display = "block";
            console.log("Create post section displayed"); // Debug log
        },

        // Chat route
        "/chat": () => {
            hideAllSections();
            document.querySelector(".chat-container").style.display = "flex";
            console.log("Chat section displayed"); // Debug log
        },

        // Login route
        "/login": () => {
            hideAllSections();
            document.querySelector(".login-section").style.display = "block";
            console.log("Login section displayed"); // Debug log
        },

        // Register route
        "/register": () => {
            hideAllSections();
            document.querySelector(".register-section").style.display = "block";
            console.log("Register section displayed"); // Debug log
        },

        // Logout route
        "/logout": () => {
            console.log("Logging out...");
            fetch("/logout", {
                method: "POST",
                credentials: "include",
            }).then(() => {
                window.location.href = "/";
            });
        },
    },

    // Navigate to a specific path and handle the route
    navigate(path) {
        console.log("Navigating to:", path); // Debug log
        history.pushState(null, "", path);
        this.handleRoute();
    },

    // Handle the current route
    handleRoute() {
        const path = window.location.pathname;
        console.log("Handling route for:", path); // Debug log
        const handler = this.routes[path] || this.routes["/"];
        if (handler) {
            handler();
        } else {
            console.error("No handler found for route:", path); // Debug log
        }
    },
};

// ==================== HIDE SECTIONS ====================
// Function to hide all sections before rendering the desired one
function hideAllSections() {
    console.log("Hiding all sections"); // Debug log
    document.querySelector(".home").style.display = "none";
    document.querySelector(".create-post").style.display = "none";
    document.querySelector(".chat-container").style.display = "none";
    document.querySelector(".login-section").style.display = "none";
    document.querySelector(".register-section").style.display = "none";
}

// ==================== FETCH AND RENDER POSTS ====================
// Fetch posts from the backend and render them dynamically
async function fetchPosts(category = "all") {
    console.log("fetchPosts called with category:", category); // Debug log
    try {
        const response = await fetch(`${category}`, {
            headers: { Accept: "application/json" },
        });

        if (!response.ok) {
            throw new Error("Failed to fetch posts");
        }

        const { posts } = await response.json();
        console.log("Posts fetched successfully:", posts); // Debug log
        renderPosts(posts); // Render the fetched posts
    } catch (error) {
        console.error("Error fetching posts:", error); // Debug log
        const postsContainer = document.getElementById("posts-container");
        postsContainer.innerHTML = "<p>Failed to load posts. Please try again later.</p>";
    }
}

// Render posts dynamically into the container
function renderPosts(posts) {
    console.log("Rendering posts..."); // Debug log
    const postsContainer = document.getElementById("posts-container");
    postsContainer.innerHTML = ""; // Clear existing posts

    if (posts.length === 0) {
        postsContainer.innerHTML = "<p>No posts available.</p>";
        console.log("No posts available to render"); // Debug log
        return;
    }

    posts.forEach((post) => {
        const postElement = document.createElement("article");
        postElement.className = "post";
        postElement.innerHTML = `
            <p>${post.postDescription || "No description available."}</p>
            <div class="post-meta">
                <span>By: ${post.username}</span>
                <span>Likes: ${post.like} | Dislikes: ${post.dislike}</span>
                <span>Categories: ${post.categoryName.join(", ")}</span>
            </div>
        `;
        postsContainer.appendChild(postElement);
        console.log("Post rendered:", post); // Debug log
    });
}

// ==================== PAGINATION ====================
// Default pagination variables
const postsPerPage = 10;
let currentPage = 1;

// Update pagination controls based on total posts
function updatePaginationControls(allPosts) {
    console.log("Updating pagination controls..."); // Debug log
    const totalPages = Math.ceil(allPosts.length / postsPerPage);
    document.querySelector(".pagination__info").textContent = `Page ${currentPage} of ${totalPages}`;
    document.querySelector(".pagination__prev").disabled = currentPage === 1;
    document.querySelector(".pagination__next").disabled = currentPage === totalPages;
}

// Display posts for the current page
function showPage(allPosts) {
    const start = (currentPage - 1) * postsPerPage;
    const end = start + postsPerPage;
    renderPosts(allPosts.slice(start, end));
}

// Pagination controls for previous and next buttons
document.querySelector(".pagination__prev").addEventListener("click", () => {
    if (currentPage > 1) {
        currentPage--;
        showPage();
        updatePaginationControls();
    }
});

document.querySelector(".pagination__next").addEventListener("click", () => {
    if (currentPage < Math.ceil(allPosts.length / postsPerPage)) {
        currentPage++;
        showPage();
        updatePaginationControls();
    }
});

// ==================== NAVIGATION ====================
// Initialize navigation links and add event listeners
document.addEventListener("DOMContentLoaded", () => {
    console.log("DOMContentLoaded event triggered"); // Debug log
        const homeSection = document.querySelector(".home");
        const postsContainer = document.getElementById("posts-container");
    
        if (!homeSection) console.error(".home element not found");
        else console.log(".home element exists:", homeSection);
    
        if (!postsContainer) console.error("#posts-container not found");
        else console.log("#posts-container exists:", postsContainer);
    
    
        router.init(); // Initialize the router

    const navLinks = document.querySelectorAll(".nav__link");
    navLinks.forEach((link) => {
        link.addEventListener("click", (e) => {
            e.preventDefault();

            // Remove active class from all links
            navLinks.forEach((l) => l.classList.remove("active"));

            // Add active class to the clicked link
            link.classList.add("active");

            // Navigate to the clicked path
            const path = link.getAttribute("data-path");
            console.log("Navigating to path from nav link:", path); // Debug log
            router.navigate(path);
        });
    });
});
