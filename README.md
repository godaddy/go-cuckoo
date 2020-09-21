# Cuckoo client library

[![Go Report Card](https://goreportcard.com/badge/github.com/godaddy/go-cuckoo)](https://goreportcard.com/report/github.com/godaddy/go-cuckoo)
[![Documentation](https://godoc.org/github.com/godaddy/go-cuckoo?status.svg)](https://godoc.org/github.com/godaddy/go-cuckoo)

Original Author: [Connor Lake](mailto:clake1@godaddy.com)

A simple go client library for the [cuckoo api](https://cuckoo.readthedocs.io/en/latest/usage/api).  See the godoc for more details and examples.

## Example Usage

```go
c := cuckoo.New(os.Getenv("API_KEY"), os.Getenv("BASEURL"))

status, _ := c.CuckooStatus(context.Background())

fmt.Println(status.Version)
```

## Auth

Cuckoo uses an API key for auth, you can see more details in the [cuckoo api documentation](https://cuckoo.readthedocs.io/en/latest/usage/api).
