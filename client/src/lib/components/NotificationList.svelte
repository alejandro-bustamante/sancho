<script lang="ts">
	import { notifications } from '$lib/stores';
	import type { Notification } from '$lib/types';

	function hideNotification(id: string) {
		notifications.update((n) => n.filter((noti) => noti.id !== id));
	}
</script>

<div class="fixed top-5 right-5 z-50 flex flex-col gap-3 w-80">
	{#each $notifications as notification (notification.id)}
		<div
			class="rounded-lg shadow-lg p-4 text-white flex justify-between items-center"
			class:bg-green-500={notification.type === 'success'}
			class:bg-blue-500={notification.type === 'info'}
			class:bg-red-500={notification.type === 'error'}
		>
			<span>{notification.message}</span>
			<button class="ml-4 hover:text-gray-200" on:click={() => hideNotification(notification.id)}>&times;</button>
		</div>
	{/each}
</div>
