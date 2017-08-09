#!/bin/bash

crtpath="certs/server.crt"
keypath="certs/server.key"

go run src/*go -crt $crtpath -key $keypath
