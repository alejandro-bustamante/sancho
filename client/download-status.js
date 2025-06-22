function downloadStatusApp() {
  return {
    downloadId: null,
    status: "Cargando estado...",

    initStatusCheck() {
      this.downloadId = this.getDownloadId();

      if (this.downloadId) {
        this.updateStatus(this.downloadId);
      }
    },

    getDownloadId() {
      const urlParams = new URLSearchParams(window.location.search);
      return urlParams.get("downloadId");
    },

    updateStatus(downloadId) {
      const proxyUrl = "http://localhost:8081/get-status";
      const url = `${proxyUrl}?downloadId=${encodeURIComponent(downloadId)}`;

      fetch(url)
        .then((response) => response.text())
        .then((status) => {
          this.status = status;

          if (
            status !== "Descarga completada" &&
            status !== "Error en la descarga"
          ) {
            setTimeout(() => this.updateStatus(downloadId), 5000); // Actualizar cada 5 segundos
          }
        })
        .catch((error) => {
          console.error("Error: ", error);
          this.status = "Error al obtener el estado de la descarga";
        });
    },
  };
}
