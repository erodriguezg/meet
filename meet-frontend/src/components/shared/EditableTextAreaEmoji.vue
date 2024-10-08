<script setup lang="ts">

import Button from 'primevue/button'
import Textarea from 'primevue/textarea'
import EmojiPicker from 'vue3-emoji-picker'
import OverlayPanel from 'primevue/overlaypanel'
import 'vue3-emoji-picker/css'
import { ref, watch } from 'vue'

const props = defineProps({
  text: String,
  placeholder: String,
  editable: Boolean,
  maxlength: Number
})
const emit = defineEmits(['accepted'])

const op = ref()
const value = ref<string>(props.text ?? '')
const oldValue = ref<string>('')
const editMode = ref<boolean>(false)

watch(() => props.text, (newValue, oldValue) => {
  value.value = newValue ?? ''
})

const toggle = (event: any) => {
  op.value.toggle(event)
}

const onSelectEmoji = (emoji: any) => {
  value.value += emoji.i
}

const updateValue = (event: any) => {
  value.value = event.target.value
}

const editAction = () => {
  oldValue.value = value.value
  editMode.value = true
}

const cancelAction = () => {
  value.value = oldValue.value
  editMode.value = false
}

const acceptAction = () => {
  editMode.value = false
  emit('accepted', value.value)
}

</script>

<template>
  <div class="EditableTextAreaEmoji">

    <div v-if="!editMode" class="no-editable">
      <div class="title">
        <template v-if="value.length === 0">
          {{ props.placeholder }}
        </template>

        <template v-if="value.length > 0">
          {{ value }}
        </template>

      </div>
      <div class="actions">
        <span title="edit" class="pi pi-pencil" @click="editAction" v-if="props.editable" />
      </div>
    </div>

    <div v-if="editMode" class="editable">
      <div class="wrapper">
        <Textarea v-model:value="value" @input="updateValue" :placeholder="props.placeholder"
          :maxlength="props.maxlength" autoResize rows="5" cols="30" />
        <Button @click="toggle">&#128578;</Button>
        <OverlayPanel ref="op" :show-close-icon="true">
          <EmojiPicker :native="true" @select="onSelectEmoji" theme="dark" />
        </OverlayPanel>
      </div>
      <div class="actions">
        <span title="cancel" class="pi pi-times" @click="cancelAction" />
        <span title="accept" class="pi pi-check" @click="acceptAction" />
      </div>
    </div>
  </div>
</template>

<style lang="scss">
.EditableTextAreaEmoji {

  .no-editable {

    .title {
      display: inline;
      text-align: justify;
    }

    .actions {
      display: inline;
      margin-left: 10px;
    }
  }

  .editable {

    .p-button {
      display: inline;
    }

    .actions {

      span {
        font-size: 1.2rem;
      }

    }

  }

}
</style>
