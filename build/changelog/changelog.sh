#!/usr/bin/env bash

# https://github.com/hIMEI29A/changelog.sh

# Simple script for easy work with
# https://github.com/mh-cbon/changelog
#
# Use this before creating a new release.

# Print usage message end exits
usage() {
   echo "Usage : changelog.sh -v|--version [VERSION NUMBER]"
}

# Wait for press any key
pressAny() {
	echo "Press any key"
	read ANY
}

# Print colorized warning message
warning() {
	echo -e "\e[1;31mWARNING!!!

	\e[1;35mTemp file CHANGELOG1.md will be created now.
	Edit it and save as CHANGELOG.md\e[0m"
}

# ClI args processing

# if there is no args provided, print usage and exit
if [ "$1" = "" ]
then
    usage
    exit
fi

# args parsing
while [ "$1" != "" ]; do
    case $1 in
        -v | --version )
        VERSION="$2"
        shift
        shift
        ;;
        * )
        usage
        exit 1
    esac
done

echo "Version = $VERSION? [y/n]"
read $REPLY

if [ $REPLY = "Y" ] || [ $REPLY = "y" ] || [ $REPLY = "Yes" ] || [ $REPLY = "yes" ]
then
	changelog prepare

	vim change.log

	changelog finalize --version=$VERSION

	warning

	pressAny

	changelog md --out=CHANGELOG1.md

	vim CHANGELOG1.md # check and save it as CHANGELOG.md

	rm CHANGELOG1.md
fi

if [ $REPLY = "N" ] || [ $REPLY = "n" ] || [ $REPLY = "No" ] || [ $REPLY = "no" ]
then
	echo "Exiting"
fi