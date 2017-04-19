#!/bin/sh

BBPATH="$PWD/BusyBee_"

if [ -f "${BBPATH}linux" ]; then
$BBPATH = "${BBPATH}linux"
elif [ -f "${BBPATH}mac" ]; then
$BBPATH = "${BBPATH}mac"
fi

while [ ! -f $BBPATH ]
do
read -p "Location of BusyBee (FullPath): " BBPATH
done

while [ -z "$EXHOST" ]
do
read -p "Exchange Host: " EXHOST
done

while [ -z "$EXEMAIL" ]
do
read -p "Exchange Email: " EXEMAIL
done

while [ -z "$EXUSER" ]
do
read -p "Exchange UserName: " EXUSER 
done

while [ -z "$EXPASS" ]
do
read -s -p "Exchange PassWord: " EXPASS
echo ""
done

while [ -z "$HCHOST" ]
do
read -p "HipChat Host: " HCHOST
done

while [ -z "$HCMENTION" ]
do
read -p "HipChat Mention: @" HCMENTION
done

while [ -z "$HCTOKEN" ]
do
read -p "HipChat Token( $HCHOST/account/api ): " HCTOKEN
done

ARGS="-exHost \"$EXHOST\" -exUser \"$EXUSER\" -exPass \"$EXPASS\" -exUID \"$EXEMAIL\" -hcHost \"$HCHOST\" -hcToken \"$HCTOKEN\" -hcUID \"$HCMENTION\""

eval "$BBPATH $ARGS"

if [ $? -gt 0 ]; then
    echo "Passed in Arguments did not allow BusyBee to execute properly. Please run install again."
    exit
fi

if crontab -l | grep -q "$BBPATH"; then
    crontab -l | grep -v "$BBPATH" | crontab -
fi
(crontab -l 2>/dev/null; echo "*/5 * * * * $BBPATH $ARGS") | crontab -