import './globals.scss'
import type { Metadata } from 'next'
import { cookies } from 'next/headers'
import Provider from './providers'

export const metadata: Metadata = {
  title: 'UF Open Source Club | OpenWebServices',
  description: 'OpenWebServices Gateway for the UF Open Source Club',
}

export default function RootLayout(props: {children: any}) {
  const jwt = cookies().get("ows-jwt")
  const random = Math.floor(Math.random() * 7)
  return (
    <html lang="en">
      <body>
	<Provider jwt={jwt} random={random}>{props.children}</Provider>
      </body>
    </html>
  )
}
