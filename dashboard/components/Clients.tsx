'use client'

import TableView from './TableView'
import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { GetClients, IsAPISuccess } from '@/APIController/API'
import { useCookies } from 'next-client-cookies'
import { Loading } from '@carbon/react'
import PaginationNav from '@carbon/react/lib/components/PaginationNav/PaginationNav'

// Data table headers.
const headers = [
  {
    key: 'name',
    header: 'Name',
  },
  {
    key: 'response_type',
    header: 'Type',
  },
  {
    key: 'redirect_uri',
    header: 'URI',
  },
  {
    key: 'scope',
    header: 'Scope',
  },
  {
    key: 'created_at',
    header: 'Created at',
  },
  {
    key: 'ttl',
    header: 'TTL',
  },
]

type ClientResponse = {
  clients: any[],
  total_count: number,
}

export default function Clients() {
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

  useEffect(() => {
    GetClients(page, token as string).then((res) => {
      if (!IsAPISuccess(res)) {
        cookies.remove('ows-access-tokens')
        router.push('/authorize')
        return
      }

      setRows((res as ClientResponse).clients.map(client => {
        let text = ""
        for (let i = 0; i < client.scope.length; ++i) {
          text += client.scope[i]
          if (i + 1 !== client.scope.length) {
            text += ", "
          }
        }
        client.scope = text
        return client
      }))

      setIsLoading(false)

      const pageCount = Math.ceil((res as ClientResponse).total_count / 20)
      if (pageCount !== numPages && pageCount > 0) {
        setNumPages(pageCount)
      }

    }).catch((err) => {
      cookies.remove('ows-access-tokens')
      router.push('/authorize')
    })
  }, [page])

  const pageChange = (newPage : number) => {
    if (newPage === page) {
      return
    }
    setIsLoading(true)
    setPage(newPage)
  }

  return (
    <>
      {
        (isLoading) ?
          (<Loading id='decoration--landing' withOverlay={true} />)
          : null
      }
      <TableView
        rows={rows}
        headers={headers}
        title="Clients"
        description="Clients are applications that have been
        authorized to use OSC accounts for authentication. They request
        users for permission and, if granted, may securely access
        their private data."
        hasAddButton={true}
      />
      <PaginationNav itemsShown={5} totalItems={numPages} onChange={pageChange} />
    </>
  )
}
