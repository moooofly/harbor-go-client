#!/bin/sh

echo "----- build harbor-go-client for packaging -----"
go build -v ../
echo "------------------------------------------------"
echo

VER=$(cat ../VERSION)

echo "===> Packaging ${VER}-bin.tar.gz"
echo
tar czvf ${VER}-bin.tar.gz harbor-go-client config.yaml rp.yaml
echo
echo "===> Packaging complete."
echo

echo "----- remove harbor-go-client -----"
rm harbor-go-client
