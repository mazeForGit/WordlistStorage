#!/bin/bash
echo "starting build.sh"

# Stops the process if something fails
# set -xe

# All of the dependencies needed/fetched for your project.
# This is what actually fixes the problem so that EB can find your dependencies. 
# FOR EXAMPLE:

echo "get go dependencies"
go get github.com/gin-gonic/gin
go get github.com/gin-contrib/static
go get github.com/sirupsen/logrus
go get github.com/mazeForGit/Wordlist

echo "create the application binary wich is referenced by Procfile"
go build -o bin/application server.go

chmod +x bin/application

echo "finished build.sh"