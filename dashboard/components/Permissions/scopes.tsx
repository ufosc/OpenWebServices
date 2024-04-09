import { Wikis, Email, Edit } from '@carbon/icons-react'
import { AccordionItem } from '@carbon/react'
import './style.scss'

const PublicScopeTitle = () =>
  <p>
    <Wikis className="perm--scope-icon" />
    {" Public Information"}
  </p>

export const PublicScope = () =>
  <AccordionItem title={PublicScopeTitle()}>
    <p style={{ fontSize: 13 }}> Read your first and last name </p>
  </AccordionItem>

const EmailScopeTitle = () =>
  <p>
    <Email className="perm--scope-icon" />
    {" Email Address"}
  </p>

export const EmailScope = () =>
  <AccordionItem title={EmailScopeTitle()}>
    <p style={{ fontSize: 13 }}> Read your email address </p>
  </AccordionItem>

const ModifyScopeTitle = () =>
  <p>
    <Edit className="perm--scope-icon" />
    {" Modify Account"}
  </p>

export const ModifyScope = () =>
  <AccordionItem title={ModifyScopeTitle()}>
    <p style={{ fontSize: 13 }}> Modify your first and last name </p>
  </AccordionItem>
