import axios from 'axios'

const confirmFileUpload = async (hash: string): Promise<void> => {
  await axios.post(`/api/v1/file/confirm/${hash}`)
}

const getStorageType = async (): Promise<string> => {
  const response = await axios.get('/api/v1/file/storage-type')
  return response.data.payload.storageType
}

const getDownloadUrl = async (hash: string): Promise<string> => {
  const response = await axios.get(`/api/v1/file/get/${hash}`)
  return response.data.payload.url
}

const getRedirectUrl = (hash: string): string => {
  let baseUrl = axios.defaults.baseURL ?? ''
  if (baseUrl === '/') {
    baseUrl = ''
  }
  return `${baseUrl}/api/v1/file/redirect/${hash}`
}

export const FileApi = {
  confirmFileUpload,
  getStorageType,
  getDownloadUrl,
  getRedirectUrl
}
