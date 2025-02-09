<template>
  <div class="send-message-container">
    <h1>Start a New Conversation</h1>

    <!-- Loading state -->
    <div v-if="loading">Sending message...</div>

    <!-- Form for text message and file upload -->
    <div v-else>
      <form @submit.prevent="submitForm">
        <div class="form-group">
          <label for="recipient_username">Recipient Username</label>
          <input 
            v-model="recipientUsername" 
            type="text" 
            id="recipient_username" 
            required
            placeholder="Enter recipient's username"
          />
        </div>

        <!-- Choose message type: text, image, or gif -->
        <div class="form-group">
          <label for="message_type">Message Type</label>
          <select v-model="messageType" id="message_type" required>
            <option disabled value="">Select message type</option>
            <option value="text">Text</option>
            <option value="image">Image</option>
            <option value="gif">GIF</option>
          </select>
        </div>

        <!-- Message content section (only shown if the user selects 'text') -->
        <div v-if="messageType === 'text'" class="form-group">
          <label for="content">Message</label>
          <textarea 
            v-model="content" 
            id="content" 
            required
            placeholder="Type your message here..."
          ></textarea>
        </div>

        <!-- File upload section (only shown if the user selects 'image' or 'gif') -->
        <div v-if="messageType === 'image' || messageType === 'gif'" class="form-group">
          <label for="file">Upload a file (Image/GIF)</label>
          <input 
            type="file" 
            id="file" 
            accept="image/*, .gif"
            @change="handleFileChange"
            required
          />
        </div>

        <div class="form-group">
          <button type="submit" class="submit-btn">Send Message</button>
        </div>
      </form>
    </div>

    <!-- Success or Error Message -->
    <div v-if="message">
      <p :class="messageType">{{ message }}</p>
    </div>
  </div>
</template>

<script>
import axios from '../services/axios.js';
export default {
  data() {
    return {
      recipientUsername: "",
      content: "",
      file: null,
      messageType: "",
      loading: false,
      message: "",
      messageClass: ""
    };
  },
  methods: {
    handleFileChange(event) { this.file = event.target.files[0]; },
    async submitForm() {
      const token = localStorage.getItem("authToken");
      const userID = localStorage.getItem("userID");
      if (!token || !userID) { this.message = "User not authenticated."; this.messageClass = "error"; return; }
      const formData = new FormData();
      formData.append("recipient_username", this.recipientUsername);
      formData.append("content_type", this.messageType);
      formData.append("content", this.content);
      if (this.file) { formData.append("file", this.file); }
      try {
        this.loading = true; this.message = ""; this.messageClass = "";
        const response = await axios.post(`/users/${userID}/conversations/first-message`, formData, {
          headers: { Authorization: `Bearer ${token}`, "Content-Type": "multipart/form-data" }
        });
        this.loading = false;
        this.message = `Message sent successfully! Conversation ID: ${response.data.c_id}`;
        this.messageClass = "success";
      } catch (error) {
        this.loading = false;
        this.message = "Failed to send message. Please try again.";
        this.messageClass = "error";
        console.error("Error sending message:", error);
      }
    }
  }
};
</script>

<style scoped>
.send-message-container {
  text-align: center;
  padding: 20px;
}

.form-group {
  margin: 15px 0;
}

input[type="text"],
textarea,
select {
  width: 100%;
  padding: 8px;
  margin: 10px 0;
  border: 1px solid #ccc;
  border-radius: 4px;
}

input[type="file"] {
  padding: 5px;
}

button.submit-btn {
  padding: 10px 20px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
}

button.submit-btn:hover {
  background-color: #0056b3;
}

.message.success {
  color: green;
}

.message.error {
  color: red;
}
</style>