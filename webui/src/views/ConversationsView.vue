<template>
  <div class="chats-container">
    <h1>Your Conversations</h1>

    <div v-if="loading">Loading chats...</div>

    <div v-else-if="conversations.length === 0">
      <p>No conversations yet.</p>
    </div>

    <ul v-else>
      <li v-for="chat in conversations" :key="chat.id" class="chat-item">
        <RouterLink :to="`/chat/${chat.id}`" class="chat-link">
          <div class="chat-info">
            <img 
              v-if="chat.photo" 
              :src="chat.photo" 
              alt="Chat Photo" 
              class="chat-photo" 
              @error="setDefaultPhoto($event)"
            />
            <span>{{ chat.name }}</span>
          </div>
          <p class="last-convo-time">{{ formatDate(chat.last_convo) }}</p>
        </RouterLink>
      </li>
    </ul>

    <!-- Add a button to navigate to SendMessageFirstView -->
    <div class="start-new-convo">
      <RouterLink to="/sendMessageFirstView">
        <button class="new-convo-btn">Start a New Conversation</button>
      </RouterLink>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "ConversationsView",
  data() {
    return {
      conversations: [],
      loading: true,
    };
  },
  async created() {
    const token = localStorage.getItem("authToken");
    const userID = localStorage.getItem("userID");

    if (!token || !userID) {
      console.warn("🚨 User not authenticated.");
      return;
    }

    try {
      console.log("🔍 Fetching user conversations...");
      const response = await axios.get(`http://localhost:3000/users/${userID}/conversations`, {
        headers: { Authorization: `Bearer ${token}` },
      });

      console.log("📝 Full API Response:", response.data); // Debugging

      this.conversations = await Promise.all(
        response.data.map(async (chat) => {
          let photoURL = "/default-profile.png"; // Default image

          if (chat.photo && chat.photo.String && chat.photo.String !== "/default-profile.png") {
            photoURL = chat.photo.String;
          }

          return {
            id: chat.id,
            name: chat.name || "Unnamed Chat",
            photo: photoURL,
            last_convo: chat.last_convo, // Capture last_convo
          };
        })
      );

      console.log("✅ Processed Conversations Data:", this.conversations);
    } catch (error) {
      console.error("❌ Error fetching conversations:", error);
    } finally {
      this.loading = false;
    }
  },
  methods: {
    setDefaultPhoto(event) {
      event.target.src = "/default-profile.png";
    },

    formatDate(dateString) {
      const date = new Date(dateString);
      return date.toLocaleString();
    }
  }
};
</script>

<style scoped>
.chats-container {
  text-align: center;
  padding: 20px;
}

.chat-item {
  list-style: none;
  margin: 10px 0;
}

.chat-link {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 8px;
  text-decoration: none;
  color: #333;
  background-color: #f9f9f9;
}

.chat-link:hover {
  background-color: #e9ecef;
}

.chat-photo {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #007bff;
}

.chat-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.last-convo-time {
  font-size: 0.8rem;
  color: #888;
  padding-left: 10px;
  padding-right: 20px;
  align-self: center;
  margin-top: 21px !important;
}

.start-new-convo {
  margin-top: 20px;
}

.new-convo-btn {
  padding: 10px 20px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
}

.new-convo-btn:hover {
  background-color: #0056b3;
}
</style>