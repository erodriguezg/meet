export interface BusinessErrorDetail {
  code: string
  message: string
  details?: Map<string, string>
}

export default class BusinessError extends Error {
  constructor (private readonly _detail: BusinessErrorDetail) {
    super(_detail.message)
  }

  public get code (): string {
    return this._detail.code
  }

  public get details (): Map<string, string> | undefined {
    return this._detail.details
  }
}
