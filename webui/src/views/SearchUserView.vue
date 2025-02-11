<template>
  <div class="search-user-container">
    <h1 class="title">Search for a User</h1>
    <div class="search-bar">
      <input
        type="text"
        v-model="searchQuery"
        placeholder="Enter username"
        class="search-input"
      />
      <button @click="searchUser" class="search-button">Search</button>
    </div>

    <div v-if="errorMessage" class="error-message">
      {{ errorMessage }}
    </div>

    <div v-if="foundUser" class="user-info">
      <p class="user-name"><strong>{{ foundUser.username }}</strong></p>
      <img
        v-if="foundUser.photo && foundUser.photo.String"
        :src="fullPhotoUrl(foundUser.photo.String)"
        alt="User Photo"
        class="user-photo"
      />
      <div class="conversation-info" v-if="conversationId">
        <p>You already have a conversation with this user.</p>
        <button @click="goToConversation" class="conversation-button">
          Go to Conversation
        </button>
      </div>
      <div class="conversation-info" v-else>
        <p>No conversation exists. Click to start a new conversation.</p>
        <button @click="startConversation" class="conversation-button">
          Start Conversation
        </button>
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
        const response = await axios.get(
          `/search/users?username=${encodeURIComponent(this.searchQuery)}`,
          {
            headers: { Authorization: `Bearer ${token}` }
          }
        );
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
        const response = await axios.post(
          `/users/${userID}/conversations/first-message`,
          formData,
          {
            headers: {
              Authorization: `Bearer ${token}`,
              "Content-Type": "multipart/form-data"
            }
          }
        );
        if (response.data && response.data.c_id) {
          this.$router.push(`/chat/${response.data.c_id}`);
        }
      } catch {
        this.errorMessage = "Error starting conversation.";
      }
    }
  }
};
</script>

<style scoped>
.search-user-container {
  max-width: 600px;
  margin: 30px auto;
  padding: 30px 20px;
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.title {
  font-size: 2em;
  margin-bottom: 20px;
  color: #333;
}

.search-bar {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-bottom: 20px;
}

.search-input {
  flex: 1;
  max-width: 400px;
  padding: 10px 15px;
  font-size: 1em;
  border: 1px solid #ccc;
  border-radius: 30px;
  outline: none;
  transition: border 0.3s ease;
}

.search-input:focus {
  border-color: #007bff;
}

.search-button {
  margin-left: 10px;
  padding: 10px 20px;
  font-size: 1em;
  background-color: #007bff;
  color: #fff;
  border: none;
  border-radius: 30px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.search-button:hover {
  background-color: #0056b3;
}

.error-message {
  color: #d9534f;
  font-weight: bold;
  margin-bottom: 20px;
}

.user-info {
  border-top: 1px solid #eee;
  padding-top: 20px;
  margin-top: 20px;
}

.user-name {
  font-size: 1.5em;
  color: #333;
  margin-bottom: 15px;
}

.user-photo {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  object-fit: cover;
  margin-bottom: 15px;
  border: 3px solid #007bff;
}

.conversation-info {
  margin-top: 20px;
}

.conversation-info p {
  font-size: 1em;
  color: #555;
  margin-bottom: 10px;
}

.conversation-button {
  padding: 10px 20px;
  font-size: 1em;
  background-color: #28a745;
  color: #fff;
  border: none;
  border-radius: 30px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.conversation-button:hover {
  background-color: #218838;
}
</style>