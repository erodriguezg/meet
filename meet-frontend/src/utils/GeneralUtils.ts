
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

export const GeneralUtils = {
  clearBlankOrNull,
  clearNull,
  stringsArrayToSelectItemArray
}

export interface SelectItem {
  name: string
  code?: string
}
