#!/bin/bash


# Set your libiio install path here
LIBIIO_VERSION="0.26"
LIBIIO_PATH="/Users/briananderson/Tools/git/libiio/build/iio.framework/Versions/0.26/iio"
LIBIIO_RPATH="@rpath/iio.framework/Versions/${LIBIIO_VERSION}/iio"


# Check if the main executable exists
if [ ! -f "main" ]; then
    echo "Error: 'main' executable not found."
    exit 1
fi

# Use otool to check for rpath
rpath=$(otool -l main | grep @rpath)

# If rpath is found, change it using install_name_tool
if [ -n "$rpath" ]; then
    echo "Found rpath: $rpath"
    old_rpath="$LIBIIO_RPATH"
    new_rpath="$LIBIIO_PATH"  # Set your desired new rpath here
    echo "Changing rpath to: $new_rpath"
    install_name_tool -change "$old_rpath" "$new_rpath" main
else
    echo "No rpath found in 'main'."
fi
