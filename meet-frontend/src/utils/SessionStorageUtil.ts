const setString = (key: string, value: string): void => {
  sessionStorage.setItem(key, value)
}

const getString = (key: string): string | null => {
  return sessionStorage.getItem(key)
}

const setItem = <T>(key: string, value: T): void => {
  sessionStorage.setItem(key, JSON.stringify({ ...value }))
}

const remove = (key: string): void => {
  sessionStorage.removeItem(key)
}

const getItem = <T>(key: string): T | null => {
  const textJson = sessionStorage.getItem(key)
  if (textJson === null) {
    return null
  }
  return JSON.parse(textJson)
}

const clearAll = (key: string): void => {
  sessionStorage.clear()
}

export const SessionStorageUtil = {
  setString,
  getString,
  setItem,
  remove,
  getItem,
  clearAll
}
