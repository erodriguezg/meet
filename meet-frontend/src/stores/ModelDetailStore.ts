import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface StoreModelDetail {
  modelNickName: string
  isSameModel: boolean
  havePermissionEditAllModels: boolean
}

export const useModelDetailStore = defineStore('modelDetail', () => {
  const modelDetail = ref<StoreModelDetail>()

  const setModelDetail = (detail: StoreModelDetail): void => {
    modelDetail.value = detail
  }

  return { modelDetail, setModelDetail }
})
