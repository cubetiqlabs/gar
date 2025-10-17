# Troubleshooting Guide

## Common Issues and Solutions

### ❌ "illegal file path" Error

**Problem:**
```bash
$ gar -xvf archive.zip
Error extracting file.txt: illegal file path: file.txt
```

**Solution:**
This error indicates a path validation issue. Try these steps:

1. **Extract to a specific directory:**
   ```bash
   mkdir extracted
   gar -xvf archive.zip extracted
   ```

2. **Verify archive is valid:**
   ```bash
   gar -tvf archive.zip
   ```

3. **Check file permissions:**
   ```bash
   ls -la
   chmod 755 .
   gar -xvf archive.zip
   ```

**Root Cause:** May occur when extracting to current directory or if there are permission issues.

**Status:** Fixed in v1.0.1 - Always extract to a named directory for best results.

---

### ❌ "Permission denied" Error

**Problem:**
```bash
$ gar -cvf archive.zip /etc/shadow
Error: permission denied: /etc/shadow
```

**Solution:**

1. **Run with appropriate permissions:**
   ```bash
   sudo gar -cvf archive.zip /etc/shadow
   ```

2. **Or, archive readable files only:**
   ```bash
   gar -cvf archive.zip ~/my-documents
   ```

3. **Check file permissions:**
   ```bash
   ls -la /path/to/files
   ```

**Root Cause:** Insufficient permissions to read files.

---

### ❌ "Archive not found" Error

**Problem:**
```bash
$ gar -xvf missing-archive.zip
Error: open missing-archive.zip: no such file or directory
```

**Solution:**

1. **Verify archive exists:**
   ```bash
   ls -la archive.zip
   ```

2. **Check full path:**
   ```bash
   gar -tvf /full/path/to/archive.zip
   ```

3. **Use absolute path in scripts:**
   ```bash
   gar -xvf $(pwd)/archive.zip
   ```

**Root Cause:** File doesn't exist or wrong path specified.

---

### ❌ "No such file or directory" Error

**Problem:**
```bash
$ gar -cvf archive.zip nonexistent-folder
Error: stat nonexistent-folder: no such file or directory
```

**Solution:**

1. **Verify folder exists:**
   ```bash
   ls -la folder-name
   ```

2. **Use correct path:**
   ```bash
   gar -cvf archive.zip ./my-folder
   ```

3. **Check working directory:**
   ```bash
   pwd
   ls -la
   ```

**Root Cause:** Input folder/file doesn't exist.

---

### ⚠️ Archive Seems Empty After Extraction

**Problem:**
```bash
$ gar -xvf archive.zip extracted/
$ ls extracted/
(no output - looks empty)
```

**Solution:**

1. **Check hidden files:**
   ```bash
   ls -la extracted/
   ```

2. **Verify archive isn't empty:**
   ```bash
   gar -tvf archive.zip
   ```

3. **Try extracting to current directory:**
   ```bash
   cd extracted
   gar -xvf ../archive.zip
   ```

**Root Cause:** Files may be hidden (start with `.`) or extraction didn't create visible files.

---

### ⚠️ Very Slow Compression/Extraction

**Problem:**
```bash
$ gar -cvf large-archive.zip huge-folder
(takes very long time)
```

**Solution:**

1. **Check system resources:**
   ```bash
   top
   df -h  # Check disk space
   ```

2. **Increase worker threads (extraction only):**
   ```bash
   gar -xvf large-archive.zip -workers=32
   ```

3. **Use fastest compression:**
   ```bash
   gar -cvf archive.zip folder -compression=fastest
   ```

4. **Check for slow disk:**
   ```bash
   # Compare with another tool
   tar -czf archive.tar.gz folder
   ```

**Root Cause:** Large files, slow disk, or system under load.

---

### ⚠️ Archive File is Too Large

**Problem:**
```bash
$ gar -cvf archive.zip large-folder
$ ls -lh archive.zip
10G archive.zip  # Very large!
```

**Solution:**

1. **Use better compression:**
   ```bash
   gar -cvf archive.zip folder -compression=best
   ```

2. **Use TAR.GZ instead of ZIP:**
   ```bash
   gar -cvfz archive.tar.gz folder
   ```

3. **Exclude unnecessary files:**
   ```bash
   cd folder
   rm -rf node_modules/ .git/
   cd ..
   gar -cvf archive.zip folder
   ```

4. **Split across multiple archives:**
   ```bash
   gar -cvf part1.zip folder/documents/
   gar -cvf part2.zip folder/media/
   ```

**Root Cause:** Uncompressed data, includes unnecessary files, or poor compression algorithm choice.

---

### 🔐 Password Protection Issues

#### Password Not Working

**Problem:**
```bash
$ gar -xvf secure.zip -password="MyPassword"
Error: decryption failed - possibly wrong password
```

**Solution:**

1. **Verify password:**
   - Double-check spelling and case sensitivity
   - Passwords are case-sensitive!

2. **Use environment variable:**
   ```bash
   export GAR_PASSWORD="CorrectPassword"
   gar -xvf secure.zip -password="$GAR_PASSWORD"
   ```

3. **Try without password (if not encrypted):**
   ```bash
   gar -xvf archive.zip  # Without password flag
   ```

**Root Cause:** Incorrect password or archive not encrypted.

#### Forgot Password

**Problem:**
```bash
$ gar -xvf secure.zip -password="WrongPassword"
Error: decryption failed
```

**Solution:**

❌ **Unfortunately:** Encrypted archives with forgotten passwords cannot be recovered.

**Prevention:**
- Store passwords securely (password manager)
- Keep copies of important archives
- Use `GAR_PASSWORD` environment variable

---

### 🔄 Format-Related Issues

#### Can't Extract TAR.GZ File

**Problem:**
```bash
$ gar -xvf archive.tar.gz
Error: gzip header
```

**Solution:**

1. **Explicitly specify format:**
   ```bash
   gar -xvf archive.tar.gz -format=targz
   ```

2. **Verify file is actually TAR.GZ:**
   ```bash
   file archive.tar.gz
   ```

3. **Try with traditional flags:**
   ```bash
   gar -action=extract -input=archive.tar.gz
   ```

**Root Cause:** File may have wrong extension or be corrupted.

#### Mixed Format Archives

**Problem:**
```bash
$ gar -tvf archive.zip
Error: cannot list mixed format archive
```

**Solution:**

1. **Determine actual format:**
   ```bash
   file archive.zip
   ```

2. **Extract with correct format:**
   ```bash
   # If it's actually TAR.GZ
   tar -xzf archive.zip
   ```

**Root Cause:** File extension doesn't match actual format.

---

## ✅ Verification Steps

### After Compression
```bash
# 1. Verify archive was created
ls -lh archive.zip

# 2. Check contents
gar -tvf archive.zip | head

# 3. Calculate checksum
md5sum archive.zip > archive.zip.md5

# 4. Test extraction
mkdir test-extract
gar -xvf archive.zip test-extract
```

### After Extraction
```bash
# 1. Verify files exist
ls -la extracted/

# 2. Check file count
find extracted/ -type f | wc -l

# 3. Verify permissions
ls -la extracted/important-file

# 4. Spot check content
head extracted/sample-file.txt
```

## 📞 Getting Help

### Debug Information to Collect

```bash
# System information
uname -a
go version

# GoArchive version
gar -help

# Reproduce the issue
# (save commands for later reporting)

# Check logs
echo "Error details here"
```

### Report an Issue

When reporting problems, include:

1. **Exact command that failed:**
   ```bash
   gar -cvf archive.zip ~/my-folder
   ```

2. **Full error message:**
   ```
   Error: permission denied
   ```

3. **System info:**
   - OS: macOS 13.4 / Ubuntu 22.04 / Windows 11
   - GoArchive version: 1.0.1
   - Go version: 1.21

4. **Steps to reproduce:**
   ```
   1. Create folder with files
   2. Run compression command
   3. Error occurs
   ```

5. **What you expected:**
   ```
   Archive should be created successfully
   ```

## 🎓 Useful Resources

- **Quick Start:** See docs/QUICK_START.md
- **Examples:** See docs/EXAMPLES.md
- **Implementation:** See docs/IMPLEMENTATION.md
- **Bug Fixes:** See docs/BUG_FIXES.md
- **Main README:** See README.md

## 💡 Prevention Tips

### Best Practices

```bash
# 1. Always verify archive after creation
gar -cvf backup.zip data/
gar -tvf backup.zip

# 2. Test extraction before deleting source
gar -xvf backup.zip test-extract/
diff -r data/ test-extract/data/

# 3. Keep backups of important archives
cp archive.zip archive-backup.zip

# 4. Use consistent file paths
cd ~/projects
gar -cvf myproject.zip myproject/  # Relative path
# NOT: gar -cvf myproject.zip ~/projects/myproject/

# 5. Document your backups
echo "Backup created on $(date)" > BACKUP.txt
gar -tvf backup.zip >> BACKUP.txt
```

### Scripting Tips

```bash
#!/bin/bash

# Safer scripting practices
set -euo pipefail

ARCHIVE="${1:-backup.zip}"
INPUT_DIR="${2:-.}"

# Verify inputs
[ -d "$INPUT_DIR" ] || { echo "Directory not found"; exit 1; }
[ -f "$ARCHIVE" ] && { echo "Archive exists, backing up"; mv "$ARCHIVE" "$ARCHIVE.old"; }

# Create archive
gar -cvf "$ARCHIVE" "$INPUT_DIR" || { echo "Compression failed"; exit 1; }

# Verify result
gar -tvf "$ARCHIVE" > /dev/null || { echo "Archive corrupt"; exit 1; }

echo "✅ Backup completed: $ARCHIVE"
```

---

**Still having issues?** Check the main README.md or open an issue on GitHub!
