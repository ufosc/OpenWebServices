'use client'

import { useContext } from 'react'
import { UserAvatar, Asleep, Awake, Logout } from '@carbon/icons-react'
import { useCookies } from 'next-client-cookies'
import { useRouter } from 'next/navigation'
import { VERSION } from '@/config'

import { Header, HeaderContainer, HeaderName, HeaderGlobalBar,
  HeaderGlobalAction, SkipToContent, useTheme } from '@carbon/react'

const NavHeader = (props : { setTheme: Function }) => {
  const router = useRouter()
  const cookies = useCookies()
  const token = cookies.get('ows-access-token')

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
    cookies.remove('ows-access-token')
    router.push("/authorize")
  }

  return (
    <HeaderContainer render={({ isSideNavExpanded, onClickSideNavExpand }) => (
      <Header aria-label="OpenWebServices">
        <SkipToContent />
        <HeaderName href="/" prefix="UF">OpenWebServices</HeaderName>
	<p>{VERSION}</p>
	<HeaderGlobalBar>
	  { themeSelector() }
	  {
	    (typeof token !== "undefined") ? (
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
