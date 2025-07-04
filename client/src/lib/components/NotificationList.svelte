<script lang="ts">
	// Importamos el store que contiene las notificaciones
	import { notifications } from '$lib/stores/stores';

	// Función que elimina una notificación por su ID
	function hideNotification(id: string) {
		// Usamos update para modificar el store: filtramos fuera la notificación con el ID dado
		notifications.update(n => n.filter(noti => noti.id !== id));
	}
</script>

<!-- Contenedor fijo donde aparecen las notificaciones -->
<div style="position: fixed; top: 1rem; right: 1rem; z-index: 50;">
	{#each $notifications as notification (notification.id)}
		<!-- Renderiza cada notificación con estilos de color según el tipo -->
		<div
			style="margin-bottom: 0.5rem; padding: 0.75rem; border-radius: 0.5rem; color: white; display: flex; justify-content: space-between;"
			class:bg-green-500={notification.type === 'success'}
			class:bg-blue-500={notification.type === 'info'}
			class:bg-red-500={notification.type === 'error'}
		>
			<!-- Texto del mensaje -->
			<span>{notification.message}</span>
			<!-- Botón para cerrar la notificación -->
			<button on:click={() => hideNotification(notification.id)} style="margin-left: 1rem;">&times;</button>
		</div>
	{/each}
</div>
