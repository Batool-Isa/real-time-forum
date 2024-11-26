document.addEventListener('DOMContentLoaded', () => {
  const postsContainer = document.getElementById('posts-container');
  const paginationContainer = document.querySelector('.pagination');
  let allPosts = [];
  let currentPage = 1;
  const postsPerPage = 10;

  // Fetch posts from the backend
  document.addEventListener('DOMContentLoaded', () => {
    const postsContainer = document.getElementById('posts-container');
    const paginationContainer = document.querySelector('.pagination');
    let allPosts = [];
    let currentPage = 1;
    const postsPerPage = 10;

    // Fetch posts from the backend
    async function fetchPosts(category = "all") {
        try {
            const response = await fetch(`/api/posts?filter-category=${category}`, {
                headers: { Accept: "application/json" },
            });
            if (!response.ok) throw new Error("Failed to fetch posts");

            const { posts } = await response.json();
            allPosts = posts; // Store posts for pagination
            renderPosts();
            updatePaginationControls();
        } catch (error) {
            console.error("Error fetching posts:", error);
            postsContainer.innerHTML = "<p>Failed to load posts. Please try again later.</p>";
        }
    }

    // Render posts for the current page
    function renderPosts() {
        postsContainer.innerHTML = ""; // Clear existing posts
        const start = (currentPage - 1) * postsPerPage;
        const end = start + postsPerPage;
        const postsToShow = allPosts.slice(start, end);

        postsToShow.forEach((post) => {
            const postElement = document.createElement("article");
            postElement.className = "post";
            postElement.dataset.category = post.categoryName.join(", ");
            postElement.innerHTML = `
                <h2>${post.title}</h2>
                <p>${post.postDescription}</p>
                <div class="post-meta">
                    <span>By: ${post.username}</span>
                    <span>Likes: ${post.like} | Dislikes: ${post.dislike}</span>
                    <span>Categories: ${post.categoryName.join(", ")}</span>
                </div>
            `;
            postsContainer.appendChild(postElement);
        });
    }

    // Update pagination controls
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

    // Fetch and render posts on page load
    fetchPosts();
});



function renderPosts(posts) {
    const postsContainer = document.getElementById("posts-container");
    postsContainer.innerHTML = ""; // Clear existing posts

    posts.forEach((post) => {
        const postElement = document.createElement("article");
        postElement.className = "post";
        postElement.dataset.category = post.categoryName.join(", ");
        postElement.innerHTML = `
            <h2>${post.title}</h2>
            <p>${post.postDescription}</p>
            <div class="post-meta">
                <span>By: ${post.username}</span>
                <span>Likes: ${post.like} | Dislikes: ${post.dislike}</span>
                <span>Categories: ${post.categoryName.join(", ")}</span>
            </div>
        `;
        postsContainer.appendChild(postElement);
    });
}


  
  // Update pagination controls
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

  // Fetch and render posts on page load
  fetchPosts();
});
