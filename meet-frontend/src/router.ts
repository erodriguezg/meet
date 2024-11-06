import { type RouteRecordRaw, createRouter, createWebHistory } from 'vue-router'
import RoutesNames from './utils/routes-names'
import { Permission } from './services/AuthService'
import { AuthGuard } from './guards/auth-guards'

const MainTemplate = (): any => import('./components/template/MainTemplate.vue')
const EmptyTemplate = (): any => import('./components/template/EmptyTemplate.vue')
const CallbackLoginPage = (): any => import('./components/pages/CallbackLoginPage.vue')
const LoginPage = (): any => import('./components/pages/LoginPage.vue')
const MeetPage = (): any => import('./components/pages/MeetPage.vue')
const UsersPage = (): any => import('./components/pages/UsersPage.vue')
const RoomsPage = (): any => import('./components/pages/RoomsPage.vue')

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: EmptyTemplate,
    children: [
      {
        name: RoutesNames.LOGIN_PAGE,
        path: '',
        component: LoginPage
      },
      {
        name: RoutesNames.MEET_PAGE,
        path: 'm/:roomId',
        component: MeetPage
      }
    ]
  },
  {
    path: '/admin',
    component: MainTemplate,
    children: [
      {
        name: RoutesNames.MANAGE_USERS,
        path: 'users',
        component: UsersPage,
        beforeEnter: (to: any, from: any) => AuthGuard.hasPermission(to, Permission.MANAGE_SYSTEM)
      }
    ]
  },
  {
    path: '/rooms',
    component: MainTemplate,
    children: [
      {
        name: RoutesNames.MANAGE_ROOMS,
        path: 'all',
        component: RoomsPage,
        meta: {
          showAllRooms: true
        },
        beforeEnter: (to: any, from: any) => AuthGuard.hasPermission(to, Permission.MANAGE_SYSTEM)
      },
      {
        name: RoutesNames.MANAGE_OWN_ROOMS,
        path: 'own',
        component: RoomsPage,
        meta: {
          showAllRooms: false
        },
        beforeEnter: (to: any, from: any) => AuthGuard.hasAnyPermission(to, [Permission.CREATE_ROOM, Permission.MANAGE_SYSTEM])
      }
    ]
  },
  { path: '/callback-login', component: CallbackLoginPage }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
