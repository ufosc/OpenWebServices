"use client"

import { FormControl, FormLabel, Input, FormHelperText, FormErrorMessage, Button } from "@chakra-ui/react"
import { useState, ChangeEvent } from "react"

export default function OSCLogin() {
  const [usernameInput, setUsernameInput] = useState('')
  const [passwordInput, setPasswordInput] = useState('')

  const handleUsernameInputChange = (e : ChangeEvent) => {
    let value = (e.target as HTMLInputElement).value;
    setUsernameInput(value);  
  }

  const handlePasswordInputChange = (e : ChangeEvent) => {
    let value = (e.target as HTMLInputElement).value;
    setPasswordInput(value);  
  }

  const isValidEmail = (email : string) => {
    const regex = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$/;
    return regex.test(email);
  };

  const isError = !isValidEmail(usernameInput)

  return (
    <FormControl isInvalid={isError} isRequired>
      <FormLabel>Username</FormLabel>
      <Input type='username' value={usernameInput} onChange={handleUsernameInputChange} />
      {!isError ? (
        <FormHelperText>
          Enter the username.
        </FormHelperText>
      ) : (
        <FormErrorMessage>Username is not a valid email.</FormErrorMessage>
      )}
      <FormLabel>Password</FormLabel>
      <Input type='password' value={passwordInput} onChange={handlePasswordInputChange} />
      <FormHelperText>
        Enter the password.
      </FormHelperText>
      <Button
        type='submit'
      >
        Submit
      </Button>
    </FormControl>
  )

}
