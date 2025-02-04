<template>
  <div class="page-wrapper">
      <!-- 1) This is our Leave Group button, shown only if it is a group. -->
    <button
      v-if="isGroup"
      class="leave-group-button"
      @click="leaveGroup"
    >
      Leave Group
    </button>
    <!-- The scrollable area with messages -->
    <div class="messages-container" ref="messagesContainer">
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
              :src="getImageUrl(message.sender_photo.String)"
              alt="Sender Photo"
              class="sender-photo"
            />
            <div class="message-content">
              <div class="message-header">
                <span>{{ message.sender_username }}</span>
                <span class="message-time">{{ formatDate(message.datetime) }}</span>

                <!-- Simple condition for displaying the delete button (for debugging) -->
                <button v-if="true" @click="deleteMessage(message.id)" class="delete-button">Delete</button>
              </div>

              <div v-if="isImage(message.content)">
                <img
                  :src="getImageUrl(message.content)"
                  alt="Image Message"
                  class="message-media"
                />
              </div>
              <div v-else-if="isGif(message.content)">
                <img
                  :src="getImageUrl(message.content)"
                  alt="Gif Message"
                  class="message-media"
                />
              </div>
              <div v-else>
                <p class="message-text">{{ message.content }}</p>
              </div>
            </div>
          </div>
        </li>
      </ul>
    </div>

    <!-- The pinned input bar at the bottom -->
    <div class="message-input-bar">
      <textarea
        v-model="messageText"
        placeholder="Type a message..."
        class="message-textarea"
        @focus="isInteracting = true"
        @blur="checkInteraction"
        @input="checkInteraction"
      ></textarea>

      <div class="media-upload">
        <input
          type="file"
          accept="image/*, .gif"
          class="file-input"
          @click="isInteracting = true"
          @change="handleFileChange"
        />
      </div>

      <button @click="sendMessage" class="send-button">Send</button>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: "GetMessagesView",
  data() {
    return {
      messages: [],
      loading: true,
      isGroup: false,
      messageText: '',
      selectedFile: null,
      isInteracting: false,
      reloadInterval: null,
    };
  },
  async created() {
    const token = localStorage.getItem("authToken");
    const conversationID = this.$route.params.c_id;

    if (!token || !conversationID) {
      console.warn("🚨 Missing token or conversation ID. Cannot load messages.");
      return;
    }

    try {
      const response = await axios.get(
        `http://localhost:3000/conversations/${conversationID}`,
        { headers: { Authorization: `Bearer ${token}` } }
      );
      this.messages = Array.isArray(response.data.messages) ? response.data.messages : [];
      this.isGroup = response.data.conversation.is_group;
    } catch (error) {
      console.error("❌ Error fetching messages:", error);
    } finally {
      this.loading = false;
    }
  },
  mounted() {
    this.$nextTick(() => {
      this.scrollToBottom();
    });

    this.reloadInterval = setInterval(() => {
      if (!this.isInteracting) {
        window.location.reload();
      }
    }, 5000);
  },
  beforeDestroy() {
    if (this.reloadInterval) {
      clearInterval(this.reloadInterval);
    }
  },
  methods: {
    scrollToBottom() {
      const container = this.$refs.messagesContainer;
      if (container) {
        container.scrollTop = container.scrollHeight;
      }
    },
    formatDate(dateString) {
      const date = new Date(dateString);
      return date.toLocaleString();
    },
    isImage(content) {
      return /\.(jpg|jpeg|png)$/i.test(content);
    },
    isGif(content) {
      return /\.(gif)$/i.test(content);
    },
    handleFileChange(event) {
      this.selectedFile = event.target.files[0];
    },
    checkInteraction() {
      if (this.messageText.trim() || this.selectedFile) {
        this.isInteracting = true;
      } else {
        this.isInteracting = false;
      }
    },
    async sendMessage() {
      const token = localStorage.getItem("authToken");
      const conversationID = this.$route.params.c_id;

      if (!token || !conversationID) {
        console.warn("🚨 Missing token or conversation ID. Cannot send message.");
        return;
      }

      let formData = new FormData();
      formData.append("content", this.messageText);
      formData.append("content_type", "text");

      if (this.selectedFile) {
        formData.append("file", this.selectedFile);
        formData.append("content_type", "photo");
      }

      try {
        const response = await axios.post(
          `http://localhost:3000/conversations/${conversationID}/messages`,
          formData,
          {
            headers: {
              Authorization: `Bearer ${token}`,
              "Content-Type": "multipart/form-data",
            },
          }
        );

        const newMessage = {
          content: response.data.content,
          sender_username: response.data.sender_username,
          sender_photo: response.data.sender_photo
            ? `http://localhost:3000${response.data.sender_photo}`
            : '/default-profile.png',
          datetime: new Date().toISOString(),
        };

        this.messages.push(newMessage);

        this.messageText = '';
        this.selectedFile = null;
        this.isInteracting = false;

        this.$nextTick(() => {
          this.scrollToBottom();
          window.location.reload();
        });
      } catch (error) {
        console.error("❌ Error sending message:", error);
      }
    },
    getImageUrl(imagePath) {
      return imagePath && imagePath.startsWith('/uploads')
        ? `http://localhost:3000${imagePath}`
        : '/default-profile.png';
    },

    // Debugging delete functionality
    async deleteMessage(messageId) {
      console.log("Attempting to delete message with ID:", messageId);
      const conversationID = this.$route.params.c_id;
      const token = localStorage.getItem("authToken");

      try {
        await axios.delete(
          `http://localhost:3000/conversations/${conversationID}/messages/${messageId}`,
          { headers: { Authorization: `Bearer ${token}` } }
        );
        console.log("Message deleted successfully.");
        this.messages = this.messages.filter((message) => message.id !== messageId);
      } catch (error) {
        console.error("❌ Error deleting message:", error);
      }
    },
        async leaveGroup() {
      const token = localStorage.getItem("authToken");
      const conversationID = this.$route.params.c_id;
      if (!token || !conversationID) {
        console.warn("🚨 Missing token or conversation ID. Cannot leave group.");
        return;
      }

      try {
        await axios.delete(
          `http://localhost:3000/groups/${conversationID}/leave`,
          { headers: { Authorization: `Bearer ${token}` } }
        );
        console.log("✅ Group left successfully.");
        // Redirect to conversation list after leaving
        this.$router.push("/conversations");
      } catch (error) {
        console.error("❌ Error leaving group:", error);
      }
    },
  },
};
</script>


<style scoped>
.leave-group-button {
  position: fixed;      /* Fix position relative to viewport */
  top: 50px;
  right: 25px;
  z-index: 999;         /* Ensure it stays on top of other elements */
  background-color: #f44336;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 10px 15px;
  cursor: pointer;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.leave-group-button:hover {
  background-color: #c62828;
}
/* Add a button for delete functionality */
.delete-button {
  background-color: red;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 5px 10px;
  cursor: pointer;
}

.delete-button:hover {
  background-color: darkred;
}
/* 
  The entire page is a column layout.
  The .messages-container will scroll if content is taller than available space.
*/
.page-wrapper {
  display: flex;
  flex-direction: column;
  height: 100vh; /* Use full viewport height */
  overflow: hidden; /* Hide any overflow beyond the container */
}

/* The scrollable area with messages occupies the main vertical space,
   leaving room at the bottom for the pinned input bar */
.messages-container {
  flex: 1; 
  overflow-y: auto;
  padding: 20px;
  background-color: #f7f7f7;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

/* The pinned bottom input bar */
.message-input-bar {
  display: flex;
  align-items: center;
  padding: 10px;
  background-color: #fdfdfd; 
  border-top: 1px solid #ccc;
  margin-bottom: 40px;  
  /* pinned to bottom, spanning full width */
  flex-shrink: 0; /* so it doesn't shrink */
}

/* Now the rest of your message styling remains the same */
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
  width: 100px;
  height: 100px;
  object-fit: cover;
  margin-top: 10px;
}

.message-item p {
  margin: 5px 0;
}

/* The input area in the bar */
.message-textarea {
  flex: 1; 
  padding: 8px;
  border-radius: 4px;
  border: 1px solid #ccc;
  resize: vertical;
  margin-right: 10px;
}

.media-upload {
  margin-right: 10px;
}

.file-input {
  font-size: 1rem;
}

.send-button {
  padding: 10px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

.send-button:hover {
  background-color: #0056b3;
}
</style>