# dt-services

## Development

```bash
# open 3 terminal windows
# port forward to consul
$ kubectl port-forward consul-server-0 8500:8500
# port forward to mongoDB
$ kubectl port-forward pw-mongodb-replicaset-0 27017:27017
# get deps
$ make init
# start application
$ make dev
```

## Deployment

```bash
# create branch
# do work
# bumpversion
$ make bumpversion-patch
# stage, commit, push
$ git add .
$ git commit -m "commit message"
$ git push origin {branch name}
# create PR on the github webapp
# merge PR, CICD only builds and deploys on the master branch
```
