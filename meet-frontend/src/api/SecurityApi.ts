import axios from 'axios'
import { type Identity } from '../services/AuthService'

const getLoginUrl = async (): Promise<string> => {
  const { data } = await axios.get('/api/v1/security/login-url')
  return data.loginUrl
}

const getToken = async (code: string, state: string): Promise<any> => {
  const params = {
    code,
    state
  }
  const { data } = await axios.post('/api/v1/security/token', null, {
    params
  })
  return data
}

const getIdentity = async (): Promise<Identity> => {
  const { data } = await axios.get('/api/v1/security/identity')
  return data.payload
}

export const SecurityApi = {
  getLoginUrl,
  getToken,
  getIdentity
}
