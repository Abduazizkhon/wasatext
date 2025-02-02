<template>
  <div class="messages-container">
    <h1>Messages</h1>

    <div v-if="loading">Loading messages...</div>

    <!-- Check if it's a group and there are no messages -->
    <div v-else-if="messages.length === 0 && isGroup">
      <p>There are no messages yet</p>
    </div>

    <!-- Default case: display messages -->
    <div v-else-if="messages.length === 0">
      <p>No messages yet.</p>
    </div>

    <ul v-else>
      <li v-for="message in messages" :key="message.id" class="message-item">
        <div class="message-info">
          <img
            v-if="message.sender_photo && message.sender_photo.String"
            :src="message.sender_photo.String || '/default-profile.png'"
            alt="Sender Photo"
            class="sender-photo"
          />
          <div class="message-content">
            <div class="message-header">
              <span>{{ message.sender_username }}</span>
              <span class="message-time">{{ formatDate(message.datetime) }}</span>
            </div>

            <!-- Check if the content is a URL (image/gif) -->
            <div v-if="isImage(message.content)">
              <img :src="message.content" alt="Image Message" class="message-media"/>
            </div>
            <div v-else-if="isGif(message.content)">
              <img :src="message.content" alt="Gif Message" class="message-media"/>
            </div>
            <div v-else>
              <p class="message-text">{{ message.content }}</p>
            </div>
          </div>
        </div>
      </li>
    </ul>
  </div>
</template>

<script>
import axios from 'axios'; // Ensure axios is imported

export default {
  name: "GetMessagesView",
  data() {
    return {
      messages: [], // Ensure it's always an array
      loading: true,
      isGroup: false, // Track if the conversation is a group
    };
  },
  async created() {
    const token = localStorage.getItem("authToken");
    const conversationID = this.$route.params.c_id; // Fetch conversation ID from route

    if (!token || !conversationID) {
      console.warn("🚨 User not authenticated or no conversation ID.");
      return;
    }

    try {
      console.log("🔍 Fetching conversation messages...");
      const response = await axios.get(`http://localhost:3000/conversations/${conversationID}`, {
        headers: { Authorization: `Bearer ${token}` },
      });

      console.log("📝 Full API Response:", response.data); // Debugging

      // Ensure messages is an array even if there are no messages
      this.messages = Array.isArray(response.data.messages) ? response.data.messages : [];
      this.isGroup = response.data.conversation.is_group; // Set the group status

      console.log("✅ Processed Messages Data:", this.messages);
    } catch (error) {
      console.error("❌ Error fetching messages:", error);
    } finally {
      this.loading = false;
    }
  },
  methods: {
    formatDate(dateString) {
      const date = new Date(dateString);
      return date.toLocaleString();
    },
    isImage(content) {
      return /\.(jpg|jpeg|png|gif)$/i.test(content);
    },
    isGif(content) {
      return /\.(gif)$/i.test(content);
    }
  }
};
</script>

<style scoped>
.messages-container {
  padding: 20px;
  background-color: #f7f7f7;
  max-width: 600px;
  margin: 0 auto;
  border-radius: 10px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.message-item {
  list-style: none;
  margin: 15px 0;
  display: flex;
  flex-direction: column;
}

.message-info {
  display: flex;
  gap: 10px;
}

.sender-photo {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
}

.message-content {
  background-color: #fff;
  border-radius: 10px;
  padding: 10px;
  margin-left: 50px;
  max-width: 80%;
}

.message-header {
  display: flex;
  justify-content: space-between;
  font-weight: bold;
  font-size: 0.9rem;
}

.message-time {
  font-size: 0.8rem;
  color: #777;
}

.message-text {
  margin-top: 10px;
  font-size: 1rem;
  line-height: 1.4;
  word-wrap: break-word;
}

.message-media {
  max-width: 100%;
  height: auto;
  margin-top: 10px;
}

.message-item p {
  margin: 5px 0;
}
</style>