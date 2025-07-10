<script lang="ts">
	// Importamos el store que contiene las notificaciones
	import { notifications } from '$lib/stores/stores';
	// Función que elimina una notificación por su ID
	function hideNotification(id: string) {
		notifications.update((n) => n.filter((noti) => noti.id !== id));
	}
</script>

<div class="fixed right-4 top-4 z-50 w-full max-w-xs space-y-2">
	{#each $notifications as notification (notification.id)}
		<div
			class="relative flex items-center justify-between rounded-lg p-4 pr-10 text-white shadow-lg"
			class:bg-green-600={notification.type === 'success'}
			class:bg-blue-600={notification.type === 'info'}
			class:bg-red-600={notification.type === 'error'}
			role="alert"
		>
			<span class="text-sm font-medium">{notification.message}</span>
			<button
				on:click={() => hideNotification(notification.id)}
				class="absolute right-2 top-1/2 -translate-y-1/2 rounded-md bg-transparent p-1 text-white/70 hover:text-white/100 focus:outline-none focus:ring-2 focus:ring-white"
				aria-label="Cerrar"
			>
				&times;
			</button>
		</div>
	{/each}
</div>
