<script setup lang="ts">
import { ref } from 'vue'
import Editor, { EditorLoadEvent } from 'primevue/editor'
import Quill from 'quill'

const richContent = ref<string>('')
const toolbarOptions = [
  ['bold', 'italic', 'underline', 'strike'], // toggled buttons
  ['blockquote', 'code-block'],
  ['link', 'image', 'video', 'formula'],

  [{ list: 'ordered' }, { list: 'bullet' }, { list: 'check' }],
  [{ script: 'sub' }, { script: 'super' }], // superscript/subscript
  [{ indent: '-1' }, { indent: '+1' }], // outdent/indent

  [{ size: ['small', false, 'large', 'huge'] }],

  [{ color: [] }, { background: [] }], // dropdown with defaults from theme
  [{ font: [] }],
  [{ align: [] }],

  ['clean'] // remove formatting button
]
const quillModules = {
  toolbar: {
    container: toolbarOptions,
    handlers: {
      image: imageHandler
    }
  }
}

let quill:Quill | null = null

function onLoadEditor (event: EditorLoadEvent) {
  quill = (event.instance as unknown) as Quill
}

function imageHandler () {
  if (quill !== null) {
    quill.focus()
    const range = quill.getSelection()
    const value = prompt('please copy paste the image url here.')
    if (value && range) {
      quill.insertEmbed(range.index, 'image', value, Quill.sources.USER)
    }
  }
}

</script>

<template>
  <div class="editPublicationPage">

    <h2 class="title">Crear Publicación</h2>

    <div class="form-wrapper">
      <!-- Título -->
      <div>
        <label for="titulo">Título</label>
        <input type="text" id="titulo" name="titulo"
          class="mt-2 block w-full px-4 py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-lg"
          placeholder="Ingrese el título">
      </div>

      <!-- URL de imagen banner -->
      <div>
        <label for="banner-url" class="block text-lg font-medium text-gray-700">URL Imagen Banner</label>
        <input type="url" id="banner-url" name="banner-url"
          class="mt-2 block w-full px-4 py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-lg"
          placeholder="Ingrese la URL de la imagen">
      </div>

      <!-- Bajada -->
      <div>
        <label for="bajada" class="block text-lg font-medium text-gray-700">Bajada</label>
        <textarea id="bajada" name="bajada" rows="4"
          class="mt-2 block w-full px-4 py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-lg"
          placeholder="Ingrese la bajada"></textarea>
      </div>

      <!-- Contenido enriquecido -->
      <div>
        <label for="contenido" class="block text-lg font-medium text-gray-700">Contenido Enriquecido</label>
        <Editor v-model="richContent" class="custom-editor" :modules="quillModules" @load="onLoadEditor">
        </Editor>
      </div>

      <!-- Botón de enviar -->
      <div class="flex justify-end">
        <button type="submit"
          class="inline-flex items-center px-6 py-3 border border-transparent text-lg font-medium rounded-lg shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">Publicar</button>
      </div>
    </div>
  </div>

</template>

<style lang="scss" scoped>
:deep(.p-editor-toolbar) {
  display: none;
}

:deep(.custom-editor .ql-editor) {
  min-height: 200px;
  max-height: 600px;
  overflow-y: auto;
  color: #f9f9f9;
  background-color: #121212;
}

iframe {
  width: 100%;
  height: 315px;
  border: none;
}

.editPublicationPage {
  @apply w-full mx-auto p-8;
}

.title {
  @apply text-3xl font-semibold mb-6;
}

.form-wrapper {
  @apply space-y-8;

  label {
    @apply block text-lg font-medium text-white;
  }
}
</style>
