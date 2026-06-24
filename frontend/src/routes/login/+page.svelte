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

<div style="max-width: 300px; margin: 50px auto; font-family: sans-serif;">
	<h2>Iniciar Sesión</h2>

	{#if mensajeError}
		<p style="color: red; font-weight: bold;">{mensajeError}</p>
	{/if}

	{#if mensajeExito}
		<p style="color: green; font-weight: bold;">{mensajeExito}</p>
	{/if}

	<form
		on:submit|preventDefault={handleLogin}
		style="display: flex; flex-direction: column; gap: 15px; mt-4"
	>
		<div>
			<label style="display: block; margin-bottom: 5px;">Usuario:</label>
			<input type="text" bind:value={usuario} required style="width: 100%; padding: 8px;" />
		</div>

		<div>
			<label style="display: block; margin-bottom: 5px;">Contraseña:</label>
			<input type="password" bind:value={contrasena} required style="width: 100%; padding: 8px;" />
		</div>

		<button
			type="submit"
			style="padding: 10px; cursor: pointer; background: #007bff; color: white; border: none;"
		>
			Entrar
		</button>
	</form>
</div>
