# talkus-listener
An http webhook listener for talkus.io chat transcripts.
Accepts JSON posted by talkus and emails the transcript to your chosen email address.
Support both authenticated and un-authenticated email.

Build with `go build`

Example usage:
`./talkus-listener --emailRecipient recieving@example.com --emailSender sending@example.com --emailServer mail.example.com`

### TODO

* Init scripts
