#!/bin/bash

TAG=$1

if [[ `git status --porcelain` ]];then
    echo "Change detected! Commit all before deploy"
    exit 1
fi

if [[ ! `git rev-parse --abbrev-ref HEAD` =~ release/[0-9]+\.[0-9]+\.[0-9]+ ]];then
    echo "You are not on release task branch!"
    exit 1
fi

if [ $(git tag -l "$TAG") ]; then
    echo "Tag $TAG is already exists on git history"
    exit 1
fi

if [[ $(gcloud config get-value project 2> /dev/null) != 'kubernetes-history-inspector' ]];then
    echo "Active project must be 'kubernetes-history-inspector'"
    exit 1
fi

exit 0