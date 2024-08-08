#!/bin/sh

set -ex
BUNDLES_DIR=./bundles
mkdir -p $BUNDLES_DIR
mkdir -p /opt/app/bundles
ls -l
(cd $BUNDLES_DIR && /opt/app/opa build . -o /opt/app/bundles/pr-bot-policy.tar.gz)
tar -ztvf /opt/app/bundles/pr-bot-policy.tar.gz
/opt/app/pr-bot-cli evaluate --config /opt/app/config/local.yaml
