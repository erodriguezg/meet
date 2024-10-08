<script setup lang="ts">
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Panel from 'primevue/panel'
import { Field, Form } from 'vee-validate'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import * as yup from 'yup'
import { ModelApi } from '../../api/ModelApi'
import { AuthService } from '../../services/AuthService'

const { t } = useI18n({ useScope: 'global' })

const loading = ref<boolean>(false)
const successDialogVisible = ref<boolean>(false)

const initialValues = {
  nickName: undefined
}

const schema = yup.object({
  nickName: yup.string().required(t('validations.required'))
})

const onSubmit = async (values: any, actions: any): Promise<void> => {
  loading.value = true
  const identity = AuthService.getIdentity()
  if (identity !== null) {
    try {
      await ModelApi.registerModel({
        nickName: values.nickName,
        personId: identity.personId
      })
      successDialogVisible.value = true
    } catch (e) {
      console.error(e)
    } finally {
      loading.value = false
    }
  }
}

const reiniciar = async (): Promise<void> => {
  successDialogVisible.value = false
  AuthService.logout()
  window.location.href = '/'
}

</script>

<template>
  <div class="ModelRegisterPage flex w-full justify-content-center mt-5">
    <Panel header="Registro de Modelo" style="width: 50rem">
      <Form :initial-values="initialValues" :validation-schema="schema" @submit="onSubmit">
        <div class="grid">
          <div class="col-12 md:col-6">
            <Field name="nickName" v-slot="{ field, errorMessage }">
              <span class="p-float-label">
                <InputText v-bind="field" :class="{ 'p-invalid': errorMessage }" />
                <label class="required">Alias de Modelo</label>
              </span>
              <small id="email-help" class="p-error">{{ errorMessage }}</small>
            </Field>
          </div>
        </div>
        <div class="button-panel">
          <Button type="button" label="Cancelar" severity="secondary" />
          <Button type="submit" label="Aceptar" severity="primary" :loading="loading" />
        </div>
      </Form>
    </Panel>

    <Dialog v-model:visible="successDialogVisible" modal header="Registro Exitoso" :style="{ width: '50vw' }" :closable="false">
      <p>
        El registro se ha realizado con éxito. Sera necesario volver iniciar sesión para que los cambios sean reflejados.
      </p>
      <template #footer>
        <Button label="Aceptar" icon="pi pi-check" @click="reiniciar" autofocus />
      </template>
    </Dialog>
  </div>
</template>

<style lang="scss" scoped></style>
