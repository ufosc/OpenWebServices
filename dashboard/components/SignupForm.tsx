'use client'

import { useState } from 'react'
import { ArrowRight } from '@carbon/icons-react'
import { useTheme, Form, Button, TextInput, Heading } from '@carbon/react'
import { PostSignup, IsAPIFailure, IsAPISuccess } from '@/APIController/API'
import { ValidateEmail } from '@/APIController/Validation'

const SignupForm = (props: { setView: Function }) => {
  const headingColor = () => {
    const { theme } = useTheme()
    return (theme == "white") ? "black" : "white"
  }

  const [hasError, setHasError] = useState("")
  const [form, setForm] = useState({
    first_name: "", last_name: "", email: "", password: "", verif: "",
  })

  const submitForm = (e : any) => {
    e.preventDefault()

    // Validate email address.
    if (!ValidateEmail(form.email)) {
      setHasError("Email address must be a valid @ufl.edu address")
      return
    }

    // Ensure password and verification match each other.
    if (form.password != form.verif) {
      setHasError("Passwords do not match")
      return
    }

    // Send post request.
    PostSignup({
      first_name: form.first_name,
      last_name: form.last_name,
      email: form.email,
      password: form.password,
      captcha: "123"
    }).then((res) => {
      if (IsAPISuccess(res)) {
	location.replace("/verify")
	return
      }

      let msg = (IsAPIFailure(res) && typeof res.error != "undefined") ?
	res.error : "An unknown error has occured. Please try again later."

      setHasError(msg)
    }).catch((err) => {
      setHasError("Server could not be reached. Please try again later")
    })
  }

  return (
    <Form className="form">
	  <Heading className="heading"
	    style={{ marginBottom: "20px", color: headingColor() }}>
	    Create an OSC Account
	  </Heading>
	  <TextInput
	    id="first_name"
	    style={{ marginBottom: "15px" }}
	    labelText="First Name"
            placeholder="Alberta"
            onChange={ (e) => setForm({
              first_name: e.target.value, last_name: form.last_name,
	      email: form.email, password: form.password, verif: form.verif })
	    }
	  />
	  <TextInput
	    id="last_name"
	    style={{ marginBottom: "15px" }}
	    labelText="Last Name"
            placeholder="Gator"
            onChange={ (e) => setForm({
	      first_name: form.first_name, last_name: e.target.value,
	      email: form.email, password: form.password, verif: form.verif })
	    }
	  />
	  <TextInput
	    id="email"
	    style={{ marginBottom: "15px" }}
	    labelText="Email Address"
            placeholder="gator@ufl.edu"
            onChange={ (e) => setForm({
	      first_name: form.first_name, last_name: form.last_name,
	      email: e.target.value, password: form.password, verif: form.verif })
	    }
	  />
	  <TextInput
	    style={{ marginBottom: "15px" }}
	    type="password"
	    labelText="Password"
	    id="password"
	    placeholder="************"
            required pattern="(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{6,}"
            onChange={ (e) => setForm({
	      first_name: form.first_name, last_name: form.last_name,
	      email: form.email, password: e.target.value, verif: form.verif })
	    }
	  />
	  <TextInput
	    style={{ marginBottom: "15px" }}
	    type="password"
	    labelText="Verify Password"
	    id="password"
	    placeholder="************"
            required pattern="(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{6,}"
            onChange={ (e) => setForm({
	      first_name: form.first_name, last_name: form.last_name,
	      email: form.email, password: form.password, verif: e.target.value })
	    }
	  />
	  <Button type="submit" className="signinform--button" onClick={submitForm}>
	    Sign Up
	    <ArrowRight className="button--arrow"/>
	  </Button>
	  {
	    (hasError != "") ? (
	      <p style={{ marginTop: 10, marginBottom: 5, color: 'red' }}>
                Error: { hasError }
	      </p>
	    ) : null
	  }
	  <hr style={{ marginTop: 30, marginBottom: 15 }}/>
	  <p style={{ color: "gray", marginBottom: 15, fontSize: 14 }}>
	    Already have an account?
	  </p>
	  <Button className="signinform--button" kind="tertiary"
	    onClick={() => { props.setView("signin") }} >
	    Sign in
	    <ArrowRight className="button--arrow"/>
	  </Button>
	</Form>
  )
}

export default SignupForm
