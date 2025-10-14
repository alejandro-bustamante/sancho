<script lang="ts">
	import { currentTrack, isPlaying } from '$lib/stores/playerStore';
	import { API_IP } from '$lib/config';

	let audioPlayer: HTMLAudioElement;

	// Reacciona a los cambios en la canción actual
	$: if ($currentTrack && audioPlayer) {
		audioPlayer.src = `${API_IP}/api/tracks/${$currentTrack.id}/stream`;
		if ($isPlaying) {
			audioPlayer.play().catch((e) => console.error('Error al reproducir audio:', e));
		}
	}

	// Reacciona a los cambios en el estado de reproducción (play/pause)
	$: if (audioPlayer) {
		if ($isPlaying) {
			audioPlayer.play().catch((e) => console.error('Error al reproducir audio:', e));
		} else {
			audioPlayer.pause();
		}
	}

	function togglePlay() {
		isPlaying.update((p) => !p);
	}
</script>

{#if $currentTrack}
	<div
		class="fixed bottom-0 left-0 right-0 z-50 flex h-20 items-center justify-between border-t border-gray-700 bg-gray-800 px-4 text-white shadow-lg"
	>
		<div class="flex items-center gap-4">
			<button on:click={togglePlay} class="rounded-full bg-purple-600 p-2 hover:bg-purple-700">
				{#if $isPlaying}
					<!-- Icono Pausa -->
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="24"
						height="24"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
						><rect x="6" y="4" width="4" height="16"></rect><rect x="14" y="4" width="4" height="16"
						></rect></svg
					>
				{:else}
					<!-- Icono Play -->
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="24"
						height="24"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"><polygon points="5 3 19 12 5 21 5 3"></polygon></svg
					>
				{/if}
			</button>
			<div>
				<p class="font-bold">{$currentTrack.title}</p>
				<p class="text-sm text-gray-400">{$currentTrack.artist}</p>
			</div>
		</div>

		<audio
			bind:this={audioPlayer}
			controls
			class="w-1/2"
			on:play={() => isPlaying.set(true)}
			on:pause={() => isPlaying.set(false)}
			on:ended={() => isPlaying.set(false)}
		>
			Tu navegador no soporta el elemento de audio.
		</audio>
	</div>
{/if}
