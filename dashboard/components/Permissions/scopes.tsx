import { Wikis, Email, Edit } from '@carbon/icons-react'
import { AccordionItem } from '@carbon/react'
import './style.scss'

const PublicScopeTitle = () => {
  return (
    <p>
      <Wikis className="perm--scope-icon" />
      {" Public Information"}
    </p>
  )
}

export const PublicScope = () => {
  return (
    <AccordionItem title={PublicScopeTitle()}>
      <p style={{ fontSize: 13 }}> Read your first and last name </p>
    </AccordionItem>
  )
}

const EmailScopeTitle = () => {
  return (
    <p>
      <Email className="perm--scope-icon" />
      {" Email Address"}
    </p>
  )
}

export const EmailScope = () => {
  return (
    <AccordionItem title={EmailScopeTitle()}>
      <p style={{ fontSize: 13 }}> Read your email address </p>
    </AccordionItem>
  )
}

const ModifyScopeTitle = () => {
  return (
    <p>
      <Edit className="perm--scope-icon" />
      {" Modify Account"}
    </p>
  )
}

export const ModifyScope = () => {
  return (
    <AccordionItem title={ModifyScopeTitle()}>
      <p style={{ fontSize: 13 }}> Modify your first and last name </p>
    </AccordionItem>
  )
}
