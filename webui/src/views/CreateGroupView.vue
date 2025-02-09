<template>
  <div class="create-group-container">
    <h1>Create a New Group</h1>

    <form @submit.prevent="createGroup">
      <div class="form-group">
        <label for="groupName">Group Name</label>
        <input type="text" id="groupName" v-model="groupName" placeholder="Enter group name" required />
      </div>

      <div class="form-group">
        <label for="groupPhoto">Group Photo (Optional)</label>
        <input type="file" id="groupPhoto" @change="handlePhotoChange" />
      </div>

      <div class="form-group">
        <label for="usernames">Add Members</label>
        <input type="text" id="usernames" v-model="usernamesInput" placeholder="Enter usernames (comma separated)" required />
      </div>

      <button type="submit" :disabled="loading">Create Group</button>
    </form>

    <div v-if="message" :class="messageClass">
      {{ message }}
    </div>
  </div>
</template>

<script>
import axios from '../services/axios.js';

export default {
  name: "CreateGroupView",
  data() {
    return {
      groupName: '',
      groupPhoto: null,
      usernamesInput: '',
      loading: false,
      message: '',
      messageClass: ''
    };
  },
  methods: {
    handlePhotoChange(event) {
      this.groupPhoto = event.target.files[0];
    },

    async createGroup() {
      this.loading = true;
      this.message = '';
      this.messageClass = '';

      const token = localStorage.getItem('authToken');
      const userID = localStorage.getItem('userID');
      if (!token || !userID) {
        this.message = 'User not authenticated';
        this.messageClass = 'error';
        return;
      }

      const usernames = this.usernamesInput.split(',').map(username => username.trim());
      const formData = new FormData();
      formData.append('group_name', this.groupName);
      if (this.groupPhoto) formData.append('photo', this.groupPhoto);
      formData.append('usernames', JSON.stringify(usernames));  // Make sure this is a JSON string

      try {
        // Use a relative URL so the axios instance's baseURL is applied automatically.
        const response = await axios.post('/groups', formData, {
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'multipart/form-data',
          },
        });
        this.loading = false;
        this.message = `Group created successfully. Group ID: ${response.data.c_id}`;
        this.messageClass = 'success';
      } catch (error) {
        this.loading = false;
        this.message = 'Failed to create group. Please try again.';
        this.messageClass = 'error';
        console.error("Error creating group:", error);
      }
    }
  }
};
</script>

<style scoped>
.create-group-container {
  text-align: center;
  padding: 20px;
}

.form-group {
  margin-bottom: 15px;
}

input[type="text"],
input[type="file"] {
  width: 100%;
  padding: 8px;
  margin: 10px 0;
  border: 1px solid #ccc;
  border-radius: 4px;
}

button {
  padding: 10px 20px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
}

button:disabled {
  background-color: #aaa;
}

button:hover {
  background-color: #0056b3;
}

.message.success {
  color: green;
}

.message.error {
  color: red;
}
</style>