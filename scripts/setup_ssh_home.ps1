<#
PowerShell helper to create an alternate HOME directory with .ssh for environments
where the Windows user folder contains non-ASCII characters that break Git/MSYS.

Usage (run as the user that will use Git):
  powershell -ExecutionPolicy Bypass -File .\scripts\setup_ssh_home.ps1

Options:
  -HomePath <path>   : path to create (default C:\git_home)
  -KeyComment <str>  : comment for generated SSH key (default timerwallet@prod)
  -Force             : regenerate key if exists

After running:
  - Copy the printed public key and add to GitHub (Settings → SSH and GPG keys).
  - To persist HOME for the current user: `setx HOME "C:\git_home"` then restart shell.
#>

param(
  [string]$HomePath = "C:\git_home",
  [string]$KeyComment = "timerwallet@prod",
  [switch]$Force
)

Write-Host "Using HOME path: $HomePath"

# Create folders
if (-not (Test-Path $HomePath)) {
  New-Item -ItemType Directory -Path $HomePath -Force | Out-Null
  Write-Host "Created $HomePath"
}
$sshDir = Join-Path $HomePath ".ssh"
if (-not (Test-Path $sshDir)) {
  New-Item -ItemType Directory -Path $sshDir -Force | Out-Null
  Write-Host "Created $sshDir"
}

# Create known_hosts (requires OpenSSH client with ssh-keyscan available)
$knownHosts = Join-Path $sshDir "known_hosts"
try {
  ssh-keyscan github.com 2>$null | Out-File -Encoding ascii $knownHosts
  Write-Host "Wrote known_hosts (github.com) to $knownHosts"
} catch {
  Write-Warning "ssh-keyscan failed — ensure OpenSSH client is installed. You can add GitHub host key manually to $knownHosts"
}

# Generate ed25519 key if missing or if Force
$keyPath = Join-Path $sshDir "id_ed25519"
if ((-not (Test-Path $keyPath)) -or $Force) {
  try {
    ssh-keygen -t ed25519 -f $keyPath -N "" -C $KeyComment | Out-Null
    Write-Host "Generated SSH key: $keyPath"
  } catch {
    Write-Warning "ssh-keygen failed. Ensure OpenSSH is available."
  }
} else {
  Write-Host "SSH key already exists at $keyPath"
}

# Start ssh-agent and add key
try {
  Start-Service ssh-agent -ErrorAction SilentlyContinue
} catch {}

try {
  ssh-add $keyPath | Out-Null
  Write-Host "Added key to ssh-agent"
} catch {
  Write-Warning "ssh-add failed — you may need to run ssh-agent or add key manually."
}

# Set HOME for current session
$env:HOME = $HomePath
Write-Host "Set HOME for current session to $HomePath"

Write-Host "\nPublic key (copy this into GitHub → Settings → SSH and GPG keys):\n"
try {
  Get-Content ($keyPath + ".pub") | ForEach-Object { Write-Host $_ }
} catch {
  Write-Warning "Public key not found at $($keyPath + '.pub')"
}

Write-Host "\nTo persist HOME for the current user run (then restart your shell):"
Write-Host "  setx HOME \"$HomePath\""
Write-Host "Then test SSH connection: ssh -T git@github.com"
