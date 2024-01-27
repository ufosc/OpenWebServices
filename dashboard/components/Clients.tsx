import TableView from './TableView'
import { useState } from 'react'

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
    key: 'ttl',
    header: 'TTL',
  },
]

export default function Clients() {
  const [rows, setRows] = useState([
    {
      id: '0',
      name: 'Dummy Client',
      response_type: 'token',
      redirect_uri: "https://google.com",
      scope: "email, public",
      ttl: "176000"
    }
  ])

  return (
    <TableView
      rows={rows}
      headers={headers}
      title="Clients"
      description="Clients are applications that have been
      authorized to use OSC accounts for authentication. They request
      users for permission and, if granted, may securely access
      their private data."
    />
  )
}
