import HelloWorld from '../components/HelloWorld.vue'
import { createRouter, createWebHistory } from "vue-router";
import Login from '../components/Login.vue'
import Profile from '../components/Profile.vue'

const routes = [
    {
        'path': '/helloworld',
        component: HelloWorld
    },
    {
        'path': '/',
        component: Login
    },
    {
        'path': '/profile',
        name: 'Profile',
        component: Profile,
        props: true
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router
