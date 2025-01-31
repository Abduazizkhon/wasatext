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

        // ✅ FORCE UI UPDATE & REDIRECT
        this.$router.push("/conversations").then(() => {
          window.location.reload(); // 🔄 Forces UI refresh
        });

      } catch (error) {
        this.errorMessage = error.response?.data?.error || "An error occurred.";
      }
    }
  }
};
</script>

<style scoped>
.login-container {
  text-align: center;
  margin-top: 50px;
}

.form-group {
  margin-bottom: 15px;
}

input {
  display: block;
  padding: 10px;
  margin: 10px auto;
  width: 250px;
  border: 1px solid #ccc;
  border-radius: 5px;
}

button {
  padding: 10px 20px;
  background-color: blue;
  color: white;
  border: none;
  cursor: pointer;
  border-radius: 5px;
}

button:hover {
  background-color: darkblue;
}

.error-message {
  color: red;
  margin-top: 10px;
}
</style>