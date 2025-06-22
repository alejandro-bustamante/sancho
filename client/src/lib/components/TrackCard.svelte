<script lang="ts">
	import { selectedUser, notifications } from '$lib/stores';
	import type { Track } from '$lib/types';

	export let track: Track;

	async function downloadTrack() {
		const id = 'starting_' + track.id;

		notifications.update((n) => [
			...n,
			{ id, message: `Iniciando descarga: ${track.title}`, type: 'info', show: true }
		]);

		try {
			const res = await fetch('http://localhost:8081/download', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					// url: track.link,
					title: track.title,
					artist: track.artist.name,
					album: track.album.title,
					user: $selectedUser
				})
			});

			if (!res.ok) throw new Error(`Error ${res.status}`);
			const data = await res.json();

			notifications.update((n) => [
				...n.filter((noti) => noti.id !== id),
				{
					id: data.downloadId,
					message: `Descarga iniciada: ${data.trackInfo.title} - ${data.trackInfo.artist}`,
					type: 'success',
					show: true
				}
			]);
		} catch (err: unknown) {
			const message = err instanceof Error ? err.message : 'Error desconocido';
			notifications.update((n) => [
				...n.filter((noti) => noti.id !== id),
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

<div class="mb-6 pb-6 border-b border-gray-200">
	<h3 class="text-xl font-semibold text-blue-600 mb-2">{track.title}</h3>
	<!-- <p class="text-sm text-gray-600">Link: {track.link}</p> -->
	<p class="text-sm text-gray-600">Artist: {track.artist.name}</p>
	<p class="text-sm text-gray-600">Album: {track.album.title}</p>
	<div class="flex flex-wrap items-center mt-3">
		<img src={track.album.cover_small} alt={track.album.title} class="rounded mr-4" />
		<audio controls class="mt-2 w-full">
			<source src={track.preview} type="audio/mpeg" />
		</audio>
	</div>
	<button
		on:click={downloadTrack}
		class="mt-3 bg-orange hover:bg-orange/80 text-white px-3 py-2 rounded text-sm transition-colors duration-300"
	>
		Descargar
	</button>
</div>
