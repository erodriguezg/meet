<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { AuthService, Profile } from '../../services/AuthService'
import RoutesNames from '../../utils/routes-names'

const route = useRoute()
const router = useRouter()

const code = route.query.code as string
const state = route.query.state as string

onMounted(async () => {
  try {
    await AuthService.processLoginCallback(code, state)
    router.push({ name: RoutesNames.HOME_PAGE })
  } catch (e) {
    console.log(e)
    router.push({ name: RoutesNames.HOME_PAGE })
  }
})

</script>

<template>
  <p>Cargando ...</p>
</template>
