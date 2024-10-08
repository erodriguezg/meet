<script setup lang="ts">
import DataView, { DataViewPageEvent } from 'primevue/dataview'
import Skeleton from 'primevue/skeleton'
import ModelSearchSideBar from '../shared/ModelSearchSideBar.vue'
import { FilterSearchModel, Model } from '../../api/domain/Model'
import { ref } from 'vue'
import { ModelApi } from '../../api/ModelApi'
import { FileApi } from '../../api/FileApi'

const pageRef = ref<{
  loading: boolean,
  first: number,
  rows: number,
  totalCount: number,
  filters?: FilterSearchModel
  models?: Model[]
}>({
  loading: false,
  first: 0,
  rows: 10,
  totalCount: 0
})

const filtersListener = async (filters: FilterSearchModel): Promise<void> => {
  pageRef.value.filters = filters
  pageRef.value.first = 0
  await loadModels()
}

const lazyLoadModels = async (event: DataViewPageEvent): Promise<void> => {
  pageRef.value.first = event.first
  await loadModels()
}

const loadModels = async (): Promise<void> => {
  pageRef.value.loading = true
  const first = pageRef.value.first
  const last = first + pageRef.value.rows
  const filters = pageRef.value.filters ?? {}
  const response = await ModelApi.searchModels(filters, first, last)
  pageRef.value.totalCount = response.totalCount
  pageRef.value.models = response.models
  pageRef.value.loading = false
}

loadModels()

</script>

<template>
  <ModelSearchSideBar @filters-changed="filtersListener" />
  <div class="ModelsPage">
    <h2>Modelos</h2>
    <DataView :value="pageRef.models" layout="grid" paginatorPosition="bottom" :paginator="true" :rows="pageRef.rows"
      data-key="nickName" :total-records="pageRef.totalCount" :lazy="true" @page="lazyLoadModels">

      <template #grid="slotProps">
        <template v-if="pageRef.loading">
          <div style="padding: .5em" class="col-12 md:col-4 lg:col-2 item">
            <div class="flex flex-column align-items-center gap-3 py-5">
              <Skeleton class="w-9 shadow-2 border-round h-10rem" />
            </div>
          </div>
        </template>

        <template v-if="!pageRef.loading">
          <div v-for="(item, index) in slotProps.items" :key="index" class="col-12 md:col-4 lg:col-2 item" >

            <a class="model-link" :href="'model/' + item.nickName" >
              <div class="img-container">

                <img v-if="item.profileImageThumbnailFileHash === undefined"
                  src="../../assets/unknown_model_profile.png" />

                <img v-if="item.profileImageThumbnailFileHash !== undefined"
                  :src="FileApi.getRedirectUrl(item.profileImageThumbnailFileHash)" />

                <div class="info">{{ item.nickName }}</div>

              </div>
            </a>
          </div>

        </template>

      </template>

      <template #empty>No se encontraron modelos.</template>
    </DataView>
  </div>
</template>

<style lang="scss" scoped>
@use '../../styles/vars' as vars;

.ModelsPage {
  width: 100%;
  padding-left: 1.5rem;
  padding-right: 1.5rem;

  @media only screen and (max-width: vars.$mobile) {
    margin-top: 2rem;
  }

  .item {
    padding: .5em;
    display: inline-table;
  }

  .model-link {

    cursor: pointer;

    div.img-container {
      position: relative;
      width: 239.500px;
      height: 179.625px;
      overflow: hidden;

      img {
        object-fit: cover;
        object-position: top;
        width: 100%;
        height: 100%;
      }

      .info {
        position: absolute;
        color: white;
        background-color: rgb(222 86 86 / 60%);
        left: 0;
        bottom: 0;
        width: 100%;
        padding-left: 10px;
      }
    }
  }

}
</style>
