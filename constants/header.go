package constants

import "net/textproto"

var (
	XServiceName  = textproto.CanonicalMIMEHeaderKey("x-Service-Name")
	XApiKey       = textproto.CanonicalMIMEHeaderKey("x-Api-Key")
	XRequestAt    = textproto.CanonicalMIMEHeaderKey("x-Request-At")
	Authorization = textproto.CanonicalMIMEHeaderKey("authorization")
)
