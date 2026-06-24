<script lang="ts">
	import { apiAuth } from '$lib/api';
	import { auth } from '$lib/authStore.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	let usuario = $state('');
	let contrasena = $state('');
	let errorMsg = $state('');
	let loading = $state(false);

	onMount(() => {
		auth.init();
		if (auth.user) {
			goto('/');
		}
	});

	async function handleLogin(e: Event) {
		e.preventDefault();
		if (!usuario || !contrasena) {
			errorMsg = 'Debe ingresar usuario y contraseña.';
			return;
		}

		errorMsg = '';
		loading = true;

		try {
			const res = await apiAuth.login(usuario, contrasena);
			auth.login(res.empleado);
			goto('/');
		} catch (err: any) {
			errorMsg = err.message || 'Error al iniciar sesión.';
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Iniciar Sesión - GPSproject</title>
	<meta name="description" content="Página de acceso." />
</svelte:head>

<div class="login-wrapper">
	<div class="glow-orb orb-1"></div>
	<div class="glow-orb orb-2"></div>

	<div class="card login-card">
		<div class="login-header">
			<div class="logo-icon">
				<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
					<path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
					<circle cx="9" cy="7" r="4" />
					<path d="M22 21v-2a4 4 0 0 0-3-3.87" />
					<path d="M16 3.13a4 4 0 0 1 0 7.75" />
				</svg>
			</div>
			<h1 id="login-title">GPSproject</h1>
			<p>Ingresa tus datos</p>
		</div>

		{#if errorMsg}
			<div class="alert alert-error" role="alert" id="login-error">
				<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
					<circle cx="12" cy="12" r="10" />
					<line x1="12" y1="8" x2="12" y2="12" />
					<line x1="12" y1="16" x2="12.01" y2="16" />
				</svg>
				<span>{errorMsg}</span>
			</div>
		{/if}

		<form onsubmit={handleLogin} aria-labelledby="login-title">
			<div class="form-group">
				<label class="form-label" for="usuario">Usuario</label>
				<input
					type="text"
					id="usuario"
					class="form-input"
					placeholder="Ej: admin"
					bind:value={usuario}
					disabled={loading}
					required
				/>
			</div>

			<div class="form-group">
				<label class="form-label" for="contrasena">Contraseña</label>
				<input
					type="password"
					id="contrasena"
					class="form-input"
					placeholder="••••••••"
					bind:value={contrasena}
					disabled={loading}
					required
				/>
			</div>

			<button type="submit" class="btn btn-primary login-btn" disabled={loading} id="btn-login-submit">
				{#if loading}
					<span class="spinner"></span> Procesando...
				{:else}
					Ingresar
				{/if}
			</button>
		</form>
	</div>
</div>

<style>
	.login-wrapper {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 100vh;
		position: relative;
		padding: 20px;
		overflow: hidden;
	}

	.glow-orb {
		position: absolute;
		border-radius: 50%;
		filter: blur(100px);
		z-index: 0;
		opacity: 0.1;
	}

	.orb-1 {
		width: 350px;
		height: 350px;
		background: #eab308;
		top: 15%;
		left: 20%;
	}

	.orb-2 {
		width: 300px;
		height: 300px;
		background: #d97706;
		bottom: 15%;
		right: 20%;
	}

	.login-card {
		width: 100%;
		max-width: 420px;
		position: relative;
		z-index: 1;
		border: 1px solid var(--border-color);
		background: #ffffff;
		box-shadow: var(--shadow-lg);
		animation: card-appear 0.6s cubic-bezier(0.16, 1, 0.3, 1);
	}

	.login-header {
		text-align: center;
		margin-bottom: 28px;
	}

	.logo-icon {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 56px;
		height: 56px;
		background: rgba(217, 119, 6, 0.1);
		border: 1.5px solid rgba(217, 119, 6, 0.2);
		color: var(--accent-color);
		border-radius: 14px;
		margin-bottom: 16px;
	}

	.logo-icon svg {
		width: 28px;
		height: 28px;
	}

	.login-header h1 {
		font-size: 1.6rem;
		font-weight: 700;
		margin-bottom: 6px;
		letter-spacing: -0.025em;
		color: var(--text-primary);
	}

	.login-header p {
		color: var(--text-secondary);
		font-size: 0.9rem;
	}

	.login-btn {
		width: 100%;
		margin-top: 10px;
		padding: 12px;
		font-size: 0.95rem;
		background: var(--accent-gradient);
		color: #ffffff;
	}

	.spinner {
		display: inline-block;
		width: 16px;
		height: 16px;
		border: 2px solid rgba(255, 255, 255, 0.3);
		border-radius: 50%;
		border-top-color: white;
		animation: spin 0.8s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	@keyframes card-appear {
		from {
			opacity: 0;
			transform: translateY(20px) scale(0.98);
		}
		to {
			opacity: 1;
			transform: translateY(0) scale(1);
		}
	}
</style>
