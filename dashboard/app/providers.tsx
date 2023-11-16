'use client'

import { useState } from 'react'
import Header from '@/components/NavHeader/navheader'
import { Content, Theme } from '@carbon/react'

const Providers = ({ children }) => {
  const [theme, setTheme] = useState("white")
  const savedTheme = localStorage.getItem("theme")
  if ((savedTheme === "white" || savedTheme === "g100") && savedTheme != theme) {
    setTheme(savedTheme)
  }

  const updateTheme = (theme) => {
    if (theme === "white") {
      localStorage.setItem("theme", "white")
      setTheme("white")
      return
    }

    if (theme === "g100") {
      localStorage.setItem("theme", "g100")
      setTheme("g100")
      return
    }
  }

  return (
    <Theme theme={theme}>
      <Header setTheme={updateTheme} />
      <Content>{children}</Content>
    </Theme>
  )
}

export default Providers
