import axios from 'axios'
import { type PackBuyDetailDto, type CapturePaymentRequest, type CreateOrderRequest, type CreateOrderResponse, type BuyPackDetailsRequest } from './domain/BuyPack'

const info = async (): Promise<any> => {
  const response = await axios.get('/api/v1/buy-pack/info')
  return response.data.payload
}

const createOrder = async (request: CreateOrderRequest): Promise<CreateOrderResponse> => {
  const response = await axios.post('/api/v1/buy-pack/create-order', request)
  return response.data.payload
}

const capturePayment = async (request: CapturePaymentRequest): Promise<void> => {
  await axios.post('/api/v1/buy-pack/capture-payment', request)
}

const details = async (request: BuyPackDetailsRequest): Promise<PackBuyDetailDto> => {
  const response = await axios.post('/api/v1/buy-pack/details', request)
  return response.data.payload
}

export const BuyPackApi = {
  info,
  createOrder,
  capturePayment,
  details
}
