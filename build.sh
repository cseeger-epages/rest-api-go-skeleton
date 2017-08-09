#!/bin/bash

params="-o bin/api"
sourcefiles="src/*.go"

go build $params $sourcefiles
