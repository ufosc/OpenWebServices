'use client'

import './page.scss'
import SigninForm from '@/components/SigninForm/signinform'
import ImageBanner from '@/components/ImageBanner/imagebanner'

const LoginPage = () => {
  return (
    <div className="loginPage">
      <div className="loginPage--form">
	<SigninForm/>
      </div>
      <ImageBanner />
    </div>
  )
}

export default LoginPage
