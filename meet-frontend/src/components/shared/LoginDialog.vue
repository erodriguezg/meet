<script setup lang="ts">

import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'

import { ref, watch } from 'vue'
import { useLoginDialogStore } from '../../stores/LoginDialogStore'
import { SecurityApi } from '../../api/SecurityApi'

const loginDialogStore = useLoginDialogStore()

const visible = ref<boolean>(false)
const loadingContinueWithGoogle = ref<boolean>(false)
const loadingContinueWithMicrosoft = ref<boolean>(false)

const continueWithGoogleAction = async () => {
  loadingContinueWithGoogle.value = true
  try {
    const url = await SecurityApi.getLoginUrl()
    window.location.href = url
  } catch (err) {
    console.error(err)
    loadingContinueWithGoogle.value = false
  }
}

const continueWithMicrosoftAction = () => {
  const reload = window.location.href
  window.location.href = reload
}

const closeDialog = () => {
  loginDialogStore.setLoginDialogData({
    dialogVisible: false
  })
}

watch(() => loginDialogStore.loginDialogData, (newValue, oldValue) => {
  if (newValue?.dialogVisible) {
    visible.value = true
  }
})

</script>

<template>

        <Dialog header="Iniciar Sesión" v-model:visible="visible" modal :style="{ width: '30vw' }" :draggable="false"
             :breakpoints="{ '960px': '75vw' }" class="LoginDialog" @hide="closeDialog">

            <div class="field">
                <label htmlFor="email">Correo Electrónico</label>
                <InputText id="email" class="form-input-text" />
            </div>
            <div class="field">
                <label htmlFor="password">Constraseña</label>
                <Password id="password" :toggleMask="true" :feedback="false"></Password>
            </div>
            <div class="field">
                <Button label="Ingresar"></Button>
            </div>

            <div class="field sign-up">
                <span>Si no tienes cuenta, registrate <a href="#">aquí</a></span>
            </div>

            <div class="field or">
                <span>O</span>
            </div>

            <div class="field buttons">

                <Button @click="continueWithGoogleAction" :loading="loadingContinueWithGoogle" outlined>
                    <img alt="logo" src="../../assets/images/google-btn.svg" style="width: 1.5rem" />
                    <span class="ml-2 font-bold">Iniciar con Google</span>
                </Button>

                <Button @click="continueWithMicrosoftAction" :loading="loadingContinueWithMicrosoft" outlined>
                    <img alt="logo" src="../../assets/images/microsoft-btn.svg" style="width: 1.5rem" />
                    <span class="ml-2 font-bold">Iniciar con Microsoft</span>
                </Button>

            </div>

        </Dialog>

</template>

<style lang="scss">
.LoginDialog {
    .p-password {
        width: 100%;

        .p-inputtext {
            width: 100%;
        }
    }

    .p-inputtext {
        width: 100%;
    }

    .p-button {
        width: 100%;
    }

    .field {
        &.sign-up {
            text-align: center;
            font-size: 1rem;

            a {
                margin-left: 10px;
            }
        }

        &.or {
            width: 100%;
            text-align: center;
            position: relative;
            border-top: 1px solid #ccc;
            margin-top: 30px;
            margin-bottom: 20px;

            span {
                background-color: #fff;
                padding: 0 10px;
                position: relative;
                top: -10px;
            }
        }

        &.buttons {
            .p-button {
                margin-bottom: 13px;
                font-size: 0.95rem;
            }
        }
    }
}
</style>
