import { RutUtils } from './RutUtils'

const testRut = (value: string | undefined | null): boolean => {
  if (value == null || value === undefined) {
    return true
  }
  return RutUtils.validate(value)
}

export const YupUtils = {
  testRut
}
