// Global state
let currentPath = '/';
const CHUNK_SIZE = 1024 * 1024; // 1MB chunks

// Initialize on page load
document.addEventListener('DOMContentLoaded', () => {
    setupUpload();
    loadFiles(currentPath);
});

// Upload functionality
function setupUpload() {
    const dropZone = document.getElementById('dropZone');
    const fileInput = document.getElementById('fileInput');

    // Click to upload
    dropZone.addEventListener('click', () => {
        fileInput.click();
    });

    // File selection
    fileInput.addEventListener('change', (e) => {
        if (e.target.files.length > 0) {
            uploadFiles(e.target.files);
        }
    });

    // Drag and drop
    dropZone.addEventListener('dragover', (e) => {
        e.preventDefault();
        dropZone.classList.add('drag-over');
    });

    dropZone.addEventListener('dragleave', () => {
        dropZone.classList.remove('drag-over');
    });

    dropZone.addEventListener('drop', (e) => {
        e.preventDefault();
        dropZone.classList.remove('drag-over');
        
        if (e.dataTransfer.files.length > 0) {
            uploadFiles(e.dataTransfer.files);
        }
    });
}

async function uploadFiles(files) {
    for (const file of files) {
        await uploadFile(file);
    }
}

async function uploadFile(file) {
    const uploadProgress = document.getElementById('uploadProgress');
    const uploadFileName = document.getElementById('uploadFileName');
    const uploadProgressBar = document.getElementById('uploadProgressBar');
    const uploadProgressText = document.getElementById('uploadProgressText');

    uploadProgress.style.display = 'block';
    uploadFileName.textContent = file.name;
    uploadProgressBar.style.width = '0%';
    uploadProgressText.textContent = '0%';

    try {
        // Read file and split into chunks
        const chunks = await splitFileIntoChunks(file);
        const remotePath = currentPath + (currentPath.endsWith('/') ? '' : '/') + file.name;

        // Upload each chunk
        for (let i = 0; i < chunks.length; i++) {
            const chunk = chunks[i];
            
            const chunkData = {
                path: remotePath,
                chunk_id: i,
                data: Array.from(new Uint8Array(chunk.data)),
                checksum: chunk.checksum,
                total: chunks.length
            };

            const response = await fetch('/upload', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(chunkData)
            });

            if (!response.ok) {
                throw new Error(`Upload failed: ${response.statusText}`);
            }

            // Update progress
            const progress = Math.round(((i + 1) / chunks.length) * 100);
            uploadProgressBar.style.width = progress + '%';
            uploadProgressText.textContent = progress + '%';
        }

        // Success - reload file list
        showMessage('Upload complete: ' + file.name, 'success');
        setTimeout(() => {
            uploadProgress.style.display = 'none';
            loadFiles(currentPath);
        }, 1500);

    } catch (error) {
        showMessage('Upload failed: ' + error.message, 'error');
        console.error('Upload error:', error);
    }
}

async function splitFileIntoChunks(file) {
    const chunks = [];
    let offset = 0;

    while (offset < file.size) {
        const chunkEnd = Math.min(offset + CHUNK_SIZE, file.size);
        const blob = file.slice(offset, chunkEnd);
        const arrayBuffer = await blob.arrayBuffer();
        const checksum = await calculateSHA256(arrayBuffer);

        chunks.push({
            data: arrayBuffer,
            checksum: checksum
        });

        offset = chunkEnd;
    }

    return chunks;
}

async function calculateSHA256(buffer) {
    const hashBuffer = await crypto.subtle.digest('SHA-256', buffer);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    return hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
}

// File browsing
async function loadFiles(path) {
    currentPath = path;
    const fileList = document.getElementById('fileList');
    const currentPathEl = document.getElementById('currentPath');

    // Update breadcrumb
    if (path !== '/') {
        const pathParts = path.split('/').filter(p => p);
        currentPathEl.innerHTML = ' / ' + pathParts.join(' / ');
    } else {
        currentPathEl.innerHTML = '';
    }

    fileList.innerHTML = '<div class="loading">Loading files...</div>';

    try {
        const response = await fetch(`/list?path=${encodeURIComponent(path)}`);
        
        if (!response.ok) {
            throw new Error('Failed to load files');
        }

        const files = await response.json();
        displayFiles(files || []);

    } catch (error) {
        fileList.innerHTML = '<div class="error">Failed to load files: ' + error.message + '</div>';
    }
}

function displayFiles(files) {
    const fileList = document.getElementById('fileList');

    if (files.length === 0) {
        fileList.innerHTML = '<div class="loading">No files in this directory</div>';
        return;
    }

    fileList.innerHTML = '';

    files.forEach(fileName => {
        const fileItem = document.createElement('div');
        fileItem.className = 'file-item';

        const isDirectory = false; // We don't have directory info yet
        const icon = isDirectory ? '📁' : '📄';

        fileItem.innerHTML = `
            <div class="file-icon">${icon}</div>
            <div class="file-info">
                <div class="file-name-text">${escapeHtml(fileName)}</div>
            </div>
            <div class="file-actions">
                <button class="btn btn-primary btn-small" onclick="downloadFile('${escapeHtml(fileName)}')">
                    Download
                </button>
            </div>
        `;

        fileList.appendChild(fileItem);
    });
}

async function downloadFile(fileName) {
    const filePath = currentPath + (currentPath.endsWith('/') ? '' : '/') + fileName;
    
    try {
        const response = await fetch(`/download?path=${encodeURIComponent(filePath)}`);
        
        if (!response.ok) {
            throw new Error('Download failed');
        }

        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = fileName;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);

        showMessage('Downloaded: ' + fileName, 'success');

    } catch (error) {
        showMessage('Download failed: ' + error.message, 'error');
    }
}

function navigateTo(path) {
    loadFiles(path);
}

function showMessage(message, type) {
    const uploadProgress = document.getElementById('uploadProgress');
    const existingMsg = document.querySelector('.message');
    
    if (existingMsg) {
        existingMsg.remove();
    }

    const msgDiv = document.createElement('div');
    msgDiv.className = type === 'error' ? 'error message' : 'success message';
    msgDiv.textContent = message;
    
    uploadProgress.parentNode.insertBefore(msgDiv, uploadProgress);
    
    setTimeout(() => {
        msgDiv.remove();
    }, 3000);
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function formatBytes(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}
