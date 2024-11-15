/*==================== SHOW NAVBAR ====================*/
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

// // Get all posts
// const posts = Array.from(postsContainer.querySelectorAll('.post-link'));

// function updatePostsPerPage() {
//   const totalPosts = posts.length;

//   // Set postsPerPage based on the number of total posts
//   postsPerPage = (totalPosts < 20) ? minPostsPerPage : defaultPostsPerPage;

//   showPage(currentPage);
// }

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

// Event listeners for pagination buttons
// prevButton.addEventListener('click', () => {
//   if (currentPage > 1) {
//     currentPage--;
//     showPage(currentPage);
//   }
// });

// nextButton.addEventListener('click', () => {
//   const totalPages = Math.ceil(posts.length / postsPerPage);
//   if (currentPage < totalPages) {
//     currentPage++;
//     showPage(currentPage);
//   }
// });

// // Initial page load
// updatePostsPerPage();
