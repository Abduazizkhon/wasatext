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

        // ✅ Store token and user info properly
        if (response.data.token && response.data.user) {
          localStorage.setItem("authToken", response.data.token);
          localStorage.setItem("username", response.data.user.username);
          localStorage.setItem("userID", response.data.user.id);
        } else {
          throw new Error("Invalid login response");
        }

        // ✅ Redirect to home page
        this.$router.push("/home");
      } catch (error) {
        this.errorMessage = error.response?.data?.error || "An error occurred.";
      }
    },
  },
};
</script>

<style scoped>
.login-container {
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

.error-message {
  margin-top: 20px;
  color: red;
  font-weight: bold;
}
</style>