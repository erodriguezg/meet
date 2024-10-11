<script setup lang="ts">
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import Checkbox from 'primevue/checkbox'
import Card from 'primevue/card'
import InputText from 'primevue/inputtext'
import Fieldset from 'primevue/fieldset'
// import PaymentPackDefinitionChiliBank from './PaymentPackDefinitionChiliBank.vue'
import { onMounted, ref } from 'vue'
import { useModelDetailStore } from '../../stores/ModelDetailStore'
import { PackPaymentMethodsApi } from '../../api/PackPaymentMethodsApi'
import { PackPaymentMethods } from '../../api/domain/PackPaymentMethods'

const modelDetailStore = useModelDetailStore()

const dialogVisible = ref<boolean>(false)
const modelNickname = ref<string>()
const packNumber = ref<number>()
const packPaymentMethods = ref<PackPaymentMethods>()

onMounted(async () => {
  const modelDetail = modelDetailStore.modelDetail
  modelNickname.value = modelDetail?.modelNickName
  packPaymentMethods.value = {
    chiliBankReceiptMethodEnabled: false,
    paypalOnlineMethodEnabled: false,
    paypalReceiptMethodEnabled: false
  }
})

const open = async (packNumberInput: number) => {
  dialogVisible.value = true
  packNumber.value = packNumberInput
  const dto = await PackPaymentMethodsApi.getPackPaymentMethods(modelNickname.value!!, packNumber.value)

  if (dto !== null) {
    packPaymentMethods.value = dto
  } else {
    packPaymentMethods.value = {
      chiliBankReceiptMethodEnabled: false,
      paypalOnlineMethodEnabled: false,
      paypalReceiptMethodEnabled: false
    }
  }
}

const save = async () => {
  await PackPaymentMethodsApi.savePackPaymentMethods(modelNickname.value!!, packNumber.value!!, packPaymentMethods.value!!)
  dialogVisible.value = false
}

const cancel = async () => {
  dialogVisible.value = false
}

defineExpose({
  open
})

</script>

<template>
    <div class="PaymentPackDefinitionDialog">
        <Dialog class="PaymentPackDefinitionDialog" v-model:visible="dialogVisible" modal header="Definir Modo Pago" :closable="false"
            :breakpoints="{ '1199px': '75vw', '575px': '90vw' }">

            <p>Selecciona los medios de pago</p>

            <Card class="PaymentPackDefinitionDialog">
                <template #title>
                    <Checkbox v-model="packPaymentMethods!.chiliBankReceiptMethodEnabled" :binary="true" />
                    <span>Comprobante Transferencia Bancaria Chile (Pesos Chilenos)</span>
                </template>
                <template #content v-if="packPaymentMethods?.chiliBankReceiptMethodEnabled">
                    <Fieldset legend="Cuenta Destino">
                        <Suspense>
                            <!-- PaymentPackDefinitionChiliBank :modelNickname="modelNickname" v-model:chiliBankAccountId="packPaymentMethods.chiliBankReceiptAccountId" -->
                        </Suspense>
                    </Fieldset>
                    <Fieldset legend="Precio Pack">
                        <div class="formgrid grid">
                            <div class="field col-12 md:col-6">
                                <label for="bankVoucherCLPPrice" class="required">Precio Pesos Chilenos</label>
                                <InputText id="bankVoucherCLPPrice"
                                    class="text-base text-color surface-overlay p-2 border-1 border-solid surface-border border-round appearance-none outline-none focus:border-primary w-full" />
                            </div>
                        </div>
                    </Fieldset>
                </template>

            </Card>

            <Card class="PaymentPackDefinitionDialog">
                <template #title>
                    <Checkbox v-model="packPaymentMethods!.paypalReceiptMethodEnabled" :binary="true" />
                    <span>Comprobante PayPal (Dolares)</span>
                </template>
                <template #content v-if="packPaymentMethods?.paypalReceiptMethodEnabled">
                    <p>Datos Destinatario</p>
                    <div class="formgrid grid">
                        <div class="field col-12 md:col-6">
                            <label for="paypalVoucherEmail">Correo electrónico destinatario transferencia</label>
                            <InputText id="paypalVoucherEmail"
                                class="text-base text-color surface-overlay p-2 border-1 border-solid surface-border border-round appearance-none outline-none focus:border-primary w-full" />
                        </div>
                    </div>
                    <p>Monto</p>
                    <div class="formgrid grid">
                        <div class="field col-12 md:col-6">
                            <label for="paypalVoucherDollarPrice">Precio Dolares Americanos</label>
                            <InputText id="paypalVoucherDollarPrice"
                                class="text-base text-color surface-overlay p-2 border-1 border-solid surface-border border-round appearance-none outline-none focus:border-primary w-full" />
                        </div>
                    </div>
                </template>
            </Card>

            <Card class="PaymentPackDefinitionDialog">
                <template #title>
                    <Checkbox v-model="packPaymentMethods!.paypalOnlineMethodEnabled" :binary="true" />
                    <span>PayPal Automático (Dolares)</span>
                </template>
                <template #content v-if="packPaymentMethods?.paypalOnlineMethodEnabled">
                    <p class="m-0">
                        Lorem ipsum dolor sit amet, consectetur adipisicing elit. Inventore sed consequuntur error
                        repudiandae numquam deserunt quisquam repellat libero asperiores earum nam nobis, culpa
                        ratione quam perferendis esse, cupiditate neque
                        quas!
                    </p>
                </template>
            </Card>

            <template #footer>
                <Button label="Cancelar" text severity="secondary" @click="cancel" />
                <Button label="Guardar" outlined severity="primary" @click="save" />
            </template>

        </Dialog>
    </div>

</template>

<style lang="scss">
.PaymentPackDefinitionDialog {

    &.p-dialog {
        max-width: 80%;
    }

    &.p-card {
        margin-bottom: 10px !important;

        .p-card-title {
            display: flex;

            .p-checkbox {
                padding-top: 8px;
                margin-right: 8px;
            }
        }
    }

}
</style>
