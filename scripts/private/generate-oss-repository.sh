#!/bin/bash
# Copyright 2024 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


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

../scripts/private/remove-internal-codes.sh ../.ossignore

rm .git
mv ../.git-escape .git
cd ..
mv oss oss-stage
git worktree remove oss
mv oss-stage oss