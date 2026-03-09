# CPTUI

CPTUI is a lightweight program written in Go to test competitive programming solutions and display them.

## Features

- View code file
- View associated devlog file
- Open link to judging platform (if included in devlog)
- Run preset test cases
- Run custom test cases

This can be a good way to systematically judge competitive programming solutions with a large array of test cases.

## Install

Download and unzip a release archive from GitHub Releases.

## Run

Run from a directory that contains `./problems`.

### macOS

```bash
./cptui
```

### Linux

```bash
./cptui
```

### Windows (PowerShell)

```powershell
.\cptui.exe
```

**A C++ compiler is required to run this program.**

Use a **GCC/g++** toolchain (needed for `bits/stdc++.h`):

- macOS (Homebrew GCC): https://formulae.brew.sh/formula/gcc
- Linux (GCC install docs): https://gcc.gnu.org/install/
- Windows (MSYS2 + MinGW-w64 GCC): https://www.msys2.org/

# Demo

The only file that is guaranteed to work without `<bits/stdc++.h>` support is `bf1easy`. Otherwise, you may need to install a compiler with support for `<bits/stdc++.h>`.

## Maintainer: Publish a New Version

1. Commit and push to `main`.
2. Create and push a semver tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

3. GitHub Actions workflow `.github/workflows/release.yml` will:
- run tests
- build multi-platform binaries
- publish release archives
