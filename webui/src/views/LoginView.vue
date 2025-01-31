<template>
  <div class="login-container">
    <h1>Login</h1>
    <form @submit.prevent="doLogin">
      <div class="form-group">
        <label for="username">Username</label>
        <input
          type="text"
          id="username"
          v-model="username"
          placeholder="Enter your username"
          required
        />
      </div>
      <button type="submit">Login</button>
    </form>
    <div v-if="errorMessage" class="error-message">
      {{ errorMessage }}
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "LoginView",
  data() {
    return {
      username: "",
      errorMessage: "",
    };
  },
  methods: {
    async doLogin() {
      try {
        const response = await axios.post("http://localhost:3000/session", {
          username: this.username,
        });

        if (response.status === 201) {
          alert("User created successfully!");
        } else if (response.status === 200) {
          alert("Login successful!");
        }

        if (response.data.token && response.data.user) {
          localStorage.setItem("authToken", response.data.token);
          localStorage.setItem("username", response.data.user.username);
          localStorage.setItem("userID", response.data.user.id);

          // ✅ Restore profile photo from localStorage if it exists
          const profilePhoto = localStorage.getItem(`profilePhoto_${response.data.user.id}`);
          if (profilePhoto) {
            localStorage.setItem("profilePhoto", profilePhoto);
          }
        } else {
          throw new Error("Invalid login response");
        }

        // ✅ Redirect to home page
        this.$router.push("/home");
      } catch (error) {
        this.errorMessage = error.response?.data?.error || "An error occurred.";
      }
    }
  }
};
</script>