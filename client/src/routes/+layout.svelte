<script lang="ts">
	import '../app.css';
	import AudioPlayer from '$lib/components/AudioPlayer.svelte';

	let { children } = $props();
</script>

<!-- 
  Añadimos un contenedor principal con padding en la parte inferior 
  para dejar espacio al reproductor de audio fijo.
-->
<div class="relative min-h-screen bg-gray-950 pb-24">
	{@render children()}
</div>

<!-- 
  El reproductor de audio se coloca aquí para que sea persistente 
  en todas las páginas de la aplicación.
-->
<AudioPlayer />
