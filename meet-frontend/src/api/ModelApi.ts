import axios from 'axios'
import { type FilterSearchModel, type SearchModelResponse, type Model, type ModelRegisterData } from './domain/Model'
import { type ResourceUploadUrlDto } from './domain/UploadFile'

const searchModels = async (filters: FilterSearchModel, first: number, last: number): Promise<SearchModelResponse> => {
  const response = await axios.post('/api/v1/model/search', filters, {
    params: {
      first,
      last
    }
  })
  return response.data.payload
}

const registerModel = async (registerData: ModelRegisterData): Promise<void> => {
  await axios.post('/api/v1/model/register', registerData)
}

const findModelByNickName = async (nickName: string): Promise<Model | undefined> => {
  const response = await axios.get(`/api/v1/model/${nickName}`)
  return response.data.payload
}

const prepareUploadProfileImg = async (nickName: string): Promise<ResourceUploadUrlDto[]> => {
  const response = await axios.post(`/api/v1/model/${nickName}/prepare-profile-img`)
  return response.data.payload
}

export const ModelApi = {
  searchModels,
  findModelByNickName,
  registerModel,
  prepareUploadProfileImg
}
