<script setup lang="ts">
import Toast from 'primevue/toast'
import FooterTemplate from './parts/FooterTemplate.vue'
import HeaderTemplate from './parts/HeaderTemplate.vue'
import LoginDialog from '../shared/LoginDialog.vue'

</script>

<template>
    <HeaderTemplate />

    <div class="center-content">

        <router-view />

    </div>

    <FooterTemplate />

    <Toast />

    <LoginDialog />
</template>

<style lang="scss" scoped>
.center-content {
    display: flex;
    width: 100%;
    min-height: 86vh;
}
</style>
