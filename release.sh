#!/bin/bash

##############################################################################
# RELEASING JAVA-OPERATOR-PLUGINS
#
# Every java-operator-plugins release should have a corresponding git
# semantic version tag begining with `v`, for example: `v1.2.3`.
#
# STEP 1: Create a release branch with the name vX.Y.x.
#         For example: git checkout -b v1.2.x
#
# STEP 2: Run this release script by providing the java-operator-plugins
#         release version as an argument in the above mentioned format.
#         For example: ./release.sh vX.Y.Z
#
# STEP 3: This script will create a release tag locally. YOU, the releaser,
#         will push the release branch and tag:
#
#           git push upstream <release-branch>
#           git push upstream <tag-name>
#
#         where <tag-name> is the release version.
#
# STEP 4: Update the release notes in github with the changes included in
#         corresponding java-operator-plugins version.
##############################################################################

set -eu

if [[ $# != 1 ]]; then
	echo "usage: $0 vX.Y.Z"
	exit 1
fi

CUR_BRANCH=`git branch --show-current`
VER=$1
NUMRE="0|[1-9][0-9]*"
PRERE="\-(alpha|beta|rc)\.[1-9][0-9]*"

# check that we are in a release branch
if ! [[ "$CUR_BRANCH" =~ ^v($NUMRE)\.($NUMRE)\.x$ ]]; then
	echo "branch: \"$CUR_BRANCH\" is not a release branch, please create a release branch in the form of vX.Y.x"
	exit 1
fi

# verify the version format
if ! [[ "$VER" =~ ^v($NUMRE)\.($NUMRE)\.($NUMRE)($PRERE)?$ ]]; then
	echo "malformed version: \"$VER\""
	exit 1
fi

# make sure we don't have anything uncommitted lying around
if ! git diff-index --quiet HEAD --; then
	echo "directory has uncommitted files"
	exit 1
fi

# run tests
echo "Running tests"
make lint test

# tag the release commit and verify its tag
echo "Creating a new tag for java-operator-plugins version $VER"
git tag --sign --message "java-operator-plugins $VER" "$VER"
git verify-tag --verbose $VER

# remind the releaser to do the next steps
echo ""
echo "Don't forget to:"
echo ""
echo "git push upstream <release-branch>"
echo "git push upstream $VER"
echo ""
echo "Also update the release notes in github for this tag."
