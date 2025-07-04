<script lang="ts">
	import { currentUser } from '$lib/stores/auth';
	import { onMount } from 'svelte';
	import { get } from 'svelte/store';

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
</script>

<!-- Al lado del resto de componentes -->
<div class="w-full flex justify-end px-4 pt-2">
	<button
		on:click={logout}
		class="text-sm text-red-500 border border-red-500 px-3 py-1 rounded hover:bg-red-100 transition"
	>
		Cerrar sesión
	</button>
</div>
{#if $currentUser === null}
	<!-- Usuario no autenticado: mostramos login -->
	<LoginPage />
{:else}
	<!-- Usuario autenticado: mostramos app principal -->
	<div class="bg-dark-gray text-light min-h-screen flex flex-col items-center pt-5 px-4">
		<NotificationList />
		<SearchBar />
		<TrackList />
	</div>
{/if}

