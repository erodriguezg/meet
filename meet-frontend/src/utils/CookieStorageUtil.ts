import Cookies from 'js-cookie'

const setString = (key: string, value: string): void => {
  Cookies.set(key, value)
}

const getString = (key: string): string | null => {
  const value = Cookies.get(key)
  if (value === undefined) {
    return null
  }
  return value
}

const setItem = <T>(key: string, value: T): void => {
  Cookies.set(key, JSON.stringify({ ...value }))
}

const remove = (key: string): void => {
  Cookies.remove(key)
}

const getItem = <T>(key: string): T | null => {
  const textJson = getString(key)
  if (textJson === null) {
    return null
  }
  return JSON.parse(textJson)
}

export const CookieStorageUtil = {
  setString,
  getString,
  setItem,
  remove,
  getItem
}
