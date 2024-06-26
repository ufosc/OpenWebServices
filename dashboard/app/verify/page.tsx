'use client'

import './page.scss'
import ImageBanner from '@/components/ImageBanner/imagebanner'
import { RandContext } from '../context'
import { useContext } from 'react'
import { Heading, Link } from '@carbon/react'

export default function Page() {
  const random = useContext(RandContext)
  return (
    <div className="verifyEmailPage">
      <div className="verifyEmailPage--prompt">
	<Heading>Awaiting Email Verification</Heading>
	<p style={{ marginTop: 15, marginBottom: 15 }}>
	  Please verify your email address to continue.
	  Make sure to check your spam folder.
	  You may safely exit this page.
	</p>
	<Link href="/authorize">Return to Sign in</Link>
      </div>
      <ImageBanner/>
    </div>
  )
}
