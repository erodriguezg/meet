import rut from 'rut.js'

const toNumber = (rutIn: string): number => {
  const rutText = rutIn.split('-')[0].replace('.', '').trim()
  return parseInt(rutText)
}

const toString = (rutNumber: number): string => {
  const digit = rut.getCheckDigit(`${rutNumber}`)
  return `${rutNumber}-${digit}`
}

const toStringFormatted = (rutNumber: number): string => {
  return rut.format(toString(rutNumber))
}

const validate = (rutIn: string): boolean => {
  return rut.validate(rutIn)
}

const format = (rutIn: string): string => {
  return rut.format(rutIn)
}

export const RutUtils = {
  toNumber,
  toString,
  toStringFormatted,
  validate,
  format
}
