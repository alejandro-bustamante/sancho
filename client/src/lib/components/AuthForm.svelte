<script lang="ts">
	import { notifications } from '$lib/stores/notifications';
    import { currentUser } from '$lib/stores/auth'

	export let mode: 'login' | 'register';

	// Campos del formulario
	let username = '';
	let password = '';
	let email = '';

	// Acción al enviar
	async function handleSubmit() {
		const payload: Record<string, string> = { username, password };
		if (mode === 'register') payload.email = email;

		const res = await fetch(
			`http://localhost:8081/${mode === 'login' ? 'auth' : 'users'}`,
			{
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(payload)
			}
		);

		const data = await res.json();

		if (!res.ok) {
			notifications.add({ type: 'error', message: data.error || 'Algo salió mal.' });
			return;
		}

        // We save the current session
        currentUser.set(username)

		notifications.add({
			type: 'success',
			message: mode === 'login'
				? 'Sesión iniciada correctamente.'
				: 'Usuario registrado con éxito.'
		});
	}
</script>

<form class="flex flex-col gap-4" on:submit|preventDefault={handleSubmit}>
	<input
		class="border p-2 rounded"
		placeholder="Usuario"
		bind:value={username}
		required
	/>

	<input
		class="border p-2 rounded"
		placeholder="Contraseña"
		type="password"
		bind:value={password}
		required
	/>

	{#if mode === 'register'}
		<input
			class="border p-2 rounded"
			placeholder="Correo electrónico"
			type="email"
			bind:value={email}
			required
		/>
	{/if}

	<button class="bg-blue-600 text-white py-2 rounded hover:bg-blue-700 transition">
		{mode === 'login' ? 'Iniciar sesión' : 'Registrarse'}
	</button>
</form>
