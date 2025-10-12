<script lang="ts">
	import { onMount } from 'svelte';
	import { currentUser } from '$lib/stores/auth';
	import { notifications } from '$lib/stores/notifications';
	import { playTrack } from '$lib/stores/playerStore';
	import type { UserLibraryTrack } from '$lib/stores/types';
	import { API_IP } from '$lib/config';

	let tracks: UserLibraryTrack[] = [];
	let filteredTracks: UserLibraryTrack[] = [];
	let isLoading = true;
	let searchTerm = '';

	// Carga las canciones del usuario al montar el componente
	onMount(async () => {
		const user = $currentUser;
		if (!user) return;

		try {
			const res = await fetch(`${API_IP}/users/${user}/tracks`);
			if (!res.ok) throw new Error('No se pudo cargar la librería.');
			const data = await res.json();
			tracks = data || [];
			filteredTracks = tracks;
		} catch (error) {
			const message = error instanceof Error ? error.message : 'Error desconocido';
			notifications.add({ type: 'error', message });
		} finally {
			isLoading = false;
		}
	});

	// Filtra las canciones según el término de búsqueda
	$: {
		if (searchTerm) {
			const lowerCaseSearch = searchTerm.toLowerCase();
			filteredTracks = tracks.filter(
				(track) =>
					track.title.toLowerCase().includes(lowerCaseSearch) ||
					(track.artist.Valid && track.artist.String.toLowerCase().includes(lowerCaseSearch)) ||
					(track.album.Valid && track.album.String.toLowerCase().includes(lowerCaseSearch))
			);
		} else {
			filteredTracks = tracks;
		}
	}

	// Formatea la duración de segundos a mm:ss
	function formatDuration(ms: number | null): string {
		if (ms === null) return 'N/A';
		const totalSeconds = Math.floor(ms / 1000);
		const minutes = Math.floor(totalSeconds / 60);
		const seconds = totalSeconds % 60;
		return `${minutes}:${seconds.toString().padStart(2, '0')}`;
	}

	// Maneja la eliminación de una canción
	async function handleDelete(trackId: number) {
		const user = $currentUser;
		if (!user) return;

		if (!confirm('¿Estás seguro de que quieres eliminar esta canción de tu librería?')) return;

		try {
			const res = await fetch(`${API_IP}/users/${user}/tracks/${trackId}`, {
				method: 'DELETE'
			});

			if (!res.ok) throw new Error('No se pudo eliminar la canción.');

			// Actualiza la UI eliminando la canción de la lista
			tracks = tracks.filter((t) => t.id !== trackId);
			notifications.add({ type: 'success', message: 'Canción eliminada correctamente.' });
		} catch (error) {
			const message = error instanceof Error ? error.message : 'Error desconocido';
			notifications.add({ type: 'error', message });
		}
	}
</script>

<div class="container mx-auto max-w-4xl py-4 text-white">
	<div class="mb-6 flex justify-end">
		<input
			type="text"
			bind:value={searchTerm}
			placeholder="Filtrar por título, artista o álbum..."
			class="w-full rounded-md border-gray-600 bg-gray-700 px-4 py-2 placeholder-gray-400 focus:border-purple-500 focus:ring-purple-500 sm:w-auto"
		/>
	</div>

	{#if isLoading}
		<p class="text-center">Cargando tu música...</p>
	{:else if filteredTracks.length === 0}
		<p class="text-center text-gray-400">No se encontraron canciones. ¡Prueba a añadir algunas!</p>
	{:else}
		<div class="overflow-hidden rounded-lg border border-gray-700 bg-gray-800">
			<table class="min-w-full divide-y divide-gray-700">
				<thead class="bg-gray-800">
					<tr>
						<th
							class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-300"
							>Título</th
						>
						<th
							class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-300"
							>Artista</th
						>
						<th
							class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-300"
							>Álbum</th
						>
						<th
							class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-300"
							>Duración</th
						>
						<th
							class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-300"
							>Acciones</th
						>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-700 bg-gray-900">
					{#each filteredTracks as track (track.id)}
						<tr class="transition-colors hover:bg-gray-800">
							<td class="whitespace-nowrap px-6 py-4 text-sm font-medium">{track.title}</td>
							<td class="whitespace-nowrap px-6 py-4 text-sm text-gray-400"
								>{track.artist.Valid ? track.artist.String : 'N/A'}</td
							>
							<td class="whitespace-nowrap px-6 py-4 text-sm text-gray-400"
								>{track.album.Valid ? track.album.String : 'N/A'}</td
							>
							<td class="whitespace-nowrap px-6 py-4 text-sm text-gray-400"
								>{formatDuration(track.duration.Valid ? track.duration.Int64 : null)}</td
							>
							<td class="whitespace-nowrap px-6 py-4 text-sm">
								<div class="flex items-center gap-2">
									<!-- svelte-ignore a11y_consider_explicit_label -->
									<button
										on:click={() => playTrack(track)}
										class="text-green-400 hover:text-green-300"
										title="Reproducir"
									>
										<svg
											xmlns="http://www.w3.org/2000/svg"
											width="20"
											height="20"
											viewBox="0 0 24 24"
											fill="none"
											stroke="currentColor"
											stroke-width="2"
											stroke-linecap="round"
											stroke-linejoin="round"><polygon points="5 3 19 12 5 21 5 3"></polygon></svg
										>
									</button>
									<!-- svelte-ignore a11y_consider_explicit_label -->
									<button
										on:click={() => handleDelete(track.id)}
										class="text-red-400 hover:text-red-300"
										title="Eliminar"
									>
										<svg
											xmlns="http://www.w3.org/2000/svg"
											width="20"
											height="20"
											viewBox="0 0 24 24"
											fill="none"
											stroke="currentColor"
											stroke-width="2"
											stroke-linecap="round"
											stroke-linejoin="round"
											><polyline points="3 6 5 6 21 6"></polyline><path
												d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"
											></path></svg
										>
									</button>
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>
