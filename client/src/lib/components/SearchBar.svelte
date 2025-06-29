<script lang="ts">
	import { trackQuery, selectedUser, tracks, isLoading, notifications } from '$lib/stores';

	async function searchTracks() {
		const query = $trackQuery.trim();
		if (!query) return;

		isLoading.set(true);
		tracks.set([]);

		// const proxyUrl = 'http://localhost:8081/proxy';
		// const apiUrl = `https://api.deezer.com/search?q=${encodeURIComponent(query)}`;
		// const url = `${proxyUrl}?url=${encodeURIComponent(apiUrl)}`;

		try {

			const res = await fetch('http://localhost:8081/search/deezer', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ title: query })
			});

			const data = await res.json();
			console.log(data)
			// const res = await fetch(url);
			// const data = await res.json();
			tracks.set(data.results || []);
		} catch (error) {
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

<div class="flex flex-wrap justify-center items-center gap-2 my-4">
	<input
		bind:value={$trackQuery}
		on:keydown={(e) => e.key === 'Enter' && searchTracks()}
		class="text-black rounded px-2 py-2 text-sm w-64"
		placeholder="Buscar canción..."
	/>
	<button on:click={searchTracks} class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">Buscar</button>
	<select bind:value={$selectedUser} class="px-3 py-2 rounded border border-gray-300 text-gray-800">
		<option value="alejandro">Alejandro</option>
		<option value="alfredo">Alfredo</option>
	</select>
</div>
