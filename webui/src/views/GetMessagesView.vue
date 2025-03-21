<template>
  <div class="page-wrapper">
    <!-- Leave Group button (shown only if it is a group) -->
    <button
      v-if="isGroup"
      class="leave-group-button"
      @click="leaveGroup"
    >
      Leave Group
    </button>

    <!-- Add User button (only for groups) -->
    <button
      v-if="isGroup"
      class="add-user-button"
      @click="toggleAddUserForm"
    >
      Add User
    </button>

    <!-- Add User form (shown when showAddUserForm === true) -->
    <div
      v-if="isGroup && showAddUserForm"
      class="add-user-form"
      @click.stop
    >
      <input
        type="text"
        placeholder="Username to add"
        v-model="usernameToAdd"
        @focus="disableAutoRefresh"
        @input="disableAutoRefresh"
        @blur="enableAutoRefresh"
      />
      <button @click="addUserToGroup">Add</button>
      <button @click="cancelAddUser">Cancel</button>
    </div>

    <!-- Set Name button (only for groups) -->
    <button
      v-if="isGroup"
      class="set-name-button"
      @click="toggleSetNameForm"
    >
      Set Name
    </button>

    <!-- Set Name form -->
    <div
      v-if="isGroup && showSetNameForm"
      class="set-name-form"
      @click.stop
    >
      <input
        type="text"
        placeholder="New Group Name"
        v-model="newGroupName"
        @focus="disableAutoRefresh"
        @input="disableAutoRefresh"
        @blur="enableAutoRefresh"
      />
      <button @click="updateGroupName">Save</button>
      <button @click="cancelSetName">Cancel</button>
    </div>

    <!-- Set Photo button (only for groups) -->
    <button
      v-if="isGroup"
      class="set-photo-button"
      @click="toggleSetPhotoForm"
    >
      Set Photo
    </button>

    <!-- Set Photo form -->
    <div
      v-if="isGroup && showSetPhotoForm"
      class="set-photo-form"
      @click.stop
    >
      <input
        type="file"
        accept="image/*"
        @click="onPhotoInputClick"
        @change="handleGroupPhotoChange"
      />
      <button @click="updateGroupPhoto">Save</button>
      <button @click="cancelSetPhoto">Cancel</button>
    </div>

    <!-- The scrollable area with messages -->
    <div class="messages-container" ref="messagesContainer">
      <h1>Messages</h1>

      <div v-if="loading">Loading messages...</div>
      <div v-else-if="messages.length === 0 && isGroup">
        <p>There are no messages yet</p>
      </div>
      <div v-else-if="messages.length === 0">
        <p>No messages yet.</p>
      </div>
      <ul v-else>
        <li v-for="message in messages" :key="message.id" class="message-item">
          <div class="message-info">
            <!-- Sender Photo -->
            <img
              v-if="message.sender_photo && message.sender_photo.String"
              :src="getImageUrl(message.sender_photo.String)"
              alt="Sender Photo"
              class="sender-photo"
            />
            <div class="message-content">
              <div class="message-header">
                <span>{{ message.sender_username }}</span>
                <!-- Show "Forwarded" label if status === 'forwarded' -->
                <span v-if="message.status === 'forwarded'" class="forwarded-label">
                  Forwarded
                </span>
                <span class="message-time">{{ formatDate(message.datetime) }}</span>

                <!-- Delete button (if message sent by current user) -->
                <button
                  v-if="message.sender_id === currentUser"
                  @click="deleteMessage(message.id)"
                  class="delete-button"
                >
                  Delete
                </button>
                <!-- Forward button -->
                <button class="forward-button" @click="toggleForwardPanel(message)">
                  {{ message.showForwardPanel ? 'Cancel Forward' : 'Forward' }}
                </button>
                <!-- Forward Panel -->
                <div v-if="message.showForwardPanel" class="forward-panel">
                  <input
                    type="text"
                    v-model="message.forwardTarget"
                    placeholder="Forward to (username or group)"
                  />
                  <button @click="forwardMessageHandler(message)">Forward</button>
                </div>
                <!-- Comment button toggles the comment section -->
                <button
                  class="comment-button"
                  @click="toggleComments(message)"
                >
                  {{ message.showComments ? 'Hide Comments' : 'Comment' }}
                </button>
              </div>

              <!-- REPLY BUTTON -->
              <button class="reply-button" @click="initiateReply(message)">
                Reply
              </button>

              <!-- Display reply snippet only if valid reply info exists -->
              <div v-if="hasReply(message)">
                <blockquote class="reply-snippet">
                  <strong>Replying to: {{ getReplySender(message.reply_to_sender) }}</strong>
                  <br />
                  "{{ replySnippet(message.reply_to_content) }}"
                </blockquote>
              </div>

              <!-- Message content (text, image, or GIF) -->
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

              <!-- Comments Section -->
              <div v-if="message.showComments" class="comments-section">
                <!-- Display each comment -->
                <div v-for="c in message.comments" :key="c.id" class="single-comment">
                  <p class="comment-header">
                    <strong>{{ c.username }}</strong>
                    <span class="comment-time">{{ formatDate(c.timestamp) }}</span>
                    <!-- Delete comment button (if owned by current user) -->
                    <button
                      v-if="c.user_id === currentUser"
                      @click="deleteComment(message, c)"
                      class="delete-comment-button"
                    >
                      Delete
                    </button>
                  </p>
                  <p class="comment-text">{{ c.content }}</p>
                </div>

                <!-- Add New Comment Form (only a ❤️ react button) -->
                <div class="add-comment-form">
                  <button @click="addEmojiComment(message)">❤️ React</button>
                </div>
              </div>
            </div>
          </div>
        </li>
      </ul>
    </div>

    <!-- REPLYING BAR (shown if we are replying to some message) -->
    <div v-if="replyToMessage" class="replying-bar">
      <p>
        Replying to: <strong>{{ replyToMessage.sender_username }}</strong><br />
        <em>{{ replySnippet(replyToMessage.content) }}</em>
        <button class="cancel-reply" @click="cancelReply">✕</button>
      </p>
    </div>

    <!-- Message Input Bar -->
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
import axios from '../services/axios.js';

export default {
  name: "GetMessagesView",
  data() {
    return {
      currentUser: localStorage.getItem("userID"),
      messages: [],
      loading: true,
      isGroup: false,
      messageText: '',
      selectedFile: null,
      isInteracting: false,
      reloadInterval: null,
      showAddUserForm: false,
      usernameToAdd: '',
      showSetNameForm: false,
      newGroupName: '',
      showSetPhotoForm: false,
      newGroupPhotoFile: null,
      // REPLY LOGIC
      replyToMessage: null
    };
  },
  async created() {
    await this.getConversation();
  },
  mounted() {
    this.$nextTick(() => {
      this.scrollToBottom();
    });
    // Auto-refresh if not interacting
    this.reloadInterval = setInterval(() => {
      if (!this.isInteracting) {
        this.getConversation();
      }
    }, 1000);
  },
  beforeUnmount() {
    if (this.reloadInterval) {
      clearInterval(this.reloadInterval);
    }
  },
  methods: {
    async getConversation() {
      const token = localStorage.getItem("authToken");
      const conversationID = this.$route.params.c_id;
      if (!token || !conversationID) {
        console.warn("Missing token or conversation ID. Cannot fetch conversation.");
        return;
      }
      try {
        const response = await axios.get(`conversations/${conversationID}`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        const fetchedMessages = Array.isArray(response.data.messages)
          ? response.data.messages
          : [];

        // Create a map of existing messages by ID for quick lookup
        const existingMessagesMap = new Map(this.messages.map(msg => [msg.id, msg]));

        // Merge fetched messages with existing messages
        const mergedMessages = fetchedMessages.map(fetchedMsg => {
          const existingMsg = existingMessagesMap.get(fetchedMsg.id);
          if (existingMsg) {
            // Preserve dynamic properties from the existing message
            return {
              ...fetchedMsg,
              showComments: existingMsg.showComments,
              showForwardPanel: existingMsg.showForwardPanel,
              forwardTarget: existingMsg.forwardTarget,
              status: existingMsg.status,
              reply_to: existingMsg.reply_to,
              reply_to_content: existingMsg.reply_to_content,
              reply_to_sender: existingMsg.reply_to_sender,
              comments: existingMsg.comments,
            };
          } else {
            // If it's a new message, initialize dynamic properties
            return {
              ...fetchedMsg,
              comments: [],
              newComment: '',
              showComments: false,
              showForwardPanel: false,
              forwardTarget: '',
              status: fetchedMsg.status || ''
            };
          }
        });

        this.messages = mergedMessages;
        this.isGroup = response.data.conversation.is_group;
      } catch (error) {
        console.error("Error fetching conversation:", error);
      } finally {
        this.loading = false;
      }
    },

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
      this.isInteracting = this.messageText.trim() || this.selectedFile ? true : false;
    },

    // Initiate a reply
    initiateReply(message) {
      this.replyToMessage = message;
      this.isInteracting = true;
    },
    // Cancel the current reply
    cancelReply() {
      this.replyToMessage = null;
      this.isInteracting = false;
    },
    // Returns a short snippet of the given text
    snippetOf(text) {
      if (!text) return '';
      return text.length > 40 ? text.substring(0, 40) + '…' : text;
    },
    // Helper: if value is an object with a String property, return that; otherwise return value.
    getPlainValue(val) {
      if (val && typeof val === 'object' && 'String' in val) {
        return val.String;
      }
      return val;
    },
    // Returns true if the message has a valid reply sender or content.
    hasReply(message) {
      const sender = this.getPlainValue(message.reply_to_sender);
      const content = this.getPlainValue(message.reply_to_content);
      return message.reply_to && (sender || content);
    },
    // Return the plain reply sender.
    getReplySender(val) {
      return this.getPlainValue(val);
    },
    // Helper to return the display for a reply's content.
    // If the content looks like a gif or image, return "gif" or "photo" respectively.
    replySnippet(content) {
      const plain = this.getPlainValue(content);
      if (!plain) return '';
      if (this.isGif(plain)) {
        return "gif";
      } else if (this.isImage(plain)) {
        return "photo";
      }
      return this.snippetOf(plain);
    },

    // Send a new message (with or without reply)
    async sendMessage() {
      const token = localStorage.getItem("authToken");
      const conversationID = this.$route.params.c_id;
      if (!token || !conversationID) {
        console.warn("Missing token or conversation ID. Cannot send message.");
        return;
      }
      try {
        let formData = new FormData();
        formData.append("content", this.messageText);
        formData.append("content_type", "text");
        // If replying, add reply_to
        if (this.replyToMessage) {
          formData.append("reply_to", this.replyToMessage.id.toString());
        }
        if (this.selectedFile) {
          formData.append("file", this.selectedFile);
          formData.append("content_type", "photo");
        }
        const response = await axios.post(
          `/conversations/${conversationID}/messages`,
          formData,
          {
            headers: {
              Authorization: `Bearer ${token}`,
              "Content-Type": "multipart/form-data"
            }
          }
        );
        const baseURL = axios.defaults.baseURL;
        // Insert the newly sent message locally
        const newMessage = {
          id: response.data.message_id || Date.now(),
          content: response.data.content,
          sender_username: response.data.sender_username,
          sender_photo: response.data.sender_photo
            ? (response.data.sender_photo.startsWith(baseURL)
                ? response.data.sender_photo
                : baseURL + response.data.sender_photo)
            : '/default-profile.png',
          datetime: new Date().toISOString(),
          sender_id: this.currentUser,
          status: '',
          comments: [],
          newComment: '',
          showComments: false,
          showForwardPanel: false,
          forwardTarget: '',
          // Store reply information as plain values
          reply_to: this.replyToMessage ? this.replyToMessage.id : null,
          reply_to_sender: this.replyToMessage ? this.replyToMessage.sender_username : '',
          reply_to_content: this.replyToMessage ? this.replyToMessage.content : ''
        };
        this.messages.push(newMessage);
        this.messageText = '';
        this.selectedFile = null;
        this.replyToMessage = null; // Clear reply
        this.isInteracting = false;
        this.$nextTick(() => {
          this.scrollToBottom();
        });
      } catch (error) {
        console.error("Error sending message:", error);
      }
    },

    // Build a full image URL if it's a local uploads path
    getImageUrl(imagePath) {
      const baseURL = axios.defaults.baseURL;
      return imagePath && imagePath.startsWith('/uploads')
        ? baseURL + imagePath
        : '/default-profile.png';
    },

    // Delete a message
    async deleteMessage(messageId) {
      const conversationID = this.$route.params.c_id;
      const token = localStorage.getItem("authToken");
      try {
        await axios.delete(`/conversations/${conversationID}/messages/${messageId}`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        this.messages = this.messages.filter((message) => message.id !== messageId);
      } catch (error) {
        console.error("Error deleting message:", error);
      }
    },

    // Delete a comment
    async deleteComment(message, comment) {
      const token = localStorage.getItem("authToken");
      const conversationID = this.$route.params.c_id;
      if (!token || !conversationID) {
        console.warn("Missing token or conversation ID. Cannot delete comment.");
        return;
      }
      try {
        await axios.delete(
          `/conversations/${conversationID}/messages/${message.id}/comments/${comment.id}`,
          { headers: { Authorization: `Bearer ${token}` } }
        );
        message.comments = message.comments.filter((c) => c.id !== comment.id);
        alert("Comment deleted successfully!");
      } catch (error) {
        console.error("Error deleting comment:", error);
        alert("Error deleting comment. Check console for details.");
      }
    },

    // Toggle comment section
    async toggleComments(message) {
      message.showComments = !message.showComments;
      if (message.showComments) {
        this.isInteracting = true;
        const token = localStorage.getItem("authToken");
        try {
          const response = await axios.get(`/messages/${message.id}/comments`, {
            headers: { Authorization: `Bearer ${token}` }
          });
          message.comments = response.data;
        } catch (error) {
          console.error("Error fetching comments:", error);
        }
      } else {
        this.$nextTick(() => {
          const anyOpen = this.messages.some((msg) => msg.showComments);
          if (!anyOpen) {
            this.isInteracting = false;
          }
        });
      }
    },

    // Add emoji comment
    async addEmojiComment(message) {
      const token = localStorage.getItem("authToken");
      const conversationID = this.$route.params.c_id;
      if (!token || !conversationID) {
        console.warn("Missing token or conversation ID. Cannot comment.");
        return;
      }
      try {
        await axios.post(
          `/conversations/${conversationID}/messages/${message.id}/comments`,
          { content_type: "emoji", content: "❤️" },
          { headers: { Authorization: `Bearer ${token}` } }
        );
        const response = await axios.get(`/messages/${message.id}/comments`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        message.comments = response.data;
      } catch (error) {
        console.error("Error adding emoji comment:", error);
        alert("Error adding comment. Check console for details.");
      }
    },

    // Toggle forward panel
    toggleForwardPanel(message) {
      message.showForwardPanel = !message.showForwardPanel;
      if (message.showForwardPanel) {
        this.isInteracting = true;
      } else {
        if (!this.messages.some((msg) => msg.showForwardPanel)) {
          this.isInteracting = false;
        }
      }
    },

    // Forward a message
    async forwardMessageHandler(message) {
      const token = localStorage.getItem("authToken");
      const conversationID = this.$route.params.c_id;
      if (!token || !conversationID) {
        console.warn("Missing token or conversation ID. Cannot forward message.");
        return;
      }
      if (!message.forwardTarget.trim()) {
        alert("Please specify a username or group name to forward to.");
        return;
      }
      let targetConversationID = "";
      if (!isNaN(message.forwardTarget)) {
        targetConversationID = message.forwardTarget;
      } else {
        targetConversationID = "new";
      }
      try {
        let payload = {};
        if (targetConversationID === "new") {
          payload = { target_username: message.forwardTarget.trim() };
        }
        await axios.post(
          `/conversations/${conversationID}/messages/${message.id}/forward/${targetConversationID}`,
          payload,
          { headers: { Authorization: `Bearer ${token}` } }
        );
        message.status = "forwarded";
        message.showForwardPanel = false;
        message.forwardTarget = "";
        alert("Message forwarded successfully!");
        if (!this.messages.some((msg) => msg.showForwardPanel)) {
          this.isInteracting = false;
        }
      } catch (error) {
        console.error("Error forwarding message:", error);
        alert("Error forwarding message. Check console for details.");
      }
    },

    // Leave group action
    async leaveGroup() {
      const token = localStorage.getItem("authToken");
      const conversationID = this.$route.params.c_id;
      if (!token || !conversationID) {
        console.warn("Missing token or conversation ID. Cannot leave group.");
        return;
      }
      try {
        await axios.delete(`/groups/${conversationID}/leave`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        this.$router.push("/conversations");
      } catch (error) {
        console.error("Error leaving group:", error);
      }
    },

    // Add user form toggles
    toggleAddUserForm() {
      this.showAddUserForm = !this.showAddUserForm;
      this.isInteracting = this.showAddUserForm;
    },
    disableAutoRefresh() {
      this.isInteracting = true;
    },
    enableAutoRefresh() {
      if (!this.usernameToAdd.trim()) {
        this.isInteracting = false;
      }
    },
    cancelAddUser() {
      this.showAddUserForm = false;
      this.usernameToAdd = '';
      this.isInteracting = false;
    },
    // Add a user to the group
    async addUserToGroup() {
      const token = localStorage.getItem("authToken");
      const conversationID = this.$route.params.c_id;
      if (!token || !conversationID) {
        console.warn("Missing token or conversation ID. Cannot add user.");
        return;
      }
      if (!this.usernameToAdd.trim()) {
        alert("Please enter a valid username.");
        return;
      }
      try {
        const data = { usernames: [this.usernameToAdd.trim()] };
        await axios.post(`/groups/${conversationID}/members`, data, {
          headers: { Authorization: `Bearer ${token}` }
        });
        alert(`User '${this.usernameToAdd.trim()}' added successfully!`);
        this.usernameToAdd = '';
        this.showAddUserForm = false;
        this.isInteracting = false;
      } catch (error) {
        console.error("Error adding user to group:", error);
        alert("Error adding user. Check console for details.");
      }
    },

    // Set group name toggles
    toggleSetNameForm() {
      this.showSetNameForm = !this.showSetNameForm;
      this.isInteracting = this.showSetNameForm;
    },
    cancelSetName() {
      this.showSetNameForm = false;
      this.newGroupName = '';
      this.isInteracting = false;
    },
    async updateGroupName() {
      const token = localStorage.getItem("authToken");
      const conversationID = this.$route.params.c_id;
      if (!token || !conversationID) {
        console.warn("Missing token or conversation ID. Cannot set group name.");
        return;
      }
      if (!this.newGroupName.trim()) {
        alert("Please enter a valid group name.");
        return;
      }
      try {
        await axios.put(`/groups/${conversationID}/name`, { new_name: this.newGroupName.trim() }, {
          headers: { Authorization: `Bearer ${token}` }
        });
        alert(`Group name updated to: "${this.newGroupName.trim()}"`);
        this.showSetNameForm = false;
        this.isInteracting = false;
      } catch (error) {
        console.error("Error setting group name:", error);
        alert("Error setting group name. Check console for details.");
      }
    },

    // Set group photo
    onPhotoInputClick() {
      this.isInteracting = true;
    },
    handleGroupPhotoChange(event) {
      this.newGroupPhotoFile = event.target.files[0] || null;
      this.isInteracting = true;
    },
    cancelSetPhoto() {
      this.showSetPhotoForm = false;
      this.newGroupPhotoFile = null;
      this.isInteracting = false;
    },
    toggleSetPhotoForm() {
      this.showSetPhotoForm = !this.showSetPhotoForm;
      this.isInteracting = this.showSetPhotoForm;
    },
    async updateGroupPhoto() {
      const token = localStorage.getItem("authToken");
      const conversationID = this.$route.params.c_id;
      if (!token || !conversationID) {
        console.warn("Missing token or conversation ID. Cannot set group photo.");
        return;
      }
      if (!this.newGroupPhotoFile) {
        alert("Please select a photo file.");
        return;
      }
      try {
        const formData = new FormData();
        formData.append("photo", this.newGroupPhotoFile);
        await axios.put(`/conversations/${conversationID}/set-group-photo`, formData, {
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "multipart/form-data"
          }
        });
        alert("Group photo updated successfully!");
        this.showSetPhotoForm = false;
        this.newGroupPhotoFile = null;
        this.isInteracting = false;
      } catch (error) {
        console.error("Error setting group photo:", error);
        alert("Error setting group photo. Check console for details.");
      }
    }
  }
};
</script>

<style scoped>
/* Position the "Leave Group" button */
.leave-group-button {
  position: fixed;
  top: 50px;
  right: 25px;
  z-index: 999;
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

/* Position the "Add User" button below Leave Group */
.add-user-button {
  position: fixed;
  top: 100px;
  right: 25px;
  z-index: 999;
  background-color: #4caf50;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 10px 15px;
  cursor: pointer;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}
.add-user-button:hover {
  background-color: #388e3c;
}

/* The form for adding a user */
.add-user-form {
  position: fixed;
  top: 150px;
  right: 25px;
  z-index: 1000;
  background-color: #fff;
  border: 1px solid #ccc;
  border-radius: 6px;
  padding: 10px;
  width: 200px;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.2);
}
.add-user-form input {
  width: 100%;
  margin-bottom: 10px;
  padding: 6px;
  box-sizing: border-box;
}
.add-user-form button {
  margin-right: 6px;
  padding: 5px 10px;
  cursor: pointer;
}

/* Set Name button */
.set-name-button {
  position: fixed;
  top: 150px;
  right: 25px;
  z-index: 999;
  background-color: #ff9800;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 10px 15px;
  cursor: pointer;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}
.set-name-button:hover {
  background-color: #f57c00;
}

/* Set Name form */
.set-name-form {
  position: fixed;
  top: 200px;
  right: 25px;
  z-index: 1000;
  background-color: #fff;
  border: 1px solid #ccc;
  border-radius: 6px;
  padding: 10px;
  width: 200px;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.2);
}
.set-name-form input {
  width: 100%;
  margin-bottom: 10px;
  padding: 6px;
  box-sizing: border-box;
}
.set-name-form button {
  margin-right: 6px;
  padding: 5px 10px;
  cursor: pointer;
}

/* Set Photo button */
.set-photo-button {
  position: fixed;
  top: 200px;
  right: 25px;
  z-index: 999;
  background-color: #2196f3;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 10px 15px;
  cursor: pointer;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}
.set-photo-button:hover {
  background-color: #1976d2;
}

/* Set Photo form */
.set-photo-form {
  position: fixed;
  top: 250px;
  right: 25px;
  z-index: 1000;
  background-color: #fff;
  border: 1px solid #ccc;
  border-radius: 6px;
  padding: 10px;
  width: 200px;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.2);
}
.set-photo-form input {
  width: 100%;
  margin-bottom: 10px;
  padding: 6px;
  box-sizing: border-box;
}
.set-photo-form button {
  margin-right: 6px;
  padding: 5px 10px;
  cursor: pointer;
}

/* "Comment" button & comment section styling */
.comment-button {
  background-color: #9c27b0;
  color: #fff;
  border: none;
  border-radius: 4px;
  padding: 5px 8px;
  margin-left: 8px;
  cursor: pointer;
}
.comment-button:hover {
  background-color: #7b1fa2;
}
.comments-section {
  margin-top: 10px;
  padding: 10px;
  border-radius: 6px;
  background-color: #f2f2f2;
}
.single-comment {
  background-color: #fff;
  margin-bottom: 6px;
  padding: 6px 8px;
  border-radius: 4px;
}
.comment-text {
  margin: 0;
  font-size: 0.9rem;
  word-wrap: break-word;
}
.add-comment-form {
  display: flex;
  margin-top: 6px;
  gap: 6px;
}
.add-comment-form button {
  background-color: #e91e63;
  color: #fff;
  border: none;
  border-radius: 4px;
  padding: 6px 12px;
  cursor: pointer;
}
.add-comment-form button:hover {
  background-color: #d81b60;
}

/* Layout: page, messages container, and input bar */
.page-wrapper {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
}
.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background-color: #f7f7f7;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}
.message-input-bar {
  display: flex;
  align-items: center;
  padding: 10px;
  background-color: #fdfdfd;
  border-top: 1px solid #ccc;
  margin-bottom: 40px;
  flex-shrink: 0;
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
  margin-left: 10px;
}
.delete-button {
  background-color: red;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 5px 10px;
  cursor: pointer;
  margin-left: 8px;
}
.delete-button:hover {
  background-color: darkred;
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
.comment-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 4px;
}
.delete-comment-button {
  background: transparent;
  border: none;
  color: #e53935;
  font-size: 0.8rem;
  cursor: pointer;
  margin-left: 10px;
  padding: 2px 6px;
  border-radius: 4px;
  transition: background-color 0.2s, color 0.2s;
}
.delete-comment-button:hover {
  background-color: #ffe6e6;
  color: #b71c1c;
}
.forward-button {
  background-color: #03a9f4;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 5px 8px;
  cursor: pointer;
  margin-left: 8px;
}
.forward-button:hover {
  background-color: #0288d1;
}
.forward-panel {
  margin-top: 10px;
  padding: 6px;
  background-color: #e3f2fd;
  border: 1px solid #90caf9;
  border-radius: 4px;
  display: flex;
  gap: 6px;
}
.forward-panel input[type="text"] {
  flex: 1;
  padding: 4px;
  border: 1px solid #ccc;
  border-radius: 4px;
}
.forwarded-label {
  font-size: 0.8rem;
  color: #555;
  margin-left: 8px;
  font-style: italic;
}

/* REPLY button */
.reply-button {
  background-color: #4a148c;
  color: #fff;
  border: none;
  border-radius: 4px;
  padding: 4px 6px;
  margin-top: 6px;
  cursor: pointer;
}
.reply-button:hover {
  background-color: #6a1b9a;
}

/* The bar showing "Replying to..." */
.replying-bar {
  background-color: #fff3cd;
  color: #856404;
  border: 1px solid #ffeeba;
  padding: 8px;
  margin: 0 10px 10px;
  border-radius: 4px;
  text-align: center;
}
.replying-bar .cancel-reply {
  background: none;
  border: none;
  color: red;
  font-weight: bold;
  margin-left: 12px;
  cursor: pointer;
  font-size: 14px;
}
</style>