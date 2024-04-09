import axios, { AxiosResponse, AxiosError } from 'axios'
import { API_ENDPOINT } from '@/config'

export const IsAPISuccess = (obj : any) => (
  typeof obj !== 'undefined' && typeof obj.message !== 'undefined'
)

export const IsAPIFailure = (obj : any) => (
  typeof obj !== 'undefined' && typeof obj.error !== 'undefined'
)

const stderror = {
  error: 'network_error',
  error_description: 'failed to reach server, ' +
    'please try again later',
}

const handleError = (rej: Function, err: AxiosError) => {
  if (IsAPIFailure(err.response)) {
    rej(err.response!.data)
    return
  }
  rej(stderror)
}

export const GetClient = (id : string) =>
  new Promise((resolve, reject) =>
    axios.get(`${API_ENDPOINT}/client/${encodeURIComponent(id)}`)
      .then((res: AxiosResponse) => resolve(res.data))
      .catch((err: AxiosError) => handleError(reject, err)))

export const DeleteClient = (id : string, token : string) =>
  new Promise((resolve, reject) =>
    axios.delete(`${API_ENDPOINT}/client/${id}`, { headers: {
      'Authorization': `Bearer ${token}`}})
      .then((res: AxiosResponse) => resolve(res.data))
      .catch((err: AxiosError) => handleError(reject, err)))

export const GetClients = (page: number, token: string) =>
  new Promise((resolve, reject) =>
    axios.get(`${API_ENDPOINT}/clients?page=${page}`, { headers: {
      'Authorization': `Bearer ${token}`}})
      .then((res: AxiosResponse) => resolve(res.data))
      .catch((err: AxiosError) => handleError(reject, err)))

export const CreateClient = (form: object, token: string) =>
  new Promise((resolve, reject) =>
    axios.post(`${API_ENDPOINT}/client`, form, { headers: {
      'Authorization': `Bearer ${token}` }})
      .then((res: AxiosResponse) => resolve(res.data))
      .catch((err: AxiosError) => handleError(reject, err)))

export const SignIn = (body: { email: string, password: string }) =>
  new Promise((resolve, reject) =>
    axios.post(`${API_ENDPOINT}/auth/signin`, body)
      .then((res: AxiosResponse) => resolve(res.data))
      .catch((err: AxiosError) => handleError(reject, err)))

export const SignUp = (body: {
  first_name: string, last_name: string,
  email: string, password: string, captcha: string
}) => new Promise((resolve, reject) =>
  axios.post(`${API_ENDPOINT}/auth/signin`, body)
    .then((res: AxiosResponse) => resolve(res.data))
    .catch((err: AxiosError) => handleError(reject, err)))

export const GetUser = (token: string) =>
  new Promise((resolve, reject) =>
    axios.get(`${API_ENDPOINT}/user`, { headers: {
      'Authorization': `Bearer ${token}`} })
      .then((res: AxiosResponse) => resolve(res.data))
      .catch((err: AxiosError) => handleError(reject, err)))

export const GetUsers = (page: number, token: string) =>
  new Promise((resolve, reject) =>
    axios.get(`${API_ENDPOINT}/users?page=${page}`, { headers: {
      'Authorization': `Bearer ${token}` } })
      .then((res: AxiosResponse) => resolve(res.data))
      .catch((err: AxiosError) => handleError(reject, err)))

export const UpdateUser = (firstName: string, lastName: string, token: string) =>
  new Promise((resolve, reject) =>
    axios.put(`${API_ENDPOINT}/user`,
      { first_name: firstName, last_name: lastName },
      { headers: { 'Authorization': `Bearer ${token}` }})
      .then((res: AxiosResponse) => resolve(res.data))
      .catch((err: AxiosError) => handleError(reject, err)))

export const UpdateUserRealms = (form: object, id: string, token: string) =>
  new Promise((resolve, reject) =>
    axios.put(`${API_ENDPOINT}/user/realms/${id}`, form,
    { headers: { 'Authorization': `Bearer ${token}` }})
    .then((res: AxiosResponse) => resolve(res.data))
    .catch((err: AxiosError) => handleError(reject, err)))

export const DeleteUser = (id: string, token: string) =>
  new Promise((resolve, reject) =>
    axios.delete(`${API_ENDPOINT}/user/${id}`, { headers: {
      'Authorization': `Bearer ${token}`}})
      .then((res: AxiosResponse) => resolve(res.data))
      .catch((err: AxiosError) => handleError(reject, err)))

export const ValidateEmail = (email : string) => {
  if (email.match(/^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$/)) {
    return true
  }
  return false
}

export const ValidateClientURLParams = (client : any) => {
  let hasDefined = false
  let hasUndefined = false
  const keys = ["response_type", "client_id", "redirect_uri", "state"]
  for (let i = 0; i < keys.length; i++) {
    if (client[keys[i]] === null) {
      hasUndefined = true
      continue
    }
    hasDefined = true
  }

  // User is either here on their own or they've been redirected by a client.
  // In the latter case, if any client parameter is defined, then all parameters
  // must be defined.
  if (hasDefined && hasUndefined) {
    return false
  }

  return true
}
