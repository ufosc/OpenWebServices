import AlertBanner from '@/components/AlertBanner'
import { Link, Accordion, AccordionItem } from '@carbon/react'

const ClientError = (props: { children: any }) => {
  return (
    <AlertBanner heading="Error: Invalid Client Authorization Request">
      <p>
	Sign-in could not be completed at this time, most likely due to an invalid
	client request. If you're the client's developer, please ensure that
	you've implemented the authorization request as specified in the OAuth2
	documentation:
      </p>
      <br />
      <Link href="https://datatracker.ietf.org/doc/html/rfc6749">
	IETF RFC 6749: The OAuth 2.0 Authorization Framework
      </Link>
      <div style={{ marginTop: "20px" }}>
	<Accordion>
          <AccordionItem title="Additional Debug Information">
            {props.children}
          </AccordionItem>
        </Accordion>
      </div>
    </AlertBanner>
  )
}

export default ClientError
