'use client'

import './page.scss'
import ImageBanner from '@/components/ImageBanner/imagebanner'
import { Heading, Link } from '@carbon/react'

const VerifyEmailPage = () => {
  return (
    <div className="verifyEmailPage">
      <div className="verifyEmailPage--prompt">
	<Heading>Awaiting Email Verification</Heading>
	<p style={{ marginTop: 15, marginBottom: 15 }}>
	  Please verify your email address to continue.
	  Make sure to check your spam folder.
	  You may safely exit this page.
	</p>
	<Link href="/signin">Return to Sign in Page</Link>
      </div>
      <ImageBanner />
    </div>
  )
}

export default VerifyEmailPage
