<script setup lang="ts">
import { onMounted } from 'vue'
import { useLoginDialogStore } from '../../stores/LoginDialogStore'
import LoginDialog from '../shared/LoginDialog.vue'

const loginDialogStore = useLoginDialogStore()

onMounted(async () => {
  loginDialogStore.setLoginDialogData({
    dialogVisible: true
  })
})

</script>

<template>
    <div class="homePage">

    </div>

    <LoginDialog />
</template>

<style lang="scss" scoped>
.publications {
    @apply grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-10 mt-10 ml-10 mr-10;
}

.publication {
    @apply w-full aspect-square bg-gray-200 p-4 text-black shadow-md;
}

.publication .title {
    @apply text-lg font-bold;
}

.publication .banner-image img {
    @apply w-full;
}

.publication .description {
    @apply text-justify;
}
</style>
