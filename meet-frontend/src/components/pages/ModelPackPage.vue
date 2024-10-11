<script setup lang="ts">

import Button from 'primevue/button'
import Checkbox from 'primevue/checkbox'
import DataView from 'primevue/dataview'
import Dialog from 'primevue/dialog'
import { getCurrentInstance, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { FileApi } from '../../api/FileApi'
import { PackApi } from '../../api/PackApi'
import { PackItemDto, TypeCodeEnum } from '../../api/domain/Pack'
import { useModelDetailStore } from '../../stores/ModelDetailStore'
import RoutesNames from '../../utils/routes-names'
import UploadDialog from '../shared/UploadDialog.vue'
import BuyPackDialog from '../shared/BuyPackDialog.vue'
import { AuthService } from '../../services/AuthService'
import EditableInputTextEmoji from '../shared/EditableInputTextEmoji.vue'
import EditableTextAreaEmoji from '../shared/EditableTextAreaEmoji.vue'
import PaymentPackDefinitionDialog from '../shared/PaymentPackDefinitionDialog.vue'

const app = getCurrentInstance()
const { t } = useI18n({ useScope: 'global' })
const modelDetailStore = useModelDetailStore()
const route = useRoute()
const router = useRouter()

const uploadModalRef = ref<InstanceType<typeof UploadDialog> | null>(null)
const buyPackDialogRef = ref<InstanceType<typeof BuyPackDialog> | null>(null)
const paymentDefinitionDialogRef = ref<InstanceType<typeof PaymentPackDefinitionDialog> | null>(null)

const modelNickName = ref<string>()
const packNumber = ref<number>()
const packItems = ref<PackItemDto[]>()
const canEditPack = ref<boolean>(false)
const canAddItem = ref<boolean>(false)
const canRequestPublish = ref<boolean>(false)
const canPublish = ref<boolean>(false)

const viewModalAdd = ref<boolean>(false)
const newItemPublic = ref<boolean>(false)

const canBuyPack = ref<boolean>(false)
const packPurchased = ref<boolean>(false)

const videoDialogVisible = ref<boolean>(false)
const videoPlayingUrl = ref<string>()
const videoContentType = ref<string>()

const titlePack = ref<string>()
const descriptionPack = ref<string>()

onMounted(async () => {
  const packNumberParam = route.params.packNumber
  if (packNumberParam) {
    packNumber.value = parseInt(route.params.packNumber as string)
  }

  const modelDetail = modelDetailStore.modelDetail
  modelNickName.value = modelDetail?.modelNickName
  const isSameModel = modelDetail?.isSameModel ?? false
  canEditPack.value = isSameModel || (modelDetail?.havePermissionEditAllModels ?? false)
  canAddItem.value = isSameModel || (modelDetail?.havePermissionEditAllModels ?? false)
  canRequestPublish.value = canAddItem.value
  canPublish.value = (modelDetail?.havePermissionEditAllModels ?? false)

  const info = await PackApi.getPackInfo(modelNickName.value!, packNumber.value!)
  titlePack.value = info.title
  descriptionPack.value = info.description

  packItems.value = await PackApi.getItemsFromPack(modelNickName.value!, packNumber.value!)
  const packLocked = packItems.value !== null && packItems.value.some(pack => pack.isLocked)
  canBuyPack.value = AuthService.isAuthenticated() && !canAddItem.value && packLocked && !isSameModel
  packPurchased.value = AuthService.isAuthenticated() && !canAddItem.value && !packLocked && !isSameModel
})

const back = () => {
  router.push({ name: RoutesNames.HOME_PAGE })
}

const openModalAdd = () => {
  viewModalAdd.value = true
  newItemPublic.value = false
}

const acceptAddItemModal = async () => {
  viewModalAdd.value = false
  await uploadModalRef.value?.open(
    async (contentType: string) => {
      let typeCode = TypeCodeEnum.IMG_JPG
      if (contentType.includes('png')) {
        typeCode = TypeCodeEnum.IMG_PNG
      } else if (contentType.includes('mp4')) {
        typeCode = TypeCodeEnum.VIDEO_MP4
      } else if (contentType.includes('ogg')) {
        typeCode = TypeCodeEnum.VIDEO_OGG
      }
      return await PackApi.prepareUploadPackItem({
        isPublic: newItemPublic.value,
        modelNickName: modelNickName.value ?? '',
        packNumber: packNumber.value ?? -1,
        typeCode
      })
    },
    async () => window.location.reload()
  )
}

const imageClickAction = async (item: PackItemDto) => {
  if (item.resourceFileHash !== undefined) {
    const url = await FileApi.getDownloadUrl(item.resourceFileHash ?? '')

    switch (item.typeCode) {
      case TypeCodeEnum.IMG_JPG:
      case TypeCodeEnum.IMG_PNG:
        app?.appContext.config.globalProperties.$viewerApi({
          images: [url]
        })
        break
      case TypeCodeEnum.VIDEO_MP4:
        videoPlayingUrl.value = url
        videoDialogVisible.value = true
        videoContentType.value = 'video/mp4'
        break
      case TypeCodeEnum.VIDEO_OGG:
        videoPlayingUrl.value = url
        videoDialogVisible.value = true
        videoContentType.value = 'video/ogg'
        break
      default:
    }
  }
}

const openBuyDialog = async () => {
  const modelNickNameVal = modelNickName.value ?? ''
  const packNumberVal = packNumber.value ?? -1
  const personIdVal = AuthService.getIdentity()?.personId ?? ''
  await buyPackDialogRef.value?.open(modelNickNameVal, packNumberVal, personIdVal)
}

const saveTitle = async (newTitle: string) => {
  await PackApi.editPackTitle(modelNickName.value!, packNumber.value!, newTitle)
  titlePack.value = newTitle
}

const saveDescription = async (newDescription: string) => {
  await PackApi.editPackDescription(modelNickName.value!, packNumber.value!, newDescription)
  descriptionPack.value = newDescription
}

const openPaymentDefinitionDialog = async () => {
  paymentDefinitionDialogRef.value!.open(packNumber.value!)
}

</script>

<template>
  <div class="Pack">
    <div class="head">
      <span class="status">
        <em v-if="packPurchased && !canEditPack" class="pi pi-lock-open" v-tooltip="'pack desbloqueado'"></em>
        <em v-if="!packPurchased && !canEditPack" class="pi pi-lock" v-tooltip="'pack bloqueado'"></em>
      </span>
      <span class="number">Pack N° {{ packNumber }} </span>
      <div class="name">
        <EditableInputTextEmoji placeholder="título del pack" :text="titlePack" :editable="canEditPack" :maxlength="30"
          @accepted="saveTitle" />
      </div>
    </div>

    <div class="description">
      <EditableTextAreaEmoji placeholder="descripción del pack" :text="descriptionPack" :editable="canEditPack"
        :maxlength="280" @accepted="saveDescription" />
    </div>

    <div class="button-panel">
      <Button v-if="canEditPack" label="Definir Modo Pago" icon="pi pi-dollar" @click="openPaymentDefinitionDialog" />
      <Button v-if="canBuyPack" label="Comprar Pack" icon="pi pi-dollar" @click="openBuyDialog" />
      <Button v-if="canAddItem" label="Agregar Imagen / Video" icon="pi pi-plus" @click="openModalAdd" />
      <Button v-if="canRequestPublish" label="Solicitar Publicación Pack" icon="pi pi-eye" />
      <Button v-if="canPublish" label="Publicar Pack" icon="pi pi-check" />
    </div>

    <DataView :value="packItems" class="pack-data-view" layout="grid" paginatorPosition="bottom" :paginator="true"
      :rows="8" data-key="order" :always-show-paginator="false">
      <template #grid="slotProps">
        <div v-for="(item, index) in slotProps.items" :key="index" class="col-12 md:col-3 item">
          <div class="img-container" @click="imageClickAction(item)">
            <img :src="FileApi.getRedirectUrl(item.thumbnailFileHash)" />
          </div>
        </div>

      </template>

      <template #empty>No se encontraron items.</template>
    </DataView>

    <Button label="Volver" @click="back" link />

    <Dialog v-model:visible="viewModalAdd" header="Agregar Foto / Video" modal>
      <div class="card flex flex-wrap justify-content-center gap-3">
        <div class="flex align-items-center">
          <Checkbox v-model="newItemPublic" inputId="itemPublic" :binary="true" />
          <label for="itemPublic" class="ml-2">Público</label>
        </div>
      </div>

      <template #footer>
        <Button :label="t('cancel')" icon="pi pi-times" text @click="() => viewModalAdd = false" />
        <Button :label="t('accept')" icon="pi pi-check" autofocus @click="acceptAddItemModal" />
      </template>
    </Dialog>

    <UploadDialog ref="uploadModalRef" />
    <BuyPackDialog ref="buyPackDialogRef" />

    <div class="video-dialog">
      <Dialog v-model:visible="videoDialogVisible" modal :style="{ width: 'auto' }"
        :breakpoints="{ '960px': '75vw', '641px': '100vw' }">

        <div class="player-container">
          <video v-if="videoPlayingUrl" height="200" controls>
            <source :src="videoPlayingUrl" :type="videoContentType" />
          </video>
        </div>

      </Dialog>
    </div>

    <PaymentPackDefinitionDialog ref="paymentDefinitionDialogRef" />

  </div>
</template>

<style lang="scss">
.Pack {

  .head {

    .status {
      float: right;

      em.pi {
        font-size: 1.5rem;
      }
    }

    .name {
      font-size: 1.5rem;
      font-weight: bold;
    }
  }

  .description {
    margin-top: 15px;
  }

  .pack-data-view {
    margin-top: 20px;
  }

  .item {
    padding: .5em;
    display: inline-table;

    .img-container {
      position: relative;
      width: 100%;
      height: 140px;
      overflow: hidden;

      img {
        object-fit: cover;
        object-position: top;
        width: 100%;
        height: 100%;
      }

    }
  }

  em.pi-lock-open {
    font-size: 2rem;
    margin-left: 5px;
    color: green;
    font-weight: bold;
  }

  em.pi-lock {
    font-size: 2rem;
    margin-left: 5px;
    color: red;
    font-weight: bold;
  }

}
</style>
