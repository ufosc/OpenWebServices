'use client'

import TableView from './TableView'
import { GetUsers, DeleteUser, IsAPISuccess, IsAPIFailure } from '@/APIController/API'
import { useState, useEffect } from 'react'
import { useCookies } from 'next-client-cookies'
import { useRouter } from 'next/navigation'
import { Loading, InlineNotification } from '@carbon/react'

// IBM Carbon is serious dogshit.
import PaginationNav from '@carbon/react/lib/components/PaginationNav/PaginationNav'

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
  total_count: number,
  count: number,
}

export default function Users() {
  const router = useRouter()
  const cookies = useCookies()
  const token = cookies.get('ows-access-token')
  if (typeof token === "undefined") {
    router.push("/authorize")
    return
  }

  const [page, setPage] = useState<number>(0)
  const [numPages, setNumPages] = useState<number>(1)
  const [rows, setRows] = useState<any>([])
  const [isLoading, setIsLoading] = useState<boolean>(true)
  const [hasNotif, setHasNotif] = useState<boolean>(false)
  const [notifData, setNotifData] = useState<{title: string,
    subtitle: string}>({title: "", subtitle: ""})

  const fetchTable = () => {
    GetUsers(page, token as string).then((res) => {
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

      setIsLoading(false)
      const pageCount = Math.ceil((res as UsersResponse).total_count / 20)
      if (pageCount !== numPages && pageCount > 0) {
        setNumPages(pageCount)
      }

    }).catch((err) => {
      cookies.remove('ows-access-tokens')
      router.push("/authorize")
    })
  }

  useEffect(fetchTable, [page])

  const pageChange = (newPage : number) => {
    if (newPage === page) {
      return
    }
    setIsLoading(true)
    setPage(newPage)
  }

  const onDelete = async (selectedRows : { id: string }[]) => {
    setIsLoading(true)
    let hasError = false
    for (let i = 0; i < selectedRows.length; ++i) {
      await DeleteUser(selectedRows[i].id, token as string).then((res) => {
        if (IsAPIFailure(res)) {
          hasError = true
        }
      }).catch(err => { hasError = true })
    }

    if (hasError) {
      setNotifData({
        title: "Error Deleting Users",
        subtitle: "You are not authorized to delete users"
      })
      setHasNotif(true)
      setTimeout(() => { setHasNotif(false) }, 5000)
    }

    fetchTable()
    setIsLoading(false)
  }

  return (
    <>
      {
        (hasNotif) ? (
          <InlineNotification
            kind="error"
            onClose={() => setHasNotif(false) }
            onCloseButtonClick={() => setHasNotif(false)}
            statusIconDescription="notification"
            subtitle={notifData.subtitle}
            title={notifData.title}
            style={{ position: "fixed", bottom: 5, left: 5}}
          />
        ) : null
      }
      {
	(isLoading) ?
	  (<Loading id="decoration--loading" withOverlay={true} />)
	  : null
      }
      <TableView
        rows={rows}
        headers={headers}
        title="Users"
        description="Users are individuals who have signed up for an OSC
        account and have and successfully verified their email address."
        hasAddButton={false}
        onDelete={onDelete}
      />
      <PaginationNav itemsShown={5} totalItems={numPages} onChange={pageChange} />
    </>
  )
}
