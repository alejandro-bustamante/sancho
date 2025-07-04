// src/lib/stores/auth.ts
import { writable } from 'svelte/store';

// Indicador de que ya estamos en el cliente y se leyó localStorage
export const isClientReady = writable(false);

export const currentUser = writable<string | null>(null);

// Solo ejecutamos en cliente
if (typeof window !== 'undefined') {
	const storedUser = localStorage.getItem('currentUser');
	currentUser.set(storedUser);

	currentUser.subscribe((value) => {
		if (value === null) {
			localStorage.removeItem('currentUser');
		} else {
			localStorage.setItem('currentUser', value);
		}
	});

	isClientReady.set(true); // <-- Ya está listo
}
