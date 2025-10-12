import { writable } from 'svelte/store';

export type View = 'search' | 'library';

export const currentView = writable<View>('search');

/**
 * A function to programmatically change the current view.
 * @param view The view to navigate to.
 */
export const navigateTo = (view: View) => {
	currentView.set(view);
};
