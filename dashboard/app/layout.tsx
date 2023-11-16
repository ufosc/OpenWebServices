import './globals.scss'
import type { Metadata } from 'next'
import Providers from './providers'

export const metadata: Metadata = {
  title: 'UF Open Source Club | OpenWebServices',
  description: 'OpenWebServices Gateway for the UF Open Source Club',
}

export default function RootLayout({children}) {
  return (
    <html lang="en">
      <body>
	<Providers>{children}</Providers>
      </body>
    </html>
  )
}
