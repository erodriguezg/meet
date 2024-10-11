import { AuthService, type Profile, type Permission } from '../services/AuthService'
import RoutesNames from '../utils/routes-names'

const hasPermission = (to: any, permission: Permission): any => {
  if (AuthService.hasPermission(permission)) {
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

export const AuthGuard = {
  hasPermission,
  hasProfile
}
