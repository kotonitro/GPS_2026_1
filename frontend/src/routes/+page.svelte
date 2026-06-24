<script lang="ts">
	import { auth } from '$lib/authStore.svelte';
	import { apiClientes, apiPromociones } from '$lib/api';
	import { onMount } from 'svelte';

	let clientCount = $state(0);
	let promocionesCount = $state(0);
	let loading = $state(true);

	onMount(async () => {
		try {
			if (auth.user) {
			const [resClientes, resPromos] = await Promise.all([
								apiClientes.getAll(),
								apiPromociones.getAll()
			]);
			clientCount = Array.isArray(resClientes) ? resClientes.length : 0;
			promocionesCount = Array.isArray(resPromos) ? resPromos.length : 0;
			}
		
		} catch (err) {
			console.error('Error al obtener la informacion de clientes:', err);
		} finally {
			loading = false;
		}
	});
</script>

<svelte:head>
	<title>GPSproject</title>
	<meta name="description" content="Principal" />
</svelte:head>

<div class="flex justify-between items-center mb-8 flex-wrap gap-4">
	<div>
		<h1 class="text-3xl font-bold tracking-tight text-text-primary">Modulos</h1>
		<p class="text-text-secondary text-sm mt-1">Bienvenido, {auth.user?.usuario || 'usuario'}</p>
	</div>
</div>

<div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
	<!-- Clientes Stat Card -->
	<a href="/clientes" class="group flex items-center gap-5 bg-bg-card border border-border-color rounded-xl p-6 shadow-md hover:border-border-color-hover hover:shadow-lg hover:-translate-y-0.5 transition-all duration-300 relative overflow-hidden text-inherit no-underline">
		<div class="flex items-center justify-center w-12 h-12 bg-accent/10 text-accent rounded-xl">
			<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
				<path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
				<circle cx="9" cy="7" r="4"/>
				<path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
				<path d="M16 3.13a4 4 0 0 1 0 7.75"/>
			</svg>
		</div>
		<div class="flex-1 flex flex-col">
			<span class="text-[0.85rem] font-medium text-text-secondary">Clientes</span>
			{#if loading}
				<span class="skeleton-loader h-8 w-14 rounded mt-1"></span>
			{:else}
				<span class="text-3xl font-bold text-text-primary mt-1 leading-none">{clientCount}</span>
			{/if}
		</div>
		<div class="text-text-muted opacity-0 -translate-x-[5px] group-hover:opacity-100 group-hover:translate-x-0 group-hover:text-accent transition-all duration-200">
			<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
				<line x1="5" y1="12" x2="19" y2="12"/>
				<polyline points="12 5 19 12 12 19"/>
			</svg>
		</div>
	</a>

	<!-- Ventas (construction) -->
	<div class="flex items-center gap-5 bg-bg-card border border-border-color rounded-xl p-6 shadow-md opacity-80 relative overflow-hidden">
		<div class="flex items-center justify-center w-12 h-12 bg-success-bg text-success-color rounded-xl">
			<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
				<circle cx="9" cy="21" r="1"/>
				<circle cx="20" cy="21" r="1"/>
				<path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"/>
			</svg>
		</div>
		<div class="flex-1 flex flex-col">
			<span class="text-[0.85rem] font-medium text-text-secondary">Ventas</span>
			<span class="text-xl font-bold text-text-muted mt-1 leading-none">falra</span>
		</div>
	</div>

	<!-- Inventario (construction) -->
	<div class="flex items-center gap-5 bg-bg-card border border-border-color rounded-xl p-6 shadow-md opacity-80 relative overflow-hidden">
		<div class="flex items-center justify-center w-12 h-12 bg-warning-bg text-warning-color rounded-xl">
			<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
				<polygon points="12 22.08 12 12 3 6.92 3 17.08 12 22.08"/>
				<polygon points="12 12 21 6.92 21 17.08 12 22.08"/>
				<polygon points="12 2 3 6.92 12 12 21 6.92 12 2"/>
			</svg>
		</div>
		<div class="flex-1 flex flex-col">
			<span class="text-[0.85rem] font-medium text-text-secondary">Inventario</span>
			<span class="text-xl font-bold text-text-muted mt-1 leading-none">falta</span>
		</div>
	</div>
</div>

<div class="max-w-[800px]">
	<div class="bg-bg-card border border-border-color rounded-xl p-6 shadow-md hover:border-border-color-hover hover:shadow-lg transition-all duration-300">
		<h2 class="text-xl font-bold text-text-primary mb-3">Gestión de Clientes</h2>
		<div class="flex gap-3">
			<a href="/clientes" class="inline-flex items-center justify-center gap-2 font-semibold text-sm px-5 py-2.5 rounded-lg cursor-pointer bg-gradient-to-r from-accent-light to-accent text-white transition-all duration-200 hover:-translate-y-[1px] hover:shadow-glow hover:text-white no-underline">
				Ir a Clientes
			</a>
		</div>
	</div>
</div>

<div class="max-w-[800px]">
	<div class="bg-bg-card border border-border-color rounded-xl p-6 shadow-md hover:border-border-color-hover hover:shadow-lg transition-all duration-300">
		<h2 class="text-xl font-bold text-text-primary mb-3">Gestión de Promociones</h2>
		<div class="flex gap-3">
			<a href="/promociones" class="inline-flex items-center justify-center gap-2 font-semibold text-sm px-5 py-2.5 rounded-lg cursor-pointer bg-gradient-to-r from-accent-light to-accent text-white transition-all duration-200 hover:-translate-y-[1px] hover:shadow-glow hover:text-white no-underline">
				Ir a Promociones
			</a>
		</div>
	</div>
</div>
