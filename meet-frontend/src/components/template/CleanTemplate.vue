<script setup lang="ts">
import Toast from 'primevue/toast'

</script>

<template>
    <div class="center-content">

        <router-view />

    </div>

    <Toast />

</template>

<style lang="scss" scoped>
.center-content {
    display: flex;
    width: 100%;
    min-height: 86vh;
}
</style>
