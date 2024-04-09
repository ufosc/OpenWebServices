'use client'

import { ValidateClientURLParams } from '@/API'
import ImageBanner from '@/components/ImageBanner/imagebanner'
import SigninForm from '@/components/SigninForm'
import SignupForm from '@/components/SignupForm'
import { useSearchParams, useRouter } from 'next/navigation'
import { useState } from 'react'
import { useCookies } from 'next-client-cookies'
import ClientError from './cerror'
import { Permissions, ClientDefinition } from './permissions'
import './style.scss'

export default function Page() {
  const router = useRouter()
  const cookies = useCookies()
  const token = cookies.get('ows-access-token')
  const [view, setView] = useState<"signin" | "signup">("signin")

  // Gather client parameters.
  const searchParams = useSearchParams()
  const client = {
    response_type: searchParams.get('response_type'),
    client_id: searchParams.get('client_id'),
    redirect_uri: searchParams.get('redirect_uri'),
    state: searchParams.get('state')
  }

  const renderForm = () => {
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

    if (typeof token === "undefined" && view === "signin") {
      return (<SigninForm setView={setView} />)
    }

    if (typeof token === "undefined" && view === "signup") {
      return (<SignupForm setView={setView} />)
    }

    // No URL parameters and user is already signed in, redirect to dashboard.
    if (typeof token !== "undefined" && client.response_type === null) {
      router.push("/")
      return null
    }

    // Display permissions page.
    return (<Permissions client={client as ClientDefinition} />)
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
