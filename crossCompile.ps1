[CmdletBinding()]
param (
    [Parameter(Mandatory=$false)][string]$Command ,
    [Parameter(Mandatory=$false)][string]$TargetOS="linux",
    [Parameter(Mandatory=$false)][string]$arch="amd64"
)

$goosList   = "aix android darwin dragonfly freebsd hurd illumos ios js linux nacl netbsd openbsd plan9 solaris windows zos".Split(" ")
$goarchList = "386 amd64 amd64p32 arm armbe arm64 arm64be ppc64 ppc64le mips mipsle mips64 mips64le mips64p32 mips64p32le ppc riscv riscv64 s390 s390x sparc sparc64".Split(" ")

function RevertChanges {
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
}

function GoBuild([string]$os,[string]$arch,[string]$additionString) {
    $env:GOOS = $os
    $env:GOARCH = $arch
    Write-Host "building for $os on $arch"
    $name = ((Get-Location).path.Split("\")[-1] + "_" + $os + "_" + $arch +"_"+$additionString)
    go build -o $name
}

if ($Command -eq "nix") {   # compile for all BSDs/linux for amd64
    foreach ($os in $goosList) {
        if ($os -match "bsd" -or $os -match "linux" -or $os -match "dragonfly") {
            GoBuild $os $arch (Get-Date -Format "MM_dd_yyyy_HH-mm")
        }
    }
} elseif ($Command -eq "bsdX64") { # compile for all BSDs on x64 (amd64)
    $arch = "amd64"
    foreach ($os in $goosList) {
        if ($os -match "bsd" -or $os -match "dragonfly") {
            GoBuild $os $arch (Get-Date -Format "MM_dd_yyyy_HH-mm")
        }
    }
} elseif ($Command -eq "bsd4All") { # compile for all BSDs on all arch
    $ErrorActionPreference = "SilentlyContinue"
    $arch = "amd64"
    foreach ($os in $goosList) {
        foreach($arch in $goarchList){
            if ($os -match "bsd" -or $os -match "dragonfly") {
                GoBuild $os $arch (Get-Date -Format "MM_dd_yyyy_HH-mm")
            }
        }
    }
} elseif ($Command -eq "openbsdX64") { # compile for openBSD on x64 (amd64)
    $arch = "amd64"
    $os = "openbsd"
    GoBuild $os $arch (Get-Date -Format "MM_dd_yyyy_HH-mm")
} elseif ($Command -eq "piBuild") { # compile for linux on arm (armv6) (for rpi 0 and  rpi 1)
    $arch = "arm"
    $os = "linux"
    $env:GOARM = 6
    GoBuild $os $arch ((Get-Date -Format "MM_dd_yyyy_HH-mm")+"armv6_pi")
} elseif ($Command -eq "pi2Build") { # compile for linux on arm (armv7) (for rpi 2 and other armv7)
    $arch = "arm"
    $os = "linux"
    $env:GOARM = 7
    GoBuild $os $arch ((Get-Date -Format "MM_dd_yyyy_HH-mm")+"armv7_pi")
} elseif ($Command -eq "pi3Build") { # compile for linux on arm (armv7) (for rpi 2 and other armv7)
    $arch = "arm64"
    $os = "linux"
    GoBuild $os $arch ((Get-Date -Format "MM_dd_yyyy_HH-mm")+"armv7_pi")
} else {
    GoBuild $TargetOS $arch (Get-Date -Format "MM_dd_yyyy_HH-mm")
}

RevertChanges




