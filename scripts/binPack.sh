#!/bin/sh

set -e

echo "[Step 1] build harbor-go-client for packaging"
echo
go build -v ../
echo

echo "[Step 2] preparing for conf/"
cp -r ../conf .
echo


VER=$(cat ../VERSION)

echo "[Step 3] Packaging ${VER}-bin.tar.gz"
echo
tar czvf ${VER}-bin.tar.gz harbor-go-client conf/*.yaml
echo

echo "[Step 4] remove harbor-go-client and conf/"
rm harbor-go-client
rm -r conf/
