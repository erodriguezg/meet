<script setup lang="ts">

import Button from 'primevue/button'
import Skeleton from 'primevue/skeleton'
import { getCurrentInstance, onMounted, ref } from 'vue'
import { FileApi } from '../../api/FileApi'
import { ModelApi } from '../../api/ModelApi'
import { Model } from '../../api/domain/Model'
import { useModelDetailStore } from '../../stores/ModelDetailStore'
import UploadDialog from './UploadDialog.vue'

const props = defineProps<{
  model: Model
}>()

const modelDetailStore = useModelDetailStore()

const app = getCurrentInstance()

const readyInit = ref<boolean>(false)
const canEdit = ref<boolean>(false)
const thumbnailImgUrl = ref<string>()
const profileImgUrl = ref<string>()
const uploadModalRef = ref<InstanceType<typeof UploadDialog> | null>(null)

onMounted(async () => {
  const isSameModel = modelDetailStore.modelDetail?.isSameModel
  const havePermissionEditAllModels = modelDetailStore.modelDetail?.havePermissionEditAllModels ?? false
  canEdit.value = isSameModel || havePermissionEditAllModels

  if (props.model.profileImageThumbnailFileHash !== undefined) {
    thumbnailImgUrl.value = await FileApi.getDownloadUrl(props.model.profileImageThumbnailFileHash ?? '')
  }

  readyInit.value = true
})

const uploadImageAction = async () => {
  await uploadModalRef.value?.open(
    async (contentType) => await ModelApi.prepareUploadProfileImg(props.model.nickName),
    async () => window.location.reload()
  )
}

const imageClickAction = async () => {
  if (props.model.profileImageFileHash !== undefined) {
    profileImgUrl.value = await FileApi.getDownloadUrl(props.model.profileImageFileHash ?? '')
  }
  app?.appContext.config.globalProperties.$viewerApi({
    images: [profileImgUrl.value ?? '']
  })
}

</script>

<template>
  <div class="ModelProfilePicture">

    <template v-if="!readyInit">
      <Skeleton size="200px" class="mr-2"></Skeleton>
    </template>

    <template v-if="readyInit">

      <div class="img-wrapper">

        <div class="info">
          <span>{{ props.model.nickName }}</span>
        </div>

        <template v-if="thumbnailImgUrl !== undefined">
          <img alt="profile-image" class="profileImg" :src="thumbnailImgUrl" @click="imageClickAction" />
        </template>

        <template v-if="thumbnailImgUrl === undefined">
          <img alt="profile-image" class="profileImg" src="../../assets/unknown_model_profile.png" />
        </template>

      </div>

      <template v-if="thumbnailImgUrl === undefined">
        <div class="upload-btn">
            <Button label="Subir Imagen" v-if="canEdit" @click="uploadImageAction" icon="pi pi-camera" />
          </div>
      </template>

      <UploadDialog ref="uploadModalRef" />
    </template>
  </div>
</template>

<style lang="scss">
.ModelProfilePicture {

  .img-wrapper {
    height: 170px;
    overflow: hidden;
    width: 230px;

    .profileImg {
      object-fit: cover;
      object-position: top;
      width: 100%;
      height: 100%;
    }

    .info {
      width: 100%;
      font-size: 20px;
      background-color: black;
      border-top-left-radius: 60%;
      border-top-right-radius: 60%;
      text-align: center;
      color: orange;
    }

  }

  .upload-btn {
      margin-top: 10px;
  }

}
</style>
