#!/bin/bash

url="https://localhost:8443/v1/projects"
#url="https://localhost:8443/help"

ab -f TLS1.2 -kc 300 -n 10000 -m POST $url
