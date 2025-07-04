<script>
	// Stores para saber si está cargando y para acceder a los resultados
	import { isLoading, tracks } from '$lib/stores/stores';
	import TrackCard from './TrackCard.svelte';
</script>

<!-- Si está cargando mostramos spinner -->
{#if $isLoading}
	<div style="text-align: center; padding: 1rem;">
		<div style="width: 2rem; height: 2rem; border: 2px solid #3b82f6; border-top-color: transparent; border-radius: 9999px; animation: spin 1s linear infinite;"></div>
		<p style="margin-top: 0.5rem;">Cargando resultados...</p>
	</div>

<!-- Si hay resultados los mostramos -->
{:else if $tracks.length > 0}
	<div style="max-width: 800px; margin: 0 auto; background: white; border-radius: 0.5rem; padding: 1rem; overflow-y: auto; max-height: 70vh;">
		{#each $tracks as track (track.track_id)}
			<TrackCard {track} />
		{/each}
	</div>
{/if}

<style>
	@keyframes spin {
		from { transform: rotate(0deg); }
		to { transform: rotate(360deg); }
	}
</style>
