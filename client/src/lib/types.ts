export type NotificationType = 'success' | 'info' | 'error';

export interface Notification {
	id: string;
	message: string;
	type: NotificationType;
	show: boolean;
}

export interface Track {
	id: string;
	title: string;
	link: string;
	preview: string;
	artist: {
		name: string;
	};
	album: {
		title: string;
		cover_small: string;
	};
}
