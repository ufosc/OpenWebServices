'use client'

import { PostSignout } from '@/APIController/API'
import { JWTContext } from '@/app/context'
import { useContext } from 'react'
import { UserAvatar, Asleep, Awake, Logout } from '@carbon/icons-react'
import { Header, HeaderContainer, HeaderName, HeaderGlobalBar, HeaderGlobalAction, SkipToContent, useTheme } from '@carbon/react'

const NavHeader = (props : { setTheme: Function }) => {
  const jwt = useContext(JWTContext)

  const themeSelector = () => {
    const { theme } = useTheme()
    if (props.setTheme === undefined) return null;

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
    PostSignout().then(() => {
      location.reload()
    })
  }

  return (
    <HeaderContainer render={({ isSideNavExpanded, onClickSideNavExpand }) => (
      <Header aria-label="OpenWebServices">
        <SkipToContent />
        <HeaderName href="/" prefix="UF">OpenWebServices</HeaderName>
	<p>v0.1.0-alpha</p>
	<HeaderGlobalBar>
	  { themeSelector() }
	  <HeaderGlobalAction aria-label="Account" tooltipAlignment="center">
	    <UserAvatar size={20} onClick={() => location.replace("/account")}/>
	  </HeaderGlobalAction>
	  {
	    (typeof jwt !== "undefined") ? (
	      <HeaderGlobalAction aria-label="Logout" tooltipAlignment="center"
		onClick={onSignout}>
	        <Logout size={20} />
	      </HeaderGlobalAction>
	    ) : null
	  }
	</HeaderGlobalBar>
      </Header>
    )}
    />
  )
}

export default NavHeader
