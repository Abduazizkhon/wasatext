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
        <!-- 1) New: "Set Name" button -->
    <button
      v-if="isGroup"
      class="set-name-button"
      @click="toggleSetNameForm"
    >
      Set Name
    </button>

    <!-- 2) New: "Set Name" form -->
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

        <!-- 1) New: Set Photo button -->
    <button
      v-if="isGroup"
      class="set-photo-button"
      @click="toggleSetPhotoForm"
    >
      Set Photo
    </button>

    <!-- 2) New: Set Photo form (shows when showSetPhotoForm === true) -->
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

      // NEW DATA PROPERTIES
      showAddUserForm: false,  // Toggle the form visibility
      usernameToAdd: '',       // The username typed by the user
      showSetNameForm: false,
      newGroupName: '',
      // NEW: for "Set Photo"
      showSetPhotoForm: false,     // Toggle the form
      newGroupPhotoFile: null,     // Hold the file user selected
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
      // Only reload if user not interacting
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

    // NEW METHODS BELOW

    // Show/hide the add user form
    toggleAddUserForm() {
      this.showAddUserForm = !this.showAddUserForm;
      // If opening the form, also temporarily disable auto-refresh
      if (this.showAddUserForm) {
        this.isInteracting = true;
      } else {
        this.isInteracting = false;
      }
    },

    // Explicitly disable auto-refresh
    disableAutoRefresh() {
      this.isInteracting = true;
    },

    // Re-enable auto-refresh
    enableAutoRefresh() {
      // Only re-enable if the username field is empty
      if (!this.usernameToAdd.trim()) {
        this.isInteracting = false;
      }
    },

    // Cancel the “Add User” action
    cancelAddUser() {
      this.showAddUserForm = false;
      this.usernameToAdd = '';
      this.isInteracting = false;
    },

    // Call POST /groups/:c_id/members endpoint to add user
    async addUserToGroup() {
      const token = localStorage.getItem("authToken");
      const conversationID = this.$route.params.c_id;
      if (!token || !conversationID) {
        console.warn("🚨 Missing token or conversation ID. Cannot add user.");
        return;
      }

      // Basic validation
      if (!this.usernameToAdd.trim()) {
        alert("Please enter a valid username.");
        return;
      }

      try {
        const data = {
          usernames: [this.usernameToAdd.trim()]
        };
        const response = await axios.post(
          `http://localhost:3000/groups/${conversationID}/members`,
          data,
          { headers: { Authorization: `Bearer ${token}` } }
        );

        console.log("✅ User(s) added:", response.data.added_users);

        // Feedback to the user
        alert(`User '${this.usernameToAdd.trim()}' added successfully!`);

        // Reset form
        this.usernameToAdd = '';
        this.showAddUserForm = false;
        this.isInteracting = false;

        // Optional: Reload page to see updated members, or do something else
        // window.location.reload();

      } catch (error) {
        console.error("❌ Error adding user to group:", error);
        if (error.response && error.response.data) {
          alert(`Failed to add user: ${error.response.data}`);
        } else {
          alert("Error adding user. Check console for details.");
        }
      }
    },
        // -----------------------------
    //  NEW: Set Name feature below 
    // -----------------------------

    // Toggle the "Set Name" form
    toggleSetNameForm() {
      this.showSetNameForm = !this.showSetNameForm;
      // Temporarily disable auto-refresh if form is open
      if (this.showSetNameForm) {
        this.isInteracting = true;
      } else {
        this.isInteracting = false;
      }
    },

    // Cancel the "Set Name" action
    cancelSetName() {
      this.showSetNameForm = false;
      this.newGroupName = '';
      this.isInteracting = false;
    },

    // PUT /groups/:c_id/name to update the group name
    async updateGroupName() {
      const token = localStorage.getItem("authToken");
      const conversationID = this.$route.params.c_id;

      if (!token || !conversationID) {
        console.warn("🚨 Missing token or conversation ID. Cannot set group name.");
        return;
      }

      // Basic validation
      if (!this.newGroupName.trim()) {
        alert("Please enter a valid group name.");
        return;
      }

      try {
        // The API expects { new_name: "<name>" }
        await axios.put(
          `http://localhost:3000/groups/${conversationID}/name`,
          { new_name: this.newGroupName.trim() },
          { headers: { Authorization: `Bearer ${token}` } }
        );

        alert(`Group name updated to: "${this.newGroupName.trim()}"`);
        this.showSetNameForm = false;
        this.isInteracting = false;

        // Optionally reload page or do something else
        // window.location.reload();

      } catch (error) {
        console.error("❌ Error setting group name:", error);
        if (error.response && error.response.data) {
          alert(`Failed to set group name: ${error.response.data}`);
        } else {
          alert("Error setting group name. Check console for details.");
        }
      }
    },

    // -----------------------------
    // NEW: Set Photo feature below
    // -----------------------------
      onPhotoInputClick() {
    this.isInteracting = true;
  },

  handleGroupPhotoChange(event) {
    this.newGroupPhotoFile = event.target.files[0] || null;
    // Keep isInteracting = true for now.
    // If user cancels the file dialog, you may not get a change event,
    // but at least the page won't reload in the middle of picking a file.
  },

  cancelSetPhoto() {
    // If user hits "Cancel" in your own form, then set isInteracting = false
    this.showSetPhotoForm = false;
    this.newGroupPhotoFile = null;
    this.isInteracting = false;  
  },

    // Toggle the "Set Photo" form
    toggleSetPhotoForm() {
      this.showSetPhotoForm = !this.showSetPhotoForm;
      if (this.showSetPhotoForm) {
        this.isInteracting = true; // Stop auto-refresh
      } else {
        this.isInteracting = false;
      }
    },

    // Cancel the "Set Photo" action
    cancelSetPhoto() {
      this.showSetPhotoForm = false;
      this.newGroupPhotoFile = null;
      this.isInteracting = false;
    },

    // Handle file selection for the group photo
    handleGroupPhotoChange(event) {
      this.newGroupPhotoFile = event.target.files[0] || null;
      // Optionally disable auto-refresh on file selection
      this.isInteracting = true;
    },

    // PUT /conversations/:c_id/set-group-photo with the chosen file
    async updateGroupPhoto() {
      const token = localStorage.getItem("authToken");
      const conversationID = this.$route.params.c_id;

      if (!token || !conversationID) {
        console.warn("🚨 Missing token or conversation ID. Cannot set group photo.");
        return;
      }

      // Ensure a file is chosen
      if (!this.newGroupPhotoFile) {
        alert("Please select a photo file.");
        return;
      }

      try {
        const formData = new FormData();
        // The API endpoint expects "photo" as the field name
        formData.append("photo", this.newGroupPhotoFile);

        await axios.put(
          `http://localhost:3000/conversations/${conversationID}/set-group-photo`,
          formData,
          {
            headers: {
              Authorization: `Bearer ${token}`,
              "Content-Type": "multipart/form-data",
            },
          }
        );

        alert("✅ Group photo updated successfully!");
        this.showSetPhotoForm = false;
        this.newGroupPhotoFile = null;
        this.isInteracting = false;

        // Optionally reload
        // window.location.reload();
      } catch (error) {
        console.error("❌ Error setting group photo:", error);
        if (error.response && error.response.data) {
          alert(`Failed to set group photo: ${error.response.data}`);
        } else {
          alert("Error setting group photo. Check console for details.");
        }
      }
    },
  },
};
</script>


<style scoped>
/* "Set Photo" button (placed below the Set Name button) */
.set-photo-button {
  position: fixed;
  top: 200px; /* Adjust as needed */
  right: 25px;
  z-index: 999;
  background-color: #2196f3; /* blue-ish */
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

/* "Set Photo" form (below the button) */
.set-photo-form {
  position: fixed;
  top: 250px; /* Just below the button */
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
/* "Set Name" button */
.set-name-button {
  position: fixed;
  top: 150px; /* Just below the "Add User" button (adjust as needed) */
  right: 25px;
  z-index: 999;
  background-color: #ff9800; /* orange-ish */
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

/* "Set Name" form */
.set-name-form {
  position: fixed;
  top: 200px; /* Just below the button */
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
/* Position the "Add User" button near the "Leave Group" button */
.add-user-button {
  position: fixed;
  top: 100px; /* Just below the "Leave Group" button (adjust as needed) */
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

/* The form container that appears for adding a user */
.add-user-form {
  position: fixed;
  top: 150px; /* Just below the "Add User" button */
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