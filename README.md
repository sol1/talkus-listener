# talkus-listener
An http webhook listener for talkus.io chat transcripts.
Accepts JSON posted by talkus and emails the transcript to your chosen email address.
Support both authenticated and un-authenticated email.

Build with `go build`

Example usage:
`./talkus-listener -config configEmail.json`

To run on startup (Requires systemd):
* `mv talkus-listener /etc/systemd/system/`
* `systemctl enable talkus-listener.service`
* `systemctl start talkus-listener.service`
