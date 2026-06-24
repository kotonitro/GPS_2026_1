<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import { auth } from '$lib/authStore.svelte';
	import { toast } from '$lib/toastStore.svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import './index.css';

	let { children } = $props();

	onMount(() => {
		auth.init();
		
		// Protection check
		if (!auth.user && $page.url.pathname !== '/login') {
			goto('/login');
		}
	});

	// React to auth status changes to redirect accordingly
	$effect(() => {
		if (!auth.loading && !auth.user && $page.url.pathname !== '/login') {
			goto('/login');
		}
	});

	function handleLogout() {
		auth.logout();
		toast.show('Sesión cerrada con éxito', 'success');
	}
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<!-- Toast Notifications -->
<div class="toast-container" aria-live="polite">
	{#each toast.toasts as t (t.id)}
		<div class="toast toast-{t.type}">
			<span>{t.message}</span>
			<button class="toast-close" onclick={() => toast.dismiss(t.id)} aria-label="Cerrar">&times;</button>
		</div>
	{/each}
</div>

{#if auth.loading}
	
	<div class="page-loader">
		<div class="spinner"></div>
		<p>Cargando ...</p>
	</div>
{:else if $page.url.pathname === '/login'}
	<!-- Login layout (no sidebar) -->
	{@render children()}
{:else}
	


	<div class="app-container">
		<aside class="sidebar">
			<div class="sidebar-logo">
				<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
					<rect x="2" y="3" width="20" height="14" rx="2" ry="2"/>
					<line x1="8" y1="21" x2="16" y2="21"/>
					<line x1="12" y1="17" x2="12" y2="21"/>
				</svg>
				<span>GPSproject</span>
			</div>

			<nav class="sidebar-menu" aria-label="Menú principal">
				<a href="/" class="sidebar-link" class:active={$page.url.pathname === '/'}>
					<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
						<rect x="3" y="3" width="7" height="9"/>
						<rect x="14" y="3" width="7" height="5"/>
						<rect x="14" y="12" width="7" height="9"/>
						<rect x="3" y="16" width="7" height="5"/>
					</svg>
					<span>Main</span>
				</a>
				<a href="/clientes" class="sidebar-link" class:active={$page.url.pathname.startsWith('/clientes')}>
					<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
						<path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
						<circle cx="9" cy="7" r="4"/>
						<path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
						<path d="M16 3.13a4 4 0 0 1 0 7.75"/>
					</svg>
					<span>Clientes</span>
				</a>
				
				<!-- construcción / placeholders para coherencia modular -->
				<div class="sidebar-separator">Módulos</div>
				<a href="#" class="sidebar-link disabled" onclick={(e) => { e.preventDefault(); toast.show('Ventas falta.', 'error'); }}>
					<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
						<circle cx="9" cy="21" r="1"/>
						<circle cx="20" cy="21" r="1"/>
						<path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"/>
					</svg>
					<span>Ventas</span>
				</a>
				<a href="#" class="sidebar-link disabled" onclick={(e) => { e.preventDefault(); toast.show('Inventario falta.', 'error'); }}>
					<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
						<line x1="16.5" y1="9.4" x2="7.5" y2="4.21"/>
						<polygon points="12 22.08 12 12 3 6.92 3 17.08 12 22.08"/>
						<polygon points="12 12 21 6.92 21 17.08 12 22.08"/>
						<polygon points="12 2 3 6.92 12 12 21 6.92 12 2"/>
						<line x1="12" y1="22.08" x2="12" y2="12"/>
					</svg>
					<span>Inventario</span>
				</a>
				<a href="#" class="sidebar-link disabled" onclick={(e) => { e.preventDefault(); toast.show('Cajas falta.', 'error'); }}>
					<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
						<rect x="2" y="2" width="20" height="8" rx="2" ry="2"/>
						<rect x="2" y="14" width="20" height="8" rx="2" ry="2"/>
						<line x1="6" y1="6" x2="6.01" y2="6"/>
						<line x1="6" y1="18" x2="6.01" y2="18"/>
					</svg>
					<span>Cajas</span>
				</a>
			</nav>

			<footer class="sidebar-footer">
				{#if auth.user}
					<div class="user-profile">
						<span class="user-name">{auth.user.usuario}</span>
						<span class="user-role badge badge-{auth.user.rol === 'Admin' ? 'admin' : 'user'}">
							{auth.user.rol}
						</span>
					</div>
				{/if}
				<button class="btn btn-secondary btn-logout" onclick={handleLogout} id="btn-logout">
					<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
						<path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
						<polyline points="16 17 21 12 16 7"/>
						<line x1="21" y1="12" x2="9" y2="12"/>
					</svg>
					<span>Cerrar Sesión</span>
				</button>
			</footer>
		</aside>

		<main class="main-content">
			{@render children()}
		</main>
	</div>
{/if}

<style>
	.page-loader {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		min-height: 100vh;
		background: var(--bg-primary);
		color: var(--text-primary);
		gap: 16px;
	}

	.page-loader .spinner {
		width: 40px;
		height: 40px;
		border: 3px solid rgba(255, 255, 255, 0.05);
		border-radius: 50%;
		border-top-color: var(--accent-color);
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	.sidebar-separator {
		font-size: 0.75rem;
		font-weight: 700;
		color: var(--text-muted);
		text-transform: uppercase;
		letter-spacing: 0.05em;
		margin: 16px 12px 6px 12px;
	}

	.sidebar-link.disabled {
		opacity: 0.5;
	}

	.btn-logout {
		width: 100%;
		justify-content: flex-start;
		font-size: 0.85rem;
		padding: 8px 12px;
		border-color: rgba(239, 68, 68, 0.1);
	}

	.btn-logout:hover {
		background: rgba(239, 68, 68, 0.1);
		color: #ef4444;
		border-color: rgba(239, 68, 68, 0.2);
	}
</style>
