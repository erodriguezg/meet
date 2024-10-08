<script setup lang="ts">
import axios from 'axios'
import Dialog from 'primevue/dialog'
import FileUpload, { FileUploadUploaderEvent } from 'primevue/fileupload'
import ProgressSpinner from 'primevue/progressspinner'
import { ref } from 'vue'
import { FileApi } from '../../api/FileApi'
import { ResourceUploadUrlDto } from '../../api/domain/UploadFile'
import { ImageVideoUtil } from '../../utils/ImageVideoUtil'

const axiosUploader = axios.create({ headers: { Authorization: undefined } })
const visible = ref<boolean>(false)
const uploading = ref<boolean>(false)
const uploadSupplierFn = ref<(contentType: string) => Promise<ResourceUploadUrlDto[]>>()
const uploadCompleteFn = ref<() => Promise<void>>()

let storageType: string | undefined

const widthThumbnail = 200

const open = async (
  uploadSupplier: (contentType: string) => Promise<ResourceUploadUrlDto[]>,
  uploadComplete: () => Promise<void>
) => {
  uploadSupplierFn.value = uploadSupplier
  uploadCompleteFn.value = uploadComplete
  uploading.value = false
  visible.value = true
  storageType = await FileApi.getStorageType()
}

defineExpose({
  open
})

const uploader = async (event: FileUploadUploaderEvent): Promise<void> => {
  uploading.value = true

  const file = event.files instanceof File ? event.files : event.files[0]

  const contentType = file.type

  if (uploadSupplierFn.value === undefined) {
    return
  }

  const toUploadsResources = await uploadSupplierFn.value(contentType)

  for (const uploadResource of toUploadsResources) {
    let blobEdited: Blob
    if (ImageVideoUtil.isImageSupportedContentType(contentType)) {
      blobEdited = await ImageVideoUtil.convertImageFileToBlob(file, widthThumbnail, uploadResource.isThumbnail, uploadResource.isBlurred)
    } else if (ImageVideoUtil.isVideoSupportedContentType(contentType)) {
      blobEdited = await ImageVideoUtil.convertVideoFileToBlob(file, widthThumbnail, uploadResource.isThumbnail, uploadResource.isBlurred)
    } else {
      throw new Error(`content type not supported: ${contentType}`)
    }
    await uploadBlob(uploadResource.uploadUrl, blobEdited)
    await FileApi.confirmFileUpload(uploadResource.fileHash)
  }

  uploading.value = false

  if (uploadCompleteFn.value !== undefined) {
    await uploadCompleteFn.value()
  }

  close()
}

const close = (): void => {
  uploadSupplierFn.value = undefined
  uploadCompleteFn.value = undefined
  visible.value = false
}

const uploadBlob = async (url: string, blob: Blob) => {
  if (storageType === 'S3') {
    await uploadBlobS3(url, blob)
  } else if (storageType === 'DROPBOX') {
    await uploadBlobDropbox(url, blob)
  }
}

const uploadBlobS3 = async (url: string, blob: Blob) => {
  try {
    const response = await axiosUploader.put(url, blob)
    console.log(response.data)
  } catch (error) {
    console.error(error)
  }
}

const uploadBlobDropbox = async (url: string, blob: Blob) => {
  try {
    const response = await axiosUploader.post(url, blob, {
      headers: {
        'Content-Type': 'application/octet-stream'
      }
    })
    console.log(response.data)
  } catch (error) {
    console.error(error)
  }
}

</script>

<template>
  <div class="UploadDialog">

    <Dialog v-model:visible="visible" modal header="Subir Archivo">

      <FileUpload v-if="!uploading" mode="basic" name="fileUpload" customUpload @uploader="uploader" auto />

      <template v-if="uploading">
        <ProgressSpinner style="width: 50px; height: 50px" strokeWidth="8" fill="var(--surface-ground)"
          animationDuration=".5s" aria-label="Custom ProgressSpinner" />
      </template>

    </Dialog>

  </div>
</template>

<style lang="scss"></style>
