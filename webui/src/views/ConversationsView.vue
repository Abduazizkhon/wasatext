<template>
  <div class="chats-container">
    <h1>My Conversations</h1>

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
            <div class="chat-text-info">
              <span class="chat-name">{{ chat.name }}</span>
              <!-- Display the preview for the last message -->
              <p class="last-message">{{ getLastMessagePreview(chat) }}</p>
            </div>
          </div>
          <p class="last-convo-time">{{ formatDate(chat.last_convo) }}</p>
        </RouterLink>
      </li>
    </ul>

    <div class="start-new-convo">
      <RouterLink to="/sendMessageFirstView">
        <button class="new-convo-btn">Start a New Conversation</button>
      </RouterLink>
      <div class="start-new-group">
        <RouterLink to="/createGroupView">
          <button class="new-group-btn">Create a New Group</button>
        </RouterLink>
      </div>
      <RouterLink to="/search/users" class="search-people-button">
        <button>Search People</button>
      </RouterLink>
    </div>
  </div>
</template>

<script>
import axios from '../services/axios.js';

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
      const response = await axios.get(`/users/${userID}/conversations`, {
        headers: { Authorization: `Bearer ${token}` },
      });

      console.log("📝 Full API Response:", response.data);

      // Process each chat object from the API and map it to our local conversation object.
      this.conversations = await Promise.all(
        response.data.map(async (chat) => {
          let photoURL = "/default-profile.png"; // Default image
          const baseURL = axios.defaults.baseURL; // Use the baseURL from axios instance

          // If a valid photo is provided (i.e. not the default), build its URL using the baseURL.
          if (chat.photo && chat.photo.String && chat.photo.String !== "/default-profile.png") {
            photoURL = chat.photo.String.startsWith(baseURL)
              ? chat.photo.String
              : `${baseURL}${chat.photo.String}`;
          }

          return {
            id: chat.id,
            name: chat.name || "Unnamed Chat",
            photo: photoURL,
            last_convo: chat.last_convo,
            // These fields remain as sql.NullString objects (with .Valid and .String)
            last_message: chat.last_message,
            last_message_type: chat.last_message_type,
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
    },
    /**
     * Returns a preview for the last message:
     * - If the type is "text", returns a truncated version (first 20 characters with ellipsis if longer).
     * - If the type is "photo" or "gif", returns the literal string "photo" or "gif".
     * - Otherwise, returns an empty string.
     *
     * This method "unwraps" the sql.NullString objects by checking their .String property.
     */
    getLastMessagePreview(chat) {
      // Unwrap the last message and type from the sql.NullString objects.
      const msg = (chat.last_message && chat.last_message.String) || "";
      const type = (chat.last_message_type && chat.last_message_type.String) || "";

      if (!msg) {
        return "";
      }
      if (type === "text") {
        return msg.length > 20 ? msg.substring(0, 20) + "..." : msg;
      } else if (type === "photo") {
        return "photo";
      } else if (type === "gif") {
        return "gif";
      } else {
        return "";
      }
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

.chat-text-info {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.chat-name {
  font-weight: bold;
  font-size: 1rem;
}

.last-message {
  font-size: 0.9rem;
  color: #777;
  margin: 0;
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

.new-group-btn {
  padding: 10px 20px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  margin-top: 10px;
}

.new-group-btn:hover {
  background-color: #0056b3;
}
</style>