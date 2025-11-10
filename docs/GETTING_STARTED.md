# Getting Started with goflux

A super simple guide to get up and running in 5 minutes! 🚀

## What is goflux?

goflux lets you upload and download files to a server. Think of it like a personal Dropbox you run yourself!

## Quick Start (3 Steps!)

### Step 1: Start the Server

Open a terminal and run:

```bash
.\bin\goflux-server.exe
```

That's it! Your server is running on `http://localhost` (port 80) 🎉

The server uses `goflux.json` for configuration (created automatically if missing).

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
3. Go to: `http://localhost`
4. Drag and drop files to upload! 📂

## Adding Security (Optional)

Want to require passwords? Here's how:

### Step 1: Create a Token

```bash
.\bin\goflux-admin.exe create --user yourname --permissions upload,download,list
```

You'll get something like:
```
Token: tok_abc123def456...
```

**IMPORTANT:** Copy this token! You won't see it again.

### Step 2: Configure Server with Authentication

Edit `goflux.json` to enable authentication:

```json
{
  "server": {
    "tokens_file": "tokens.json"
  }
}
```

Then start the server:

```bash
.\bin\goflux-server.exe
```

### Step 3: Configure Client with Token

Edit your `goflux.json` to add the token:

```json
{
  "client": {
    "token": "tok_abc123def456..."
  }
}
```

Or use environment variable:

```bash
$env:GOFLUX_TOKEN = "tok_abc123def456..."
.\bin\goflux.exe put myfile.txt /myfile.txt
```

## Troubleshooting

### "Server failed to start"
- Make sure nothing else is using port 80
- Edit `goflux.json` to change the port in `server.address`

### "Connection refused"
- Is the server running? Check the first terminal window
- Are you using the right port? Default is 80
- Check `server_url` in your client config

### "Authentication failed"
- Did you set the token in config or environment variable?
- Is your token correct? (no extra spaces!)
- Make sure server has `tokens_file` configured

### "File not found"
- Use forward slashes: `/photos/file.jpg` not `\photos\file.jpg`
- Server paths start with `/`

## Configuration Files

### Use Different Configs for Different Environments

**Development (goflux-dev.json):**
```json
{
  "client": {
    "server_url": "http://localhost"
  }
}
```

**Production (goflux-prod.json):**
```json
{
  "client": {
    "server_url": "http://myserver.com",
    "token": "tok_production_token"
  }
}
```

Then use:
```bash
.\bin\goflux.exe --config goflux-prod.json put file.txt /file.txt
```

## Full Example Session

Here's a complete example of uploading vacation photos:

```bash
# Terminal 1: Start the server
.\bin\goflux-server.exe

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

- **Share with friends:** Give them your server address and a token, or share your config file
- **Use from anywhere:** Deploy the server to a cloud provider
- **Organize files:** Create folders like `/photos`, `/documents`, `/backups`
- **Stay secure:** Always use tokens when the server is accessible from the internet!
- **Multiple environments:** Create different config files for dev/staging/prod

## Cheat Sheet

| What do you want to do? | Command |
|-------------------------|---------|
| Start server | `.\bin\goflux-server.exe` |
| Upload file | `.\bin\goflux.exe put file.txt /file.txt` |
| Download file | `.\bin\goflux.exe get /file.txt file.txt` |
| List files | `.\bin\goflux.exe ls /` |
| Use different config | `.\bin\goflux.exe --config prod.json ls` |
| Use web UI | Open browser to `http://localhost` |
| Create token | `.\bin\goflux-admin.exe create --user yourname` |
| Stop server | Press `Ctrl+C` in the server terminal |

## Need Help?

- 📖 Full documentation: [README.md](../README.md)
- 🔒 Security guide: [AUTHENTICATION.md](AUTHENTICATION.md)
- ⚙️ Configuration guide: [CONFIGURATION.md](CONFIGURATION.md)
- 🏗️ Architecture details: [architecture.md](architecture.md)

---

**That's it!** You're now a goflux expert! 🎓
