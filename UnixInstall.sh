#!/bin/sh

BBPATH="$PWD/BusyBee_"

if [ -f "${BBPATH}linux" ]; then
BBPATH="${BBPATH}linux"
elif [ -f "${BBPATH}mac" ]; then
BBPATH="${BBPATH}mac"
fi

while [ ! -f $BBPATH ]
do
read -p "Location of BusyBee (FullPath): " BBPATH
done

eval "$BBPATH"

if [ $? -gt 0 ]; then
    echo "Passed in Arguments did not allow BusyBee to execute properly. Please check the configuration file run install again."
    exit
fi

if crontab -l | grep -q "$BBPATH"; then
    crontab -l | grep -v "$BBPATH" | crontab -
fi
(crontab -l 2>/dev/null; echo "*/5 * * * * cd $PWD;$BBPATH") | crontab -