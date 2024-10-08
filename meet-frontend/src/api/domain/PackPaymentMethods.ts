export interface PackPaymentMethods {
  chiliBankReceiptAccountId?: string
  chiliBankReceiptCLPPrice?: number
  chiliBankReceiptMethodEnabled: boolean
  paypalOnlineMethodEnabled: boolean
  paypalOnlineRecipientEmail?: string
  paypalOnlineUSDPrice?: number
  paypalReceiptMethodEnabled: boolean
  paypalReceiptRecipientEmail?: string
  paypalReceiptUSDPrice?: string
}
