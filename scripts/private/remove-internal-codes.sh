#!/bin/bash

# This is temporal adhoc script to extract the internal KHI implementation for OSSing.
# Eventually the OSSed version should be the base and internal code should be amended in internal repository. But for this time, extract the OSS part from internal repository.
# The dependency graph will be flipped once we done OSSing.

set -x

git clone sso://khi/kubernetes-history-inspector-oss oss && (cd oss && f=`git rev-parse --git-dir`/hooks/commit-msg ; mkdir -p $(dirname $f) ; curl -Lo $f https://gerrit-review.googlesource.com/tools/hooks/commit-msg ; chmod +x $f)
cd oss

# Safe guard not to destroy parent directory on the previous command error
CURRENT_DIRECTORY=${PWD##*/}
if [ $CURRENT_DIRECTORY != "oss" ]; then
    echo "unexpected status: directory must be inside of oss"
    exit 1
fi

mv .git ../.git-escape
cd ..
git worktree add oss origin/beta
cd oss

# Safe guard not to destroy parent directory on the previous command error
CURRENT_DIRECTORY=${PWD##*/}
if [ $CURRENT_DIRECTORY != "oss" ]; then
    echo "unexpected status: directory must be inside of oss"
    exit 1
fi
VERSION=$(cat VERSION)
OSSIGNORE=$(sed '/^[ \t]*#/d' ../.ossignore) # Removing comment lines

echo "$OSSIGNORE" | while read -r pattern;
do
    rm -rf $pattern -r
done
rm .git
mv ../.git-escape .git
cd ..
mv oss oss-stage
git worktree remove oss
mv oss-stage oss