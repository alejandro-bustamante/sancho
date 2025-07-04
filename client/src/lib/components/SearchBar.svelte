<script lang="ts">
	// Importamos los stores que manejan la búsqueda y estado general
	import { trackQuery, selectedUser, tracks, isLoading, notifications } from '$lib/stores/stores';

	// Función que realiza la búsqueda
	async function searchTracks() {
		const query = $trackQuery.trim(); // Accedemos al valor reactivo del store
		if (!query) return;

		isLoading.set(true);  // Indicamos que está cargando
		tracks.set([]);       // Limpiamos resultados anteriores

		try {
			const res = await fetch('http://localhost:8081/search', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					 title: query,
				})
			});

			const data = await res.json();
			tracks.set(data.results || []);  // Actualizamos store con resultados
		} catch (error) {
			// Agregamos una notificación de error
			notifications.update(n => [...n, {
				id: Date.now().toString(),
				message: 'Error al buscar canciones',
				type: 'error',
				show: true
			}]);
			console.error(error);
		} finally {
			isLoading.set(false);
		}
	}
</script>

<!-- Contenedor de barra de búsqueda -->
<div style="display: flex; gap: 0.5rem; margin: 1rem 0; flex-wrap: wrap; justify-content: center;">
	<!-- Campo de entrada de texto. Está enlazado al store trackQuery -->
	<input
		bind:value={$trackQuery}
		on:keydown={(e) => e.key === 'Enter' && searchTracks()}
		style="padding: 0.5rem; border-radius: 0.25rem;"
		placeholder="Buscar canción..."
	/>

	<!-- Botón para iniciar búsqueda -->
	<button on:click={searchTracks} style="padding: 0.5rem 1rem; background-color: #3b82f6; color: white; border-radius: 0.25rem;">
		Buscar
	</button>

	<!-- Selector de usuario -->
	<select bind:value={$selectedUser} style="padding: 0.5rem; border-radius: 0.25rem;">
		<option value="alejandro">Alejandro</option>
		<option value="alfredo">Alfredo</option>
	</select>
</div>
