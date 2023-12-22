'use client'

import './style.scss'
import { ArrowRight, Wikis } from '@carbon/icons-react'
import { useTheme, Button, Form, Heading, Accordion, AccordionItem } from '@carbon/react'
import { TypeAuthGrant } from '@/APIController/types'
import { PublicScope, EmailScope, ModifyScope } from './scopes'
import { useCookies } from 'next-client-cookies'

const PermissionsForm = (props: { client: TypeAuthGrant, state: string }) => {
  const cookies = useCookies()

  const headingColor = () => {
    const { theme } = useTheme()
    return (theme == "white") ? "black" : "white"
  }

  const onAccept = () => {
    location.replace(`/auth/authorize?response_type=${props.client.response_type}` +
      `&client_id=${props.client.id}&redirect_uri=${encodeURIComponent(props.client.redirect_uri)}` +
      `&state=${props.state}&assertion=${cookies.get('ows-jwt')}`)
  }

  const onReject = () => {
    location.replace("/authorize")
  }

  return (
    <Form className="form">
      <Heading className="heading"
	style={{ marginBottom: "20px", color: headingColor() }}>
	Authorize {props.client.name}
      </Heading>
      <p style={{ marginBottom: "20px" }}>
	You are attempting to sign-in to '{props.client.name}' using your Open Source Club account.
	The client is requesting your permission to access the following:
      </p>
      <Accordion style={{ marginBottom: "20px" }}>
	<AccordionItem title="Client Description" open={true}>
	  <p>{props.client.description}</p>
	</AccordionItem>
	{ (!props.client.scope?.includes("public")) ? null : (<PublicScope />) }
	{ (!props.client.scope?.includes("email")) ? null : (<EmailScope />) }
	{ (!props.client.scope?.includes("modify")) ? null : (<ModifyScope />) }
      </Accordion>
      <p style={{ marginBottom: "20px" }}>Your password will never be shared</p>
      <Button className="perm-button" onClick={onAccept}>
	Accept
	<ArrowRight className="button--arrow" />
      </Button>
      <Button className="perm-button" kind="danger--tertiary" onClick={onReject}>Reject</Button>
    </Form>
  )
}

export default PermissionsForm
