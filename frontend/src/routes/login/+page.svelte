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

<div class="flex items-center justify-center min-h-screen relative p-5 overflow-hidden">
	<!-- Glow Orbs -->
	<div class="absolute rounded-full filter blur-[100px] z-0 opacity-10 w-[350px] h-[350px] bg-accent-light top-[15%] left-[20%]"></div>
	<div class="absolute rounded-full filter blur-[100px] z-0 opacity-10 w-[300px] h-[300px] bg-accent bottom-[15%] right-[20%]"></div>

	<div class="w-full max-w-[420px] relative z-10 border border-border-color bg-white rounded-xl p-6 shadow-lg animate-fade-in">
		<div class="text-center mb-7">
			<div class="inline-flex items-center justify-center w-14 h-14 bg-accent/10 border border-accent/20 text-accent rounded-2xl mb-4">
				<svg xmlns="http://www.w3.org/2000/svg" class="w-7 h-7" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
					<path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
					<circle cx="9" cy="7" r="4" />
					<path d="M22 21v-2a4 4 0 0 0-3-3.87" />
					<path d="M16 3.13a4 4 0 0 1 0 7.75" />
				</svg>
			</div>
			<h1 id="login-title" class="text-2xl font-bold mb-1.5 tracking-tight text-text-primary">GPSproject</h1>
			<p class="text-text-secondary text-sm">Ingresa tus datos</p>
		</div>

		{#if errorMsg}
			<div class="p-4 rounded-lg flex gap-3 text-sm mb-5 bg-danger-bg text-danger-color border border-red-500/15" role="alert" id="login-error">
				<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
					<circle cx="12" cy="12" r="10" />
					<line x1="12" y1="8" x2="12" y2="12" />
					<line x1="12" y1="16" x2="12.01" y2="16" />
				</svg>
				<span>{errorMsg}</span>
			</div>
		{/if}

		<form onsubmit={handleLogin} aria-labelledby="login-title">
			<div class="flex flex-col gap-1.5 mb-5">
				<label class="text-[0.85rem] font-semibold text-text-secondary" for="usuario">Usuario</label>
				<input
					type="text"
					id="usuario"
					class="bg-white border border-[rgba(15,30,54,0.15)] rounded-lg px-4 py-3 text-text-primary text-sm outline-none transition-all duration-200 focus:border-accent focus:ring-2 focus:ring-accent/15 placeholder:text-text-muted/80 disabled:opacity-50 disabled:cursor-not-allowed"
					placeholder="Ej: admin"
					bind:value={usuario}
					disabled={loading}
					required
				/>
			</div>

			<div class="flex flex-col gap-1.5 mb-5">
				<label class="text-[0.85rem] font-semibold text-text-secondary" for="contrasena">Contraseña</label>
				<input
					type="password"
					id="contrasena"
					class="bg-white border border-[rgba(15,30,54,0.15)] rounded-lg px-4 py-3 text-text-primary text-sm outline-none transition-all duration-200 focus:border-accent focus:ring-2 focus:ring-accent/15 placeholder:text-text-muted/80 disabled:opacity-50 disabled:cursor-not-allowed"
					placeholder="••••••••"
					bind:value={contrasena}
					disabled={loading}
					required
				/>
			</div>

			<button type="submit" class="w-full mt-2.5 p-3 text-sm font-semibold text-white bg-gradient-to-r from-accent-light to-accent rounded-lg cursor-pointer hover:shadow-glow hover:-translate-y-[1px] disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none disabled:shadow-none transition-all duration-200" disabled={loading} id="btn-login-submit">
				{#if loading}
					<span class="inline-block w-4 h-4 border-2 border-white/30 rounded-full border-t-white animate-spin"></span> Procesando...
				{:else}
					Ingresar
				{/if}
			</button>
		</form>
	</div>
</div>
