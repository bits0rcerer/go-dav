# Go-DAV

Minimal WebDAV server written in go

---

## Environment Variables

| Key                       | Explanation                                  | Default |      Example      |
|---------------------------|----------------------------------------------|:-------:|:-----------------:|
| GODAV_ROOT                | Path to root directory served by this server |    -    |      "/data"      |
| GODAV_PREFIX              | URL path prefix                              |   ""    | "/path/to/webdav" |
| GODAV_NO_AUTH             | Disables authorization                       | "False" |      "True"       |
| PORT                      | Port to listen on                            | "8080"  |      "42069"      |
| GODAV_USER_\<user name\>¹ | SHA256 hash of \<user name\>                 |    -    |         -         |

#### ¹ Example

You are able to create multiple users

```bash
# user silverhand with password "saka sucks"
GODAV_USER_silverhand=4822d7069138c1975cc8fb4453e41255b971c3c7483f2b063a76932b230a6564

# another user lucy with password "moon"
GODAV_USER_lucy=9e78b43ea00edcac8299e0cc8df7f6f913078171335f733a21d5d911b6999132
```

---

## Notice

- This server is meant to run behind a reverse proxy which should take care of SSL
- This server trusts all common proxy headers

> [NOTE](https://pkg.go.dev/github.com/gorilla/handlers?utm_source=godoc#ProxyHeaders): This middleware should only be
> used when behind a reverse proxy like nginx, HAProxy or Apache. Reverse proxies that don't (or are configured not to)
> strip these headers from client requests, or where these headers are accepted "as is" from a remote client (e.g. when
> Go
> is not behind a proxy), can manifest as a vulnerability if your application uses these headers for validating the '
> trustworthiness' of a request.