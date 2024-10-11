
const backendUrl: string = import.meta.env.VITE_APP_BACKEND_URL as string

const clearBlankOrNull = (input: any): any | undefined => {
  if (input === undefined || input === null) {
    return undefined
  }
  if (typeof input === 'string' || input instanceof String) {
    const text: string = input as string
    if (text.trim().length === 0) {
      return undefined
    }
  }
  return input
}

const clearNull = (input: any): any | undefined => {
  if (input === undefined || input === null) {
    return undefined
  }
  return input
}

const stringsArrayToSelectItemArray = (inputArray: string[]): SelectItem[] => {
  return inputArray.map(inString => {
    return {
      name: inString,
      code: inString
    }
  })
}

const getWebSocketBaseUrl = (): string => {
  if (backendUrl === '/') {
    return `${window.location.protocol}//${window.location.hostname}`
  } else {
    return backendUrl
  }
}

export const GeneralUtils = {
  clearBlankOrNull,
  clearNull,
  stringsArrayToSelectItemArray,
  getWebSocketBaseUrl
}

export interface SelectItem {
  name: string
  code?: string
}
