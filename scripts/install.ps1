$ErrorActionPreference = "Stop"

$repo = "https://github.com/bokshi-gh/http-file-server.git"
$tmpDir = "$env:TEMP\http-file-server-build"
$exeName = "httpfs.exe"
$installRoot = "C:\Program Files\httpfs"
$installBin = "$installRoot\bin"

if (Test-Path $tmpDir) {
    Remove-Item -Recurse -Force $tmpDir
}

git clone $repo $tmpDir

Set-Location "$tmpDir\cmd\httpfs"
go build -o $exeName

if (-not (Test-Path $installBin)) {
    New-Item -ItemType Directory -Path $installBin -Force | Out-Null
}

Move-Item $exeName "$installBin\$exeName" -Force

Copy-Item "$tmpDir\README.md" "$installRoot\README.md" -Force
Copy-Item "$tmpDir\LICENSE" "$installRoot\LICENSE" -Force

$oldPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if (-not ($oldPath -split ";" | Where-Object { $_ -eq $installBin })) {
    [Environment]::SetEnvironmentVariable("PATH", "$oldPath;$installBin", "User")
    Write-Host "Added $installBin to user PATH. You may need to restart PowerShell to use it."
} else {
    Write-Host "$installBin already exists in PATH."
}

Set-Location $env:TEMP
Remove-Item -Recurse -Force $tmpDir

Write-Host "Build complete!"
Write-Host "Executable: $installBin\$exeName"
Write-Host "Documentation: $installRoot\README.md"
Write-Host "License: $installRoot\LICENSE"
Write-Host "You can now run 'httpfs.exe' from any PowerShell session."
