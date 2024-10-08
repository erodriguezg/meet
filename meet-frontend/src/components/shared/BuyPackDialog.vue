<script setup lang="ts">

import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import { ref } from 'vue'
import { CreateOrderData, OnApproveData, loadScript } from '@paypal/paypal-js'
import { BuyPackApi } from '../../api/BuyPackApi'
import BlockUI from 'primevue/blockui'
import { PackBuyDetailDto } from '../../api/domain/BuyPack'

const loading = ref<boolean>(false)
const visible = ref<boolean>(false)
const visibleSuccess = ref<boolean>(false)
const buyPackDto = ref<PackBuyDetailDto>()
const paypalBtn = ref<HTMLElement>()

const open = async (modelNickName: string, packNumber: number, personId: string) => {
  const info = await BuyPackApi.info()

  buyPackDto.value = await BuyPackApi.details({
    modelNickName,
    packNumber
  })

  visible.value = true

  try {
    const paypal = await loadScript({
      clientId: info.clientId
    })
    if (paypal !== undefined && paypal !== null && paypal.Buttons !== undefined) {
      await paypal.Buttons({
        createOrder: async (data: CreateOrderData) => {
          loading.value = true
          const response = await BuyPackApi.createOrder({
            modelNickName,
            packNumber,
            personId
          })
          loading.value = false
          return response.orderId
        },
        onApprove: async (data: OnApproveData) => {
          loading.value = true
          await BuyPackApi.capturePayment({
            orderId: data.orderID
          })
          loading.value = false
          visible.value = false
          visibleSuccess.value = true
        }
      }).render(paypalBtn.value ?? '')
    }
  } catch (err) {
    console.error(err)
  }
}

defineExpose({
  open
})

const acceptSuccessDialog = () => {
  visibleSuccess.value = false
  const reload = window.location.href
  window.location.href = reload
}

</script>

<template>
    <div class="BuyPackDialog">

        <Dialog v-model:visible="visible" modal header="Adquirir Pack" :style="{ width: '50vw' }"
        :breakpoints="{ '960px': '75vw', '641px': '100vw' }" >

          <BlockUI :blocked="loading">
            <div class="priceSection">
                <p>Valor del Pack: <em>${{buyPackDto?.packDollarValue}} US</em></p>
            </div>

            <div class="paypalButtons" ref="paypalBtn" ></div>
          </BlockUI>

        </Dialog>

        <Dialog v-model:visible="visibleSuccess" modal header="Transación Existosa" :style="{ width: '50vw' }"
        :breakpoints="{ '960px': '75vw', '641px': '100vw' }" >

          <p>Se ha adquirido el Pack</p>

          <template #footer>
            <Button label="Aceptar" icon="pi pi-check" @click="acceptSuccessDialog" autofocus />
          </template>

        </Dialog>

    </div>
</template>

<style lang="scss">
.BuyPackDialog {

    .priceSection {
        color: red;
    }

}
</style>
