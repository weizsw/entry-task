<template>
    <div>Username: {{ username }}</div>
    <div>Nickname: {{ nicknameDisplay }}</div>
    <div><input type="text" v-model="nickname"> <button @click="changeNickName">Change</button></div>
    <div class="dp"><img :src="pic" alt=""></div>
    <div><input type="file" ref="fileInt" accept="image/png, image/jpeg" @change="changePic">
    </div>
</template>
<style>
img {
    width: 300px;
    height: 300px;
}
</style>
<script>
export default {
    name: "Profile",
    data() {
        return {
            username: '',
            nicknameDisplay: '',
            pic: '',
        }
    },
    created() {
        this.$watch(
            () => localStorage.getItem('username'),
            () => {
                this.getUserInfo()
            },
            { immediate: true }
        )

    },
    methods: {
        getUserInfo() {
            this.axios({
                method: 'post',
                url: '/profile',
                data: ({
                    username: localStorage.getItem('username'),
                    token: localStorage.getItem('token'),
                })
            }).then(response => {
                this.username = response.data['data'].username
                this.nicknameDisplay = response.data['data'].nickname
                this.pic = response.data['data'].pic
            }).catch(error => {
                console.log(error)
            });
        },
        changeNickName() {
            this.axios({
                method: 'post',
                url: '/nickname',
                data: ({
                    username: localStorage.getItem('username'),
                    nickname: this.nickname,
                    token: localStorage.getItem('token'),
                })
            }).then(response => {
                this.nicknameDisplay = this.nickname
                this.nickname = ""
            }).catch(error => {
                console.log(error)
            });
        },
        changePic() {
            const file = this.$refs.fileInt.files[0]
            const param = new FormData()
            param.append('pic', file)
            param.append('username', localStorage.getItem('username'))
            param.append('token', localStorage.getItem('token'))

            this.axios({
                method: 'post',
                url: '/pic',
                headers: {
                    'Content-Type': 'multipart/form-data',
                },
                data: param,
            }).then(response => {
                this.pic = response.data['data']['pic']
                console.log(this.pic)
                location.reload()
            }).catch(error => {
                console.log(error)
            });
        }
    }
}
</script>