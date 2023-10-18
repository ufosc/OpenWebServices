"use client"

import { FormControl, FormLabel, Input, FormErrorMessage, Button, Box, Heading, Flex } from "@chakra-ui/react"
import { useState, ChangeEvent } from "react"
import { ArrowForwardIcon } from '@chakra-ui/icons'

export const AuthForm = ({variant}: {variant: string}) => {
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
  
    const isError = !isValidEmail(usernameInput) && usernameInput !== ''
  
    return (
        <Box w={'100%'} h={'100%'} display={'flex'} justifyContent={'center'} alignItems={'center'} >
            <Box w={{base: '80%', md: '400px'}} bg={'blackAlpha.50'} p={3}>
                <Heading variant={'h5'} mb={5}>{variant === 'signup' ? 'Sign up' : 'Log in'}</Heading>
                <FormControl isInvalid={isError} isRequired  w={'100%'} h={'100%'}>
                    <FormLabel>Username</FormLabel>
                    <Input type='username' value={usernameInput} onChange={handleUsernameInputChange} />
                    {!isError ? (
                    <></>
                    ) : (
                    <FormErrorMessage>Username is not a valid email.</FormErrorMessage>
                    )}
                    <FormLabel mt={3}>Password</FormLabel>
                    <Input type='password' value={passwordInput} onChange={handlePasswordInputChange} />
                    <Flex justifyContent={'flex-end'}>
                        <Button
                        mt={3}
                        type='submit'
                        colorScheme="blue"
                        rightIcon={<ArrowForwardIcon />} 
                        >
                        Submit
                        </Button>
                    </Flex>
                </FormControl>
            </Box>
        </Box>
    )
}
  