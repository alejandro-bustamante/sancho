function searchApp() {
  return {
    trackQuery: "",
    tracks: [],
    isLoading: false,
    selectedUser: "alejandro",
    notifications: [],

    searchTracks() {
      const query = this.trackQuery.trim();
      if (!query) return;
      this.isLoading = true;
      this.tracks = [];
      const proxyUrl = "http://localhost:8081/proxy";
      const apiUrl = `https://api.deezer.com/search?q=${encodeURIComponent(
        query
      )}`;
      const url = `${proxyUrl}?url=${encodeURIComponent(apiUrl)}`;
      fetch(url)
        .then((response) => response.json())
        .then((data) => {
          this.tracks = data.data;
          this.isLoading = false;
        })
        .catch((error) => {
          console.error("Error: ", error);
          this.isLoading = false;
          this.showNotification("Error al buscar canciones", "error");
        });
    },

    downloadTrack(track) {
      // Mostrar notificación de "iniciando descarga"
      this.showNotification(
        `Iniciando descarga: ${track.title}...`,
        "info",
        "starting_" + track.id
      );

      fetch("http://localhost:8081/download", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          url: track.link,
          title: track.title,
          artist: track.artist.name,
          album: track.album.title,
          user: this.selectedUser,
        }),
      })
        .then((response) => {
          if (!response.ok) {
            throw new Error(`Error en la respuesta: ${response.status}`);
          }
          return response.json();
        })
        .then((data) => {
          // Ocultar notificación de "iniciando"
          this.hideNotification("starting_" + track.id);

          // Mostrar notificación de éxito con la información detallada
          this.showNotification(
            `Descarga iniciada: ${data.trackInfo.title} - ${data.trackInfo.artist}`,
            "success",
            data.downloadId
          );

          console.log("Descarga iniciada:", data);
        })
        .catch((error) => {
          // Ocultar notificación de "iniciando"
          this.hideNotification("starting_" + track.id);

          console.error("Error calling download in the backend", error);
          this.showNotification(
            "Error al iniciar la descarga: " + error.message,
            "error"
          );
        });
    },

    // Notification system
    showNotification(message, type = "info", id = null) {
      const notificationId = id || Date.now().toString();

      // Comprobar si ya existe una notificación con este ID
      const existingIndex = this.notifications.findIndex(
        (n) => n.id === notificationId
      );
      if (existingIndex !== -1) {
        // Actualizar la notificación existente
        this.notifications[existingIndex].message = message;
        this.notifications[existingIndex].type = type;
        return;
      }

      const notification = {
        id: notificationId,
        message,
        type,
        show: true,
      };

      this.notifications.push(notification);

      // Auto-hide after 5 seconds for success and error messages
      if (type !== "info") {
        setTimeout(() => {
          this.hideNotification(notification.id);
        }, 5000);
      }
    },

    hideNotification(id) {
      const index = this.notifications.findIndex((n) => n.id === id);
      if (index !== -1) {
        this.notifications.splice(index, 1);
      }
    },
  };
}
