'use client'

import { GetUser, IsAPISuccess } from '@/APIController/API'
import { redirect } from 'next/navigation'
import { useState } from 'react'
import { useCookies } from 'next-client-cookies'
import Modify from './modify'

import { Loading } from '@carbon/react'

export default function MyAccount() {
  const cookies = useCookies()
  const jwt = cookies.get('ows-jwt')
  if (typeof jwt === "undefined") {
    redirect("/authorize")
  }

  const [data, setData] = useState<Object | null>(null)
  if (data === null) {
    GetUser(jwt).then((res) => {
      if (!IsAPISuccess(res)) {
	cookies.remove('ows-jwt')
	location.replace("/authorize")
	return
      }
      setData(res)
    }).catch((err) => {
      cookies.remove('ows-jwt')
      location.replace("/authorize")
    })
    return (<Loading withOverlay={true} />)
  }

  return (<Modify data={data} />)
}
