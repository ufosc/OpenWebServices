import {
  Notification,
  UserAvatar,
  Asleep,
  Awake
} from '@carbon/icons-react'

import {
  Header,
  HeaderContainer,
  HeaderName,
  HeaderGlobalBar,
  HeaderGlobalAction,
  SkipToContent,
  HeaderNavigation,
  HeaderMenuItem,
  HeaderMenuButton,
  SideNav,
  SideNavItems,
  HeaderSideNavItems,
  Link,
  useTheme
} from '@carbon/react'

const NavHeader = ({ setTheme }) => {
  const themeSelector = () => {
    const { theme } = useTheme()

    if (setTheme === undefined) return null;

    return (theme == "white") ? (
      <HeaderGlobalAction aria-label="Theme Selector"
	onClick={() => setTheme("g100")}>
	    <Asleep size={20} />
      </HeaderGlobalAction>
    ) : (
      <HeaderGlobalAction aria-label="Theme Selector"
        onClick={()=> setTheme("white")} >
	  <Awake size={20} />
      </HeaderGlobalAction>
    )
  }

  return (
    <HeaderContainer render={({ isSideNavExpanded, onClickSideNavExpand }) => (
      <Header aria-label="OpenWebServices">
        <SkipToContent />
	<HeaderMenuButton
	  aria-label="Open menu"
	  onClick={onClickSideNavExpand}
	  isActive={isSideNavExpanded}
	/>
        <HeaderName href="/" prefix="UF OSC">OpenWebServices</HeaderName>
	<HeaderNavigation aria-label="OpenWebServices">
          <HeaderMenuItem href="/">Dashboard</HeaderMenuItem>
	  <HeaderMenuItem href="https://github.com/ufosc">Github</HeaderMenuItem>
	  <HeaderMenuItem href="https://docs.ufosc.org/">Docs</HeaderMenuItem>
        </HeaderNavigation>
	<SideNav
          aria-label="Side navigation"
          expanded={isSideNavExpanded}
          isPersistent={false}
        >
	  <HeaderSideNavItems>
          <HeaderMenuItem href="/">Dashboard</HeaderMenuItem>
  	  <HeaderMenuItem href="https://github.com/ufosc">Github</HeaderMenuItem>
	  <HeaderMenuItem href="https://docs.ufosc.org/">Docs</HeaderMenuItem>
        </HeaderSideNavItems>
        </SideNav>
	<HeaderGlobalBar>
	  { themeSelector() }
	  <HeaderGlobalAction aria-label="Account" tooltipAlignment="center">
	    <UserAvatar size={20} />
	  </HeaderGlobalAction>
	</HeaderGlobalBar>
      </Header>
    )}
    />
  )
}

export default NavHeader
