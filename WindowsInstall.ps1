$User = $(whoami)
$TaskName = 'BusyBee'
$BBPath = "$PSScriptRoot/BusyBee.exe"
while(!(Test-Path $BBPath)){
    $BBPath = Read-Host -Prompt 'Location of BusyBee.exe (FullPath)'
}

while([string]::IsNullOrEmpty($ExHost)){
    $ExHost = Read-Host -Prompt 'Exchange Host'
}

while([string]::IsNullOrEmpty($ExEmail)){
    $ExEmail = Read-Host -Prompt 'Exchange Email'
}

while([string]::IsNullOrEmpty($ExUser)){
    $ExUser = Read-Host -Prompt 'Exchange UserName'
}

while([string]::IsNullOrEmpty($ExPass)){
    $ExPass = Read-Host -Prompt 'Exchange PassWord' -AsSecureString
}

while([string]::IsNullOrEmpty($HcHost)){
    $HcHost = Read-Host -Prompt 'Hipchat Host'
}

while([string]::IsNullOrEmpty($HcMention)){
    $HcMention = Read-Host -Prompt 'Hipchat Mention @'
}

while([string]::IsNullOrEmpty($HcToken)){
    $HcToken = Read-Host -Prompt "Hipchat Token ( $HcHost/account/api )"
}

$RawPass = [Runtime.InteropServices.Marshal]::PtrToStringAuto([Runtime.InteropServices.Marshal]::SecureStringToBSTR($ExPass))

$Args = "-exHost `"$ExHost`" -exUser `"$ExUser`" -exPass `"$RawPass`" -exUID `"$ExEmail`" -hcHost `"$HcHost`" -hcToken `"$HcToken`" -hcUID `"$HcMention`""

$Test = (Start-Process -FilePath $BBPath -ArgumentList ($Args -split ' ') -WindowStyle Hidden -PassThru -Wait -RedirectStandardError 'blah.txt')

if($Test.ExitCode -gt 0){
    Write-Error 'Passed in Arguments did not allow BusyBee to execute properly. Please run install again.'
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