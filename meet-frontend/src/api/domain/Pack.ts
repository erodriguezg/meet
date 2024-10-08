export enum TypeCodeEnum {
  IMG_JPG = 'img-jpg',
  IMG_PNG = 'img-png',
  VIDEO_MP4 = 'video-mp4',
  VIDEO_OGG = 'video-ogg'
}

export interface PrepareUploadPackItemDto {
  modelNickName: string
  packNumber: number
  typeCode: TypeCodeEnum
  isPublic: boolean
}

export interface PackKeyDto {
  modelNickName: string
  packNumber: number
}

export interface PackDto {
  packNumber: number
  title?: string
  coverImageFileHash?: string
  isLocked: boolean
}

export interface PackItemDto {
  typeCode: TypeCodeEnum
  itemNumber: number
  resourceFileHash?: string
  thumbnailFileHash: string
  isLocked: boolean
}

export interface PackInfoDto {
  title?: string
  description?: string
}
