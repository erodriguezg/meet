<script setup lang="ts">
import { ref } from 'vue'
import Button from 'primevue/button'
import TreeTable from 'primevue/treetable'
import Column from 'primevue/column'

const nodes = ref<any>([
  {
    data: {
      name: 'category a'
    },
    children: [
      {
        data: {
          name: 'category a.1'
        }
      }
    ]
  },
  {
    data: {
      name: 'category b'
    },
    children: [
      {
        data: {
          name: 'category b.1'
        }
      }
    ]
  }
])
</script>

<template>
    <div class="categoriesPage">
        <h1>Categories</h1>

        <Button label="Add Category" icon="fa fa-plus" />
        <TreeTable :value="nodes" tableStyle="min-width: 50rem">
            <Column field="name" header="Category Name" expander style="width: 34%"></Column>
            <Column header="Actions" style="width: 20%"></Column>
        </TreeTable>
    </div>
</template>

<style lang="scss" scoped>
</style>
