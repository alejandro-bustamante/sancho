<script lang="ts">
	import { API_IP } from '$lib/config';

	// Importamos los stores que manejan la búsqueda y estado general
	import { trackQuery, selectedUser, tracks, isLoading, notifications } from '$lib/stores/stores';

	// Función que realiza la búsqueda
	async function searchTracks() {
		const query = $trackQuery.trim(); // Accedemos al valor reactivo del store
		if (!query) return;

		isLoading.set(true); // Indicamos que está cargando
		tracks.set([]); // Limpiamos resultados anteriores

		try {
			let api = `${API_IP}/search?q=${encodeURIComponent(query)}`;
			// const res = await fetch(`${API_IP}/search?q=${encodeURIComponent(query)}`);
			const res = await fetch(api);
			const data = await res.json();
			tracks.set(data.results || []); // Actualizamos store con resultados
		} catch (error) {
			// Agregamos una notificación de error
			notifications.update((n) => [
				...n,
				{
					id: Date.now().toString(),
					message: 'Error al buscar canciones',
					type: 'error',
					show: true
				}
			]);
			console.error(error);
		} finally {
			isLoading.set(false);
		}
	}
</script>

<div class="mb-8 flex flex-wrap items-center justify-center gap-4 rounded-lg bg-gray-800 p-4">
	<input
		bind:value={$trackQuery}
		on:keydown={(e) => e.key === 'Enter' && searchTracks()}
		class="rounded-md border border-gray-600 bg-gray-700 px-4 py-2 text-gray-200 placeholder-gray-400 transition focus:border-transparent focus:outline-none focus:ring-2 focus:ring-purple-500"
		placeholder="Buscar canción..."
	/>

	<button
		on:click={searchTracks}
		class="rounded-md bg-purple-600 px-6 py-2 font-semibold text-white transition-colors hover:bg-purple-700 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 focus:ring-offset-gray-800"
	>
		Buscar
	</button>
</div>
