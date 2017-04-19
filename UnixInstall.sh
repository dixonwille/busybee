#!/bin/sh
BBPATH="$PWD/BusyBee"
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

ARGS="-exHost '$EXHOST' -exUser '$EXUSER' -exPass '$EXPASS' -exUID '$EXEMAIL' -hcHost '$HCHOST' -hcToken '$HCTOKEN' -hcUID '$HCMENTION'"

if ! crontab -l | grep -q "$BBPATH"; then
   (crontab -l 2>/dev/null; echo "*/5 * * * * $BBPATH $ARGS") | crontab -
fi
