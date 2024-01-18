# WebSMTP
[![Go Reference](https://pkg.go.dev/badge/github.com/ufosc/OpenWebServices/websmtp.svg)](https://pkg.go.dev/github.com/ufosc/OpenWebServices/websmtp)

WebSTMP implements an SMTP client that can be interfaced via HTTP. Its purpose is to provide emailing capabilities to web/public clients that may otherwise not have access to low-level TCP APIs. Since each request can send a variable number of emails (which consequently takes a variable amount of time), a request-reference model is employed: clients submit a request to send emails, and receive a reference that they can use to track the status of their request. This means that clients receive immediate feedback for variable-complexity operations.

## Dependencies
1. [GoLang](https://go.dev/doc/install).

## Install
Clone the parent repository:
```bash
git clone https://github.com/ufosc/OpenWebServices.git
```

Navigate to the project directory:
```bash
cd OpenWebServices/websmtp
```

Install Go mod dependencies:
```
go get
```

## Usage
Within the project directory (`OpenWebServices/websmtp`), start the server:
```bash
go run ./...
```

### Submitting mailer jobs
To submit a job to the mailer, send a POST request to `http://localhost:8080/mail/send` with the following JSON body:
```json
{
    "from": "<source>@<domain>.com",
    "to": ["example@example.com", "foo@bar.com"],
    "subject": "Hello World",
    "body": "hello, how are you sir?"
}
```

You should receive a response similar to:
```json
{"ref":"c1de3135-fbaf-4319-9382-bc23362deac2"}
```

`ref` is a reference number which you can use to periodically query the status of the mailer job.

### Querying mailer job status
To query the status of your job, send a GET request to `http://localhost:8080/mail/status/:REF`, where `:REF` is the relevant reference number. If the job has succeeded, you should eventually receive a response that looks as follows:
```json
{
    "id": "f683a80c-ba2d-4285-8e1d-4aeaeb37a729",
    "status": "completed",
    "failed": [],
    "time_completed": 1705589319
}
```

If the `:REF` is invalid, unknown to the server, or simply hasn't started, the response will look as follows:
```json
{
    "id": "8e766178-bd19-4387-aefa-c9756585e6b6",
    "status": "not started",
    "failed": [],
    "time_completed": 0
}
```

The server uses an in-memory cache to keep track of job requests, which is periodically wiped. Pending jobs will remain in queue even if their corresponding cache entry has been wiped, meaning that the status endpoint has no knowledge of whether an arbitrary entry exists. If you know your reference ID, then the status response is accurate; otherwise, it is not.

## Configure
The server is configured by the following environment variables:
```go
type Config struct {
	GIN_MODE string // "release" or "debug".
	PORT     string // server port.
	THREADS  string // number of worker threads.
}
```

The environment variables are pulled from the OS. Additionally, you may create a `.env` file in the current working directory for conveniently specifying options. The `.env` file format is as follows:
```
GIN_MODE="debug"
PORT="8080"
THREADS="1"
<KEY>=<VALUE>
```

## Security (None!!)
Until OAuth2 is finalized, the websmtp server is <b>not intended for public use</b>. Routes are unauthorized, so the server is vulnerable to abuse. For the time being, it is recommended to use websmtp behind a firewall - this defeats the purpose of the project, but it only a temporary measure.

## Sender Policy Framework

A lot of SMTP servers will decline emails from source addresses that do not comply with [Sender Policy Framework](https://en.wikipedia.org/wiki/Sender_Policy_Framework). This is a receiver-side filtering mechanism; as far as the websmtp server is concerned, the email has been delivered succesfully.

The only way to resolve this issue is to update the DNS records assosciated with your desired sending addresses.
