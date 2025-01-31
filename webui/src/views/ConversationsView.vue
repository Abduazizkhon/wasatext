<template>
  <div class="conversations-container">
    <h1>Conversations</h1>

    <div v-if="isLoading" class="loading">Loading...</div>
    <div v-else>
      <div v-if="conversations.length > 0">
        <h2>Private Chats</h2>
        <ul class="conversation-list">
          <li v-for="convo in privateConversations" :key="convo.id" @click="openConversation(convo.id)">
            <img v-if="convo.photo" :src="'http://localhost:3000' + convo.photo" alt="User Photo" class="profile-photo" />
            <span>{{ convo.name }}</span>
          </li>
        </ul>

        <h2>Group Chats</h2>
        <ul class="conversation-list">
          <li v-for="convo in groupConversations" :key="convo.id" @click="openConversation(convo.id)">
            <img v-if="convo.photo" :src="'http://localhost:3000' + convo.photo" alt="Group Photo" class="profile-photo" />
            <span>{{ convo.name }}</span>
          </li>
        </ul>
      </div>
      <div v-else class="no-conversations">No conversations found.</div>
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
      isLoading: true,
      isLoggedIn: false,
    };
  },
  computed: {
    privateConversations() {
      return this.conversations.filter(convo => !convo.is_group);
    },
    groupConversations() {
      return this.conversations.filter(convo => convo.is_group);
    }
  },
  created() {
    const token = localStorage.getItem("authToken");
    const userID = localStorage.getItem("userID");

    if (!token || !userID) {
      this.$router.push("/"); // Redirect to login if not authenticated
      return;
    }

    this.isLoggedIn = true;
    this.fetchConversations(userID, token);
  },
  methods: {
    async fetchConversations(userID, token) {
      try {
        const response = await axios.get(`http://localhost:3000/users/${userID}/conversations`, {
          headers: { Authorization: `Bearer ${token}` },
        });
        this.conversations = response.data;
      } catch (error) {
        console.error("Error fetching conversations:", error);
      } finally {
        this.isLoading = false;
      }
    },
    openConversation(conversationId) {
      this.$router.push(`/conversations/${conversationId}`);
    }
  }
};
</script>

<style scoped>
.conversations-container {
  text-align: center;
  margin-top: 20px;
}

.loading {
  font-size: 18px;
  color: gray;
}

.no-conversations {
  color: red;
  font-weight: bold;
}

.conversation-list {
  list-style: none;
  padding: 0;
}

.conversation-list li {
  display: flex;
  align-items: center;
  padding: 10px;
  cursor: pointer;
  border-bottom: 1px solid #ddd;
  transition: background 0.3s;
}

.conversation-list li:hover {
  background: #f0f0f0;
}

.profile-photo {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  margin-right: 10px;
  border: 2px solid #007bff;
}
</style>