# Examples & Use Cases

## 🚀 Quick Examples

### Basic Compress
```bash
gar -cvf my-archive.zip my-folder
```

**Output:**
```
Compressing my-folder to my-archive.zip...
  Adding: file1.txt
  Adding: file2.txt
  Adding: subdir/file3.txt
Operation completed in 2.491084ms
```

### Basic Extract
```bash
gar -xvf my-archive.zip
```

**Output:**
```
Extracting my-archive.zip to .../
  Extracting: file1.txt
  Extracting: file2.txt
  Extracting: subdir/file3.txt
Operation completed in 1.186875ms
```

### List Contents
```bash
gar -tvf my-archive.zip
```

**Output:**
```
Archive contents:
  file1.txt (12 bytes)
  file2.txt (12 bytes)
  subdir/ (0 bytes)
  subdir/file3.txt (12 bytes)
Operation completed in 307.708µs
```

## 📦 Format Examples

### TAR.GZ Compression
```bash
gar -cvfz my-archive.tar.gz my-folder
gar -xvf my-archive.tar.gz
```

### bzip2 Compression
```bash
gar -cvfj my-archive.tar.bz2 my-folder
gar -xvf my-archive.tar.bz2
```

### 7-zip Compression
```bash
gar -cvfZ my-archive.7z my-folder
gar -xvf my-archive.7z
```

## 🔒 Encryption Examples

### Create Encrypted Archive
```bash
gar -cvf secure.zip my-folder -password="MySecurePassword123"
```

### Extract Encrypted Archive
```bash
gar -xvf secure.zip -password="MySecurePassword123"
```

### Using Environment Variables (Recommended)
```bash
export GAR_PASSWORD="MySecurePassword123"
gar -cvf secure.zip my-folder -password="$GAR_PASSWORD"
gar -xvf secure.zip -password="$GAR_PASSWORD"
```

## ⚡ Performance Examples

### Multi-threaded Extraction
```bash
# Use 16 worker threads for faster extraction
gar -xvf large-archive.zip -workers=16
```

### Compression Levels
```bash
# Fastest compression (larger file)
gar -cvf archive.zip folder -compression=fastest

# Balanced compression (default)
gar -cvf archive.zip folder -compression=normal

# Best compression (slower)
gar -cvf archive.zip folder -compression=best
```

## 📂 Common Scenarios

### Project Backup
```bash
# Create timestamped backup
gar -cvf project-backup-$(date +%Y%m%d).zip ~/my-project

# Extract when needed
gar -xvf project-backup-20250117.zip ~/restored-project
```

### Log File Archive
```bash
# Compress logs in TAR.GZ with best compression
gar -cvfz logs-archive.tar.gz /var/log -compression=best

# Extract to temp for inspection
mkdir /tmp/logs-extracted
gar -xvf logs-archive.tar.gz /tmp/logs-extracted
```

### Software Distribution
```bash
# Create encrypted archive for secure distribution
gar -cvf software-release.zip build/ docs/ README.md -password="release123"

# Recipients extract securely
gar -xvf software-release.zip -password="release123"
```

### Cross-Format Conversion
```bash
# Extract from RAR (old system)
gar -xvf data.rar extracted-files

# Re-archive in modern TAR.GZ format
gar -cvfz data.tar.gz extracted-files

# Clean up
rm -rf extracted-files
```

### Selective Listing
```bash
# List archive contents
gar -tvf archive.zip

# Find specific files in archive
gar -tvf archive.zip | grep -i ".jpg"

# Count files
gar -tvf archive.zip | wc -l
```

## 🔄 Comparison: Old vs New Syntax

### Compress
```bash
# New Unix-style (short & clean)
gar -cvf archive.zip my-folder

# Old traditional (verbose but explicit)
gar -action=compress -input=my-folder -output=archive.zip -verbose
```

### Extract
```bash
# New Unix-style
gar -xvf archive.zip output-dir

# Old traditional
gar -action=extract -input=archive.zip -output=output-dir -verbose
```

### List
```bash
# New Unix-style
gar -tvf archive.zip

# Old traditional
gar -action=list -input=archive.zip -verbose
```

## 🛠️ Advanced Usage

### Shell Scripts
```bash
#!/bin/bash

# Backup script
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_NAME="backup_${TIMESTAMP}.zip"

echo "Creating backup: $BACKUP_NAME"
gar -cvf "$BACKUP_NAME" ~/important-data

echo "Backup created: $(ls -lh $BACKUP_NAME | awk '{print $5}')"
```

### Cron Job
```bash
# Daily backup (add to crontab)
0 2 * * * /usr/local/bin/gar -cvfz /backups/daily-$(date +\%Y\%m\%d).tar.gz /home/data
```

### Pipeline Usage
```bash
# Find and archive files modified in last 7 days
find ~/projects -mtime -7 -type f | xargs gar -cvf recent-changes.zip

# List and filter
gar -tvf archive.zip | grep -E "\.(go|py|js)$"
```

## 💾 Real-World Examples

### Developer Workflow
```bash
# Before committing code, create backup
gar -cvf pre-commit-backup-$(date +%Y%m%d-%H%M%S).zip src/ tests/

# Later, extract if needed
gar -xvf pre-commit-backup-20250117-143022.zip
```

### Data Migration
```bash
# Export data
gar -cvfz export-$(date +%Y%m%d).tar.gz database-dump/

# Transfer to new system
scp export-*.tar.gz user@newhost:/tmp/

# Extract on new system
ssh user@newhost "gar -xvf /tmp/export-*.tar.gz"
```

### Archive Verification
```bash
# Create and verify backup
gar -cvf backup.zip important-data
gar -tvf backup.zip > backup-manifest.txt

# Later, verify backup is intact
gar -tvf backup.zip > backup-manifest-verify.txt
diff backup-manifest.txt backup-manifest-verify.txt
```

## ✨ Tips & Tricks

### Quick Check Archive Size
```bash
gar -tvf archive.zip | tail -1
```

### Extract Without Creating Directory
```bash
# All files in current directory
gar -xvf archive.zip

# With specific output
gar -xvf archive.zip .
```

### Create Multiple Archives from Pattern
```bash
for dir in project_*/; do
  gar -cvf "${dir%/}.zip" "$dir"
done
```

### Parallel Archive Creation
```bash
# Process multiple folders in parallel
find ~/projects -maxdepth 1 -type d | \
  xargs -I {} -P 4 gar -cvf {}.zip {}
```
