#!/usr/bin/env bash

PROGNAME=md-file-viewer
mkdir -p builds

function do_build()
{
    if [[ "$GOOS" != "darwin" || "$GOARCH" != "386" ]]; then
        echo "Building $PROGNAME-$GOOS-$GOARCH"
        go build
        OUTPROGNAME=$PROGNAME
        OUTGOARCH=$GOARCH
        if [[ "$GOOS" == "windows" ]]; then
            OUTPROGNAME=$OUTPROGNAME.exe
            OUTGOARCH=$OUTGOARCH.exe
        fi
        mv $OUTPROGNAME builds/$PROGNAME-$GOOS-$OUTGOARCH
    fi
}

function do_arch_builds()
{
    GOARCH=amd64 do_build
    GOARCH=386 do_build
}

GOOS=linux do_arch_builds
GOOS=darwin do_arch_builds
GOOS=windows do_arch_builds

echo "Copying css, templates..."
cp -R css templates builds/
