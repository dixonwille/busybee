#!/bin/sh
if ! crontab -l | grep -q '/path/to/job'; then
    (crontab -l 2>/dev/null; echo "*/5 * * * * /path/to/job -with args") | crontab -
fi
