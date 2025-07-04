<script lang="ts">
	import AuthForm from './AuthForm.svelte';

	// Estado para saber si estamos en modo login o registro
	let isLogin = true;

	function toggleMode() {
		isLogin = !isLogin;
	}
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-100">
	<div class="w-full max-w-md p-6 bg-white rounded shadow">
		<h2 class="text-2xl font-semibold text-center mb-4">
			{isLogin ? 'Iniciar sesión' : 'Registrarse'}
		</h2>

		<AuthForm mode={isLogin ? 'login' : 'register'} />

		<p class="text-center mt-4 text-sm text-gray-600">
			{isLogin ? '¿No tienes una cuenta?' : '¿Ya tienes una cuenta?'}
			<button class="text-blue-500 ml-2" on:click={toggleMode}>
				{isLogin ? 'Regístrate' : 'Inicia sesión'}
			</button>
		</p>
	</div>
</div>
