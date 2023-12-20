'use client'

import { redirect } from 'next/navigation'
import { useCookies } from 'next-client-cookies'
import MyAccount from '@/components/MyAccount/MyAccount'
import { Heading, Tabs, TabList, Tab, TabPanels, TabPanel } from '@carbon/react'

export default function Page() {
  const cookies = useCookies()
  const jwt = cookies.get('ows-jwt')
  if (typeof jwt === "undefined") {
    redirect("/authorize")
  }

  return (
    <div className="account">
      <Heading> Dashboard </Heading>
      <Tabs>
	<TabList contained className="account--tablist">
	  <Tab>Analytics</Tab>
	  <Tab>Clients</Tab>
	  <Tab>Users</Tab>
	  <Tab disabled>DNS</Tab>
	  <Tab disabled>CDN</Tab>
	  <Tab>My Account</Tab>
	</TabList>
	<TabPanels>
	  <TabPanel> Tab Panel 1 </TabPanel>
	  <TabPanel> Tab Panel 2 </TabPanel>
	  <TabPanel> Tab Panel 3 </TabPanel>
	  <TabPanel> Tab Panel 4 </TabPanel>
	  <TabPanel> Tab Panel 5 </TabPanel>
	  <TabPanel><MyAccount /></TabPanel>
	</TabPanels>
      </Tabs>
    </div>
  )
}
