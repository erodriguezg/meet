<template>

    <ModelSearchSideBar />

    <div class="homePage">

        <div class="publications">

            <div v-for="i in 10" :key="i" class="publication">

                    <div class="title">
                        Title of publication {{ i }}
                    </div>
                    <div class="banner-image">
                        <img
                            src="https://t4.ftcdn.net/jpg/03/37/92/49/240_F_337924964_fF6I0iZkv9lvoIEn71f0OsgawArIPPDu.jpg" />
                    </div>

                    <div class="description">
                        orem Ipsum es simplemente el texto de relleno de las
                        imprentas y archivos de texto. Lorem Ipsum ha sido el
                        texto de relleno estándar de las industrias desde el
                        año 1500, cuando un impresor (N. del T. persona que
                        se dedica a la imprenta) desconocido usó una g
                        alería de textos y los mezcló de tal manera que
                        logró hacer un libro de textos especimen. No
                        sólo sobrevivió 500 años, sino que tambien i
                        ngresó como texto de relleno en documentos
                    </div>

                    <div class="footer"></div>

            </div>

        </div>

    </div>
</template>

<script setup lang="ts">
import ModelSearchSideBar from '../shared/ModelSearchSideBar.vue'

</script>

<style lang="scss" scoped>
.publications {
    @apply grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-10 mt-10 ml-10 mr-10;
}

.publication {
    @apply w-full aspect-square bg-gray-200 p-4 text-black shadow-md;
}

.publication .title {
    @apply text-lg font-bold;
}

.publication .banner-image img {
    @apply w-full;
}

.publication .description {
    @apply text-justify;
}
</style>
