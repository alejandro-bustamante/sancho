<script lang="ts">
	import { currentUser } from '$lib/stores/auth';
	import { onMount } from 'svelte';
	import { get } from 'svelte/store';

	import IndexFolderForm from '$lib/components/IndexFolderForm.svelte';
	import NotificationList from '$lib/components/NotificationList.svelte';
	import SearchBar from '$lib/components/SearchBar.svelte';
	import TrackList from '$lib/components/TrackList.svelte';
	import LoginPage from '$lib/components/LoginPage.svelte';

	function logout() {
		currentUser.set(null); // Esto también borra localStorage
	}

	let user: string | null = null;

	// Nos suscribimos al store para saber si hay sesión
	$currentUser; // esto vuelve reactivo automáticamente

	let showMenu = false;
</script>

<!-- Al lado del resto de componentes -->
<div class="flex w-full justify-end px-4 pt-2">
	<button
		on:click={logout}
		class="rounded border border-red-900 px-3 py-1 text-sm text-red-50 transition hover:bg-red-800"
	>
		Cerrar sesión
	</button>
</div>
{#if $currentUser === null}
	<!-- Usuario no autenticado: mostramos login -->
	<LoginPage />
{:else}
	<!-- Usuario autenticado: mostramos app principal -->
	<div class="bg-dark-gray text-light flex min-h-screen flex-col items-center px-4 pt-5">
		<NotificationList />
		<SearchBar />
		<TrackList />
	</div>
{/if}
