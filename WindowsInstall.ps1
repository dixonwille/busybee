$User = $(whoami)
$TaskName = 'BusyBee'
$BBPath = "$PSScriptRoot/BusyBee.exe"
while(!(Test-Path $BBPath)){
    $BBPath = Read-Host -Prompt 'Location of BusyBee.exe (FullPath)'
}

if (Get-ScheduledTask -TaskName $TaskName -ErrorAction 'Ignore') {
    Unregister-ScheduledTask -TaskName $TaskName -Confirm:$false
}

$Action = New-ScheduledTaskAction -Execute $BBPath
$Trigger = New-ScheduledTaskTrigger -AtLogOn -User $User
$Settings = New-ScheduledTaskSettingsSet -Hidden -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -StartWhenAvailable -RunOnlyIfNetworkAvailable
$Task = New-ScheduledTask -Action $Action -Trigger $Trigger -Settings $Settings
Register-ScheduledTask -TaskName $TaskName -User $User -InputObject $Task > $null
$RegTask = Get-ScheduledTask -TaskName $TaskName
$RegTask.Triggers.Repetition.Interval = "PT5M"
Set-ScheduledTask -User $User -InputObject $RegTask > $null