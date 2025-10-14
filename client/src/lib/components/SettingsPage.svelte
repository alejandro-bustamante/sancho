<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { notifications } from '$lib/stores/notifications';

	let isRunning = false;
	let processed = 0;
	let total = 0;
	let progress = 0;
	let intervalId: ReturnType<typeof setInterval>;

	async function getStatus() {
		try {
			const res = await fetch(`http://localhost:5400/api/library/thumbnails/status`);
			if (!res.ok) {
				throw new Error('Failed to fetch status');
			}
			const data = await res.json();
			isRunning = data.isRunning;
			processed = data.processed;
			total = data.total;

			if (total > 0) {
				progress = (processed / total) * 100;
			} else {
				progress = 0;
			}

			if (!isRunning && intervalId) {
				if (total > 0 && processed === total) {
					notifications.add({
						type: 'success',
						message: 'Generación de miniaturas completada.'
					});
				}
				clearInterval(intervalId);
			}
		} catch (error) {
			notifications.add({
				type: 'error',
				message: 'Error al consultar estado de la generación de miniaturas.'
			});
			if (intervalId) clearInterval(intervalId);
		}
	}

	async function startGeneration() {
		try {
			const res = await fetch(`http://localhost:5400/api/library/thumbnails`, { method: 'POST' });
			if (!res.ok) {
				throw new Error('Failed to start generation process');
			}
			notifications.add({
				type: 'info',
				message: 'Iniciando generación de miniaturas...'
			});

			// Start polling
			await getStatus(); // Get initial status right away
			if (isRunning && !intervalId) {
				intervalId = setInterval(getStatus, 1000);
			}
		} catch (error) {
			notifications.add({
				type: 'error',
				message: 'Error al iniciar la generación de miniaturas.'
			});
		}
	}

	onMount(async () => {
		await getStatus(); // Check status on component mount in case it's already running
		if (isRunning && !intervalId) {
			intervalId = setInterval(getStatus, 1000);
		}
	});

	onDestroy(() => {
		if (intervalId) {
			clearInterval(intervalId);
		}
	});
</script>

<div class="container mx-auto max-w-4xl py-4 text-white">
	<h2 class="mb-6 text-2xl font-bold">Configuración de la Librería</h2>

	<div class="rounded-lg border border-gray-700 bg-gray-800 p-6">
		<h3 class="mb-4 text-lg font-semibold">Miniaturas de Álbumes</h3>
		<p class="mb-4 text-gray-400">
			Genera las imágenes de portada para los álbumes de tu librería que no la tengan. Este proceso
			puede tardar varios minutos dependiendo del tamaño de tu colección.
		</p>

		{#if isRunning}
			<div class="w-full">
				<p class="mb-2 text-center">{processed} / {total} álbumes procesados</p>
				<div class="h-4 w-full rounded-full bg-gray-700">
					<div class="h-4 rounded-full bg-purple-600" style="width: {progress}%"></div>
				</div>
			</div>
		{:else}
			<button
				on:click={startGeneration}
				class="rounded-lg bg-purple-600 px-6 py-2 font-semibold text-white transition-colors hover:bg-purple-700 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 focus:ring-offset-gray-800 disabled:opacity-50"
				disabled={isRunning}
			>
				Generar Miniaturas
			</button>
		{/if}
	</div>
</div>
