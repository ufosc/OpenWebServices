'use client'

import TableView from './TableView'
import { GetUsers, IsAPISuccess } from '@/APIController/API'
import { useState, useEffect } from 'react'
import { useCookies } from 'next-client-cookies'
import { useRouter } from 'next/navigation'

// Data table headers.
const headers = [
  {
    key: 'email',
    header: 'Email'
  },
  {
    key: 'first_name',
    header: 'First Name'
  },
  {
    key: 'last_name',
    header: 'Last Name'
  },
  {
    key: 'realms',
    header: 'Realms'
  },
  {
    key: 'created_at',
    header: 'Created at'
  }
]

type UsersResponse = {
  users: any[],
}

export default function Users() {
  const router = useRouter()
  const cookies = useCookies()
  const token = cookies.get('ows-access-token')
  if (typeof token === "undefined") {
    router.push("/authorize")
    return
  }

  const [rows, setRows] = useState<any>([])
  useEffect(() => {
    GetUsers(0, token as string).then((res) => {
      if (!IsAPISuccess(res)) {
        cookies.remove('ows-access-tokens')
        router.push('/authorize')
        return
      }

      setRows((res as UsersResponse).users.map(user => {
        let text = ""
        for (let i = 0; i < user.realms.length; ++i) {
          text += user.realms[i]
          if (i + 1 !== user.realms.length) {
            text += ", "
          }
        }
        user.realms = text
        return user
      }))
    }).catch((err) => {
      cookies.remove('ows-access-tokens')
      router.push("/authorize")
    })
  }, [])

  return (
    <TableView
      rows={rows}
      headers={headers}
      title="Users"
      description="Users are individuals who have signed up for an OSC
      account and have and successfully verified their email address."
      hasAddButton={false}
    />
  )
}
