# Run only on a disposable Windows host; never use an existing user installation.
param([Parameter(Mandatory)][string]$PackageDirectory)
$ErrorActionPreference = 'Stop'
$version = (Get-Content "$PSScriptRoot/../VERSION" -Raw).Trim()
$installDir = Join-Path $env:LOCALAPPDATA 'Programs/scrcpy-tui'
if (Test-Path $installDir) { throw 'MSI test requires a clean host; installation already exists' }
$originalPath = [Environment]::GetEnvironmentVariable('Path', 'User')
$current = (Resolve-Path "$PackageDirectory/scrcpy-tui_${version}_windows_amd64.msi").Path
$prior = (Resolve-Path "$PackageDirectory/upgrade-fixture.msi").Path
$logDir = Join-Path $env:RUNNER_TEMP 'msi-logs'
New-Item -ItemType Directory -Force $logDir | Out-Null
function Invoke-Msi([string]$Action, [string]$Package, [string]$Name) {
    $log = Join-Path $logDir "$Name.log"
    $p = Start-Process msiexec.exe -ArgumentList @($Action, "`"$Package`"", '/qn', '/norestart', '/l*v', "`"$log`"") -Wait -PassThru
    if ($p.ExitCode -ne 0) { throw "MSI $Name failed: $($p.ExitCode); log: $log" }
}
function Assert-Installed([string]$ProductVersion) {
    $exe = Join-Path $installDir 'scrcpy-tui.exe'
    $actual = & $exe --version
    if ($LASTEXITCODE -ne 0 -or $actual -ne $version) { throw 'Installed binary version mismatch' }
    $value = (Get-ItemProperty 'HKCU:/Software/e-novatisHQ/scrcpy-tui').InstalledVersion
    if ($value -ne $ProductVersion) { throw "Installed product version mismatch: $value" }
    $entries = [Environment]::GetEnvironmentVariable('Path', 'User').Split(';')
    $matches = @($entries | Where-Object { $_.TrimEnd('\') -ieq $installDir.TrimEnd('\') })
    if ($matches.Count -ne 1) { throw 'User PATH must contain exactly one installation directory' }
    foreach ($entry in ($originalPath -split ';' | Where-Object { $_ })) {
        if ($entries -notcontains $entry) { throw 'MSI changed an unrelated PATH entry' }
    }
}
try {
    Invoke-Msi '/i' $prior 'install-prior'
    Assert-Installed '0.0.1'
    Invoke-Msi '/i' $current 'upgrade-current'
    Assert-Installed $version
    Invoke-Msi '/i' $current 'reinstall-current'
    Assert-Installed $version
    Invoke-Msi '/x' $current 'uninstall-current'
    if (Test-Path $installDir) { throw 'Uninstall left the installation directory' }
    if (Test-Path 'HKCU:/Software/e-novatisHQ/scrcpy-tui') { throw 'Uninstall left product registration' }
    $after = [Environment]::GetEnvironmentVariable('Path', 'User')
    # Windows Installer can normalize a trailing separator while removing its
    # appended entry. Require every actual original entry, in the same order.
    $beforeEntries = @(([string]$originalPath).TrimEnd(';') -split ';')
    $afterEntries = @(([string]$after).TrimEnd(';') -split ';')
    if (($beforeEntries -join ';') -cne ($afterEntries -join ';')) {
        throw "Uninstall changed user PATH entries (before: $($beforeEntries.Count), after: $($afterEntries.Count))"
    }
    Write-Output 'MSI install, upgrade, reinstall, PATH preservation and uninstall: PASS'
} finally {
    # This host is disposable. Attempt owned product cleanup even after an assertion fails.
    if (Test-Path (Join-Path $installDir 'scrcpy-tui.exe')) {
        $null = Start-Process msiexec.exe -ArgumentList @('/x', "`"$current`"", '/qn', '/norestart') -Wait -PassThru
        $null = Start-Process msiexec.exe -ArgumentList @('/x', "`"$prior`"", '/qn', '/norestart') -Wait -PassThru
    }
}
