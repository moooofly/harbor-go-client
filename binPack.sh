#!/usr/bin/env bash

echo "===> Packaging v$(cat VERSION)-bin.tar.gz ..."
echo
tar czvf v$(cat VERSION)-bin.tar.gz harbor-go-client config.yaml rp.yaml
echo
echo "===> Packaging complete."
