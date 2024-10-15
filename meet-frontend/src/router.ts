import { type RouteRecordRaw, createRouter, createWebHistory } from 'vue-router'
import RoutesNames from './utils/routes-names'
import { Permission } from './services/AuthService'
import { AuthGuard } from './guards/auth-guards'

const MainTemplate = (): any => import('./components/template/MainTemplate.vue')
const CleanTemplate = (): any => import('./components/template/CleanTemplate.vue')
const CallbackLoginPage = (): any => import('./components/pages/CallbackLoginPage.vue')
const HomePage = (): any => import('./components/pages/HomePage.vue')
const EditPublicationPage = (): any => import('./components/pages/EditPublicationPage.vue')
const CategoriesPage = (): any => import('./components/pages/CategoriesPage.vue')
const MeetPage = (): any => import('./components/pages/MeetPage.vue')

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: CleanTemplate,
    children: [
      {
        name: RoutesNames.HOME_PAGE,
        path: '',
        component: HomePage
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
        name: RoutesNames.ADMINISTRATION_CATEGORIES,
        path: 'categories',
        component: CategoriesPage,
        beforeEnter: (to: any, from: any) => AuthGuard.hasPermission(to, Permission.MANAGE)
      },
      {
        name: RoutesNames.NEW_PUBLICATION_PAGE,
        path: '/publication/new',
        component: EditPublicationPage,
        beforeEnter: (to: any, from: any) => AuthGuard.hasPermission(to, Permission.PUBLISH)
      }
    ]
  },
  {
    path: '/publication',
    component: MainTemplate,
    children: [
      {
        name: RoutesNames.NEW_PUBLICATION_PAGE,
        path: '/new',
        component: EditPublicationPage,
        beforeEnter: (to: any, from: any) => AuthGuard.hasPermission(to, Permission.PUBLISH)
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
