<script setup lang="ts">
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterView, useRoute } from 'vue-router'
import { ModelApi } from '../../api/ModelApi'
import { Model } from '../../api/domain/Model'
import { useModelDetailStore } from '../../stores/ModelDetailStore'
import ModelProfileInfo from '../shared/ModelProfileInfo.vue'
import ModelProfilePicture from '../shared/ModelProfilePicture.vue'
import { AuthService, Permission } from '../../services/AuthService'

const { t } = useI18n({ useScope: 'global' })

const modelDetailStore = useModelDetailStore()

const readyInit = ref<boolean>(false)

const route = useRoute()

const nick = ref<string>()

const modelData = ref<Model>()

onMounted(async () => {
  const nickParam = route.params.nick
  if (nickParam) {
    nick.value = route.params.nick as string
  }

  modelData.value = await ModelApi.findModelByNickName(nick.value ?? '')
  readyInit.value = true

  if (modelData.value !== undefined) {
    modelDetailStore.setModelDetail({
      modelNickName: modelData.value.nickName,
      isSameModel: modelData.value.nickName === AuthService.getIdentity()?.modelNickName,
      havePermissionEditAllModels: AuthService.hasPermission(Permission.MANAGE)
    })
  }
})

const acceptNotFoundAction = () => {
  window.location.href = '/'
}
</script>

<template>
    <div class="modelsPage">
        <template v-if="readyInit">
            <template v-if="modelData !== undefined">

                <div class="desktop">
                    <div class="container grid">
                        <div class="col-3">
                            <div class="title">
                                <ModelProfilePicture :model="modelData" />
                            </div>
                            <ModelProfileInfo />
                        </div>
                        <div class="col-9">
                            <RouterView />
                        </div>
                    </div>
                </div>
                <div class="mobile">
                    <div class="container">
                        <div class="row">
                            <div class="title">
                                <ModelProfilePicture :model="modelData" />
                            </div>
                        </div>
                        <div class="row">
                            <ModelProfileInfo />
                        </div>
                        <div class="row">
                            <RouterView />
                        </div>
                    </div>
                </div>

            </template>

            <template v-if="modelData === undefined">
                <Dialog :visible="true" modal header="Error">
                    <p>
                        {{ t("modelNotFound") }}
                    </p>
                    <template #footer>
                        <Button :label="t('accept')" icon="pi pi-check" @click="acceptNotFoundAction" />
                    </template>
                </Dialog>
            </template>
        </template>
    </div>
</template>

<style lang="scss">
.modelsPage {
    padding: 1rem 5rem 1rem 5rem;

    width: 100%;

    .title {
        display: flex;
    }
}
</style>
