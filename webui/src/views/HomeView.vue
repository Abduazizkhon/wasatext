<template>
  <div class="home-container">
    <h1>Home Page</h1>
    <p>Welcome, {{ username }}!</p>
	<li class="nav-item">
		<RouterLink v-if="isLoggedIn" to="/users/me/username">
			<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#file-text"/></svg>
			Profile
		</RouterLink>
	</li>
  </div>
</template>

<script>
export default {
  name: "HomeView",
  data() {
    return {
      username: "",
      isLoggedIn: false, // Track if the user is logged in
    };
  },
  created() {
    const storedUsername = localStorage.getItem("username");
    const token = localStorage.getItem("authToken");

    if (token) {
      this.isLoggedIn = true;
      this.username = storedUsername || "Guest";
    } else {
      this.isLoggedIn = false;
      this.username = "Guest";
    }
  },
};
</script>

<style scoped>
.home-container {
  text-align: center;
  margin-top: 50px;
}
</style>