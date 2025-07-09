<script>
	import { isLoading, tracks } from '$lib/stores/stores';
	import TrackCard from './TrackCard.svelte';
</script>

{#if $isLoading}
	<div class="flex flex-col items-center justify-center py-8 text-gray-300">
		<div
			class="h-8 w-8 animate-spin rounded-full border-4 border-purple-500 border-t-transparent"
		></div>
		<p class="mt-4 text-sm">Cargando resultados...</p>
	</div>
{:else if $tracks.length > 0}
	<div class="max-h-[90vh] overflow-y-auto px-4 py-4">
		<div class="mx-auto flex max-w-4xl flex-col gap-4">
			{#each $tracks as track (track.track_id)}
				<TrackCard {track} />
			{/each}
		</div>
	</div>
{/if}
