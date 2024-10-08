export interface CreateOrderRequest {
  personId: string
  modelNickName: string
  packNumber: number
}

export interface BuyPackDetailsRequest {
  modelNickName: string
  packNumber: number
}

export interface PackBuyDetailDto {
  modelNickName: string
  packTitle?: string
  packDollarValue: number
}

export interface CreateOrderResponse {
  orderId: string
}

export interface CapturePaymentRequest {
  orderId: string
}
