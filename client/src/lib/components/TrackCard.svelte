<script lang="ts">
	import { selectedUser, notifications } from '$lib/stores/stores';
	import { currentUser, isClientReady } from '$lib/stores/auth';
	import type { Track } from '$lib/stores/types';
	import { onMount } from 'svelte';
	import { API_IP } from '$lib/config';

	export let track: Track;

	let showPlayer = false;
	let sampleUrl: string | null = null;
	let isLoadingSample = false;

	// Formatea duración (segundos → mm:ss)
	function formatDuration(seconds: number): string {
		const mins = Math.floor(seconds / 60);
		const secs = seconds % 60;
		return `${mins}:${secs.toString().padStart(2, '0')}`;
	}

	// Función que obtiene el sample
	async function loadSample() {
		if (sampleUrl || isLoadingSample) return;
		isLoadingSample = true;

		try {
			const res = await fetch(`${API_IP}/api/search/${track.isrc}/sample`);
			const data = await res.json();
			if (!res.ok || !data.sample_url) throw new Error(data.error || 'Error al obtener muestra');
			sampleUrl = data.sample_url;
			showPlayer = true;
		} catch (error) {
			notifications.update((n) => [
				...n,
				{
					id: Date.now().toString(),
					message: `No se pudo cargar el sample: ${error instanceof Error ? error.message : 'Error desconocido'}`,
					type: 'error',
					show: true
				}
			]);
		} finally {
			isLoadingSample = false;
		}
	}

	// Función que inicia la descarga
	async function downloadTrack() {
		const tempId = 'starting_' + track.track_id;
		if (!$isClientReady || !$selectedUser) {
			notifications.update((n) => [
				...n,
				{
					id: Date.now().toString(),
					message: 'No hay usuario seleccionado. Por favor, inicia sesión.',
					type: 'error',
					show: true
				}
			]);
			return;
		}

		notifications.update((n) => [
			...n,
			{
				id: tempId,
				message: `Iniciando descarga: ${track.title}`,
				type: 'info',
				show: true
			}
		]);

		try {
			const res = await fetch(`${API_IP}/api/downloads`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					id: track.track_id,
					isrc: track.isrc,
					user: $selectedUser,
					quality: 2
				})
			});

			const data = await res.json();
			if (!res.ok) throw new Error(data.error || `Error ${res.status}`);

			// Reemplaza la notificación de inicio por la respuesta del servidor
			notifications.update((n) => [
				...n.filter((noti) => noti.id !== tempId),
				{
					id: data.downloadId,
					message: data.message,
					type: data.status === 'exists' ? 'success' : 'info',
					show: true
				}
			]);

			// Solo iniciar seguimiento si no es "exists"
			if (data.status === 'downloading' || data.status === 'linking') {
				pollDownloadStatus(data.downloadId, track.title);
			}
		} catch (err: unknown) {
			const message = err instanceof Error ? err.message : 'Error desconocido';
			notifications.update((n) => [
				...n.filter((noti) => noti.id !== tempId),
				{
					id: Date.now().toString(),
					message: `Error al iniciar la descarga: ${message}`,
					type: 'error',
					show: true
				}
			]);
		}
	}

	function pollDownloadStatus(downloadId: string, title: string, attempts = 0) {
		const maxAttempts = 30; // 30 intentos → 1 minuto

		const interval = setInterval(async () => {
			try {
				const res = await fetch(`${API_IP}/api/downloads/${downloadId}/status`);
				if (!res.ok) throw new Error(`Error ${res.status}`);
				const data = await res.json();
				const status = data.status;

				let type: 'info' | 'success' | 'error' = 'info';
				let message = `Descargando ${title}...`;

				if (status === 'indexing') {
					message = `Indexando ${title}...`;
				} else if (status === 'success') {
					type = 'success';
					message = `Canción lista: ${title}`;
					clearInterval(interval);
				} else if (status === 'failed') {
					type = 'error';
					message = `Descarga fallida: ${title}`;
					clearInterval(interval);
				} else if (status === 'canceled') {
					type = 'error';
					message = `Descarga cancelada: ${title}`;
					clearInterval(interval);
				}

				// Actualiza la notificación
				notifications.update((n) => [
					...n.filter((noti) => noti.id !== downloadId),
					{
						id: downloadId,
						message,
						type,
						show: true
					}
				]);

				if (++attempts >= maxAttempts) {
					clearInterval(interval);
					notifications.update((n) => [
						...n.filter((noti) => noti.id !== downloadId),
						{
							id: downloadId,
							message: `Sin respuesta del servidor para: ${title}`,
							type: 'error',
							show: true
						}
					]);
				}
			} catch (error) {
				clearInterval(interval);
				notifications.update((n) => [
					...n.filter((noti) => noti.id !== downloadId),
					{
						id: downloadId,
						message: `Error consultando estado de descarga: ${title}`,
						type: 'error',
						show: true
					}
				]);
			}
		}, 1000);
	}
</script>

<div
	class="flex flex-col items-center gap-4 rounded-lg border border-gray-700 bg-gray-900 p-4 transition-all hover:bg-gray-800 sm:flex-row"
>
	<img
		src={track.image}
		alt="Cover de {track.album}"
		class="h-24 w-24 flex-shrink-0 rounded-md object-cover"
	/>

	<div class="flex-grow text-center sm:text-left">
		<h3 class="text-lg font-bold text-gray-100">{track.title}</h3>
		<p class="text-sm text-gray-400">Artista: {track.artist}</p>
		<p class="text-sm text-gray-400">Álbum: {track.album}</p>
		<p class="text-sm text-gray-400">Duración: {formatDuration(track.duration)}</p>

		{#if sampleUrl && showPlayer}
			<audio class="mt-2 w-full" controls src={sampleUrl}></audio>
		{/if}
	</div>

	<div class="ml-auto flex flex-col gap-2">
		<button
			on:click={downloadTrack}
			class="rounded-lg bg-purple-600 px-4 py-2 font-semibold text-white transition-colors hover:bg-purple-700 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 focus:ring-offset-gray-900"
		>
			Descargar
		</button>

		<button
			on:click={loadSample}
			class="rounded-lg border border-purple-500 px-4 py-2 text-purple-400 transition-colors hover:bg-purple-800 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 focus:ring-offset-gray-900 disabled:opacity-50"
			disabled={isLoadingSample}
		>
			{isLoadingSample ? 'Cargando...' : 'Escuchar muestra'}
		</button>
	</div>
</div>
