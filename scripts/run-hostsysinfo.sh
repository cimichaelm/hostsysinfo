#!/bin/sh
# file: run-hostsysinfo.sh
#
defaults()
{
    bindir=./bin
    config=config/site-config.yaml
    prog=$bindir/hostsysinfo
}

run_prog()
{
    Largs=$*
    $prog -config $config $Largs
}

args=$*
defaults

run_prog $*

