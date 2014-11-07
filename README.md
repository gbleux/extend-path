# Utility for PATH manipulation

This program is a small utility for working with PATH, or PATH-like
environment variables. It can append, prepend and delete entries, while
taking care of duplicates and already existing entries. It can also be
used to test for the existence of a directory in the PATH.
The program does *not* modify any environment variables, but merely
writes the results to _standard output_. The user is responsible for
making use of the output (i.e. exporting and overwriting the previous
value). The entries in the path list are split based on the platform
specific separator value.

## Installation

Assuming your go development environment is set up properly, `go get`
should fetch the sources, build the binary and place it in your GOPATH
for future use.

    go get http://github.com/gbleux/extend-path

Alternatively you can build it locally and copy the binary manually

>$> export GOPATH="$PWD"  
>$> go build .  
>$> cp extend-path "$HOME/bin"

## Usage

    extend-path [-d|-a|-p] [-s] [-v] [-e VAR] DIR...

The program was writte with the intend of setting of environments for
interactive shells. As an example one could extend the PATH variable
with additional directories in the _.profile_ initialization script.

    PATH=$(extend-path -relocate -append \
        /opt/android/tools \
        /opt/android/build-tools \
        /opt/android/platform-tools)
    PATH=$(extend-path -delete /usr/local/bin)
    PATH=$(extend-path -validate -prepend "$HOME/bin")
