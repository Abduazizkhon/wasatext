
<template>
  <div class="profile-container">
    <h1>Profile</h1>

    <div v-if="isLoggedIn">
      <!-- Display Profile Picture -->
      <img v-if="profilePhoto" :src="'http://localhost:3000' + profilePhoto" alt="Profile Photo" class="profile-photo" />

      <input type="file" @change="handleFileUpload" accept="image/*" />
      <button @click="updatePhoto" :disabled="!selectedFile" class="update-button">Update Photo</button>

      <p v-if="errorMessage" class="error-message">{{ errorMessage }}</p>
    </div>

    <div v-else class="error-message">
      User not logged in
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "ProfileView",
  data() {
    return {
      selectedFile: null,
      profilePhoto: "",
      isLoggedIn: false,
      errorMessage: "",
    };
  },
  created() {
    const token = localStorage.getItem("authToken");
    if (token) {
      this.isLoggedIn = true;
      this.fetchProfilePhoto();
    } else {
      this.$router.push("/");
    }
  },
  methods: {
    handleFileUpload(event) {
      this.selectedFile = event.target.files[0];
    },
    async updatePhoto() {
        if (!this.selectedFile) {
            this.errorMessage = "Please select an image first.";
            return;
        }

        const token = localStorage.getItem("authToken");
        const userID = localStorage.getItem("userID"); // ✅ Get user ID

        if (!token || !userID) {
            this.errorMessage = "You must be logged in.";
            return;
        }

        const formData = new FormData();
        formData.append("photo", this.selectedFile);

        try {
            const response = await axios.put("http://localhost:3000/users/me/photo", formData, {
                headers: { Authorization: `Bearer ${token}`, "Content-Type": "multipart/form-data" }
            });

            if (response.data.photo) {
                this.profilePhoto = response.data.photo; // ✅ Update component state

                // ✅ Store photo in localStorage with user ID
                localStorage.setItem(`profilePhoto_${userID}`, response.data.photo);

                // Redirect to home page
                this.$router.push("/home");
            }
        } catch (error) {
            this.errorMessage = error.response?.data?.error || "An error occurred.";
        }
    },
    async fetchProfilePhoto() {
      const token = localStorage.getItem("authToken");
      try {
        const response = await axios.get("http://localhost:3000/users/me", {
          headers: { Authorization: `Bearer ${token}` }
        });

        if (response.data.photo) {
          this.profilePhoto = response.data.photo;
          localStorage.setItem("profilePhoto", response.data.photo); // ✅ Store it in localStorage
        }
      } catch (error) {
        console.error("Error fetching profile photo:", error);
      }
    }
  }
};
</script>

<style scoped>
.profile-container {
  text-align: center;
  margin-top: 50px;
}

.profile-photo {
  width: 150px;
  height: 150px;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid #007bff;
}

.update-button {
  margin-top: 15px;
  padding: 10px 20px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

.update-button:hover {
  background-color: #0056b3;
}

.error-message {
  color: red;
  font-weight: bold;
}
</style>