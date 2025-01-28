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
      newUsername: "",
      isLoggedIn: false,
      successMessage: "",
      errorMessage: "",
    };
  },
  created() {
    const token = localStorage.getItem("authToken");
    if (token) {
      this.isLoggedIn = true;
    } else {
      this.isLoggedIn = false;
      this.$router.push("/"); // Redirect if not logged in
    }
  },
  methods: {
    async updateUsername() {
        this.errorMessage = "";
        const token = localStorage.getItem("authToken");

        if (!token) {
            this.errorMessage = "You must be logged in.";
            return;
        }

        try {
            console.log("Updating username:", this.newUsername); // Debugging log
            console.log("Auth Token:", token); // Debugging log

            const response = await axios.put(`http://localhost:3000/users/me/username?t=${new Date().getTime()}`, 
                { newname: this.newUsername },
                {
                    headers: {
                        "Authorization": `Bearer ${token}`, // ✅ Ensure correct token format
                        "Content-Type": "application/json"
                    }
                }
            );

            console.log("Response:", response); // Debugging log
            localStorage.setItem("username", this.newUsername);
            this.successMessage = "Username updated successfully!";
            this.errorMessage = "";

            setTimeout(() => {
                this.$router.push("/home");
            }, 1000);
        } catch (error) {
            console.error("Error updating username:", error);
            this.errorMessage = error.response?.data?.error || "An error occurred. Please try again.";
        }
    }
  }
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