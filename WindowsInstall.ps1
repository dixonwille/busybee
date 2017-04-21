$User = $(whoami)
$TaskName = 'BusyBee'
$BBPath = "$PSScriptRoot/BusyBee.exe"
while(!(Test-Path $BBPath)){
    $BBPath = Read-Host -Prompt 'Location of BusyBee.exe (FullPath)'
}

$Test = (Start-Process -FilePath $BBPath -PassThru -Wait)

if($Test.ExitCode -gt 0){
    Write-Error 'Passed in Arguments did not allow BusyBee to execute properly. Please check the configuration file run install again.'
    exit
}

if (Get-ScheduledTask -TaskName $TaskName -ErrorAction 'Ignore') {
    Unregister-ScheduledTask -TaskName $TaskName -Confirm:$false
}

$Action = New-ScheduledTaskAction -Execute $BBPath -Argument $Args
$Trigger = New-ScheduledTaskTrigger -AtLogOn -User $User
$Settings = New-ScheduledTaskSettingsSet -Hidden -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -StartWhenAvailable -RunOnlyIfNetworkAvailable
$Task = New-ScheduledTask -Action $Action -Trigger $Trigger -Settings $Settings
Register-ScheduledTask -TaskName $TaskName -User $User -InputObject $Task > $null
$RegTask = Get-ScheduledTask -TaskName $TaskName
$RegTask.Triggers.Repetition.Interval = "PT5M"
Set-ScheduledTask -User $User -InputObject $RegTask > $null