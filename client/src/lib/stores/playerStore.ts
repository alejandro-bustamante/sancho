import { writable } from 'svelte/store';
import type { UserLibraryTrack } from './types';

export interface NowPlaying {
	id: number;
	title: string;
	artist: string;
}

export const currentTrack = writable<NowPlaying | null>(null);

export const isPlaying = writable(false);

export function playTrack(track: UserLibraryTrack) {
	currentTrack.set({
		id: track.id,
		title: track.title,
		artist: track.artist.Valid ? track.artist.String : 'Artista Desconocido'
	});
	isPlaying.set(true);
}
