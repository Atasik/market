document.addEventListener('DOMContentLoaded', function () {
    var dragDropArea = document.getElementById('drag-drop-area');
    var fileInput = document.getElementById('file');
    var filePreview = document.getElementById('image-container');
  
    function previewFile(file) {
      var reader = new FileReader();
      reader.onload = function (e) {
        var img = document.createElement('img');
        img.src = e.target.result;
        img.classList.add('img-fluid');
        filePreview.innerHTML = '';
        filePreview.appendChild(img);
        var closeButton = document.createElement('span');
        closeButton.classList.add('close-button');
        closeButton.innerHTML = '&times;';
        closeButton.addEventListener('click', function () {
          filePreview.innerHTML = '';
          filePreview.appendChild(dragDropArea);
          fileInput.value = '';
          dragDropArea.style.display = 'flex';
        });
        filePreview.appendChild(closeButton);
      };
      reader.readAsDataURL(file);
      dragDropArea.style.display = 'none';
    }
  
    function handleDrop(e) {
      e.preventDefault();
      dragDropArea.classList.remove('highlight');
      var files = e.dataTransfer.files;
      fileInput.files = files;
      if (files.length > 0) {
        previewFile(files[0]);
      }
    }
  
    function handleFileChange() {
      var files = fileInput.files;
      if (files.length > 0) {
        previewFile(files[0]);
      }
    }
  
    dragDropArea.addEventListener('dragover', function (e) {
      e.preventDefault();
      dragDropArea.classList.add('highlight');
    });
  
    dragDropArea.addEventListener('dragleave', function (e) {
      e.preventDefault();
      dragDropArea.classList.remove('highlight');
    });
  
    dragDropArea.addEventListener('drop', handleDrop);
  
    dragDropArea.addEventListener('click', function () {
      fileInput.click();
    });
  
    fileInput.addEventListener('change', handleFileChange);
  });