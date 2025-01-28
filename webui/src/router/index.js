import { createRouter, createWebHashHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from "../views/LoginView.vue"
import UsernameView from "../views/UsernameView.vue"

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/home', component: HomeView },
    { path: '/', component: LoginView },
    { 
      path: '/users/me/username', 
      component: UsernameView,
      beforeEnter: (to, from, next) => {
        const token = localStorage.getItem("authToken");
        console.log("Router Guard Check: Token =", token); // DEBUG
        if (!token) {
          next("/"); // Redirect to login if not authenticated
        } else {
          next(); // Allow access if logged in
        }
      }
    }
  ]
});

export default router;