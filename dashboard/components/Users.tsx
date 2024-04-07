'use client'

import TableView from './TableView'
import { useState, useEffect } from 'react'
import { useCookies } from 'next-client-cookies'
import { createPortal } from 'react-dom'
import { useRouter } from 'next/navigation'

import {
  Loading, InlineNotification, Modal, ModalBody,
  Form, TextInput, FormGroup, Checkbox, Button,
} from '@carbon/react'

import {
  GetUsers, DeleteUser,
  IsAPISuccess, IsAPIFailure, UpdateUserRealms
} from '@/APIController/API'

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
    subtitle: string, kind: 'success' | 'error'}>({
      title: "", subtitle: "", kind: "error",
    })

  const [userModal, setUserModal] = useState<boolean>(false)
  const [modifyUserForm, setModifyUserForm] = useState({
    id: "",
    first_name: "",
    last_name: "",
    scope: {
      clients_read: false,
      clients_delete: false,
      clients_create: false,
      users_read: false,
      users_delete: false,
      users_update: false,
    }
  })

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
      const pageCount = Math.ceil((res as UsersResponse).total_count / 10)
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
        subtitle: "You are not authorized to delete users",
        kind: "error",
      })
      setHasNotif(true)
      setTimeout(() => { setHasNotif(false) }, 5000)
    }

    fetchTable()
    setIsLoading(false)
  }

  const onEdit = (cells : Array<{id: string, value: any}>) => {
    let id = cells[0].id.split(':')
    if (id.length !== 2)
      return null

    let scope = cells[3].value.split(", ")
    setModifyUserForm({
      id: id[0],
      first_name: cells[1].value,
      last_name: cells[2].value,
      scope: {
        clients_read: scope.includes("clients.read"),
        clients_delete: scope.includes("clients.delete"),
        clients_create: scope.includes("clients.create"),
        users_read: scope.includes("users.read"),
        users_delete: scope.includes("users.delete"),
        users_update: scope.includes("users.update"),
      },
    })

    setUserModal(true)
  }

  const submitModifyForm = (e : any) => {
    e.preventDefault()
    if (typeof token === "undefined")
      return null

    let realms = []
    if (modifyUserForm.scope.clients_read) {
      realms.push("clients.read")
    }

    if (modifyUserForm.scope.clients_delete) {
      realms.push("clients.delete")
    }

    if (modifyUserForm.scope.clients_create) {
      realms.push("clients.create")
    }

    if (modifyUserForm.scope.users_read) {
      realms.push("users.read")
    }

    if (modifyUserForm.scope.users_delete) {
      realms.push("users.delete")
    }

    if (modifyUserForm.scope.users_update) {
      realms.push("users.update")
    }

    let form = {
      first_name: modifyUserForm.first_name,
      last_name: modifyUserForm.last_name,
      realms: realms,
    }

    UpdateUserRealms(form, modifyUserForm.id, token).then((res) => {
      setModifyUserForm({
        id: "",
        first_name: "",
        last_name: "",
        scope: {
          clients_read: false,
          clients_delete: false,
          clients_create: false,
          users_read: false,
          users_delete: false,
          users_update: false,
        }
      })

      if (IsAPISuccess(res)) {
        setUserModal(false)
        setNotifData({
          title: "Success",
          subtitle: "User modified succesfully",
          kind: "success",
        })
        setHasNotif(true)
        setTimeout(() => { setHasNotif(false) }, 5000)
        fetchTable()
        return
      }

      if (IsAPIFailure(res) && res.error == 'insufficient_scope') {
        setUserModal(false)
        setNotifData({
          title: "Error Modifying User",
          subtitle: "You are not authorized to modify users",
          kind: "error",
        })
        setHasNotif(true)
        setTimeout(() => { setHasNotif(false) }, 5000)
        return
      }

      let msg = (IsAPIFailure(res) && typeof res.error !== "undefined") ?
        res.error : "An unknown error has occured. Please try again later."

      setUserModal(false)
      setNotifData({
        title: "Error",
        subtitle: msg,
        kind: "error",
      })
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
	  (<Loading id="decoration--loading" withOverlay={true} />)
	  : null
      }
      {
        (userModal) ? createPortal(
          <Modal
            open={userModal}
            onRequestClose={() => { setUserModal(false) }}
            modalHeading="Modify User"
            modalLabel="Users"
            primaryButtonText="Add"
            secondaryButtonText="Cancel"
            onRequestSubmit={() => {
              let submit = document.getElementById('submit-modify-user')
              if (typeof submit !== 'undefined' && submit !== null) {
                submit.click()
              }
            }}
          >
            <ModalBody>
              <Form className="form" id='modify-user-form'
                onSubmit={submitModifyForm}>
                <TextInput
                  id="first-name"
                  style={{ marginBottom: "15px" }}
                  labelText="First Name"
                  value={modifyUserForm.first_name}
                  required
                  onChange = {(e) => setModifyUserForm({
                    id: modifyUserForm.id,
                    first_name: e.target.value,
                    last_name: modifyUserForm.last_name,
                    scope: modifyUserForm.scope,
                  })}
                />
                <TextInput
                  id="last-name"
                  style={{ marginBottom: "15px" }}
                  value={modifyUserForm.last_name}
                  labelText="Last Name"
                  required
                  onChange = {(e) => setModifyUserForm({
                    id: modifyUserForm.id,
                    first_name: modifyUserForm.first_name,
                    last_name: e.target.value,
                    scope: modifyUserForm.scope,
                  })}
                />
                <FormGroup legendText="User Realms">
                  <Checkbox labelText="clients.read" id="realm-clients-read"
                    checked={modifyUserForm.scope.clients_read}
                    onChange = {(e, { checked, id }) => setModifyUserForm({
                      id: modifyUserForm.id,
                      first_name: modifyUserForm.first_name,
                      last_name: modifyUserForm.last_name,
                      scope: {
                        clients_read: checked,
                        clients_delete: modifyUserForm.scope.clients_delete,
                        clients_create: modifyUserForm.scope.clients_create,
                        users_read: modifyUserForm.scope.users_read,
                        users_delete: modifyUserForm.scope.users_delete,
                        users_update: modifyUserForm.scope.users_update,
                      },
                    })}
                  />
                  <Checkbox labelText="clients.delete" id="realm-clients-delete"
                    checked={modifyUserForm.scope.clients_delete}
                    onChange = {(e, { checked, id }) => setModifyUserForm({
                      id: modifyUserForm.id,
                      first_name: modifyUserForm.first_name,
                      last_name: modifyUserForm.last_name,
                      scope: {
                        clients_read: modifyUserForm.scope.clients_read,
                        clients_delete: checked,
                        clients_create: modifyUserForm.scope.clients_create,
                        users_read: modifyUserForm.scope.users_read,
                        users_delete: modifyUserForm.scope.users_delete,
                        users_update: modifyUserForm.scope.users_update,
                      },
                    })}
                  />
                  <Checkbox labelText="clients.create" id="realm-clients-create"
                    checked={modifyUserForm.scope.clients_create}
                    onChange = {(e, { checked, id }) => setModifyUserForm({
                      id: modifyUserForm.id,
                      first_name: modifyUserForm.first_name,
                      last_name: modifyUserForm.last_name,
                      scope: {
                        clients_read: modifyUserForm.scope.clients_read,
                        clients_delete: modifyUserForm.scope.clients_delete,
                        clients_create: checked,
                        users_read: modifyUserForm.scope.users_read,
                        users_delete: modifyUserForm.scope.users_delete,
                        users_update: modifyUserForm.scope.users_update,
                      },
                    })}
                  />
                  <Checkbox labelText="users.read" id="realm-users-read"
                    checked={modifyUserForm.scope.users_read}
                    onChange = {(e, { checked, id }) => setModifyUserForm({
                      id: modifyUserForm.id,
                      first_name: modifyUserForm.first_name,
                      last_name: modifyUserForm.last_name,
                      scope: {
                        clients_read: modifyUserForm.scope.clients_read,
                        clients_delete: modifyUserForm.scope.clients_delete,
                        clients_create: modifyUserForm.scope.clients_create,
                        users_read: checked,
                        users_delete: modifyUserForm.scope.users_delete,
                        users_update: modifyUserForm.scope.users_update,
                      },
                    })}
                  />
                  <Checkbox labelText="users.delete" id="realm-users-delete"
                    checked={modifyUserForm.scope.users_delete}
                    onChange = {(e, { checked, id }) => setModifyUserForm({
                      id: modifyUserForm.id,
                      first_name: modifyUserForm.first_name,
                      last_name: modifyUserForm.last_name,
                      scope: {
                        clients_read: modifyUserForm.scope.clients_read,
                        clients_delete: modifyUserForm.scope.clients_delete,
                        clients_create: modifyUserForm.scope.clients_create,
                        users_read: modifyUserForm.scope.users_read,
                        users_delete: checked,
                        users_update: modifyUserForm.scope.users_update,
                      },
                    })}
                  />
                  <Checkbox labelText="users.update" id="realm-users-update"
                    checked={modifyUserForm.scope.users_update}
                    onChange = {(e, { checked, id }) => setModifyUserForm({
                      id: modifyUserForm.id,
                      first_name: modifyUserForm.first_name,
                      last_name: modifyUserForm.last_name,
                      scope: {
                        clients_read: modifyUserForm.scope.clients_read,
                        clients_delete: modifyUserForm.scope.clients_delete,
                        clients_create: modifyUserForm.scope.clients_create,
                        users_read: modifyUserForm.scope.users_read,
                        users_delete: modifyUserForm.scope.users_delete,
                        users_update: checked,
                      },
                    })}
                  />
                </FormGroup>
                <Button id='submit-modify-user' type='submit'
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
        title="Users"
        description="Users are individuals who have signed up for an OSC
        account and have and successfully verified their email address."
        hasCreateButton={false}
        hasModifyButton={true}
        onDelete={onDelete}
        onCreate={() => {}}
        onEdit={onEdit}
      />
      <PaginationNav itemsShown={5} totalItems={numPages} onChange={pageChange} />
    </>
  )
}
