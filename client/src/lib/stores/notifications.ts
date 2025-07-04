import { writable } from 'svelte/store';

export type Notification = {
	type: 'success' | 'error' | 'info';
	message: string;
};

function createNotificationStore() {
	const { subscribe, update, set } = writable<Notification[]>([]);

	return {
		subscribe,
		set,
		update,
		add: (notification: Notification) => update((n) => [...n, notification]),
		clear: () => set([])
	};
}

export const notifications = createNotificationStore();
