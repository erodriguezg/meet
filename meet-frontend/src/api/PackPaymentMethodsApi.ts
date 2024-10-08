import axios from 'axios'
import { type PackPaymentMethods } from './domain/PackPaymentMethods'

const getPackPaymentMethods = async (modelNickName: string, packNumber: number): Promise<PackPaymentMethods> => {
  const response = await axios.get(`/api/v1/pack-payment-methods/${modelNickName}/${packNumber}`)
  return response.data.payload
}

const savePackPaymentMethods = async (modelNickName: string, packNumber: number, data: PackPaymentMethods): Promise<PackPaymentMethods> => {
  const response = await axios.post(`/api/v1/pack-payment-methods/${modelNickName}/${packNumber}`, data)
  return response.data.payload
}

export const PackPaymentMethodsApi = {
  getPackPaymentMethods,
  savePackPaymentMethods
}
