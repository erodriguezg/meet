<script setup lang="ts">
import { ref } from 'vue'
import RoutesNames from '../../../utils/routes-names'
import { AuthService, Permission } from '../../../services/AuthService'

const seeModelRegisterLink = ref<boolean>(true)

seeModelRegisterLink.value = AuthService.hasPermission(Permission.MANAGE)

</script>

<template>
    <div class="footerTemplate">
        <p>version: 0.4.0</p>
        <div v-if="seeModelRegisterLink">
            <RouterLink :to="{name: RoutesNames.NEW_PUBLICATION_PAGE}"  >registro de modelo</RouterLink>
        </div>
    </div>
</template>

<style lang="scss" scoped>
    .footerTemplate {
        bottom: 10px;
        width: 100%;
        font-size: 0.9rem;
        text-align: center;
    }
</style>
