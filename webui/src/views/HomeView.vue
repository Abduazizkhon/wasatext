<template>
  <div class="home-container">
    <h1>Home Page</h1>
    <p>Welcome, {{ username }}!</p>
  </div>
</template>

<script>
export default {
  name: "HomeView",
  data() {
    return {
      username: "", // Username to display
    };
  },
  created() {
    // Retrieve username from localStorage
    const storedUsername = localStorage.getItem("username");
    console.log("Stored username:", storedUsername); // Debug log
    this.username = storedUsername || "Guest";
  },
};
</script>

<style scoped>
.home-container {
  text-align: center;
  margin-top: 50px;
}
</style>