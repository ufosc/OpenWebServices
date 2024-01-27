import TableView from './TableView'
import { useState } from 'react'

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
]

export default function Users() {
    const [rows, setRows] = useState([
      {
        id: '0',
        email: 'testing@ufosc.org',
        first_name: 'Dummy',
        last_name: 'User',
        realms: ['clients.read', 'clients.create']
      },
      {
        id: '1',
        email: 'testing2@ufosc.org',
        first_name: 'Dummy',
        last_name: 'User 2',
        realms: ['clients.read', 'clients.create']
      }
    ])

  return (
    <TableView
      rows={rows}
      headers={headers}
      title="Users"
      description="Users are individuals who have signed up for an OSC
      account and have and successfully verified their email address."
    />
  )
}
