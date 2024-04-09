'use client'

import { ArrowRight } from '@carbon/icons-react'
import { SignIn, ValidateEmail } from '@/API'
import { useTheme, Button, Form, Heading, TextInput, Link } from '@carbon/react'
import { useState } from 'react'
import { useCookies } from 'next-client-cookies'
import { useRouter } from 'next/navigation'

const SigninForm = (props: { setView: Function }) => {
  const router = useRouter()
  const cookies = useCookies()
  const headingColor = () => {
    const { theme } = useTheme()
    return (theme == "white") ? "black" : "white"
  }

  const [hasError, setHasError] = useState("")
  const [form, setForm] = useState({ email: "", password: "" })
  const submitForm = async (e : any) => {
    e.preventDefault()

    // Validate email address.
    if (!ValidateEmail(form.email)) {
      setHasError("Email address is invalid")
      return
    }

    // Make API call.
    SignIn(form).then(_res => {
      let res = _res as { token: string }
      if (typeof res.token !== "undefined") {
	cookies.set('ows-access-token', res.token)
      }
      router.refresh()
    }).catch((err) => {
      setHasError(err.error_description)
    })

  }

  return (
    <Form className="form">
      <Heading className="heading"
	style={{ marginBottom: "20px", color: headingColor() }}>
	Sign in to Open Source Club
      </Heading>
      <TextInput
	id="email"
	style={{ marginBottom: "15px" }}
	placeholder="gator@ufl.edu"
	labelText="Email Address"
	onChange={(e) => setForm({ email: e.target.value, password: form.password })}/>
      <TextInput
	style={{ marginBottom: "15px" }}
	type="password"
	labelText="Password"
	id="password"
	placeholder="************"
	required pattern="(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{12,}"
	onChange={(e) => setForm({ password: e.target.value, email: form.email })}/>
      <Button type="submit" className="signinform--button" onClick={submitForm}>
	Continue
	<ArrowRight className="button--arrow" />
      </Button>
      {
	(hasError != "") ? (
	  <p style={{ marginTop: 10, marginBottom: 5, color: 'red' }}>
	    Error: { hasError }
	  </p>) : null
      }
      <Link style={{ fontSize: 13 }} href="/reset"> Forgot Password? </Link>
      <hr style={{ marginTop: 30, marginBottom: 15 }} />
      <p style={{ color: "gray", marginBottom: 15, fontSize: 14 }}>
	Don't have an account?
      </p>
      <Button className="signinform--button" kind="tertiary"
	onClick={() => { props.setView("signup") }}>
	Create an account
	<ArrowRight className="button--arrow" />
      </Button>
    </Form>
  )
}

export default SigninForm
