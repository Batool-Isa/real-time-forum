//sidebar
const menuItems = document.querySelectorAll(".menu-item");
const messagesNotification = document.querySelector("#messages-notifications");
const messages = document.querySelector(".messages");
const message = messages.querySelectorAll(".message");
const messageSearch = document.querySelector("#message-search");

//remove active class from all menu items
const changeActiveItem = () => {
    menuItems.forEach((item) => {
        item.classList.remove("active");
    });
};

// Left sidebar menu items handler
menuItems.forEach((item) => {
    item.addEventListener("click", () => {
        changeActiveItem();
        item.classList.add("active");

        if (item.querySelector('h3').textContent === 'Home') {
            // Restore the feeds content
            document.querySelector('.middle').innerHTML = feedsContent;
        }

        if (item.id != "notifications") {
            document.querySelector(".notifications-popup").style.display = "none";
        } else {
            document.querySelector(".notifications-popup").style.display = "block";
            document.querySelector(
                "#notifications .notification-count"
            ).style.display = "none";
        }
    });
});


messagesNotification.addEventListener("click", () => {
    messages.style.boxShadow = "0 0 1rem var(--color-primary)";
    messagesNotification.querySelector(".notification-count").style.display =
        "none";
    setTimeout(() => {
        messages.style.boxShadow = "none";
    }, 2000);
});

// Make messages clickable
const messageItems = document.querySelectorAll('.messages .message');
messageItems.forEach(messageItem => {
    messageItem.addEventListener('click', () => {
        const username = messageItem.querySelector('h5').textContent;
        const profilePic = messageItem.querySelector('.profile-pic img').src;
        showChat(username, profilePic);
                // Set messages menu item as active
                menuItems.forEach(item => {
                    if(item.querySelector('h3').textContent === 'Messages') {
                        changeActiveItem();
                        item.classList.add('active');
                    }
                });
    });
});

// Get all message elements
const messageElements = document.querySelectorAll('.messages .message');

// Add click handler to each message
messageElements.forEach(messageItem => {
    messageItem.style.cursor = 'pointer';
    messageItem.addEventListener('click', () => {
        // Get user info from the clicked message
        const username = messageItem.querySelector('h5').textContent;
        const profilePic = messageItem.querySelector('.profile-pic img').src;
        
        // Show chat interface
        showChat(username, profilePic);
        
        // Visual feedback for selected chat
        messageElements.forEach(item => item.classList.remove('active'));
        messageItem.classList.add('active');
    });
});


function showChat(username, profilePic) {
    // Get the middle section and clear the feeds
    const middleSection = document.querySelector('.middle');
    middleSection.innerHTML = `
        <div class="chat-container">
            <div class="chat-header">
                <div class="profile-pic">
                    <img src="${profilePic}" />
                </div>
                <div class="chat-user-info">
                    <h4>${username}</h4>
                    <p class="status">online</p>
                </div>
            </div>

            <div class="chat-messages" id="chat-messages">
                <!-- Sample messages -->
                <div class="message received">
                    Hey, how are you?
                </div>
                <div class="message sent">
                    I'm good! How about you?
                </div>
            </div>

            <div class="chat-input">
                <form id="message-form">
                    <input type="text" id="message-input" placeholder="Type a message...">
                    <button type="submit" class="btn btn-primary">Send</button>
                </form>
            </div>
        </div>
    `;

    // Add message form handler
    document.getElementById('message-form').addEventListener('submit', handleMessageSubmit);
}

function handleMessageSubmit(e) {
    e.preventDefault();
    const input = document.getElementById('message-input');
    if (input.value.trim()) {
        addMessage(input.value, 'sent');
        input.value = '';
    }
}

function addMessage(content, type) {
    const messagesContainer = document.getElementById('chat-messages');
    const messageDiv = document.createElement('div');
    messageDiv.className = `message ${type}`;
    messageDiv.textContent = content;
    messagesContainer.appendChild(messageDiv);
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
}

// Store the original feeds content when page loads
let feedsContent;
document.addEventListener('DOMContentLoaded', () => {
    feedsContent = document.querySelector('.middle').innerHTML;
});
