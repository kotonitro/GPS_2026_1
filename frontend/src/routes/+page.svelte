<script lang="ts">
	import { auth } from '$lib/authStore.svelte';
	import { apiClientes } from '$lib/api';
	import { onMount } from 'svelte';

	let clientCount = $state(0);
	let loading = $state(true);

	onMount(async () => {
		try {
			if (auth.user) {
				const res = await apiClientes.getAll();
				clientCount = Array.isArray(res) ? res.length : 0;
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

<div class="page-header">
	<div>
		<h1 class="page-title">Modulos</h1>
		<p class="page-subtitle">Bienvenido, {auth.user?.usuario || 'usuario'}</p>
	</div>
</div>

<div class="dashboard-grid">
	

	<a href="/clientes" class="card stat-card">
		<div class="stat-icon">
			<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
				<path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
				<circle cx="9" cy="7" r="4"/>
				<path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
				<path d="M16 3.13a4 4 0 0 1 0 7.75"/>
			</svg>
		</div>
		<div class="stat-content">
			<span class="stat-label">Clientes</span>
			{#if loading}
				<span class="skeleton skeleton-text" style="width: 60px; height: 2rem; margin-top: 8px;"></span>
			{:else}
				<span class="stat-value">{clientCount}</span>
			{/if}
		</div>
		<div class="stat-arrow">
			<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
				<line x1="5" y1="12" x2="19" y2="12"/>
				<polyline points="12 5 19 12 12 19"/>
			</svg>
		</div>
	</a>





	<div class="card stat-card construction">
		<div class="stat-icon icon-ventas">
			<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
				<circle cx="9" cy="21" r="1"/>
				<circle cx="20" cy="21" r="1"/>
				<path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"/>
			</svg>
		</div>
		<div class="stat-content">
			<span class="stat-label">Ventas</span>
			<span class="stat-value text-construction">falra</span>
		</div>
	</div>



	
	<div class="card stat-card construction">
		<div class="stat-icon icon-inventario">
			<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
				<polygon points="12 22.08 12 12 3 6.92 3 17.08 12 22.08"/>
				<polygon points="12 12 21 6.92 21 17.08 12 22.08"/>
				<polygon points="12 2 3 6.92 12 12 21 6.92 12 2"/>
			</svg>
		</div>
		<div class="stat-content">
			<span class="stat-label">Inventario</span>
			<span class="stat-value text-construction">falta</span>
		</div>
	</div>
</div>

<div class="info-section">
	<div class="card info-card">
		<h2>Gestión de Clientes</h2>
		
		<div class="action-shortcuts">
			<a href="/clientes" class="btn btn-primary">Ir a Clientes</a>
		</div>
	</div>
</div>

<style>
	.dashboard-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
		gap: 24px;
		margin-bottom: 32px;
	}

	.stat-card {
		display: flex;
		align-items: center;
		text-decoration: none;
		color: inherit;
		gap: 20px;
		position: relative;
		overflow: hidden;
		transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
	}

	.stat-card:hover {
		transform: translateY(-2px);
		border-color: var(--accent-color);
	}

	.stat-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 48px;
		height: 48px;
		background: rgba(99, 102, 241, 0.1);
		color: var(--accent-color);
		border-radius: 12px;
	}

	.icon-ventas {
		background: rgba(16, 185, 129, 0.1);
		color: var(--success-color);
	}

	.icon-inventario {
		background: rgba(245, 158, 11, 0.1);
		color: var(--warning-color);
	}

	.stat-content {
		flex: 1;
		display: flex;
		flex-direction: column;
	}

	.stat-label {
		font-size: 0.85rem;
		font-weight: 500;
		color: var(--text-secondary);
	}

	.stat-value {
		font-size: 1.8rem;
		font-weight: 700;
		color: var(--text-primary);
		line-height: 1.2;
		margin-top: 4px;
	}

	.text-construction {
		font-size: 1.25rem;
		color: var(--text-muted);
	}

	.stat-card.construction {
		cursor: default;
		opacity: 0.8;
	}

	.stat-card.construction:hover {
		transform: none;
		border-color: var(--border-color);
	}

	.stat-arrow {
		color: var(--text-muted);
		opacity: 0;
		transform: translateX(-5px);
		transition: all 0.2s ease;
	}

	.stat-card:hover .stat-arrow {
		opacity: 1;
		transform: translateX(0);
		color: var(--accent-color);
	}

	.info-section {
		max-width: 800px;
	}

	.info-card h2 {
		font-size: 1.25rem;
		font-weight: 600;
		margin-bottom: 12px;
		color: var(--text-primary);
	}

	.info-card p {
		color: var(--text-secondary);
		font-size: 0.95rem;
		margin-bottom: 20px;
		line-height: 1.6;
	}

	.action-shortcuts {
		display: flex;
		gap: 12px;
	}
</style>
