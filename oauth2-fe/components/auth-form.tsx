"use client"

import { FormControl, FormLabel, Input, FormErrorMessage, Button, Box, Heading, Flex, Text, Link, Alert, AlertIcon } from "@chakra-ui/react"
import { useState, ChangeEvent } from "react"
import { ArrowForwardIcon } from '@chakra-ui/icons'
import axios from "axios"
import { redirect, useSearchParams } from "next/navigation"

//http://localhost:3000/signin?response_type=code&client_id=s6BhdRkqt3&state=xyz&redirect_uri=https%3A%2F%2Fclient.example.com%2Fcb

export const AuthForm = ({variant}: {variant: string}) => {
    const [firstNameInput, setFirstNameInput] = useState('')
    const [lastNameInput, setLastNameInput] = useState('')
    const [usernameInput, setUsernameInput] = useState('')
    const [passwordInput, setPasswordInput] = useState('')

    const [submissionSuccess, setSubmissionSuccess] = useState('')
    const [submissionError, setSubmissionError] = useState('')

    const searchParams = useSearchParams();
    const response_type = searchParams.get('response_type')
    const client_id = searchParams.get('client_id')
    const redirect_uri = searchParams.get('redirect_uri')
    const state = searchParams.get('state')
    console.log(response_type, client_id, redirect_uri, state)

    const handleFirstNameInputChange = (e : ChangeEvent) => {
      let value = (e.target as HTMLInputElement).value;
      setFirstNameInput(value);  
    }
  
    const handleLastNameInputChange = (e : ChangeEvent) => {
      let value = (e.target as HTMLInputElement).value;
      setLastNameInput(value);  
    }

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

    const SignInPostRequest = () => {
      axios.post('/auth/signin', {username: usernameInput, password: passwordInput})
        .then(response => {
          console.log(response)
          redirect(`/signin/permissions/?response_type=${response_type}&client_id=${client_id}&state=${state}&redirect_uri=${redirect_uri}`)
        })
        .catch(error => {
          console.log(`/signin/permissions/?response_type=${response_type}&client_id=${client_id}&state=${state}&redirect_uri=${redirect_uri}`)
          console.error(error.message)
          setSubmissionError(error.message)
        })
    }

    const SignUpPostRequest = () => {
      axios.post('/auth/signup', {first_name: firstNameInput, last_name: lastNameInput, username: usernameInput, password: passwordInput})
        .then(response => {
          console.log(response)
          setSubmissionSuccess('Success. Please log in.')
          setSubmissionError('')
        })
        .catch(error => {
          console.error(error.message)
          setSubmissionError(error.message)
          setSubmissionSuccess('');
        })
    }
  
    const isError = !isValidEmail(usernameInput) && usernameInput !== ''
  
    return (
        <Box w={'100%'} h={'100%'} display={'flex'} justifyContent={'center'} alignItems={'center'} >
            <Box w={{base: '80%', md: '400px'}} bg={'blackAlpha.50'} p={3}>
                <Heading variant={'h5'} mb={5}>{variant === 'signup' ? 'Join OSC' : 'Sign in'}</Heading>
                <FormControl isInvalid={isError} isRequired  w={'100%'} h={'100%'}>
                    {variant === 'signin' ? (<></>) : 
                    <>
                      <FormLabel mt={3}>First Name</FormLabel>
                      <Input type="first-name" value={firstNameInput} onChange={handleFirstNameInputChange} />
                      <FormLabel mt={3} >Last Name</FormLabel>
                      <Input mb={3} type="last-name" value={lastNameInput} onChange={handleLastNameInputChange} />
                    </>}
                    <FormLabel>Username</FormLabel>
                    <Input type='username' value={usernameInput} onChange={handleUsernameInputChange} />
                    {!isError ? (
                    <></>
                    ) : (
                    <FormErrorMessage>Username is not a valid email.</FormErrorMessage>
                    )}
                    <Flex justifyContent={'space-between'} alignItems={'center'}>
                      <FormLabel mt={3} w={'fit-content'}>Password</FormLabel>
                      {variant === 'signup' ? <></> : <Text color={'blue'}>Forgot password?</Text>}
                    </Flex>
                    <Input type='password' value={passwordInput} onChange={handlePasswordInputChange} />
                    <Flex justifyContent={'flex-end'}>
                      <Button
                        mt={3}
                        type='submit'
                        colorScheme="blue"
                        rightIcon={<ArrowForwardIcon />} 
                        onClick={variant === 'signup' ? SignUpPostRequest : SignInPostRequest}
                        >
                        Submit
                      </Button>
                    </Flex>
                    {variant === 'signup' ? <></> : 
                      <Flex justifyContent={'center'} mt={2}>
                        <Text mr={1}>New to OSC?</Text>
                        <Link href={'/signup'} color="#3182CE">Create an account</Link>
                      </Flex>}
                </FormControl>
                {!submissionSuccess ? (<></>) : (
                  <Alert status='success' mt={2}>
                    <AlertIcon />{submissionSuccess}
                  </Alert>
                )}
                {!submissionError ? (<></>) : (
                  <Alert status='error' mt={2}>
                    <AlertIcon />{submissionError}
                  </Alert>
                )}
            </Box>
        </Box>
    )
}
  