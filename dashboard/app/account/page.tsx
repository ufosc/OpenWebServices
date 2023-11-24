'use client'

import { redirect } from 'next/navigation'
import { cookies } from 'next/headers'
import { useContext } from 'react'
import { JWTContext } from '@/app/context'

export default function Page() {
  const jwt = useContext(JWTContext)
  if (typeof jwt === undefined) {
    redirect("/authorize")
  }

  return (<p>Account Page: Currently unavailable (work in progress)</p>)
}
