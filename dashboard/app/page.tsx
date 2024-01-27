'use client'

import { useState } from 'react'
import { useCookies } from 'next-client-cookies'
import MyAccount from '@/components/MyAccount'
import { GetUser, IsAPISuccess } from '@/APIController/API'
import Users from '@/components/Users'
import Clients from '@/components/Clients'
import { Heading, Tabs, TabList, Tab,
  TabPanels, TabPanel, Loading } from '@carbon/react'

type User = {
  realms: string[],
}

export default function Page() {
  const cookies = useCookies()
  const jwt = cookies.get('ows-jwt')
  if (typeof jwt === "undefined" && typeof window !== "undefined") {
    window.location.replace("/authorize")
  }

  const [user, setUser] = useState<User | null>(null)
  if (user === null) {
    GetUser(jwt as string).then((res) => {
        if (!IsAPISuccess(res)) {
          cookies.remove('ows-jwt')
          if (typeof window !== "undefined") {
            window.location.replace("/authorize")
          }
        }
        setUser(res as User)
      }).catch((err) => {
        cookies.remove('ows-jwt')
        if (typeof window !== "undefined") {
          window.location.replace("/authorize")
        }
      })
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
	  <TabPanel><Clients /></TabPanel>
	  <TabPanel><Users /></TabPanel>
	</TabPanels>
      </Tabs>
    </div>
  )
}
