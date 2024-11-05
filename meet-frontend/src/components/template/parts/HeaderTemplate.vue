<script setup lang="ts">
import Button from 'primevue/button'
import Menubar from 'primevue/menubar'
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import RoutesNames from '../../../utils/routes-names'
import { AuthService, Permission } from '../../../services/AuthService'
import { useLoginDialogStore } from '../../../stores/LoginDialogStore'

const loginDialogStore = useLoginDialogStore()

const router = useRouter()

const loadingLogin = ref<boolean>(false)
const loadingLogout = ref<boolean>(false)

const isLoggedIn = ref<boolean>(false)

isLoggedIn.value = AuthService.isAuthenticated()

const items = ref<any>([
  {
    label: 'My Profile',
    visible: AuthService.hasAnyPermission([Permission.EDIT_OWN_PROFILE, Permission.MANAGE_SYSTEM]),
    icon: 'fa-regular fa-address-card',
    command: () => {
      router.push({
        name: RoutesNames.EDIT_OWN_PROFILE
      })
    }
  },
  {
    label: 'Administration',
    icon: 'pi pi-fw pi-user',
    visible: AuthService.hasPermission(Permission.MANAGE_SYSTEM),
    items: [
      {
        label: 'Users',
        icon: 'pi pi-fw pi-user-plus',
        command: () => {
          router.push({
            name: RoutesNames.MANAGE_USERS
          })
        }
      }
    ]
  },
  {
    label: 'Rooms',
    icon: 'fa-regular fa-comments-o',
    visible: AuthService.hasAnyPermission([Permission.CREATE_ROOM, Permission.MANAGE_SYSTEM]),
    items: [
      {
        label: 'Create Room',
        icon: 'fa-regular fa-plus',
        command: () => {
          router.push({
            name: RoutesNames.CREATE_ROOM
          })
        }
      },
      {
        label: 'Manage Own Rooms',
        icon: 'fa-regular fa-smile-o',
        command: () => {
          router.push({
            name: RoutesNames.MANAGE_OWN_ROOMS
          })
        }
      },
      {
        label: 'Manage All Rooms',
        icon: 'fa-regular fa-globe',
        command: () => {
          router.push({
            name: RoutesNames.MANAGE_ROOMS
          })
        }
      }
    ]
  },
  {
    label: 'Close Session',
    icon: 'pi pi-fw pi-power-off',
    command: () => { logoutAction() }
  }
])

const logoAction = () => {
  router.push({
    name: RoutesNames.LOGIN_PAGE
  })
}

const loginAction = async () => {
  loginDialogStore.setLoginDialogData({
    dialogVisible: true
  })
}

const logoutAction = async () => {
  loadingLogout.value = true
  try {
    AuthService.logout()
    window.location.href = '/'
  } catch (err) {
    console.error(err)
  } finally {
    loadingLogout.value = false
  }
}

</script>

<template>
  <div class="headerTemplate">

    <div class="desktop">
      <div class="content">
        <img src="../../../assets/logo3.png" alt="logo2" class="logo" @click="logoAction" />
        <Menubar v-if="isLoggedIn" :model="items" />
        <Button v-if="!isLoggedIn" label="Ingresar" @click="loginAction" :loading="loadingLogin"></Button>
      </div>
    </div>

    <div class="mobile">
      <div class="content">
        <Menubar v-if="isLoggedIn" :model="items" />
        <img src="../../../assets/logo3.png" alt="logo2" class="logo" @click="logoAction" />
        <Button v-if="!isLoggedIn" icon="pi pi-sign-in" rounded @click="loginAction" :loading="loadingLogin"></Button>
      </div>
    </div>

  </div>
</template>

<style lang="scss">
@use '../../../styles/vars' as vars;

.headerTemplate {

  .desktop {

    .content {
      min-height: 4rem;
      background-color: black;
      display: flex;
      justify-content: flex-end;
      padding-right: 3rem;
      align-items: center;

      .logo {
        width: 18rem;
        position: absolute;
        left: 1rem;
        cursor: pointer;
        top: 5px;
      }

      .p-button {
        height: 2rem;
        background-color: vars.$colorYellow;
        color: vars.$colorBlack;
        border: none;

        &:hover {
          background-color: vars.$colorRed;
          color: vars.$colorBlack;
          border: none;
        }
      }

      .p-menubar {
        background: none;
        border: none;

        .p-menuitem-text {
          color: vars.$colorWhite !important;
        }

        .p-menuitem-icon {
          color: vars.$colorWhite !important;
        }

        .p-menuitem-content:hover {
          background: none !important;
        }

        .p-submenu-list {
          background: vars.$colorBlack !important;
        }

        .p-menuitem.p-highlight>.p-menuitem-content {
          background: none !important;

          .p-menuitem-text {
            color: vars.$colorRed !important;
          }

          .p-menuitem-icon {
            color: vars.$colorRed !important;
          }

        }

      }
    }

  }

  .mobile {

    .content {
      min-height: 4rem;
      background-color: black;
      display: flex;
      justify-content: flex-start;
      padding-right: 3rem;
      align-items: center;

      .logo {
        max-width: 16rem;
        min-width: 12rem;
        cursor: pointer;
        margin-left: 15px;
      }

      .p-button {
        background-color: rgb(221 121 53);
        color: vars.$colorBlack;
        border: none;
        margin-right: 15px;
        position: absolute;
        right: 0;

        &:hover {
          background-color: vars.$colorRed;
          color: vars.$colorBlack;
          border: none;
        }
      }

      ul.p-menubar-root-list {
        height: 160vw;
        width: 90vw;
      }

      div.p-menubar {
        margin-right: 10px;
        background-color: black;
        color: white;
        border-color: black;

        .p-menubar-button {
          background-color: black;
          color: white;
        }

        .p-submenu-list,
        .p-menubar-root-list {
          background-color: black;
        }

        .p-menuitem-icon,
        .p-menuitem-text {
          color: white;
        }

        .p-menuitem.p-highlight>.p-menuitem-content {
          background-color: black;

          .p-menuitem-icon,
          .p-menuitem-text {
            color: #f3943d;
          }
        }

        .p-menuitem:not(.p-highlight):not(.p-disabled).p-focus>.p-menuitem-content,
        .p-menubar-root-list>.p-menuitem:not(.p-highlight):not(.p-disabled)>.p-menuitem-content:hover {
          background-color: black;
        }

      }
    }

  }

}
</style>
