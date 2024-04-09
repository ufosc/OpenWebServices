'use client'

import { useState, useEffect } from 'react'
import { useCookies } from 'next-client-cookies'
import MyAccount from '@/components/MyAccount'
import { GetUser, IsAPISuccess } from '@/API'
import Users from '@/components/Users'
import Clients from '@/components/Clients'
import { useRouter } from 'next/navigation'
import { Heading, Tabs, TabList, Tab,
  TabPanels, TabPanel, Loading } from '@carbon/react'

type User = {
  realms: string[],
}

export default function Page() {
  const router = useRouter()
  const cookies = useCookies()
  const token = cookies.get('ows-access-token')
  const [user, setUser] = useState<User | null>(null)

  useEffect(() => {
    if (typeof token === "undefined") {
      router.push("/authorize")
    }

    if (user === null) {
      GetUser(token as string)
        .then((res) => setUser(res as User))
        .catch((err) => {
          cookies.remove('ows-access-token')
          router.push("/authorize")
        })
    }
  }, [])

  if (user === null) {
    return (<Loading withOverlay={true} />)
  }

  return (
    <div className="account">
      <Heading> Dashboard </Heading>
      <Tabs>
	<TabList contained className="account--tablist" aria-label="dashboard">
	  <Tab>My Account</Tab>
	  <Tab disabled={!user.realms?.includes("clients.read")}>Clients</Tab>
	  <Tab disabled={!user.realms?.includes("users.read")}>Users</Tab>
	</TabList>
	<TabPanels>
	  <TabPanel><MyAccount user={user} /></TabPanel>
	  <TabPanel>{(user.realms?.includes("clients.read")) ? (<Clients />) : null }</TabPanel>
	  <TabPanel>{ (user.realms?.includes("users.read")) ? (<Users />) : null }</TabPanel>
	</TabPanels>
      </Tabs>
    </div>
  )
}
