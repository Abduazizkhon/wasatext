<template>
  <div class="username-container">
    <h1>Update Your Username</h1>
    <form @submit.prevent="updateUsername">
      <div class="form-group">
        <label for="newname">New Username</label>
        <input
          type="text"
          id="newname"
          v-model="newUsername"
          placeholder="Enter your new username"
          required
        />
      </div>
      <!-- ✅ Button will only be disabled if user is NOT logged in OR the input is empty -->
      <button type="submit" :disabled="!isLoggedIn || newUsername.trim() === ''">Update</button>
    </form>

    <div v-if="!isLoggedIn" class="error-message">
      You must be logged in to update your username.
    </div>

    <div v-if="successMessage" class="success-message">
      {{ successMessage }}
    </div>
    <div v-if="errorMessage" class="error-message">
      {{ errorMessage }}
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "UsernameView",
  data() {
    return {
      newUsername: "", // Holds the new username input
      isLoggedIn: false, // Tracks login status
      successMessage: "", // Message to show on success
      errorMessage: "", // Message to show on error
    };
  },
  created() {
    // ✅ Ensure the user is logged in by checking for an auth token
    const token = localStorage.getItem("authToken");
    console.log("Checking login status: Token =", token); // Debugging

    if (token) {
      this.isLoggedIn = true;
    } else {
      this.isLoggedIn = false;
      this.$router.push("/"); // Redirect to login if not logged in
    }
  },
  methods: {
    async updateUsername() {
      if (!this.isLoggedIn) {
        this.errorMessage = "You must be logged in to update your username.";
        return;
      }

      try {
        const token = localStorage.getItem("authToken");
        if (!token) {
          this.errorMessage = "User not authenticated. Please log in.";
          return;
        }

        console.log("Attempting to update username:", this.newUsername); // Debugging

        // ✅ Send API request to update the username
        const response = await axios.put(
          "http://localhost:3000/users/me/username",
          { newname: this.newUsername },
          {
            headers: {
              Authorization: `Bearer ${token}`, // Include token in Authorization header
            },
          }
        );

        // ✅ Handle successful response
        this.successMessage = response.data.message;
        this.errorMessage = ""; // Clear any previous errors
        alert("Username updated successfully!");

        // ✅ Update localStorage with the new username
        localStorage.setItem("username", this.newUsername);

        // ✅ Redirect to home page after successful update
        this.$router.push("/home");
      } catch (error) {
        // Handle errors from the API
        this.successMessage = "";
        if (error.response && error.response.data) {
          this.errorMessage = error.response.data.error;
        } else {
          this.errorMessage = "An error occurred. Please try again.";
        }
      }
    },
  },
};
</script>

<style scoped>
.username-container {
  max-width: 400px;
  margin: 0 auto;
  padding: 20px;
  text-align: center;
  border: 1px solid #ddd;
  border-radius: 8px;
  box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.1);
}

.form-group {
  margin-bottom: 20px;
}

label {
  display: block;
  margin-bottom: 8px;
  font-weight: bold;
}

input {
  width: 100%;
  padding: 10px;
  font-size: 16px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

button {
  width: 100%;
  padding: 10px;
  font-size: 16px;
  color: #fff;
  background-color: #007bff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.success-message {
  margin-top: 20px;
  color: green;
  font-weight: bold;
}

.error-message {
  margin-top: 20px;
  color: red;
  font-weight: bold;
}
</style>