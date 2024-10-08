import axios from 'axios'
import { type PackInfoDto, type PackDto, type PackItemDto, type PackKeyDto, type PrepareUploadPackItemDto } from './domain/Pack'
import { type ResourceUploadUrlDto } from './domain/UploadFile'

const prepareUploadPackItem = async (data: PrepareUploadPackItemDto): Promise<ResourceUploadUrlDto[]> => {
  const response = await axios.post('/api/v1/pack/prepare-upload-item', data)
  return response.data.payload
}

const publishPack = async (data: PackKeyDto): Promise<void> => {
  await axios.post('/api/v1/pack/publish', data)
}

const readyToPublishPack = async (data: PackKeyDto): Promise<void> => {
  await axios.post('/api/v1/pack/ready-to-publish', data)
}

const getPacksFromModel = async (modelNickName: string): Promise<PackDto[]> => {
  const response = await axios.get(`/api/v1/pack/${modelNickName}`)
  return response.data.payload
}

const getPackInfo = async (modelNickName: string, packNumber: number): Promise<PackInfoDto> => {
  const response = await axios.get(`/api/v1/pack/${modelNickName}/${packNumber}/info`)
  return response.data.payload
}

const getItemsFromPack = async (modelNickName: string, packNumber: number): Promise<PackItemDto[]> => {
  const response = await axios.get(`/api/v1/pack/${modelNickName}/${packNumber}/items`)
  return response.data.payload
}

const createNewPack = async (modelNickName: string): Promise<PackDto> => {
  const response = await axios.put(`/api/v1/pack/${modelNickName}/new`)
  return response.data.payload
}

const deletePackItem = async (modelNickName: string, packNumber: number, packItem: number): Promise<void> => {
  await axios.delete(`/api/v1/pack/${modelNickName}/${packNumber}`)
}

const deletePack = async (modelNickName: string, packNumber: number): Promise<void> => {
  await axios.delete(`/api/v1/pack/${modelNickName}`)
}

const editPackTitle = async (modelNickName: string, packNumber: number, title: string): Promise<void> => {
  await axios.post(`/api/v1/pack/${modelNickName}/${packNumber}/title`, { title })
}

const editPackDescription = async (modelNickName: string, packNumber: number, description: string): Promise<void> => {
  await axios.post(`/api/v1/pack/${modelNickName}/${packNumber}/description`, { description })
}

export const PackApi = {
  prepareUploadPackItem,
  publishPack,
  readyToPublishPack,
  getPacksFromModel,
  getPackInfo,
  getItemsFromPack,
  createNewPack,
  deletePackItem,
  deletePack,
  editPackTitle,
  editPackDescription
}
