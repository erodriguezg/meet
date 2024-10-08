import { createI18n } from 'vue-i18n'

/*
 * All i18n resources specified in the plugin `include` option can be loaded
 * at once using the import syntax
 */
import messages from '@intlify/unplugin-vue-i18n/messages'

const defaultLocale = navigator.languages[0].split('-')[0]

console.log(`defaultLocale: ${defaultLocale}`)

export const i18n = createI18n({
  locale: defaultLocale,
  fallbackLocale: 'es', // set fallback locale
  messages
})
