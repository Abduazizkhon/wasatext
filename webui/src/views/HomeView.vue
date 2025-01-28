<template>
  <div class="home-container">
    <h1>Home Page</h1>
    <p>Welcome, {{ username }}!</p>

    <!-- Profile Picture -->
    <img v-if="profilePhoto" :src="'http://localhost:3000' + profilePhoto" alt="Profile Photo" class="profile-photo" />

    <RouterLink to="/users/me/username" class="nav-button">Profile</RouterLink>
    <RouterLink to="/users/me/photo" class="photo-button">Upload Profile Photo</RouterLink>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "HomeView",
  data() {
    return {
      username: "",
      profilePhoto: localStorage.getItem("profilePhoto") || "", // ✅ Load from localStorage immediately
    };
  },
  created() {
    const storedUsername = localStorage.getItem("username");
    this.username = storedUsername || "Guest";
    this.fetchProfilePhoto();
  },
  methods: {
    async fetchProfilePhoto() {
      const token = localStorage.getItem("authToken");
      try {
        const response = await axios.get("http://localhost:3000/users/me", {
          headers: { Authorization: `Bearer ${token}` }
        });

        if (response.data.photo) {
          this.profilePhoto = response.data.photo;
          localStorage.setItem("profilePhoto", response.data.photo); // ✅ Update localStorage
        }
      } catch (error) {
        console.error("Error fetching profile photo:", error);
      }
    }
  }
};
</script>

<style scoped>
.home-container {
  text-align: center;
  margin-top: 50px;
}

.profile-photo {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid #007bff;
}

.nav-button, .photo-button {
  display: block;
  margin: 15px auto;
  padding: 10px 15px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  text-align: center;
  text-decoration: none;
  width: 200px;
}

.nav-button:hover, .photo-button:hover {
  background-color: #0056b3;
}
</style>