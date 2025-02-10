<template>
  <div class="search-user-container">
    <h1>Search for a User</h1>
    <div class="search-bar">
      <input type="text" v-model="searchQuery" placeholder="Enter username" />
      <button @click="searchUser">Search</button>
    </div>

    <div v-if="errorMessage" class="error-message">
      {{ errorMessage }}
    </div>

    <div v-if="foundUser" class="user-info">
      <p><strong>{{ foundUser.username }}</strong></p>
      <img
        v-if="foundUser.photo && foundUser.photo.String"
        :src="fullPhotoUrl(foundUser.photo.String)"
        alt="User Photo"
        class="user-photo"
      />
      <div v-if="conversationId">
        <p>You already have a conversation with this user.</p>
        <button @click="goToConversation">Go to Conversation</button>
      </div>
      <div v-else>
        <p>No conversation exists. Click to start a new conversation.</p>
        <button @click="startConversation">Start Conversation</button>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "../services/axios.js";
export default {
  name: "SearchUserView",
  data() {
    return {
      searchQuery: "",
      foundUser: null,
      conversationId: null,
      errorMessage: ""
    };
  },
  methods: {
    fullPhotoUrl(photoPath) {
      return axios.defaults.baseURL + photoPath;
    },
    async searchUser() {
      this.errorMessage = "";
      this.foundUser = null;
      this.conversationId = null;

      const token = localStorage.getItem("authToken");
      if (!token) {
        this.errorMessage = "User not authenticated.";
        return;
      }
      try {
        const response = await axios.get(`/search/users?username=${encodeURIComponent(this.searchQuery)}`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        this.foundUser = response.data.user;
        this.conversationId = response.data.conversation_id;
      } catch (err) {
        if (err.response && err.response.status === 404) {
          this.errorMessage = "No such user found.";
        } else {
          this.errorMessage = "Error searching for user.";
        }
      }
    },
    goToConversation() {
      this.$router.push(`/chat/${this.conversationId}`);
    },
    async startConversation() {
      const token = localStorage.getItem("authToken");
      const userID = localStorage.getItem("userID");
      if (!token || !userID) {
        this.errorMessage = "User not authenticated.";
        return;
      }
      try {
        const formData = new FormData();
        formData.append("recipient_username", this.foundUser.username);
        formData.append("content_type", "text");
        formData.append("content", "Hi!");
        const response = await axios.post(`/users/${userID}/conversations/first-message`, formData, {
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "multipart/form-data"
          }
        });
        if (response.data && response.data.c_id) {
          this.$router.push(`/chat/${response.data.c_id}`);
        }
      } catch (err) {
        this.errorMessage = "Error starting conversation.";
      }
    }
  }
};
</script>

<style scoped>
.search-user-container {
  text-align: center;
  padding: 20px;
}
.search-bar {
  margin-bottom: 20px;
}
.search-bar input {
  width: 60%;
  padding: 8px;
  font-size: 16px;
}
.search-bar button {
  padding: 8px 16px;
  margin-left: 10px;
  font-size: 16px;
  cursor: pointer;
}
.user-info {
  margin-top: 20px;
}
.user-photo {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  object-fit: cover;
  margin-top: 10px;
}
.error-message {
  color: red;
  font-weight: bold;
  margin-top: 10px;
}
</style>