<script lang="ts">
	import { notifications } from '$lib/stores/notifications';
	import { currentUser } from '$lib/stores/auth';
	import { API_IP } from '$lib/config';

	export let mode: 'login' | 'register';
	// Campos del formulario
	let username = '';
	let password = '';
	let email = '';
	// Acción al enviar
	async function handleSubmit() {
		const payload: Record<string, string> = { username, password };
		if (mode === 'register') payload.email = email;

		const res = await fetch(`${API_IP}/api/${mode === 'login' ? 'auth' : 'users'}`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(payload)
		});
		const data = await res.json();

		if (!res.ok) {
			notifications.add({ type: 'error', message: data.error || 'Algo salió mal.' });
			return;
		}

		// We save the current session
		currentUser.set(username);

		notifications.add({
			type: 'success',
			message: mode === 'login' ? 'Sesión iniciada correctamente.' : 'Usuario registrado con éxito.'
		});
	}
</script>

<form class="space-y-4" on:submit|preventDefault={handleSubmit}>
	<input
		class="w-full rounded-md border border-gray-600 bg-gray-700 px-4 py-2 text-gray-200 placeholder-gray-400 transition focus:border-transparent focus:outline-none focus:ring-2 focus:ring-purple-500"
		placeholder="Usuario"
		bind:value={username}
		required
	/>

	<input
		class="w-full rounded-md border border-gray-600 bg-gray-700 px-4 py-2 text-gray-200 placeholder-gray-400 transition focus:border-transparent focus:outline-none focus:ring-2 focus:ring-purple-500"
		placeholder="Contraseña"
		type="password"
		bind:value={password}
		required
	/>

	{#if mode === 'register'}
		<input
			class="w-full rounded-md border border-gray-600 bg-gray-700 px-4 py-2 text-gray-200 placeholder-gray-400 transition focus:border-transparent focus:outline-none focus:ring-2 focus:ring-purple-500"
			placeholder="Correo electrónico"
			type="email"
			bind:value={email}
			required
		/>
	{/if}

	<button
		class="w-full rounded-md bg-purple-600 px-4 py-2 font-bold text-white transition-colors hover:bg-purple-700 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 focus:ring-offset-gray-800"
	>
		{mode === 'login' ? 'Iniciar sesión' : 'Registrarse'}
	</button>
</form>
