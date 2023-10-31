"use client"

import { Box, Flex, Text, Button } from "@chakra-ui/react"
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
      <Box>
        <Text>Grant Permissions?</Text>
        <Flex>
          <Button onClick={() => window.close()}>Cancel</Button>
          <Button colorScheme="green" onClick={redirectToURI}>Authorize OSC</Button>
        </Flex>
      </Box>
    )
}