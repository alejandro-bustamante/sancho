<script lang="ts">
	import { selectedUser, notifications } from '$lib/stores/stores';
	import type { Track } from '$lib/stores/types';

	// Recibimos una pista como prop
	export let track: Track;

	// Función que inicia la descarga de una pista
	async function downloadTrack() {
		const id = 'starting_' + track.track_id;

		// Mostramos notificación de inicio
		notifications.update(n => [...n, {
			id,
			message: `Iniciando descarga: ${track.title}`,
			type: 'info',
			show: true
		}]);

		try {
			const res = await fetch('http://localhost:8081/downloads', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					id: track.track_id,
					isrc: track.isrc,
					user: $selectedUser,
					quality: 2,
				})
			});

			if (!res.ok) throw new Error(`Error ${res.status}`);
			const data = await res.json();

			// Reemplazamos la notificación temporal por una de éxito
			notifications.update(n => [
				...n.filter(noti => noti.id !== id),
				{
					id: data.downloadId,
					message: `${data.message}`,
					type: 'success',
					show: true
				}
			]);
		} catch (err: unknown) {
			const message = err instanceof Error ? err.message : 'Error desconocido';
			notifications.update(n => [
				...n.filter(noti => noti.id !== id),
				{
					id: Date.now().toString(),
					message: `Error al iniciar la descarga: ${message}`,
					type: 'error',
					show: true
				}
			]);
		}
	}
</script>

<!-- Tarjeta que muestra información de la pista -->
<div style="margin-bottom: 1rem; padding-bottom: 1rem; border-bottom: 1px solid #ccc;">
	<h3 style="font-size: 1.2rem; font-weight: bold; color: #2563eb;">{track.title}</h3>
	<p style="font-size: 0.9rem; color: #666;">Artista: {track.artist}</p>
	<p style="font-size: 0.9rem; color: #666;">Álbum: {track.album}</p>

	<!-- Imagen y preview -->
	<div style="display: flex; align-items: center; margin-top: 0.5rem;">
		<img src={track.image} alt={track.artist} style="border-radius: 0.25rem; margin-right: 1rem;" />
		<!-- <audio controls style="width: 100%;">
			<source src={track.preview} type="audio/mpeg" />
		</audio> -->
	</div>

	<!-- Botón de descarga -->
	<button on:click={downloadTrack} style="margin-top: 0.5rem; background-color: orange; color: white; padding: 0.5rem 1rem; border-radius: 0.25rem;">
		Descargar
	</button>
</div>
