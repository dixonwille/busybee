# BusyBee

A bot (or server depending on how far this goes) that will update hipchat status based on exchange calendar. It is being built in a way that it should be able to update any Chatting services status with any calendars events (using interfaces).

## Hipchat User Token

Go to [Account API Access](https://www.hipchat.com/account/api) to grab a token for BusyBee. The token must have `Administer Group` and `View Group` scopes in order for BusyBee to update your users status.

## Install

Go to the [Latest Release](https://github.com/dixonwille/busybee/releases/latest) and download `BusyBee_linux`, `BusyBee_mac` or `BusyBee.exe` depending on if you are on a Linux, Mac or Windows machine respectivley. Also download the appropriate install script for your machine as well. This install script will setup you machine to run BusyBee every 5 minutes.

Once everything is downloaded it is recommended to move both the files somewhere they will be out of the way. A good location may be `%HOME%/BusyBee/` and put them in that folder. 

Once files are moved, if you are on a Unix box (Mac or Linux) we need to make the files executable.

```bash
cd ~/BusyBee #Assuming this is where you moved your files too
chmod +x UnixInstall.sh
chmod +x WindowsInstall.ps1
```

> If you are on a Mac and try to double click the files you will get a warning. To bypass, navigate to this file in `Finder`, Control-Click the file, then from the shortcut menu, select `Open`. It will warn you one more time but with a different option now. Click `Open` again.

Now run the install script for your machine (`UnixInstall.sh` or `WindowsInstall.ps1`). It will prompt you to answer a few questions.

> `Exchange UserName` is not your email address. It is the AD Auth username you use to sign in Exchange.

> Do not include the `@` symbol in the `Hipchat Mention` prompt as it is a special character and may throw the script off.

> `http://` and `https://` is not required as BusyBee will check. If it does not exist it uses `https://`.

## Edit Task

There may be times when you need to update the task because something has changed. Examples include:

* Moved executable
* Change Hipchat @Mention name
* Change AD Auth Password

To edit the task, simply rerun the install script as it should remove the old task and add the updated new task in its place.