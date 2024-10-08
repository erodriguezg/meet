<script setup lang="ts">
import { usePrimeVue } from 'primevue/config'
import Toast from 'primevue/toast'
import { mergeProps, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { getPrimeVueLocale } from './utils/PrimeVueLocaleUtil'

const { t } = useI18n({ useScope: 'global' })

onMounted(() => {
  const primevue = usePrimeVue()
  if (primevue.config.locale !== undefined) {
    const i18nLocale = getPrimeVueLocale(t)
    primevue.config.locale = mergeProps(primevue.config.locale as any, i18nLocale) as any
  }
})

</script>

<template>
  <router-view />
  <Toast />
</template>

<style lang="scss" scoped></style>
