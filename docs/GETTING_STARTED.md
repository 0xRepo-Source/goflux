# Getting Started with goflux

A super simple guide to get up and running in 5 minutes! 🚀

## What is goflux?

goflux lets you upload and download files to a server. Think of it like a personal Dropbox you run yourself!

## Quick Start (3 Steps!)

### Step 1: Start the Server

Open a terminal and run:

```bash
.\bin\goflux-server.exe --storage ./myfiles
```

That's it! Your server is running on `http://localhost:8080` 🎉

### Step 2: Upload a File

Open **another** terminal (keep the server running!) and upload a file:

```bash
.\bin\goflux.exe put myfile.txt /myfile.txt
```

This copies `myfile.txt` from your computer to the server.

### Step 3: Download it Back

Download the file to make sure it worked:

```bash
.\bin\goflux.exe get /myfile.txt downloaded.txt
```

Now you have `downloaded.txt` with the same content! ✨

## Common Commands

### Upload a file
```bash
.\bin\goflux.exe put <local-file> <server-path>
```

**Example:**
```bash
.\bin\goflux.exe put photo.jpg /photos/vacation.jpg
```

**You'll see a nice progress bar:**
```
Uploading photo.jpg (2450000 bytes, 3 chunks)...
Uploading... 100% [========================================] (1.2 MB/s)
✓ Upload complete: photo.jpg → /photos/vacation.jpg
```

### Download a file
```bash
.\bin\goflux.exe get <server-path> <local-file>
```

**Example:**
```bash
.\bin\goflux.exe get /photos/vacation.jpg my-photo.jpg
```

### See what's on the server
```bash
.\bin\goflux.exe ls /
```

**Example:**
```bash
.\bin\goflux.exe ls /photos
```

## 🔄 What if the Upload Fails?

Don't worry! goflux automatically resumes interrupted uploads.

If your upload is interrupted (Wi-Fi drops, computer crashes, etc.), just run the same command again:

```bash
# Upload starts...
.\bin\goflux.exe put bigfile.zip /bigfile.zip
# ... oops, internet disconnects! ...

# Just run it again - it will resume automatically!
.\bin\goflux.exe put bigfile.zip /bigfile.zip
```

**What you'll see:**
```
Uploading bigfile.zip (50000000 bytes, 48 chunks)...
🔄 Resuming upload: 25/48 chunks already uploaded
  Uploaded chunk 26/48
  ... continues from where it left off ...
```

**Magic!** ✨ It remembers what was already uploaded and only sends the missing parts.

## Using the Web Browser UI

Don't want to use the terminal? No problem!

1. Start the server (same as Step 1 above)
2. Open your web browser
3. Go to: `http://localhost:8080`
4. Drag and drop files to upload! 📂

## Adding Security (Optional)

Want to require passwords? Here's how:

### Step 1: Create a Token

```bash
.\bin\goflux-admin.exe create --user yourname --permissions upload,download,list
```

You'll get something like:
```
Token: abc123def456...
```

**IMPORTANT:** Copy this token! You won't see it again.

### Step 2: Start Server with Authentication

```bash
.\bin\goflux-server.exe --storage ./myfiles --tokens tokens.json
```

### Step 3: Use Your Token

Now you need the token for every command:

```bash
.\bin\goflux.exe --token "abc123def456..." put myfile.txt /myfile.txt
```

**Pro tip:** Set it once and forget it:
```bash
$env:GOFLUX_TOKEN = "abc123def456..."
.\bin\goflux.exe put myfile.txt /myfile.txt
```

## Troubleshooting

### "Server failed to start"
- Make sure nothing else is using port 8080
- Try a different port: `.\bin\goflux-server.exe --addr :9000 --storage ./myfiles`

### "Connection refused"
- Is the server running? Check the first terminal window
- Are you using the right port? Default is 8080

### "Authentication failed"
- Did you forget the `--token`?
- Is your token correct? (no extra spaces!)
- Try setting `$env:GOFLUX_TOKEN` instead

### "File not found"
- Use forward slashes: `/photos/file.jpg` not `\photos\file.jpg`
- Server paths start with `/`

## Full Example Session

Here's a complete example of uploading vacation photos:

```bash
# Terminal 1: Start the server
.\bin\goflux-server.exe --storage ./myfiles

# Terminal 2: Upload some photos
.\bin\goflux.exe put beach.jpg /vacation/beach.jpg
.\bin\goflux.exe put sunset.jpg /vacation/sunset.jpg
.\bin\goflux.exe put family.jpg /vacation/family.jpg

# List them to make sure they're there
.\bin\goflux.exe ls /vacation

# Download one back
.\bin\goflux.exe get /vacation/beach.jpg beach-copy.jpg
```

## What Next?

- **Share with friends:** Give them your server address (like `http://192.168.1.100:8080`) and a token
- **Use from anywhere:** Deploy the server to a cloud provider
- **Organize files:** Create folders like `/photos`, `/documents`, `/backups`
- **Stay secure:** Always use tokens when the server is accessible from the internet!

## Cheat Sheet

| What do you want to do? | Command |
|-------------------------|---------|
| Start server | `.\bin\goflux-server.exe --storage ./myfiles` |
| Upload file | `.\bin\goflux.exe put file.txt /file.txt` |
| Download file | `.\bin\goflux.exe get /file.txt file.txt` |
| List files | `.\bin\goflux.exe ls /` |
| Use web UI | Open browser to `http://localhost:8080` |
| Create token | `.\bin\goflux-admin.exe create --user yourname` |
| Stop server | Press `Ctrl+C` in the server terminal |

## Need Help?

- 📖 Full documentation: [README.md](../README.md)
- 🔒 Security guide: [AUTHENTICATION.md](AUTHENTICATION.md)
- 🏗️ Architecture details: [architecture.md](architecture.md)

---

**That's it!** You're now a goflux expert! 🎓
