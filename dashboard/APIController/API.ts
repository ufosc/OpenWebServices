import axios, { AxiosResponse, AxiosError } from 'axios'
import { TypeSigninBody, TypeSignupBody, APIResponse } from './types'

const RootURL = "http://localhost:8080"

export const GetClient = (id : string) => {
  return new Promise((resolve: Function, reject: Function) => {
    axios.get(`${RootURL}/client/${encodeURIComponent(id)}`).then((res : AxiosResponse) => {
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

export const DeleteClient = (id : string, token : string) => {
  return new Promise((resolve: Function, reject: Function) => {
    axios.delete(`${RootURL}/client/${id}`, { headers: {
      'Authorization': `Bearer ${token}`}}).then((res : AxiosResponse) => {
        resolve(res.data)
      })
      .catch((err: AxiosError) => {
        if (err.response && IsAPIFailure(err.response.data)) {
          resolve(err.response.data)
          return
        }
        reject(err)
      })
  })
}

export const GetClients = (page : number, token : string) => {
  return new Promise((resolve: Function, reject: Function) => {
    axios.get(`${RootURL}/clients?page=${page}`, {
      headers: { 'Authorization': `Bearer ${token}`}})
      .then((res : AxiosResponse) => {
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

export const PostCreateClient = (form: object, token : string) => {
  return new Promise((resolve: Function, reject: Function) => {
    axios.post(`${RootURL}/client`, form, {
      headers: { 'Authorization': `Bearer ${token}` }})
      .then((res : AxiosResponse) => {
        resolve(res.data);
      }).catch((err: AxiosError) => {
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
    axios.post(`${RootURL}/auth/signin`, body)
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
    axios.post(`${RootURL}/auth/signup`, body)
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

export const GetUser = (token : string) => {
  return new Promise((resolve: Function, reject: Function) => {
    axios.get(`${RootURL}/user`, { headers: { 'Authorization': `Bearer ${token}`} })
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

export const GetUsers = (page : number, token : string) => {
  return new Promise((resolve: Function, reject: Function) => {
    axios.get(`${RootURL}/users?page=${page}`, { headers: { 'Authorization': `Bearer ${token}` } })
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

export const UpdateUser = (firstName : string, lastName : string,
  token : string) => {
  return new Promise((resolve: Function, reject: Function) => {
    axios.put(`${RootURL}/user`, { first_name: firstName, last_name: lastName },
      { headers: { 'Authorization': `Bearer ${token}` }})
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

export const UpdateUserRealms = (form: object, id: string, token: string) => {
  return new Promise((resolve: Function, reject: Function) => {
    axios.put(`${RootURL}/user/realms/${id}`, form,
      { headers: { 'Authorization': `Bearer ${token}` }})
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

export const DeleteUser = (id : string, token : string) => {
  return new Promise((resolve: Function, reject: Function) => {
    axios.delete(`${RootURL}/user/${id}`, { headers: {
      'Authorization': `Bearer ${token}`}}).then((res : AxiosResponse) => {
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
