'use client'

import { useContext } from 'react'
import { UserAvatar, Asleep, Awake, Logout } from '@carbon/icons-react'
import { useCookies } from 'next-client-cookies'

import { Header, HeaderContainer, HeaderName, HeaderGlobalBar,
  HeaderGlobalAction, SkipToContent, useTheme } from '@carbon/react'

const NavHeader = (props : { setTheme: Function }) => {
  const cookies = useCookies()
  const jwt = cookies.get('ows-jwt')

  const themeSelector = () => {
    const { theme } = useTheme()
    return (theme == "white") ? (
      <HeaderGlobalAction aria-label="Theme Selector"
	onClick={() => props.setTheme("g100")}>
	    <Awake size={20} />
      </HeaderGlobalAction>
    ) : (
      <HeaderGlobalAction aria-label="Theme Selector"
        onClick={()=> props.setTheme("white")} >
	  <Asleep size={20} />
      </HeaderGlobalAction>
    )
  }

  const onSignout = () => {
    cookies.remove('ows-jwt')
    if (typeof window !== "undefined") {
      window.location.replace("/")
    }
  }

  return (
    <HeaderContainer render={({ isSideNavExpanded, onClickSideNavExpand }) => (
      <Header aria-label="OpenWebServices">
        <SkipToContent />
        <HeaderName href="/" prefix="UF">OpenWebServices</HeaderName>
	<p>v0.1.0-alpha</p>
	<HeaderGlobalBar>
	  { themeSelector() }
	  {
	    (typeof jwt !== "undefined") ? (
	      <>
		<HeaderGlobalAction aria-label="Logout" tooltipAlignment="center"
		  onClick={onSignout}>
		  <Logout size={20} />
		</HeaderGlobalAction>
	      </>
	    ) : null
	  }
	</HeaderGlobalBar>
      </Header>
    )}
    />
  )
}

export default NavHeader
