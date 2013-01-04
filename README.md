webchat
=======

go based websockets chats, served up in a bootstrap website

![nothing fancy](https://raw.github.com/dearing/webchat/master/www/img/ss.png)

build
----

Nothing to it.

```
git clone https://github.com/dearing/webchat

go build

webchat --help

Usage of webchat:
  -conf="webchat.conf": json configuration
  -gen=false: generate a default conf in working directory

webchat -gen=true
webchat -conf=webchat.conf

```

default config:
----
```
{
  "Certificate": "cert.pem",
	"CertificateKey": "cert.key",
	"WWWHost": ":9000",
	"WWWRoot": "www",
	"Tag": "/chat",
	"TLS": false,
	"Verbose": true,
	"EscapeData": false
}
```
