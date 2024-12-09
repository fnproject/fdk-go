#!/usr/bin/env bash

user="fnproject"
image="go"
goversion="1"
stretch="stretch"
alpine="alpine"


docker push ${user}/${image}:${goversion}-${stretch}
docker push ${user}/${image}:${goversion}-${stretch}-dev

docker push ${user}/${image}:${goversion}-${alpine}
docker push ${user}/${image}:${goversion}-${alpine}-dev
