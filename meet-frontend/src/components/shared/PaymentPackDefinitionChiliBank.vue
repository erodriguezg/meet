<script setup lang="ts">

import Button from 'primevue/button'
import Dropdown from 'primevue/dropdown'
import InputText from 'primevue/inputtext'
import { Field, Form } from 'vee-validate'
import * as yup from 'yup'
import { ref, onMounted, watch } from 'vue'
import { ChiliBankApi } from '../../api/ChiliBankApi'
import { ChiliBankAccountDTO } from '../../api/domain/ChiliBank'
import { SelectItem, GeneralUtils } from '../../utils/GeneralUtils'
import { useI18n } from 'vue-i18n'
import { YupUtils } from '../../utils/YupUtils'
import { RutUtils } from '../../utils/RutUtils'

const { t } = useI18n({ useScope: 'global' })

const props = defineProps<{
  modelNickname: string,
  chiliBankAccountId?: string
}>()

const emit = defineEmits<{(e: 'update:chiliBankAccountId', chiliBankAccountId: string | undefined): void}>()

const initialValues = {
  rut: null,
  holderName: null,
  bankName: null,
  accountType: null,
  accountNumber: null
}

const schema = yup.object({
  rut: yup.string()
    .required(`${t('validations.required')}`)
    .test('rut', `${t('validations.rut')}`, YupUtils.testRut),
  holderName: yup.string().required(t('validations.required')),
  bankName: yup.string().required(t('validations.required')),
  accountType: yup.string().required(t('validations.required')),
  accountNumber: yup.number().positive()
    .transform((v, o) => o === '' ? null : v)
    .required(t('validations.required'))
    .typeError(t('validations.number'))
})

const account = ref<ChiliBankAccountDTO>()

const accounts = ref<ChiliBankAccountDTO[]>([])

const editMode = ref<boolean>(false)

const chiliBanks = ref<SelectItem[]>([])

const chiliBanksTypeAccounts = ref<SelectItem[]>([])

onMounted(async () => {
  chiliBanks.value = GeneralUtils.stringsArrayToSelectItemArray(await ChiliBankApi.getBanks())
  chiliBanksTypeAccounts.value = GeneralUtils.stringsArrayToSelectItemArray(await ChiliBankApi.getAccountTypes())
  accounts.value = await ChiliBankApi.getModelAccounts(props.modelNickname!)
  account.value = accounts.value.find(ac => ac.id === props.chiliBankAccountId)
})

const addAccountBtnClick = () => {
  editMode.value = true
}

const acceptEditAccountClick = async (values: any) => {
  const accountAux = {
    rut: RutUtils.toNumber(values.rut),
    accountNumber: parseInt(values.accountNumber),
    accountType: values.accountType,
    bankName: values.bankName,
    holderName: values.holderName
  } as unknown as ChiliBankAccountDTO

  const savedAccount = await ChiliBankApi.saveModelAccount(props.modelNickname!, accountAux)
  accounts.value = await ChiliBankApi.getModelAccounts(props.modelNickname!)
  account.value = savedAccount

  editMode.value = false
}

const cancelEditAccountClick = () => {
  editMode.value = false
}

watch(account, (newAccount) => {
  emit('update:chiliBankAccountId', newAccount?.id)
})

</script>
<template>
    <div class="PaymentPackDefinitionChiliBank">

        <div class="formgrid grid" v-if="!editMode">
            <div class="field col-12 md:col-6">
                <label for="chileBankAccount">Cuenta</label>
                <Dropdown id="chileBankAccount" v-model="account" :options="accounts" placeholder="Seleccione..."
                    :show-clear="true"
                    class="text-base text-color surface-overlay p-2 border-1 border-solid surface-border border-round appearance-none outline-none focus:border-primary w-full">
                    <template #value="slotProps">
                        <div v-if="slotProps.value" class="flex align-items-center">
                            <span>{{ slotProps.value.rut }} {{ slotProps.value.holderName }}</span>
                            <span>, {{ slotProps.value.bankName }}</span>
                            <span>, {{ slotProps.value.accountType }} {{ slotProps.value.accountNumber }}</span>
                        </div>
                        <span v-else>
                            {{ slotProps.placeholder }}
                        </span>
                    </template>
                    <template #option="slotProps">
                        <div class="flex align-items-center">
                            <span>{{ slotProps.option.rut }} {{ slotProps.option.holderName }}</span>
                            <span>, {{ slotProps.option.bankName }}</span>
                            <span>, {{ slotProps.option.accountType }} {{ slotProps.option.accountNumber }}</span>
                        </div>
                    </template>
                </Dropdown>

                <Button label="Agregar Cuenta" severity="secondary" text @click="addAccountBtnClick" />

            </div>
        </div>

        <div class="formgrid grid" v-if="editMode">

            <Form :initial-values="initialValues" :validation-schema="schema" @submit="acceptEditAccountClick"  >

                <Field name="rut" v-slot="{ field, errorMessage, handleChange }">
                    <div class="field col-12 md:col-6">
                        <label for="rut">RUT</label>
                        <InputText :model-value="field.value" @update:modelValue="handleChange"
                            placeholder="Ej. 12345678-5"
                            class="text-base text-color surface-overlay p-2 border-1 border-solid surface-border border-round appearance-none outline-none focus:border-primary w-full"
                            :class="{ 'p-invalid': errorMessage }" />
                        <small v-show="errorMessage" id="rut-help" class="p-error">{{ errorMessage }}</small>
                    </div>
                </Field>

                <Field name="holderName" v-slot="{ field, errorMessage, handleChange }">
                    <div class="field col-12 md:col-6">
                        <label for="holderName">Nombre Titular</label>
                        <InputText :model-value="field.value" @update:modelValue="handleChange"
                            class="text-base text-color surface-overlay p-2 border-1 border-solid surface-border border-round appearance-none outline-none focus:border-primary w-full"
                            :class="{ 'p-invalid': errorMessage }" />
                        <small v-show="errorMessage" id="holder-name-help" class="p-error">{{ errorMessage }}</small>
                    </div>
                </Field>

                <Field name="bankName" v-slot="{ field, errorMessage, handleChange }">
                    <div class="field col-12 md:col-6">
                        <label for="chiliBankName">Banco</label>
                        <Dropdown :model-value="field.value" @update:modelValue="handleChange" :options="chiliBanks"
                            showClear optionLabel="name" optionValue="code"
                            class="text-base text-color surface-overlay p-2 border-1 border-solid surface-border border-round appearance-none outline-none focus:border-primary w-full"
                            :class="{ 'p-invalid': errorMessage }" />
                        <small v-show="errorMessage" id="holder-name-help" class="p-error">{{ errorMessage }}</small>
                    </div>
                </Field>

                <Field name="accountType" v-slot="{ field, errorMessage, handleChange }">
                    <div class="field col-12 md:col-6">
                        <label for="chileBankTypeAccount">Tipo de Cuenta</label>
                        <Dropdown :model-value="field.value" @update:modelValue="handleChange"
                            :options="chiliBanksTypeAccounts" showClear optionLabel="name" optionValue="code"
                            class="text-base text-color surface-overlay p-2 border-1 border-solid surface-border border-round appearance-none outline-none focus:border-primary w-full"
                            :class="{ 'p-invalid': errorMessage }" />
                        <small v-show="errorMessage" id="holder-name-help" class="p-error">{{ errorMessage }}</small>
                    </div>
                </Field>

                <Field name="accountNumber" v-slot="{ field, errorMessage, handleChange }">
                    <div class="field col-12 md:col-6">
                        <label for="numeroCuenta">Numero de Cuenta</label>
                        <InputText :model-value="field.value" @update:modelValue="handleChange"
                            class="text-base text-color surface-overlay p-2 border-1 border-solid surface-border border-round appearance-none outline-none focus:border-primary w-full"
                            :class="{ 'p-invalid': errorMessage }" />
                        <small v-show="errorMessage" id="holder-name-help" class="p-error">{{ errorMessage }}</small>
                    </div>
                </Field>

                <div class="field col-12">
                    <Button label="Aceptar" type="submit" severity="primary" />
                    <Button label="Cancelar" severity="secondary" @click="cancelEditAccountClick" />
                </div>

            </Form>

        </div>

    </div>
</template>
<style lang="scss">
.PaymentPackDefinitionChiliBank {}
</style>
