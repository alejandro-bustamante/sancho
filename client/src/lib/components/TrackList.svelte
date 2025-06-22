<!-- src/lib/components/TrackList.svelte -->
<script>
  import { isLoading, tracks } from '$lib/stores';
  import TrackCard from './TrackCard.svelte';
</script>

{#if $isLoading}
  <div class="text-center py-4">
    <div class="inline-block animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-500"></div>
    <p class="mt-2">Cargando resultados...</p>
  </div>
{:else if $tracks.length > 0}
  <div class="w-full max-w-3xl mx-auto bg-white rounded shadow-lg p-6 text-gray-800 overflow-y-auto max-h-[70vh]">
    {#each $tracks as track (track.id)}
      <TrackCard {track} />
    {/each}
  </div>
{/if}
