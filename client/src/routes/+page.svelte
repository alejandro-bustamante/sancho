<script lang="ts">
	import { currentUser } from '$lib/stores/auth';
	import { onMount } from 'svelte';
	import { get } from 'svelte/store';

	import NotificationList from '$lib/components/NotificationList.svelte';
	import SearchBar from '$lib/components/SearchBar.svelte';
	import TrackList from '$lib/components/TrackList.svelte';
	import LoginPage from '$lib/components/LoginPage.svelte';
	import UserLibrary from '$lib/components/UserLibrary.svelte';

	function logout() {
		currentUser.set(null); // Esto también borra localStorage
	}

	let user: string | null = null;

	// Nos suscribimos al store para saber si hay sesión
	$currentUser; // esto vuelve reactivo automáticamente
	let currentView: 'search' | 'library' = 'search';

	let showMenu = false;
</script>

<!-- Al lado del resto de componentes -->
<div class="flex w-full items-center justify-end gap-4 px-4 pt-2">
	{#if $currentUser}
		<button
			on:click={() => (currentView = 'library')}
			class="rounded-md bg-purple-600 px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-purple-700 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 focus:ring-offset-gray-900"
		>
			Mi Librería
		</button>
		<button
			on:click={logout}
			class="rounded border border-red-900 px-3 py-1 text-sm text-red-50 transition hover:bg-red-800"
		>
			Cerrar sesión
		</button>
	{/if}
</div>

{#if $currentUser === null}
	<!-- Usuario no autenticado: mostramos login -->
	<LoginPage />
{:else}
	<!-- Usuario autenticado: mostramos app principal -->
	<div class="text-light flex min-h-screen flex-col items-center bg-gray-950 px-4 pt-5">
		<h1 class="mb-4 text-4xl font-bold text-white">Sancho</h1>
		<NotificationList />
		{#if currentView === 'search'}
			<SearchBar />
			<TrackList />
		{:else if currentView === 'library'}
			<UserLibrary />
		{/if}
	</div>
{/if}
