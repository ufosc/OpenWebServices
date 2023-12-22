'use client'

import { useState } from 'react'
import { redirect } from 'next/navigation'
import { useCookies } from 'next-client-cookies'
import MyAccount from '@/components/MyAccount'
import { GetUser, IsAPISuccess } from '@/APIController/API'
import Users from '@/components/Users'
import { Heading, Tabs, TabList, Tab,
  TabPanels, TabPanel, Loading } from '@carbon/react'

export default function Page() {
  const cookies = useCookies()
  const jwt = cookies.get('ows-jwt')
  if (typeof jwt === "undefined") {
    redirect("/authorize")
  }

  const [user, setUser] = useState<Object | null>(null)
  if (user === null) {
    GetUser(jwt).then((res) => {
      if (!IsAPISuccess(res)) {
	cookies.remove('ows-jwt')
	location.replace("/authorize")
	return
      }
      setUser(res)
    }).catch((err) => {
      cookies.remove('ows-jwt')
      location.replace("/authorize")
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
	  <Tab disabled>Analytics</Tab>
	  <Tab disabled>DNS</Tab>
	  <Tab disabled>CDN</Tab>
	</TabList>
	<TabPanels>
	  <TabPanel><MyAccount user={user} /></TabPanel>
	  <TabPanel> Tab Panel 1 </TabPanel>
	  <TabPanel> Tab Panel 2 </TabPanel>
	</TabPanels>
      </Tabs>
    </div>
  )
}
