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

        // Redirect to the homepage or another page after login
        this.$router.push("/home");
      } catch (error) {
        if (error.response && error.response.data) {
          this.errorMessage = error.response.data.error || "An error occurred.";
        } else {
          this.errorMessage = "Unable to connect to the server.";
        }
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

button:hover {
  background-color: #0056b3;
}

.error-message {
  margin-top: 20px;
  color: red;
  font-weight: bold;
}
</style>