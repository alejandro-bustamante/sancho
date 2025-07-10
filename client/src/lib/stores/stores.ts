// src/lib/stores/stores.ts
import { writable, type Writable } from 'svelte/store';
import type { Notification, Track } from './types';
import { currentUser } from './auth';

export const trackQuery = writable('');
export const selectedUser: Writable<string | null> = writable(null);
export const tracks = writable<Track[]>([]);
export const isLoading = writable(false);
export const notifications = writable<Notification[]>([]);

// Sincronizar selectedUser con currentUser
currentUser.subscribe((value) => {
	selectedUser.set(value);
});
