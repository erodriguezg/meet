<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { AuthService, Permission } from '../../services/AuthService'
import RoutesNames from '../../utils/routes-names'

const route = useRoute()
const router = useRouter()

const code = route.query.code as string
const state = route.query.state as string

onMounted(async () => {
  try {
    await AuthService.processLoginCallback(code, state)
    if (AuthService.hasPermission(Permission.MANAGE_SYSTEM)) {
      router.push({ name: RoutesNames.MANAGE_ROOMS })
    } else if (AuthService.hasPermission(Permission.CREATE_ROOM)) {
      router.push({ name: RoutesNames.MANAGE_OWN_ROOMS })
    } else {
      router.push({ name: RoutesNames.LOGIN_PAGE })
    }
  } catch (e) {
    console.log(e)
    router.push({ name: RoutesNames.LOGIN_PAGE })
  }
})

</script>

<template>
  <p>Cargando ...</p>
</template>
