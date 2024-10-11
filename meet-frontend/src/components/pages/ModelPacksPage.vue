<script setup lang="ts">
import Button from 'primevue/button'
import DataView from 'primevue/dataview'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { FileApi } from '../../api/FileApi'
import { PackApi } from '../../api/PackApi'
import { PackDto } from '../../api/domain/Pack'
import { useModelDetailStore } from '../../stores/ModelDetailStore'
import RoutesNames from '../../utils/routes-names'

const modelDetailStore = useModelDetailStore()

const packs = ref<PackDto[]>()

const router = useRouter()

const modelNickName = ref<string>()
const canCreatePack = ref<boolean>(false)

onMounted(async () => {
  modelNickName.value = modelDetailStore.modelDetail?.modelNickName ?? ''

  const isSameModel = modelDetailStore.modelDetail?.isSameModel
  const havePermissionEditAllModels = modelDetailStore.modelDetail?.havePermissionEditAllModels ?? false
  canCreatePack.value = isSameModel || havePermissionEditAllModels

  packs.value = await PackApi.getPacksFromModel(modelNickName.value)
})

const goToPack = (packNumber: number) => {
  router.push({
    name: RoutesNames.HOME_PAGE,
    params: {
      packNumber
    }
  })
}

const createNewPack = async () => {
  const newPack = await PackApi.createNewPack(modelNickName.value ?? '')
  goToPack(newPack.packNumber)
}

</script>

<template>
  <div class="ModelPacks">

    <div>
      <Button v-if="canCreatePack" label="Nuevo Pack" @click="createNewPack" />
    </div>
    <DataView :value="packs" layout="grid" paginatorPosition="bottom" :paginator="true" :rows="8" data-key="order"
      :always-show-paginator="false">
      <template #grid="slotProps">
        <div v-for="(item, index) in slotProps.items" :key="index" class="col-12 md:col-4 lg:col-2 item">
          <p>Pack {{ item.packNumber }}</p>
          <div class="pack-img" @click="() => { goToPack(item.packNumber) }">

            <img v-if="item.coverImageFileHash === undefined"
              src="https://i.discogs.com/7Yic1tgwSvxVB3yb1tLcds_zLtvscwchE8M0gjmO9uc/rs:fit/g:sm/q:90/h:594/w:600/czM6Ly9kaXNjb2dz/LWRhdGFiYXNlLWlt/YWdlcy9SLTExMzAy/NjItMTU0MDU4MDA4/OS03OTMyLmpwZWc.jpeg" />

            <img v-if="item.coverImageFileHash !== undefined" :src="FileApi.getRedirectUrl(item.coverImageFileHash)" />

          </div>
        </div>
      </template>
      <template #empty>No se encontraron packs.</template>
    </DataView>

  </div>
</template>

<style lang="scss">
.ModelPacks {

  .item {
    padding: .5em;

    .pack-img img {
      width: 200px;
    }
  }

}
</style>
