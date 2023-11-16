'use client'

import './signinform.scss'

import { ArrowRight, Information } from '@carbon/icons-react'

import {
  Form,
  TextInput,
  Heading,
  Button,
  Checkbox,
  Link,
  Tooltip,
  useTheme
} from '@carbon/react'

const SigninForm = () => {
  const headingColor = () => {
    const { theme } = useTheme()
    return (theme == "white") ? "black" : "white"
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
	labelText="Email Address"
	placeholder="gator@ufl.edu"
      />
      <TextInput
	style={{ marginBottom: "15px" }}
	type="password"
	labelText="Password"
	id="password"
	placeholder="************"
	required pattern="(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{6,}"/>
      <Button className="signinform--button">
	Continue
	<ArrowRight className="button--arrow"/>
      </Button>
      <Link style={{ fontSize: 13 }} href="/reset"> Forgot Password? </Link>
      <hr style={{ marginTop: 30, marginBottom: 15 }}/>
      <p style={{ color: "gray", marginBottom: 15, fontSize: 14 }}>
	Don't have an account?
      </p>
      <Button className="signinform--button" kind="tertiary" href="/signup">
	Create an account
	<ArrowRight className="button--arrow"/>
      </Button>
    </Form>
  )
}

export default SigninForm
