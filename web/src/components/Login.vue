<template>
    <div>
        <div>Username:</div>
        <div><input type="text" v-model="username" /></div>
        <div>Password</div>
        <div><input type="password" v-model="password" /></div>
        <div class="sb"><button @click="login">login</button></div>
    </div>
</template>
<style>
.sb {
    text-align: center;
    padding-top: 10px;
}
</style>
<script>
import md5 from 'js-md5'
export default {
    name: "Login",
    data() {
        return {
            username: '',
            password: '',
        }
    },

    methods: {
        login() {
            this.axios({
                method: 'post',
                url: '/login',
                data: ({
                    username: this.username,
                    password: md5(this.password),
                })
            }).then(response => {
                localStorage.setItem('token', response.data['data']['token'])
                localStorage.setItem('username', this.username)
                this.$router.push({
                    name: 'Profile',
                    params: {
                        username: this.username,
                    }
                })
            }).catch(error => {
                alert('Wrong username or password')
                this.password = ""
            });
        }
    }
}

</script>