#!/bin/sh

set -ex
PRBOT_DIR=.prbot
mkdir -p $PRBOT_DIR
mkdir -p /opt/app/bundles
ls -l
(cd $PRBOT_DIR && /opt/app/opa build . -o /opt/app/bundles/pr-bot-policy.tar.gz)
tar -ztvf /opt/app/bundles/pr-bot-policy.tar.gz
/opt/app/pr-bot-cli evaluate --config /opt/app/config/local.yaml
