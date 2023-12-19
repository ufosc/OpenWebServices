import axios, { AxiosResponse, AxiosError } from 'axios'
import { TypeSigninBody, TypeSignupBody, APIResponse } from './types'

export const GetClient = (id : string) => {
  return new Promise((resolve: Function, reject: Function) => {
    axios.get(`client/${encodeURIComponent(id)}`).then((res : AxiosResponse) => {
      resolve(res.data)
    }).catch((err : AxiosError) => {
      if (err.response && IsAPIFailure(err.response.data)) {
	resolve(err.response.data)
	return
      }
      reject(err)
    })
  })
}

export const PostSignin = (body: TypeSigninBody) => {
  return new Promise((resolve: Function, reject: Function) => {
    axios.post('/auth/signin', body)
      .then((res : AxiosResponse) => {
	resolve(res.data)
      })
      .catch((err : AxiosError) => {
	if (err.response && IsAPIFailure(err.response.data)) {
	  resolve(err.response.data)
	  return
	}
	reject(err)
      })
  })
}

export const PostSignup = (body: TypeSignupBody) => {
  return new Promise((resolve: Function, reject: Function) => {
    axios.post('/auth/signup', body)
      .then((res : AxiosResponse) => {
	resolve(res.data)
      })
      .catch((err : AxiosError) => {
	if (err.response && IsAPIFailure(err.response.data)) {
	  resolve(err.response.data)
	  return
	}
	reject(err)
      })
  })
}

export const GetUser = (jwt : string) => {
  return new Promise((resolve: Function, reject: Function) => {
    axios.get('/user', { headers: { 'Authorization': `Bearer ${jwt}`} })
      .then((res : AxiosResponse) => {
	resolve(res.data)
      })
      .catch((err : AxiosError) => {
	if (err.response && IsAPIFailure(err.response.data)) {
	  resolve(err.response.data)
	  return
	}
	reject(err)
      })
  })
}

export function IsAPISuccess(obj : any): obj is APIResponse {
  return typeof (obj as APIResponse).message !== 'undefined'
}

export function IsAPIFailure(obj : any): obj is APIResponse {
  return typeof (obj as APIResponse).error !== 'undefined'
}
