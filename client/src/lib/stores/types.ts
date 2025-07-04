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
