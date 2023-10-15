# WebSMTP

WebSTMP implements an SMTP client that can be interfaced via HTTP. Its purpose is to provide emailing capabilities to web/public clients that may otherwise not have access to low-level TCP APIs. Since each request can send a variable number of emails (which consequently takes a variable amount of time), a request-reference model is employed: clients submit a request to send emails, and receive a reference that they can use to track the status of their request. This means that clients receive immediate feedback for variable-complexity operations.

A lot of SMTP servers decline addresses without [SPF](https://en.wikipedia.org/wiki/Sender_Policy_Framework), so the source address is always pre-specified by the server. In the future, either add restrictions to what emails can be sent (to prevent abuse of the organization email) or allow all email sources altogether.
