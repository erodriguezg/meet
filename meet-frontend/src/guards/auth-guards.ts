import { AuthService, type Profile, type Permission } from '../services/AuthService'
import RoutesNames from '../utils/routes-names'

const hasPermission = (to: any, permission: Permission): any => {
  if (AuthService.hasPermission(permission)) {
    return true
  } else {
    return { name: RoutesNames.HOME_PAGE }
  }
}

const hasAnyPermission = (to: any, permissionList: Permission[]): any => {
  if (AuthService.hasAnyPermission(permissionList)) {
    return true
  } else {
    return { name: RoutesNames.HOME_PAGE }
  }
}

const hasProfile = (to: any, profile: Profile): any => {
  if (AuthService.hasProfile(profile)) {
    return to
  } else {
    return { name: RoutesNames.HOME_PAGE }
  }
}

const hasAnyProfile = (to: any, profileList: Profile[]): any => {
  if (AuthService.hasAnyProfile(profileList)) {
    return to
  } else {
    return { name: RoutesNames.HOME_PAGE }
  }
}

export const AuthGuard = {
  hasPermission,
  hasAnyPermission,
  hasProfile,
  hasAnyProfile
}
