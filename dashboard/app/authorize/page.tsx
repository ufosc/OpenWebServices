'use client'

import { GetClient, IsAPIFailure, IsAPISuccess } from '@/APIController/API'
import { ValidateClientURLParams } from '@/APIController/Validation'
import { TypeAuthGrant, TypeClient, TypeGetClientResponse } from '@/APIController/types'
import AlertBanner from '@/components/AlertBanner'
import ImageBanner from '@/components/ImageBanner/imagebanner'
import PermissionsForm from '@/components/Permissions/Form'
import SigninForm from '@/components/SigninForm'
import SignupForm from '@/components/SignupForm'
import { Accordion, AccordionItem, Link, Loading } from '@carbon/react'
import { redirect, useSearchParams } from 'next/navigation'
import { useState } from 'react'
import { useCookies } from 'next-client-cookies'
import './style.scss'

const ClientError = (props: { children: any }) => {
  return (
    <AlertBanner heading="Error: Invalid Client Authorization Request">
      <p>
	Sign-in could not be completed at this time, most likely due to an invalid
	client request. If you're the client's developer, please ensure that
	you've implemented the authorization request as specified in the OAuth2
	documentation:
      </p>
      <br />
      <Link href="https://datatracker.ietf.org/doc/html/rfc6749">
	IETF RFC 6749: The OAuth 2.0 Authorization Framework
      </Link>
      <div style={{ marginTop: "20px" }}>
	<Accordion>
          <AccordionItem title="Additional Debug Information">
            {props.children}
          </AccordionItem>
        </Accordion>
      </div>
    </AlertBanner>
  )
}

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
  const [clientData, setClientData] = useState<TypeClient | null>(null)
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

export default function Page() {
  const cookies = useCookies()
  const jwt = cookies.get('ows-jwt')
  const [view, setView] = useState<"signin" | "signup">("signin")

  // Gather client parameters.
  const searchParams = useSearchParams()
  const client : TypeAuthGrant = {
    response_type: searchParams.get('response_type'),
    client_id: searchParams.get('client_id'),
    redirect_uri: searchParams.get('redirect_uri'),
    state: searchParams.get('state')
  }

  const renderForm = () => {
    // Validate URL params.
    if (!ValidateClientURLParams(client)) {
      return (
	<ClientError>
	  <p>
            The authorization request is missing required URL parameters. Please
            ensure you're passing all required parameters in URL encoded format,
            using the "application/x-www-form-urlencoded" format.
          </p>
	</ClientError>
      )
    }

    if (typeof jwt === "undefined" && view === "signin") {
      return (<SigninForm setView={setView} />)
    }

    if (typeof jwt === "undefined" && view === "signup") {
      return (<SignupForm setView={setView} />)
    }

    // No URL parameters and user is already signed in, redirect to account page.
    if (typeof jwt !== "undefined" && client.response_type === null) {
      redirect("/account")
      return null
    }

    // Display permissions page.
    return (<Permissions client={client} />)
  }

  return (
    <div className="loginPage">
      <div className="loginPage--form">
	{ renderForm() }
      </div>
      <ImageBanner/>
    </div>
  )
}
