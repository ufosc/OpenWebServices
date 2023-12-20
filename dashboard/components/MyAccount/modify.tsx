'use client'

import { useState } from 'react'
import { UpdateUser, IsAPISuccess, IsAPIFailure } from '@/APIController/API'
import { useCookies } from 'next-client-cookies'
import { Edit, EditOff } from '@carbon/icons-react'
import { Form, useTheme, TextInput, Button } from '@carbon/react'

const editButton = (isEditing : boolean, setIsEditing : Function) => {
  if (isEditing) {
    return (
      <Button kind="danger" onClick={() => setIsEditing(false) }>
	  Cancel
	  <EditOff className="button--arrow" />
      </Button>
    )
  }

  return (
    <Button onClick={() => setIsEditing(true) }>
	  Edit
	  <Edit className="button--arrow" />
    </Button>
  )
}

export default function Modify({ data } : any) {
  const cookies = useCookies()

  const headingColor = () => {
    const { theme } = useTheme()
    return (theme == "white") ? "black" : "white"
  }

  const [ isEditing, setIsEditing ] = useState(false)
  const [ newData, setNewData ] = useState(data)
  const [ hasError, setHasError ] = useState("")
  const [ hasSuccess, setHasSuccess ] = useState(false)
  const submitForm = async (e : any) => {
    e.preventDefault()
    setHasSuccess(false)
    setHasError("")

    if (newData.first_name.length > 20) {
      setHasError("first name cannot be longer than 20 characters");
      return
    }

    if (newData.last_name.length > 20) {
      setHasError("last name cannot be longer than 20 characters");
      return
    }

    if (newData.first_name.length < 2) {
      setHasError("first name cannot be less than 2 characters");
      return
    }

    if (newData.last_name.length < 2) {
      setHasError("last name cannot be less than 2 characters");
      return
    }

    UpdateUser(newData.first_name, newData.last_name, cookies.get('ows-jwt'))
      .then((res) => {
	if (IsAPISuccess(res)) {
	  setHasSuccess(true)
	  setIsEditing(false)
	  return
	}

	let msg = (IsAPIFailure(res) && typeof res.error != "undefined") ?
	  res.error : "An unknown error has occurred. Please try again later."

	setHasError(msg)
      }).catch((err) => {
	setHasError("Server could not be reached. Please try again later")
      })
  }

  return (
    <div className="accountPage">
      <Form style={{ width: "100%", maxWidth: "700px" }}>
	<TextInput
	  id="email"
	  style={{ marginBottom: "15px" }}
	  disabled
	  value={data.email}
	  labelText="Email Address"
	/>
	<TextInput
	  id="first_name"
	  style={{ marginBottom: "15px" }}
	  value={newData.first_name}
	  labelText="First Name"
	  onChange={(e) => setNewData({ first_name: e.target.value,
	    last_name: newData.last_name }) }
	  disabled={ !isEditing }
	/>
	<TextInput
	  id="last_name"
	  style={{ marginBottom: "35px" }}
	  value={newData.last_name}
	  labelText="Last Name"
	  onChange={(e) => setNewData({ first_name: newData.first_name,
	    last_name: e.target.value }) }
	  disabled={ !isEditing }
	/>
	{ editButton(isEditing, setIsEditing) }
	{
	  (isEditing) ? (
	    <Button type="submit" onClick={submitForm}
	      style={{ marginLeft: 10 }}>
	      Update
	    </Button>
	  ) : null
	}
	{
	  (hasError != "") ? (
	    <p style={{ marginTop: 10, marginBottom: 5, color: 'red' }}>
	      Error: { hasError }
	    </p>) : null
	}
	{
	  (hasSuccess) ? (
	    <p style={{ marginTop: 10, marginBottom: 5, color: 'green' }}>
	      Account updated succesfully
	    </p>) : null
	}
      </Form>
    </div>
  )
}
