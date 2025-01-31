<script setup>
import { RouterLink, RouterView, useRouter } from 'vue-router'
import { ref, onMounted } from 'vue'

const router = useRouter();
const isAuthenticated = ref(false);

onMounted(() => {
  const token = localStorage.getItem("authToken");
  isAuthenticated.value = !!token; // ✅ Checks login status
});

const logout = () => {
  localStorage.clear();
  isAuthenticated.value = false;
  router.push("/login"); // ✅ Redirects to login
  window.location.reload(); // 🚀 **Forces UI update**
};
</script>

<template>
	<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
		<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="/">WASAText</a>
	</header>

	<div class="container-fluid">
		<div class="row">
			<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
				<div class="position-sticky pt-3 sidebar-sticky">
					<h6 class="sidebar-heading px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>General</span>
					</h6>
					<ul class="nav flex-column">
						<li class="nav-item">
							<RouterLink to="/home" class="nav-link">Home</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink 
								to="/conversations" 
								class="nav-link" 
								:class="{ disabled: !isAuthenticated }">
								Chats
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink 
								to="/groups" 
								class="nav-link" 
								:class="{ disabled: !isAuthenticated }">
								Groups
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink to="/users/me/username" class="nav-link" :class="{ disabled: !isAuthenticated }">
								Profile
							</RouterLink>
						</li>
						<li v-if="!isAuthenticated" class="nav-item">
							<RouterLink to="/login" class="nav-link">Login</RouterLink>
						</li>
						<li v-if="isAuthenticated" class="nav-item">
							<button class="nav-link logout-button" @click="logout">Logout</button>
						</li>
					</ul>
				</div>
			</nav>

			<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
				<RouterView />
			</main>
		</div>
	</div>
</template>

<style>
.disabled {
  pointer-events: none;
  opacity: 0.5;
}
.logout-button {
  background: none;
  border: none;
  color: red;
  cursor: pointer;
  font-size: 16px;
  padding: 10px;
  width: 100%;
  text-align: left;
}
.logout-button:hover {
  text-decoration: underline;
}
</style>