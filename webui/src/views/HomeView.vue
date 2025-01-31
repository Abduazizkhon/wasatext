<template>
  <div class="home-container">
    <h1>Home Page</h1>
    <p>Welcome, {{ username }}!</p>

    <!-- ✅ Display Profile Photo if Available -->
    <img v-if="profilePhoto" :src="profilePhoto" alt="Profile Photo" class="profile-photo" />

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
      username: localStorage.getItem("username") || "Guest",
      profilePhoto: "",
    };
  },
  async created() {
    await this.fetchProfilePhoto();
  },
  methods: {
    async fetchProfilePhoto() {
      const token = localStorage.getItem("authToken");
      const userID = localStorage.getItem("userID");

      if (!userID || !token) {
        console.warn("🚨 Missing user ID or token.");
        return;
      }

      try {
        console.log("🔍 Fetching profile photo from API...");
        const response = await axios.get(`http://localhost:3000/users/${userID}`, {
          headers: { Authorization: `Bearer ${token}` },
        });

        console.log("📝 Full API Response:", response.data); // 🚀 Debugging step

        if (response.status === 404) {
          console.error("🚨 User not found in database.");
          return;
        }

        if (response.data.photo) {
          console.log("📸 Raw photo value from API:", response.data.photo); // 🚀 Debugging step

          if (typeof response.data.photo === "string") {
            this.profilePhoto = `http://localhost:3000${response.data.photo}`;
            localStorage.setItem(`profilePhoto_${userID}`, this.profilePhoto);
            console.log("✅ Final profile photo URL:", this.profilePhoto);
          } else {
            console.error("❌ API photo format is incorrect (expected a string but got an object)");
            console.log("🧐 Actual type:", typeof response.data.photo, "| Value:", response.data.photo);
            this.profilePhoto = "";
          }
        } else {
          console.warn("🚨 No profile photo found for user.");
          this.profilePhoto = ""; // Default empty state
        }
      } catch (error) {
        console.error("❌ Error fetching profile photo:", error);
        this.profilePhoto = ""; // Prevent UI from breaking
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
</style>