import { writable } from 'svelte/store';
import type { Notification, Track } from './types';

export const trackQuery = writable('');
export const selectedUser = writable('alejandro');
export const tracks = writable<Track[]>([]);
export const isLoading = writable(false);
export const notifications = writable<Notification[]>([]);
