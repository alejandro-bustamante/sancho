<script lang="ts">
	import AuthForm from './AuthForm.svelte';

	// Estado para saber si estamos en modo login o registro
	let isLogin = true;
	function toggleMode() {
		isLogin = !isLogin;
	}
</script>

<div class="flex min-h-screen items-center justify-center bg-gray-900 text-gray-200">
	<div
		class="w-full max-w-md space-y-6 rounded-lg border border-gray-700 bg-gray-800 p-8 shadow-lg"
	>
		<h2 class="text-center text-3xl font-bold text-white">
			{isLogin ? 'Iniciar sesión' : 'Registrarse'}
		</h2>

		<AuthForm mode={isLogin ? 'login' : 'register'} />

		<p class="text-center text-sm text-gray-400">
			{isLogin ? '¿No tienes una cuenta?' : '¿Ya tienes una cuenta?'}
			<button
				class="font-semibold text-purple-400 transition-colors hover:text-purple-300"
				on:click={toggleMode}
			>
				{isLogin ? 'Regístrate' : 'Inicia sesión'}
			</button>
		</p>
	</div>
</div>
