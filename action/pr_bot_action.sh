#!/bin/sh

set -ex
mkdir -p ./bundles
mkdir -p /opt/app/bundles
ls -l
/opt/app/opa build ./bundles -o /opt/app/bundles/pr-bot-policy.tar.gz
tar -ztvf /opt/app/bundles/pr-bot-policy.tar.gz
/opt/app/pr-bot-cli evaluate --config /opt/app/config/local.yaml
