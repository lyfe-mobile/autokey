#!/bin/bash

notify() {
    COLOR=$1
    MESSAGE=$2
    declare -A colors
    colors[red]=1
    colors[green]=2

    if [ -n "${HIPCHAT_ROOM}" -a -n "${HIPCHAT_TOKEN}" ]; then
        curl -sS \
             -H 'Content-type: application/json' \
             -d "{\"color\":\"$COLOR\", \"message\":\"$MESSAGE\", \"message_format\":\"html\", \"notify\":true}" \
             https://api.hipchat.com/v2/room/$HIPCHAT_ROOM/notification?auth_token=$HIPCHAT_TOKEN
    else
        HC=" (HipChat disabled)"
    fi

    tput setaf ${colors[$COLOR]}; echo $MESSAGE$HC
}

buildstatus() {
    COLOR=$1
    STATUS=$2
    if [ -n "${TRAVIS_COMMIT}" ]; then
        COMMIT=${TRAVIS_COMMIT:0:7}
        BRANCH=${TRAVIS_BRANCH}
        URL="https://travis-ci.com/${TRAVIS_REPO_SLUG}/builds/${TRAVIS_BUILD_ID}"
        CI_INFO="<A HREF=$URL>${TRAVIS_REPO_SLUG}#${TRAVIS_BUILD_NUMBER}</A>"
    else # We're not on Travis
        COMMIT=`git rev-parse --short HEAD`
        BRANCH=`git rev-parse --abbrev-ref HEAD`
        CI_INFO=`hostname`
    fi

    AUTHOR=`git --no-pager show -s --format='%an'`
    notify $COLOR "${CI_INFO} (${BRANCH} - ${COMMIT} : ${AUTHOR}): ${STATUS}"
}

do_test () {
    PROFILE=coverage.txt
    # For reporting to coveralls.io
    go get github.com/mattn/goveralls
    # For merging coverage profiles
    go get github.com/wadey/gocovmerge
    # For running the tests
    go get github.com/onsi/ginkgo/ginkgo

    export PKGS=$(go list ./... | egrep -v '/vendor/')
    export PKGS_DELIM=$(echo "$PKGS" | paste -sd "," -)

    set -o pipefail # To detect when a test fails.

    if eval "`go list -f '{{if or (len .TestGoFiles) (len .XTestGoFiles)}}ginkgo '"$TAGS"' '"$V"' --randomizeAllSpecs --randomizeSuites  -gcflags=-l -covermode count -coverprofile {{.Name}}_{{len .Imports}}_{{len .Deps}}.coverprofile -coverpkg '"$PKGS_DELIM"' {{.Dir}} |& grep -v "warning: no packages being tested depend on github.com/lyfe-mobile" && \{{end}}' $PKGS; echo :` "; then
        # Build was successful.
        gocovmerge `ls **/*.coverprofile` > $PROFILE
        COVERAGE=`go tool cover -func=$PROFILE | awk '/(statements)/ {print $NF}'`
        go tool cover -html=$PROFILE -o cover.htm
        if [[ $TRAVIS_COMMIT ]];  then
            goveralls -service=travis-ci
            rm -f cover.htm
        fi
        buildstatus green "build successful. Coverage: $COVERAGE."
    else
        EXITCODE=$?
        buildstatus red "build has failed!"
        exit $EXITCODE
    fi
}

do_build() {
    echo -n Building...
    if ! git diff-index --quiet HEAD -- ; then
        dirty="-X main.GitDirty=true"
    fi

    go build $TAGS -ldflags "-X main.BuildDate=`date -u '+%Y-%m-%d_%I:%M:%S%p'` \
      $dirty \
      -X main.GitHash=`git rev-parse --short HEAD` \
      -X main.GitBranch=`git rev-parse --abbrev-ref HEAD`" .
    EXITCODE=$?
    echo .
    if [[ $EXITCODE -ne 0 ]]; then
      exit $EXITCODE
    fi
}

shopt -s globstar

trap "rm -f $PROFILE; rm -f **/*.coverprofile; tput sgr0" EXIT

# Make sure all Go code is formatted.

FILES=$(go fmt -n ./... |sed 's/ -w//' | bash)

if [[ -n "$FILES" ]]; then
    notify red "Build cannot continue until these files are gofmtd: ${FILES//$'\n'/, }"
    exit 1
fi

while [ $# -gt 0 ]; do
    case "$1" in
        -tags) shift; TAGS="-tags $1";;
        -v) V="-v";;
        -build) BUILD=1;;
        *) usage;;
    esac
    shift
done

if [[ $BUILD ]] ; then
    do_build
else
    do_test
fi
