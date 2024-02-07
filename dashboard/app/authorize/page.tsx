'use client'

import { ValidateClientURLParams } from '@/APIController/Validation'
import ImageBanner from '@/components/ImageBanner/imagebanner'
import SigninForm from '@/components/SigninForm'
import SignupForm from '@/components/SignupForm'
import { useSearchParams, useRouter } from 'next/navigation'
import { useState } from 'react'
import { useCookies } from 'next-client-cookies'
import { TypeAuthGrant } from '@/APIController/types'
import ClientError from './cerror'
import Permissions from './permissions'
import './style.scss'

export default function Page() {
  const cookies = useCookies()
  const router = useRouter()
  const token = cookies.get('ows-access-token')
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
