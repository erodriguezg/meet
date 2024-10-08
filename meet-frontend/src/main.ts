import PrimeVue from 'primevue/config'
import Lara from '@primevue/themes/lara'
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'

import 'viewerjs/dist/viewer.css'
import VueViewer from 'v-viewer'

import 'primevue/resources/themes/lara-light-indigo/theme.css'
import 'primeicons/primeicons.css'

import './styles/main.scss'

import ConfirmationService from 'primevue/confirmationservice'
import ToastService from 'primevue/toastservice'
import Tooltip from 'primevue/tooltip'
import { AxiosUtil } from './utils/AxiosUtil'
import { AuthService } from './services/AuthService'
import { i18n } from './i18n'

const initAsync = async (): Promise<void> => {
  AxiosUtil.initialize()
  await AuthService.initialize()

  const pinia = createPinia()
  const app = createApp(App)
  app.use(PrimeVue, {
    pt: Lara,
    ripple: true,
    zIndex: {
      modal: 1000, // dialog, sidebar
      overlay: 1100, // dropdown, overlaypanel
      menu: 1100, // overlay menus
      tooltip: 1200 // tooltip
    }
  })
  app.use(i18n)
  app.use(pinia)
  app.use(ConfirmationService)
  app.use(ToastService)
  app.use(router)
  app.use(VueViewer as any)
  app.directive('tooltip', Tooltip)
  app.mount('#app')
}

// support async top level initialization

initAsync().catch(err => {
  console.error(err)
})
