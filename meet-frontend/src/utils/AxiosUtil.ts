import axios, { type AxiosResponse } from 'axios'
import BusinessError from '../errors/BusinessError'

const initialize = (): void => {
  axios.defaults.baseURL = import.meta.env.VITE_APP_BACKEND_URL
  axios.defaults.headers.post['Content-Type'] = 'application/json'
  axios.interceptors.response.use(businessExceptionInterceptor)
}

const setAuthorization = (token: string): void => {
  axios.defaults.headers.Authorization = `Bearer ${token}`
}

const businessExceptionInterceptor = (response: AxiosResponse): AxiosResponse => {
  if (response?.data?.status === 'BUSINESS_ERROR') {
    throw new BusinessError(response.data.respuesta)
  }
  return response
}

export const AxiosUtil = {
  initialize,
  setAuthorization
}
