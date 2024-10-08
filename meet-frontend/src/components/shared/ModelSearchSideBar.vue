<script setup lang="ts">
import { ref } from 'vue'
import { FilterSearchModel } from '../../api/domain/Model'
import ModelSearchFilters from './ModelSearchFilters.vue'
const emits = defineEmits<{(e: 'filters-changed', filters: FilterSearchModel): void }>()

const open = ref(true)
const openMobile = ref(false)

const switchOpenClose = () => {
  open.value = !open.value
}

const switchOpenCloseMobile = () => {
  openMobile.value = !openMobile.value
}

const filtersListener = async (filters: FilterSearchModel): Promise<void> => {
  emits('filters-changed', filters)
}

</script>

<template>
    <div class="modelSearchSideBar">

        <div class="desktop" :class="open ? 'open' : 'close'">
            <div class="switch" @click="switchOpenClose"><i class="pi"
                    :class="open ? 'pi-angle-left' : 'pi-angle-right'"></i>
            </div>
            <div class="content">
                <div class="title"><em class="pi pi-filter"></em><span>FILTROS</span></div>
                <ModelSearchFilters @filters-changed="filtersListener" />
            </div>

        </div>

        <div class="mobile" :class="openMobile ? 'open' : 'close'">
            <div class="title" @click="switchOpenCloseMobile">{{ openMobile ? 'CERRAR FILTROS' : 'FILTROS'}}</div>
            <div class="content">
                <ModelSearchFilters @filters-changed="filtersListener" />
            </div>
        </div>

    </div>
</template>

<style lang="scss" scoped>
@use '@/styles/_vars' as vars;

.modelSearchSideBar {

    .mobile {
        background-color: rgb(0 0 0 / 60%);
        width: 100%;
        height: 2rem;
        position: absolute;
        text-align: right;
        padding-right: 2rem;

        &.open {
            .content {
                display: block;
            }
        }

        &.close {
            .content {
                display: none;
            }
        }

        .content {
            background-color: rgb(0 0 0 / 85%);
            z-index: 1;
            position: absolute;
            width: 100%;
            padding-right: 12px;
        }
    }

    .desktop {
        position: relative;
        height: 40rem;
        width: 300px;
        background-color: black;
        margin-top: 3rem;
        margin-right: 1.5rem;

        .switch {
            width: 1.5rem;
            position: relative;
            float: right;
            left: 1.5rem;
            background-color: red;
            font-size: 1.2rem;
            cursor: pointer;
        }

        &.open {
            transition: 0.5s;

            .content {
                overflow-x: hidden;
                display: block;
            }

        }

        &.close {
            width: 0px;
            transition: 0.5s;

            .content {
                display: none;
            }
        }

        .content {

            display: flex;

            .title {
                text-align: center;
                font-weight: bold;
                margin-top: 10px;
                text-transform: uppercase;

                background: -webkit-linear-gradient(vars.$colorYellow, vars.$colorRed);
                -webkit-background-clip: text;
                background-clip: text;
                -webkit-text-fill-color: transparent;

                span {
                    padding-left: 5px;
                }
            }

        }
    }

}
</style>
