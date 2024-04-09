'use client'

import PermissionsForm from '@/components/Permissions/Form'
import { useState } from 'react'
import ClientError from './cerror'
import { Loading } from '@carbon/react'
import { GetClient } from '@/API'

export type ClientDefinition = {
  response_type: string;
  client_id:     string;
  redirect_uri:  string;
  state:         string;
}

export const Permissions = (props: { client: ClientDefinition }) => {

  const client = props.client

  // Validate response_type.
  if (client.response_type !== "code" && client.response_type !== "token") {
    return <ClientError><p>
        The authorization request's response_type was not one of
        'authorization_code' or 'token'.
      </p></ClientError>
  }

  // Ensure client_id is not null before we request the API for it.
  if (client.client_id === "" || client.client_id === null) {
    return <ClientError><p>Missing client_id parameter.</p></ClientError>
  }

  // Ensure state is not null.
  if (client.state === "" || client.state === null) {
    return <ClientError><p>The state parameter cannot be an empty string.</p>
      </ClientError>
  }

  // Ensure redirect_uri is not blank.
  if (client.redirect_uri === "" || client.redirect_uri === null) {
    return <ClientError>
	<p>The redirect_uri parameter cannot be an empty string.</p>
      </ClientError>
  }

  // Fetch and verify client information.
  const [clientError, setClientError] = useState("")
  const [clientData, setClientData] = useState<ClientDefinition | null>(null)

  if (clientData === null && clientError === "") {
    GetClient(props.client.client_id).then((_res) => {
      let res = (_res as ClientDefinition)

      // Ensure response types are same.
      if (client.response_type !== res.response_type) {
	setClientError("URL parameter response_type and client response_type do not match.")
	return
      }

      // Ensure redirect_uri's are the same.
      if (client.redirect_uri !== res.redirect_uri) {
	setClientError("URL parameter redirect_uri and client redirect_uri do not match.")
	return
      }

      setClientData(res)
    }).catch((err) => setClientError(err.error_description))
  }

  if (clientError === "" && clientData === null) {
    return (<Loading style={{margin: "auto auto auto auto"}} withOverlay={false} />)
  }

  if (clientError !== "") {
    return <ClientError>
        <p>The client could not be verified. {clientError}</p>
      </ClientError>
  }

  return <PermissionsForm client={clientData} state={props.client.state}/>
}
