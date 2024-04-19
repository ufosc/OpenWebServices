'use client'

import TableView from './TableView'
import PaginationNav from '@carbon/react/lib/components/PaginationNav/PaginationNav'
import { DeleteClient, GetClients, CreateClient } from '@/API'
import { useEffect, useState } from 'react'
import { createPortal } from 'react-dom'
import { useRouter } from 'next/navigation'
import { useCookies } from 'next-client-cookies'

import {
  Loading, InlineNotification, Modal, ModalBody, Form,
  TextInput, TextArea, Select, SelectItem, FormGroup, Checkbox,
  Button,
} from '@carbon/react'

// Data table headers.
const headers = [
  {
    key: 'id',
    header: 'Client ID',
  },
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
    header: 'TTL (sec)',
  },
]

type ClientResponse = {
  clients: any[],
  total_count: number,
  count: number,
}

export default function Clients() {
  const router = useRouter()
  const cookies = useCookies()
  const token = cookies.get('ows-access-token')
  useEffect(() => {
    if (typeof token === "undefined") {
      router.push("/authorize")
    }
  }, [])

  const [page, setPage] = useState<number>(0)
  const [numPages, setNumPages] = useState<number>(1)
  const [rows, setRows] = useState<any>([])
  const [isLoading, setIsLoading] = useState<boolean>(true)
  const [hasNotif, setHasNotif] = useState<boolean>(false)
  const [notifData, setNotifData] = useState<{title: string,
    subtitle: string, kind: 'success' | 'error' }>({
      title: "", subtitle: "", kind: "error"
    })

  const [clientModal, setClientModal] = useState<boolean>(false)
  const [createClientForm, setCreateClientForm] = useState({
    name: "",
    description: "",
    response_type: "code",
    redirect_uri: "",
    email: false,
  })

  const fetchTable = () => {
    GetClients(page, token as string).then((res) => {
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
      const pageCount = Math.ceil((res as ClientResponse).total_count / 10)
      if (pageCount !== numPages && pageCount > 0) {
        setNumPages(pageCount)
      }

    }).catch((err) => {
      cookies.remove('ows-access-tokens')
      router.push('/authorize')
    })
  }

  // Fetch table every time page is changed.
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
      await DeleteClient(selectedRows[i].id, token as string)
        .catch(err => { hasError = true })
    }

    if (hasError) {
      setNotifData({
        title: "Error Deleting Clients",
        subtitle: "You are not authorized to delete clients",
        kind: "error",
      })
      setHasNotif(true)
      setTimeout(() => { setHasNotif(false) }, 5000)
    }

    fetchTable()
    setIsLoading(false)
  }

  const submitCreateClientForm = (e : any) => {
    e.preventDefault()
    if (typeof token === "undefined")
      return

    let form = {
      name: createClientForm.name,
      description: createClientForm.name,
      response_type: createClientForm.response_type,
      redirect_uri: createClientForm.redirect_uri,
      scope: ['public']
    }

    if (createClientForm.email && createClientForm.response_type === 'code')
      form.scope.push('email')

    CreateClient(form, token)
      .then((res) => {
        setCreateClientForm({
          name: "",
          description: "",
          response_type: "code",
          redirect_uri: "",
          email: false,
        })

        setClientModal(false)
        setNotifData({
          title: "Success",
          subtitle: "Client created successfully",
          kind: "success"
        })
        setHasNotif(true)
        setTimeout(() => { setHasNotif(false) }, 5000)
        fetchTable()
      })
      .catch((err) => {
        setClientModal(false)
        if (err.error == 'insufficient_scope') {
          setNotifData({
            title: "Error Creating Client",
            subtitle: "You are not authorized to create clients",
            kind: "error"
          })
        } else {
          setNotifData({
            title: "Error",
            subtitle: err.error_description,
            kind: "error"
          })
        }
        setHasNotif(true)
        setTimeout(() => { setHasNotif(false) }, 5000)
      })
  }

  return (
    <>
      {
        (hasNotif) ? (
          <InlineNotification
            kind={notifData.kind}
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
          (<Loading id='decoration--landing' withOverlay={true} />)
          : null
      }
      {
        (clientModal) ?
          createPortal(
            <Modal
              open={clientModal}
              onRequestClose={() => { setClientModal(false) }}
              modalHeading="Create a new client"
              modalLabel="Clients"
              primaryButtonText="Add"
              secondaryButtonText="Cancel"
              onRequestSubmit={() => {
                let submit = document.getElementById('submit-create-client')
                if (typeof submit !== 'undefined' && submit !== null) {
                  submit.click()
                }
              }}
            >
              <ModalBody>
                <Form className="form" id='create-client-form'
                  onSubmit={submitCreateClientForm}>
                  <TextInput
                    id="name"
                    style={{ marginBottom: "15px" }}
                    placeholder="My Web App"
                    labelText="Name"
                    required
                    onChange = {(e) => setCreateClientForm({
                      name: e.target.value,
                      description: createClientForm.description,
                      response_type: createClientForm.response_type,
                      redirect_uri: createClientForm.redirect_uri,
                      email: createClientForm.email,
                    })}
                  />
                  <TextArea
                    id="description"
                    style={{ marginBottom: "15px" }}
                    placeholder="Description"
                    labelText="Description"
                    required
                    onChange = {(e) => setCreateClientForm({
                      name: createClientForm.name,
                      description: e.target.value,
                      response_type: createClientForm.response_type,
                      redirect_uri: createClientForm.redirect_uri,
                      email: createClientForm.email,
                    })}
                  />
                  <Select id="response-type"
                    style={{ marginBottom: "15px" }}
                    defaultValue="placeholder-item"
                    required
                    onChange = {(e) => setCreateClientForm({
                      name: createClientForm.name,
                      description: createClientForm.description,
                      response_type: e.target.value,
                      redirect_uri: createClientForm.redirect_uri,
                      email: createClientForm.email,
                    })}
                  >
                    <SelectItem value="code" text="Authorization code" />
                    <SelectItem value="token" text="Implicit Token" />
                  </Select>
                  <TextInput
                    id="redirect_uri"
                    style={{ marginBottom: "15px" }}
                    placeholder="https://example.com/oauth"
                    labelText="Redirect URI"
                    required
                    onChange = {(e) => setCreateClientForm({
                      name: createClientForm.name,
                      description: createClientForm.description,
                      response_type: createClientForm.response_type,
                      redirect_uri: e.target.value,
                      email: createClientForm.email,
                    })}
                  />
                  <FormGroup legendText="Client Scope"
                    style={{ marginBottom: "15px" }}
                  >
                    <Checkbox labelText="Public" checked id="scope-check-public" />
                    { (createClientForm.response_type === 'code') ? (
                    <Checkbox labelText="Email" id="scope-check-email"
                      onChange={ (e, { checked, id }) => setCreateClientForm({
                        name: createClientForm.name,
                        description: createClientForm.description,
                        response_type: createClientForm.response_type,
                        redirect_uri: createClientForm.redirect_uri,
                        email: checked,
                      })}
                    />) : null
                    }
                  </FormGroup>
                  <Button id='submit-create-client' type='submit'
                    style={{ display: 'none' }}>
                    Submit
                  </Button>
                </Form>
              </ModalBody>
            </Modal>,
            document.body
          )
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
        hasCreateButton={true}
        hasModifyButton={false}
        onCreate={() => { setClientModal(true) }}
        onDelete={onDelete}
        onEdit={() => {}}
      />
      <PaginationNav itemsShown={5} totalItems={numPages} onChange={pageChange} />
    </>
  )
}
