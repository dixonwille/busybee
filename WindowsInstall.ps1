$User = $(whoami)
$TaskName = 'BusyBee'
if (Get-ScheduledTask -TaskName $TaskName) {
    exit
}
$Action = New-ScheduledTaskAction -Execute "C:/Path/To/Exe" -Argument "-any args -passed here"
$Trigger = New-ScheduledTaskTrigger -AtLogOn -User $User
$Settings = New-ScheduledTaskSettingsSet -Hidden -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -StartWhenAvailable -RunOnlyIfNetworkAvailable
$Task = New-ScheduledTask -Action $Action -Trigger $Trigger -Settings $Settings
Register-ScheduledTask -TaskName $TaskName -User $User -InputObject $Task > $null
$RegTask = Get-ScheduledTask -TaskName $TaskName
$RegTask.Triggers.Repetition.Interval = "PT5M"
Set-ScheduledTask -User $User -InputObject $RegTask > $null
