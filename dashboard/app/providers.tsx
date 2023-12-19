'use client'

import { useState } from 'react'
import { Content, Theme } from '@carbon/react'
import Header from '@/components/NavHeader'
import { RandContext } from './context'

const Provider = (props: { children : any; random: number }) => {
  const [theme, setTheme] = useState("white")

  const savedTheme = localStorage.getItem("theme")
  if ((savedTheme === "white" || savedTheme === "g100") && savedTheme != theme) {
    setTheme(savedTheme)
  }

  const updateTheme = (newTheme : string) => {
    if ((newTheme === "white" || newTheme === "g100") && newTheme != theme) {
      localStorage.setItem("theme", newTheme)
      setTheme(newTheme)
    }
  }

  return (
    <Theme theme={theme}>
      <RandContext.Provider value={props.random}>
	<Header setTheme={updateTheme} />
	<Content>{props.children}</Content>
      </RandContext.Provider>
    </Theme>
  )
}

export default Provider
