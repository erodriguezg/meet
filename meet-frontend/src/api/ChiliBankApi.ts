import axios from 'axios'
import { type ChiliBankAccountDTO } from './domain/ChiliBank'

const getBanks = async (): Promise<string[]> => {
  const response = await axios.get('/api/v1/chili-bank/banks')
  return response.data.payload
}

const getAccountTypes = async (): Promise<string[]> => {
  const response = await axios.get('/api/v1/chili-bank/account-types')
  return response.data.payload
}

const getModelAccounts = async (modelNickname: string): Promise<ChiliBankAccountDTO[]> => {
  const response = await axios.get(`/api/v1/chili-bank/${modelNickname}/accounts`)
  return response.data.payload
}

const saveModelAccount = async (modelNickname: string, account: ChiliBankAccountDTO): Promise<ChiliBankAccountDTO> => {
  const response = await axios.post(`/api/v1/chili-bank/${modelNickname}/accounts`, account)
  return response.data.payload
}

const deleteModelAccount = async (modelNickname: string, accountId: string): Promise<void> => {
  await axios.delete(`/api/v1/chili-bank/${modelNickname}/${accountId}`)
}

export const ChiliBankApi = {
  getBanks,
  getAccountTypes,
  getModelAccounts,
  saveModelAccount,
  deleteModelAccount
}
