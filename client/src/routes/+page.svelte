<script lang="ts">
	import { currentUser } from '$lib/stores/auth';
	import NotificationList from '$lib/components/NotificationList.svelte';
	import SearchBar from '$lib/components/SearchBar.svelte';
	import TrackList from '$lib/components/TrackList.svelte';
	import LoginPage from '$lib/components/LoginPage.svelte';
	import UserLibrary from '$lib/components/UserLibrary.svelte';
	import SettingsPage from '$lib/components/SettingsPage.svelte';

	// Import the new store for view management
	import { currentView, navigateTo } from '$lib/stores/navigation';

	function logout() {
		currentUser.set(null);
		navigateTo('search'); // Reset the view to 'search' on logout
	}
</script>

{#if $currentUser === null}
	<LoginPage />
{:else}
	<div class="text-light flex min-h-screen flex-col items-center bg-gray-950 px-4 pt-5">
		<div class="flex w-full max-w-4xl items-center justify-between">
			<h1 class="text-4xl font-bold text-white">Sancho</h1>
			<button
				on:click={logout}
				class="rounded border border-red-900 px-3 py-1 text-sm text-red-50 transition hover:bg-red-800"
			>
				Cerrar sesión
			</button>
		</div>

		<NotificationList />

		<div class="mb-4 mt-8 w-full max-w-4xl border-b border-gray-700">
			<nav class="-mb-px flex space-x-4">
				<button
					on:click={() => navigateTo('search')}
					class="whitespace-nowrap border-b-2 px-4 py-2 text-lg font-medium transition-colors"
					class:border-purple-500={$currentView === 'search'}
					class:text-purple-400={$currentView === 'search'}
					class:border-transparent={$currentView !== 'search'}
					class:text-gray-400={$currentView !== 'search'}
					class:hover:text-white={$currentView !== 'search'}
					class:hover:border-gray-500={$currentView !== 'search'}
				>
					Búsqueda
				</button>
				<button
					on:click={() => navigateTo('library')}
					class="whitespace-nowrap border-b-2 px-4 py-2 text-lg font-medium transition-colors"
					class:border-purple-500={$currentView === 'library'}
					class:text-purple-400={$currentView === 'library'}
					class:border-transparent={$currentView !== 'library'}
					class:text-gray-400={$currentView !== 'library'}
					class:hover:text-white={$currentView !== 'library'}
					class:hover:border-gray-500={$currentView !== 'library'}
				>
					Librería
				</button>
				<button
					on:click={() => navigateTo('settings')}
					class="whitespace-nowrap border-b-2 px-4 py-2 text-lg font-medium transition-colors"
					class:border-purple-500={$currentView === 'settings'}
					class:text-purple-400={$currentView === 'settings'}
					class:border-transparent={$currentView !== 'settings'}
					class:text-gray-400={$currentView !== 'settings'}
					class:hover:text-white={$currentView !== 'settings'}
					class:hover:border-gray-500={$currentView !== 'settings'}
				>
					Configuración
				</button>
			</nav>
		</div>

		<div class="w-full max-w-4xl">
			{#if $currentView === 'search'}
				<SearchBar />
				<TrackList />
			{:else if $currentView === 'library'}
				<UserLibrary />
			{:else if $currentView === 'settings'}
				<SettingsPage />
			{/if}
		</div>
	</div>
{/if}
