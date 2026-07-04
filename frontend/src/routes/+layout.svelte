<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import { auth } from '$lib/authStore.svelte';
	import { toast } from '$lib/toastStore.svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import './layout.css';

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
<div class="toast-container fixed top-6 right-6 z-50 flex flex-col gap-2.5 max-w-[380px] w-[calc(100%-48px)]" aria-live="polite">
	{#each toast.toasts as t (t.id)}
		<div class="toast p-4 rounded-lg flex items-center justify-between gap-3 text-white font-semibold text-sm shadow-lg animate-fade-in {t.type === 'success' ? 'bg-[#10b981]' : 'bg-[#ef4444]'}">
			<span>{t.message}</span>
			<button class="toast-close cursor-pointer opacity-70 border-none bg-transparent text-white text-lg hover:opacity-100" onclick={() => toast.dismiss(t.id)} aria-label="Cerrar">&times;</button>
		</div>
	{/each}
</div>

{#if auth.loading}
	<!-- General Page Loader -->
	<div class="flex flex-col items-center justify-center min-h-screen bg-bg-primary text-text-primary gap-4">
		<div class="w-10 h-10 border-3 border-black/5 rounded-full border-t-accent animate-spin"></div>
		<p class="font-medium text-sm text-text-secondary">Cargando GPSproject...</p>
	</div>
{:else if $page.url.pathname === '/login'}
	<!-- Login layout (no sidebar) -->
	{@render children()}
{:else}
	<!-- Main Dashboard Layout -->
	<div class="flex min-h-screen">
		<aside class="w-[260px] bg-bg-secondary border-r border-border-color flex flex-col fixed top-0 bottom-0 left-0 z-10 transition-transform duration-300 md:translate-x-0">
			<div class="p-6 text-xl font-bold tracking-tight bg-gradient-to-r from-text-primary to-accent bg-clip-text text-transparent flex items-center gap-2.5 border-b border-border-color">
				<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
					<rect x="2" y="3" width="20" height="14" rx="2" ry="2"/>
					<line x1="8" y1="21" x2="16" y2="21"/>
					<line x1="12" y1="17" x2="12" y2="21"/>
				</svg>
				<span>GPSproject</span>
			</div>

			<nav class="flex-1 py-5 px-3 list-none flex flex-col gap-1.5" aria-label="Menú principal">
				<a href="/" class="flex items-center gap-3 p-3 text-text-secondary no-underline rounded-xl font-semibold text-sm transition-all duration-200 hover:bg-text-primary/4 hover:text-text-primary {$page.url.pathname === '/' ? 'bg-gradient-to-r from-accent-light to-accent text-white shadow-glow hover:text-white hover:bg-gradient-to-r' : ''}">
					<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
						<rect x="3" y="3" width="7" height="9"/>
						<rect x="14" y="3" width="7" height="5"/>
						<rect x="14" y="12" width="7" height="9"/>
						<rect x="3" y="16" width="7" height="5"/>
					</svg>
					<span>Main</span>
				</a>
				<a href="/clientes" class="flex items-center gap-3 p-3 text-text-secondary no-underline rounded-xl font-semibold text-sm transition-all duration-200 hover:bg-text-primary/4 hover:text-text-primary {$page.url.pathname.startsWith('/clientes') ? 'bg-gradient-to-r from-accent-light to-accent text-white shadow-glow hover:text-white hover:bg-gradient-to-r' : ''}">
					<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
						<path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
						<circle cx="9" cy="7" r="4"/>
						<path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
						<path d="M16 3.13a4 4 0 0 1 0 7.75"/>
					</svg>
					<span>Clientes</span>
				</a>
				<a href="/promociones" class="flex items-center gap-3 p-3 text-text-secondary no-underline rounded-xl font-semibold text-sm transition-all duration-200 hover:bg-text-primary/4 hover:text-text-primary {$page.url.pathname.startsWith('/promociones') ? 'bg-gradient-to-r from-accent-light to-accent text-white shadow-glow hover:text-white hover:bg-gradient-to-r' : ''}">
					<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
						<path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
						<circle cx="9" cy="7" r="4"/>
						<path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
						<path d="M16 3.13a4 4 0 0 1 0 7.75"/>
					</svg>
					<span>Promociones</span>
				</a>
				
				<!-- construcción / placeholders para coherencia modular -->
				<div class="text-[0.75rem] font-bold text-text-muted uppercase tracking-wider my-4 mx-3">Módulos</div>
				<a href="#" class="flex items-center gap-3 p-3 text-text-secondary/50 no-underline rounded-xl font-semibold text-sm transition-all duration-200 cursor-not-allowed" onclick={(e) => { e.preventDefault(); toast.show('Ventas falta.', 'error'); }}>
					<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
						<circle cx="9" cy="21" r="1"/>
						<circle cx="20" cy="21" r="1"/>
						<path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"/>
					</svg>
					<span>Ventas</span>
				</a>
				<a href="#" class="flex items-center gap-3 p-3 text-text-secondary/50 no-underline rounded-xl font-semibold text-sm transition-all duration-200 cursor-not-allowed" onclick={(e) => { e.preventDefault(); toast.show('Inventario falta.', 'error'); }}>
					<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
						<line x1="16.5" y1="9.4" x2="7.5" y2="4.21"/>
						<polygon points="12 22.08 12 12 3 6.92 3 17.08 12 22.08"/>
						<polygon points="12 12 21 6.92 21 17.08 12 22.08"/>
						<polygon points="12 2 3 6.92 12 12 21 6.92 12 2"/>
						<line x1="12" y1="22.08" x2="12" y2="12"/>
					</svg>
					<span>Inventario</span>
				</a>
				<a href="#" class="flex items-center gap-3 p-3 text-text-secondary/50 no-underline rounded-xl font-semibold text-sm transition-all duration-200 cursor-not-allowed" onclick={(e) => { e.preventDefault(); toast.show('Cajas falta.', 'error'); }}>
					<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
						<rect x="2" y="2" width="20" height="8" rx="2" ry="2"/>
						<rect x="2" y="14" width="20" height="8" rx="2" ry="2"/>
						<line x1="6" y1="6" x2="6.01" y2="6"/>
						<line x1="6" y1="18" x2="6.01" y2="18"/>
					</svg>
					<span>Cajas</span>
				</a>
			</nav>

			<footer class="p-4 border-t border-border-color flex flex-col gap-3 bg-text-primary/2">
				{#if auth.user}
					<div class="flex flex-col gap-0.5">
						<span class="font-semibold text-sm text-text-primary">{auth.user.usuario}</span>
						<span class="inline-flex items-center w-max px-2.5 py-0.5 rounded-full text-[0.75rem] font-bold border {auth.user.rol === 'Admin' ? 'bg-[#d97706]/10 text-[#b45309] border-[#d97706]/20' : 'bg-[#3b82f6]/10 text-[#1d4ed8] border-[#3b82f6]/20'}">
							{auth.user.rol}
						</span>
					</div>
				{/if}
				<button class="w-full flex items-center gap-2 font-semibold text-[0.85rem] py-2 px-3 border border-red-500/10 rounded-lg cursor-pointer bg-transparent text-[#3c4f6b] hover:bg-red-500/10 hover:text-red-500 hover:border-red-500/20 transition-all duration-200" onclick={handleLogout} id="btn-logout">
					<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
						<path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
						<polyline points="16 17 21 12 16 7"/>
						<line x1="21" y1="12" x2="9" y2="12"/>
					</svg>
					<span>Cerrar Sesión</span>
				</button>
			</footer>
		</aside>

		<main class="flex-1 ml-[260px] p-10 min-w-0 max-md:ml-0 max-md:p-5">
			{@render children()}
		</main>
	</div>
{/if}
