import PermissionsForm from '@/components/Permissions/Form'
import { useState } from 'react'
import { TypeAuthGrant, TypeClient, TypeGetClientResponse } from '@/APIController/types'
import { GetClient, IsAPIFailure, IsAPISuccess } from '@/APIController/API'
import { APIResponse } from '@/APIController/types'
import ClientError from './cerror'
import { Loading } from '@carbon/react'

const Permissions = (props: { client : TypeAuthGrant }) => {

  // Validate response_type.
  if (props.client.response_type !== "code" && props.client.response_type !== "token") {
    return (
      <ClientError>
	<p>
	  The authorization request's response_type was not one of
	  'authorization_code' or 'token'.
	</p>
      </ClientError>
    )
  }

  // Ensure client_id is not null before we request the API for it.
  if (props.client.client_id === "" || props.client.client_id === null) {
    return (
      <ClientError>
	<p>Missing client_id parameter.</p>
      </ClientError>
    )
  }

  // Ensure state is not blank.
  if (props.client.state === "" || props.client.state === null) {
    return (
      <ClientError>
	<p>The state parameter cannot be an empty string.</p>
      </ClientError>
    )
  }

  // Ensure redirect_uri is not blank.
  if (props.client.redirect_uri === "" || props.client.redirect_uri === null) {
    return (
      <ClientError>
	<p>The redirect_uri parameter cannot be an empty string.</p>
      </ClientError>
    )
  }

  // Fetch and verify client information.
  const [clientError, setClientError] = useState("")
  const [clientData, setClientData] = useState<APIResponse | TypeClient | null>(null)
  if (clientData === null && clientError === "") {
    GetClient(props.client.client_id).then((res) => {
      if (!IsAPISuccess(res)) {
	setClientError((IsAPIFailure(res) && typeof res.error != "undefined") ?
	  res.error : clientError)
	return
      }

      // Ensure response types are same.
      if (props.client.response_type !== (res as TypeGetClientResponse).response_type) {
	setClientError("URL parameter response_type and client response_type do not match.")
	return
      }

      // Ensure redirect_uri's are the same.
      if (props.client.redirect_uri !== (res as TypeGetClientResponse).redirect_uri) {
	setClientError("URL parameter redirect_uri and client redirect_uri do not match.")
	return
      }

      setClientData(res)

    }).catch((err) => {
      setClientError("Client could not be found or server failed to respond.")
    })
  }

  if (clientError === "" && clientData === null) {
    return (<Loading style={{margin: "auto auto auto auto"}} withOverlay={false} />)
  }

  if (clientError !== "") {
    return (
      <ClientError>
	<p>The client could not be verified. {clientError}</p>
      </ClientError>
    )
  }

  return (<PermissionsForm client={clientData} state={props.client.state}/>)
}

export default Permissions
