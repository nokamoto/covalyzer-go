#/bin/bash

cd $(mktemp -d)

# https://magefile.org/
git clone https://github.com/magefile/mage
cd mage
go run bootstrap.go
