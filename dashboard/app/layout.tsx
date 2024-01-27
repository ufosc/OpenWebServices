import './globals.scss'
import type { Metadata } from 'next'
import Provider from './providers'
import { CookiesProvider } from 'next-client-cookies/server'

export const metadata: Metadata = {
  title: 'UF Open Source Club | OpenWebServices',
  description: 'OpenWebServices Gateway for the UF Open Source Club',
}

export default function RootLayout(props: {children: any}) {
  const random = Math.floor(Math.random() * 7)
  return (
    <html lang="en">
      <body>
	<CookiesProvider>
	  <Provider random={random}>{props.children}</Provider>
	</CookiesProvider>
      </body>
    </html>
  )
}
