'use client'

import { Heading, useTheme } from '@carbon/react'

const AlertBanner = (props: { heading: any, children: any }) => {
  const headingColor = () => {
    const { theme } = useTheme()
    return (theme == "white") ? "black" : "white"
  }

  return (
    <div>
      <Heading className="heading"
	style={{ marginBottom: "20px", color: headingColor() }}>
	{props.heading}
      </Heading>
      {props.children}
    </div>
  )
}

export default AlertBanner
