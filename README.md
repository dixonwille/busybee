# BusyBee 

[![GitHub release](https://img.shields.io/github/release/dixonwille/busybee.svg)](https://github.com/dixonwille/busybee/releases/latest) [![Build Status](https://travis-ci.org/dixonwille/busybee.svg?branch=master)](https://travis-ci.org/dixonwille/busybee) [![GoDoc](https://godoc.org/github.com/dixonwille/busybee?status.svg)](https://godoc.org/github.com/dixonwille/busybee)

> **WINDOWS NOT RECOMMENDED** I am having trouble finding a reliable way to run a task every 5 minutes without annoying the user with a popup each time. If a solution is found and working please let me know!

A bot (or server depending on how far this goes) that will update hipchat status based on exchange calendar. It is being built in a way that it should be able to update any Chatting services status with any calendars events (using interfaces).

## Hipchat User Token

Go to [Account API Access](https://www.hipchat.com/account/api) to grab a token for BusyBee. The token must have `Administer Group` and `View Group` scopes in order for BusyBee to update your users status.

## Install

Go to the [Latest Release](https://github.com/dixonwille/busybee/releases/latest) and download `BusyBee_linux`, `BusyBee_mac` or `BusyBee.exe` depending on if you are on a Linux, Mac or Windows machine respectivley. Also download the appropriate install script for your machine as well. This install script will setup you machine to run BusyBee every 5 minutes.

Once everything is downloaded it is recommended to move both the files somewhere they will be out of the way. A good location may be `%HOME%/BusyBee/` and put them in that folder. 

Once files are moved, if you are on a Unix box (Mac or Linux) we need to make the files executable.

```bash
cd ~/BusyBee #Assuming this is where you moved your files too
chmod +x BusyBee_linux #Or BusyBee_mac
chmod +x UnixInstall.sh
```

> If you are on a Mac and try to double click the files you will get a warning. To bypass, navigate to this file in `Finder`, Control-Click the file, then from the shortcut menu, select `Open`. It will warn you one more time but with a different option now. Click `Open` again.

For Mac or Linux machines it should be as simple as running `UnixInstall.sh` and everything will be set up.

As for Windows you will need to run `BusyBee.exe` first and answer all the questions. Make sure that it did not fail after all prompts were answered. Then run `WindowsInstall.ps1` to install it to the Task Scheduler. You will have to log out then back in so that the schedular triggers the newly added task.

## Edit Task

There may be times when you need to update the task because something has changed. Examples include:

* Moved executable
* Change Hipchat @Mention name
* Change AD Auth Password

To edit the task, simply edit the configuration file and on next run BusyBee will use those values.

If the value you are trying to edit is encrypted, you have two options:

1. Remove it from the configuration file and then run BusyBee to reprompt for that field
2. Run `./BusyBee_linux encrypt string` replacing string with the value you want to encrypt. Copy the output and paste it into the config

> Using the latter, make sure you run your version of BusyBee.

## Uninstalling

If you want to uninstall or upgrade to a new version of BusyBee follow the steps bellow depending on you operating system to remove BusyBee. It is recommended to Uninstall everytime you want to upgrade for right now. After you uninstall using one of the methods below, make sure to remove the executable and the install script that was downloaded.

### Unix (Mac or Linux)

Open a terminal and run `crontab -e`. This will open an editor. Remove the line that has `/path/to/BusyBee` in it. This will stop it from running every 5 minutes.

### Windows

1. Search for Task Scheduler in Windows 8 or Windows 10.
2. On the left side click `Task Scheduler Library`
3. In the middle find the Task named `BusyBee`
4. Click on it
5. On the left there is a delete button. Click it.
