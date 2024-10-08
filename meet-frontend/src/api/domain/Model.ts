
export interface ModelRegisterData {
  nickName: string
  personId: string
}

export interface Model {
  id?: string
  personId: string
  nickName: string
  profileImageFileHash?: string
  profileImageThumbnailFileHash?: string
  aboutMe?: string
  countryCode?: string
  city?: string
  zodiacSignCode?: string
}

export interface FilterSearchModel {
  nickName?: string
  countryCode?: string
  cityName?: string
  zodiacSignCode?: string
}

export interface SearchModelResponse {
  totalCount: number
  models?: Model[]
}
