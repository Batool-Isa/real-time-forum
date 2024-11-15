document.addEventListener('DOMContentLoaded', () => {
    const toggleButton = document.getElementById('navbar-toggle');
    const sidebar = document.getElementById('sidebar');
    const mainContent = document.getElementById('main-content');

    toggleButton.addEventListener('click', () => {
        sidebar.classList.toggle('active');
        mainContent.classList.toggle('active');
    });
});
