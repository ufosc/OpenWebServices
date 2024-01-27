'use client'

import { useState, useEffect } from 'react'
import { Content, Theme } from '@carbon/react'
import Header from '@/components/NavHeader'
import { RandContext } from './context'

const Provider = (props: { children : any; random: number }) => {
  const [theme, setTheme] = useState<"white" | "g100">("white")
  const savedTheme = (typeof window !== "undefined") ?
    localStorage.getItem("theme") : theme

  const updateTheme = (newTheme : string) => {
    if ((newTheme === "white" || newTheme === "g100") && newTheme != theme) {
      setTheme(newTheme)
    }
  }

  useEffect(() => {
    if ((savedTheme === "white" || savedTheme === "g100") &&
      savedTheme != theme) {
        setTheme(savedTheme)
      }
  }, [])

  useEffect(() => {
    if (typeof window !== "undefined") {
      localStorage.setItem("theme", theme)
    }
  }, [theme])

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
