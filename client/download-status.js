function getDownloadId() {
  const urlParams = new URLSearchParams(window.location.search);
  return urlParams.get('downloadId');
}

function updateStatus(downloadId) {
  // const proxyUrl = 'http://192.168.0.28:8081/get-status';
  // using localhost to avoid ip changes
  const proxyUrl = 'http://localhost:8081/get-status';
  const url = `${proxyUrl}?downloadId=${encodeURIComponent(downloadId)}`;

  fetch(url)
    .then(response => response.text())
    .then(status => {
      document.getElementById('downloadStatus').textContent = `Estado actual: ${status}`;

      if (status !== 'Descarga completada' && status !== 'Error en la descarga') {
        setTimeout(() => updateStatus(downloadId), 5000); // Actualizar cada 5 segundos
      }
    })
    .catch(error => {
      console.error('Error: ', error);
      document.getElementById('downloadStatus').textContent = 'Error al obtener el estado de la descarga';
    });
}

document.addEventListener('DOMContentLoaded', function () {
  const downloadId = getDownloadId();

  if (!downloadId) {
    document.getElementById('downloadStatus').textContent = 'Error: No se proporcionó un ID de descarga';
    return;
  }

  updateStatus(downloadId);
});
