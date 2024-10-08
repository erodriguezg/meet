import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface StoreLoginDialog {
  dialogVisible: boolean
}

export const useLoginDialogStore = defineStore('loginDialog', () => {
  const loginDialogData = ref<StoreLoginDialog>()

  const setLoginDialogData = (inData: StoreLoginDialog): void => {
    loginDialogData.value = inData
  }

  return { loginDialogData, setLoginDialogData }
})
