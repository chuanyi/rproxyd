# RProxyd

A tiny easy-to-use reverse proxy written by golang.

## Config

RProxyd using ./rproxyd as its config file, and supported commands as follows:

+ listen ip:port

  > listen address, `listen :8888` for example.

+ proxy /route  http://ip:port

  > proxy rule, `proxy /api  http://127.0.0.1:15000/` for example.
  >
  > RProxyd will try to match-router-rule one by one, if failed, return 404. 

## Build

`>go build rproxyd.go`

## Service

You can use [NSSM](http://www.nssm.cc/) to let rproxyd run as a windows service.