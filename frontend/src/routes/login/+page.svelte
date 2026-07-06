<script lang="ts">
	import { goto } from '$app/navigation';
	import { login } from '$lib/api';
	import { Store, Eye, EyeOff } from '@lucide/svelte';

	let usuario = $state('');
	let contrasena = $state('');

	let mensajeError = $state('');
	let mensajeExito = $state('');
	
	let loading = $state(false);
	let showPassword = $state(false);

	async function handleLogin() {
		mensajeError = '';
		mensajeExito = '';
		loading = true;

		try {
			await login(usuario, contrasena);
			goto('/dashboard');
		} catch (error: any) {
			mensajeError = error.message;
		} finally {
			loading = false;
		}
	}
	
	function togglePasswordVisibility(e: Event) {
		e.preventDefault();
		showPassword = !showPassword;
	}
</script>

<svelte:head>
	<title>Iniciar Sesión - GPS_2026_1</title>
</svelte:head>

<div class="flex min-h-screen w-full bg-bg-primary">
	<!-- Left Side - Image/Banner -->
	<div class="relative hidden w-1/2 bg-black lg:flex">
		<!-- Placeholder para imagen de fondo (Supermercado) -->
		<img 
			src="https://images.unsplash.com/photo-1542838132-92c53300491e?q=80&w=1974&auto=format&fit=crop" 
			alt="Fondo de supermercado" 
			class="absolute inset-0 h-full w-full object-cover opacity-50"
		/>
		<div class="absolute inset-0 bg-gradient-to-t from-black/80 to-transparent"></div>
		
		<div class="relative z-10 flex h-full flex-col justify-end p-12 pb-20 text-white">
			<h1 class="mb-4 text-4xl font-bold">El mejor sistema para<br/>tu tienda de barrio</h1>
			<p class="max-w-md text-lg text-[#c5b8ad]">
				Control de inventario, fiados, ventas y más.<br/>Todo en un solo lugar.
			</p>
		</div>
	</div>

	<!-- Right Side - Login Form -->
	<div class="flex w-full items-center justify-center bg-bg-primary p-8 lg:w-1/2">
		<div class="w-full max-w-md">
			<!-- Logo -->
			<div class="mb-10 flex items-center gap-3">
				<div class="flex h-12 w-12 items-center justify-center rounded-xl bg-accent/20 text-accent">
					<Store size={28} />
				</div>
				<div class="flex flex-col">
					<span class="tracking-tight text-2xl font-bold text-text-primary">GPS_2026_1</span>
					<span class="tracking-wider text-xs font-semibold uppercase text-accent">Gestión comercial</span>
				</div>
			</div>

			<h2 class="mb-2 text-3xl font-bold text-text-primary">Bienvenido</h2>
			<p class="mb-8 text-sm text-text-secondary">Ingresa tus datos para continuar</p>

			{#if mensajeError}
				<div class="mb-5 flex gap-3 rounded-lg border border-red-500/15 bg-danger-bg p-4 text-sm text-danger-color">
					<span>{mensajeError}</span>
				</div>
			{/if}

			{#if mensajeExito}
				<div class="mb-5 flex gap-3 rounded-lg border border-green-500/15 bg-exito/10 p-4 text-sm text-exito">
					<span>{mensajeExito}</span>
				</div>
			{/if}

			<form onsubmit={handleLogin} class="flex flex-col gap-5">
				<div>
					<label class="mb-2 block text-[0.65rem] font-bold uppercase tracking-wider text-text-secondary" for="usuario">
						USUARIO
					</label>
					<input
						id="usuario"
						type="text"
						bind:value={usuario}
						oninput={(e) => { e.currentTarget.value = e.currentTarget.value.replace(/[^a-zA-Z]/g, ''); usuario = e.currentTarget.value; }}
						placeholder="Ej: admin"
						required
						class="w-full rounded-lg border border-border-color bg-bg-card px-4 py-3 text-text-primary placeholder:text-text-muted/50 transition-all focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
					/>
				</div>

				<div>
					<label class="mb-2 block text-[0.65rem] font-bold uppercase tracking-wider text-text-secondary" for="contrasena">
						CONTRASEÑA
					</label>
					<div class="relative flex items-center">
						<!-- Solución para el toggle en Svelte: renderizar condicionalmente inputs separados -->
						{#if showPassword}
							<input
								id="contrasena-text"
								type="text"
								bind:value={contrasena}
								oninput={(e) => { e.currentTarget.value = e.currentTarget.value.replace(/\s/g, ''); contrasena = e.currentTarget.value; }}
								placeholder="Ej: Admin123."
								required
								class="w-full rounded-lg border border-border-color bg-bg-card px-4 py-3 pr-12 text-text-primary placeholder:text-text-muted/50 transition-all focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
							/>
						{:else}
							<input
								id="contrasena"
								type="password"
								bind:value={contrasena}
								oninput={(e) => { e.currentTarget.value = e.currentTarget.value.replace(/\s/g, ''); contrasena = e.currentTarget.value; }}
								placeholder="••••••••"
								required
								class="w-full rounded-lg border border-border-color bg-bg-card px-4 py-3 pr-12 text-text-primary placeholder:text-text-muted/50 transition-all focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
							/>
						{/if}
						<button 
							type="button" 
							class="absolute right-4 cursor-pointer text-text-muted transition-colors hover:text-text-primary"
							onclick={togglePasswordVisibility}
							aria-label={showPassword ? "Ocultar contraseña" : "Ver contraseña"}
						>
							{#if showPassword}
								<EyeOff size={18} />
							{:else}
								<Eye size={18} />
							{/if}
						</button>
					</div>
				</div>

				<button
					type="submit"
					disabled={loading}
					class="mt-4 w-full rounded-lg bg-primario px-4 py-3.5 font-bold text-white shadow-sm transition-all hover:bg-primario-hover disabled:cursor-not-allowed disabled:opacity-70"
				>
					{#if loading}
						Cargando...
					{:else}
						Iniciar sesión
					{/if}
				</button>
			</form>
		</div>
	</div>
</div>