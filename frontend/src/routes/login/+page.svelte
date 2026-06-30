<script lang="ts">
	import { apiFetch } from '$lib/api';

	let usuario = $state('');
	let contrasena = $state('');

	let mensajeError = $state('');
	let mensajeExito = $state('');

	async function handleLogin() {
		mensajeError = '';
		mensajeExito = '';

		try {
			const respuesta = await apiFetch('/auth/login', {
				method: 'POST',
				body: JSON.stringify({ usuario, contrasena })
			});

			mensajeExito = `¡Bienvenido, ${respuesta.empleado?.usuario || 'admin'}!`;
			console.log('Datos recibidos:', respuesta);
		} catch (error: any) {
			mensajeError = error.message;
		}
	}
</script>

<div class="max-w-sm mx-auto mt-20 p-6 bg-white rounded-lg shadow-md">
	<h2 class="text-2xl font-bold text-primario text-center mb-6">Iniciar Sesión</h2>

	{#if mensajeError}
		<p class="text-error font-bold text-center mb-4">{mensajeError}</p>
	{/if}

	{#if mensajeExito}
		<p class="text-exito font-bold text-center mb-4">{mensajeExito}</p>
	{/if}

	<form on:submit|preventDefault={handleLogin} class="flex flex-col gap-4">
		<div>
			<label class="block text-sm font-semibold mb-1">Usuario:</label>
			<input
				type="text"
				bind:value={usuario}
				required
				class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:border-primario focus:ring-1 focus:ring-primario"
			/>
		</div>

		<div>
			<label class="block text-sm font-semibold mb-1">Contraseña:</label>
			<input
				type="password"
				bind:value={contrasena}
				required
				class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:border-primario focus:ring-1 focus:ring-primario"
			/>
		</div>

		<button
			type="submit"
			class="w-full bg-primario hover:bg-primario-hover text-white font-bold py-2 px-4 rounded-md transition-colors mt-2"
		>
			Entrar
		</button>
	</form>
</div>
