export type NotificationType = 'success' | 'info' | 'error';

export interface Notification {
	id: string;
	message: string;
	type: NotificationType;
	show: boolean;
}

export interface Track {
	// id: string;
	title: string;
	artist: string;
	album: string;
	duration: number;
	image: string;
	track_id: number;
	source: string;
	isrc: string;
}

export interface UserLibraryTrack {
	id: number;
	title: string;
	duration: { Int64: number; Valid: boolean };
	artist: { String: string; Valid: boolean };
	album: { String: string; Valid: boolean };
	album_art_path: { String: string; Valid: boolean };
}
