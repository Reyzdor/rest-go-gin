let selectedFiles = [];
        
        function handleFiles(files) {
            for (let file of files) {
                selectedFiles.push(file);
            }
            updateFileList();
        }
        
        function updateFileList() {
            const fileList = document.getElementById('file-list');
            fileList.innerHTML = '';
            
            if (selectedFiles.length === 0) {
                fileList.innerHTML = '<p>Нет выбранных файлов</p>';
                return;
            }
            
            selectedFiles.forEach((file, index) => {
                const reader = new FileReader();
                reader.onload = function(e) {
                    const div = document.createElement('div');
                    div.className = 'file-item';
                    div.innerHTML = `
                        <div style="display: flex; align-items: center;">
                            <img src="${e.target.result}" class="file-preview">
                            <div>
                                <strong>${file.name}</strong><br>
                                <small>${Math.round(file.size / 1024)} KB</small>
                                ${index === 0 ? '<br><small style="color:#872BFF;">★ Главное фото</small>' : ''}
                            </div>
                        </div>
                        <button type="button" onclick="removeFile(${index})" 
                                style="background:#ff4444; color:white; border:none; padding:5px 10px; border-radius:3px; cursor:pointer;">
                            Удалить
                        </button>
                    `;
                    fileList.appendChild(div);
                };
                reader.readAsDataURL(file);
            });
        }
        
        function removeFile(index) {
            selectedFiles.splice(index, 1);
            updateFileList();
        }
        
        const dropArea = document.querySelector('.file-upload-area');
        
        ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
            dropArea.addEventListener(eventName, preventDefaults, false);
        });
        
        function preventDefaults(e) {
            e.preventDefault();
            e.stopPropagation();
        }
        
        ['dragenter', 'dragover'].forEach(eventName => {
            dropArea.addEventListener(eventName, highlight, false);
        });
        
        ['dragleave', 'drop'].forEach(eventName => {
            dropArea.addEventListener(eventName, unhighlight, false);
        });
        
        function highlight() {
            dropArea.style.background = '#e9e1ff';
        }
        
        function unhighlight() {
            dropArea.style.background = '#f8f5ff';
        }
        
        dropArea.addEventListener('drop', handleDrop, false);
        
        function handleDrop(e) {
            const dt = e.dataTransfer;
            const files = dt.files;
            handleFiles(files);
        }