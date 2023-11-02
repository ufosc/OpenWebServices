"use client"

import { Box, Flex, Text, Button, Heading, Link } from "@chakra-ui/react"
import { useRouter, useSearchParams } from "next/navigation"

export default function PermissionsPage() {
    const searchParams = useSearchParams();
    const redirect_uri = searchParams.get('redirect_uri')
    const router = useRouter();

    const redirectToURI = () => {
      if(redirect_uri) {
        console.log(redirect_uri)
        router.push(redirect_uri)
      }
    }
  
    return (
      <Box w={'100%'} h={'100%'} display={'flex'} justifyContent={'center'} alignItems={'center'} >
        <Box w={{base: '80%', md: '400px'}} bg={'blackAlpha.50'} p={3}>
          <Heading variant={'h6'} mb={5}>OSC would like permission to:</Heading>
          <Text>Verify your OSC identity and view necessary resources.</Text>
          <Text w={'fit-content'} borderBottom={'1px solid black'}>Resources</Text>
          <Text>- First and last name</Text>
          <Text>- Email addresses</Text>
          <Flex>
            <Button onClick={() => window.close()} w={'50%'}>Cancel</Button>
            <Button colorScheme="green" onClick={redirectToURI} w={'50%'}>Authorize OSC</Button>
          </Flex>
          <Text w={'fit-content'} mt={3} textAlign={'center'}>Authorizing will redirect to</Text>
          <Link color="blue">{redirect_uri}</Link>
        </Box>
      </Box>
    )
}