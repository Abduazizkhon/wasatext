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
          <img 
            v-if="chat.photo" 
            :src="chat.photo" 
            alt="Chat Photo" 
            class="chat-photo" 
            @error="setDefaultPhoto($event)"
          />
          <span>{{ chat.name }}</span>
        </RouterLink>
      </li>
    </ul>
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

          // ✅ Access the photo URL correctly from the nested 'String' field
          if (chat.photo && chat.photo.String && chat.photo.String !== "/default-profile.png") {
            // If the photo path is valid, directly use it
            photoURL = chat.photo.String;
          }

          // Returning the processed chat with photo
          return {
            id: chat.id,
            name: chat.name || "Unnamed Chat",
            photo: photoURL,
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
    // ✅ Fallback if image fails to load
    setDefaultPhoto(event) {
      event.target.src = "/default-profile.png";
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
  gap: 10px;
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
</style>