#!/bin/bash

url="https://localhost:8443/v1/projects"
#url="https://localhost:8443/help"

ab -f ALL -kc 1000 -n 10000 -m POST $url
