#!/bin/sh

docker rmi registry.cn-hangzhou.aliyuncs.com/bullyun/stateful-operator:1.0 -f
docker tag deploy-html-controller:1.0 registry.cn-hangzhou.aliyuncs.com/bullyun/stateful-operator:1.0
docker push registry.cn-hangzhou.aliyuncs.com/bullyun/stateful-operator:1.0

