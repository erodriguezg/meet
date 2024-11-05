import { SecurityApi } from '../api/SecurityApi'
import { AxiosUtil } from '../utils/AxiosUtil'
import { CookieStorageUtil } from '../utils/CookieStorageUtil'
import { SessionStorageUtil } from '../utils/SessionStorageUtil'

const jwtCookieKey = 'jwt'
const identitySessionKey = 'identity'

export interface Identity {
  personId: string
  email: string
  firstName: string
  lastName: string
  profileCode: number
  profileName: string
  permissionsCodes: number[]
  modelId?: string
  modelNickName?: string
}

export enum Profile {
  ADMINISTRATOR = 1,
  USER = 2
}

export enum Permission {
  MANAGE_SYSTEM = 1,
  EDIT_OWN_PROFILE = 2,
  CREATE_ROOM = 3
}

const processLoginCallback = async (code: string, state: string): Promise<void> => {
  const tokenResponse = await SecurityApi.getToken(code, state)
  CookieStorageUtil.setString(jwtCookieKey, tokenResponse.jwt)
  await refreshIdentity(tokenResponse.jwt)
}

const initialize = async (): Promise<void> => {
  const jwt = CookieStorageUtil.getString(jwtCookieKey)
  if (jwt !== null) {
    await refreshIdentity(jwt)
  }
}

const getIdentity = (): Identity | null => {
  return SessionStorageUtil.getItem(identitySessionKey)
}

const isAuthenticated = (): boolean => {
  return getIdentity() !== null
}

const hasProfile = (profile: Profile): boolean => {
  return getIdentity()?.profileCode === profile
}

const hasAnyProfile = (profileList: Profile[]): boolean => {
  return profileList.some(p => hasProfile(p)) ?? false
}

const hasPermission = (permission: Permission): boolean => {
  return getIdentity()?.permissionsCodes?.some(pc => Number(pc) === Number(permission)) ?? false
}

const hasAnyPermission = (permissionList: Permission[]): boolean => {
  return permissionList.some(p => hasPermission(p)) ?? false
}

const logout = (): void => {
  CookieStorageUtil.remove(jwtCookieKey)
  SessionStorageUtil.remove(identitySessionKey)
}

const refreshIdentity = async (jwtToken: string): Promise<void> => {
  try {
    AxiosUtil.setAuthorization(jwtToken)
    const identity = await SecurityApi.getIdentity()
    SessionStorageUtil.setItem(identitySessionKey, identity)
  } catch (e) {
    cleanStorage()
  }
}

const cleanStorage = (): void => {
  CookieStorageUtil.remove(jwtCookieKey)
  SessionStorageUtil.remove(identitySessionKey)
}

export const AuthService = {
  processLoginCallback,
  initialize,
  getIdentity,
  logout,
  isAuthenticated,
  hasProfile,
  hasAnyProfile,
  hasPermission,
  hasAnyPermission
}
