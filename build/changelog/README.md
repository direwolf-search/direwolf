# changelog.sh

Simple bash script for easy work with [changelog](https://github.com/mh-cbon/changelog).

## Install & usage

You must install [changelog](https://github.com/mh-cbon/changelog) first.

Then

`git clone github.com/hIMEI29A/changelog.sh`

`make`

When new git repo started, init [changelog](https://github.com/mh-cbon/changelog) with

`changelog init`

After new release final commit and before create new tag call **this** script:

`changelog.sh -v | --version`

where version is a new tag.