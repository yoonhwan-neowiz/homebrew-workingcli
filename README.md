# homebrew-workingcli

Homebrew tap for Git Assistant (ga) - A smart Git workflow optimizer for large repositories

## 🍺 Installation

```bash
# Add the tap
brew tap yoonhwan-neowiz/workingcli

# Install ga
brew install ga
```

## 🔄 Update

```bash
# Update tap
brew update

# Upgrade ga to latest version
brew upgrade ga
```

## 🗑️ Uninstall

```bash
brew uninstall ga
brew untap yoonhwan-neowiz/workingcli
```

## 📦 What is Git Assistant (ga)?

Git Assistant is a powerful CLI tool designed to optimize and streamline Git operations for large-scale repositories. It provides:

- **Smart Repository Optimization**: Reduce 79GB+ repositories to manageable sizes (30MB)
- **AI-Powered Commit Messages**: Generate conventional commit messages automatically
- **Interactive Staging**: User-friendly file selection interface
- **Performance Workflows**: Optimized commands for faster Git operations

### Quick Start

```bash
# Interactive staging
ga

# AI commit message generation
ga commit

# Check repository optimization status
ga opt quick status

# Optimize large repository
ga opt quick to-slim

# Get help
ga --help
ga opt help workflow
```

## 🔧 Development

### Building from Source

If you want to build ga from source:

```bash
# Clone the main repository
git clone https://github.com/yoonhwan-neowiz/WorkingCli
cd WorkingCli

# Build
./build.command
```

### Creating a New Release

For maintainers, to create a new release:

```bash
# Run the release script with version number
cd scripts
./release.sh 0.1.1

# Follow the instructions to:
# 1. Push Formula changes
# 2. Create GitHub release
# 3. Upload binary archives
```

### Testing Formula Locally

```bash
# Build test archives
cd scripts
./build-release.sh

# Install from local Formula
brew install --build-from-source ./Formula/ga.rb
```

## 📄 Formula Details

The Formula supports:
- macOS (Intel & Apple Silicon)
- Linux (amd64 & arm64)
- Automatic platform detection
- SHA256 verification

## 🐛 Issues & Support

- Main Project: [WorkingCli Issues](https://github.com/yoonhwan-neowiz/WorkingCli/issues)
- Tap Issues: [homebrew-workingcli Issues](https://github.com/yoonhwan-neowiz/homebrew-workingcli/issues)

## 📜 License

MIT License - See [LICENSE](LICENSE) for details

---

<div align="center">
  
**Git Assistant (ga)** - Making Git workflows smarter and faster 🚀

</div>