<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { checkSession, logout, apiClientes } from '$lib/api';
	import { auth } from '$lib/authStore.svelte';
	import { toast } from '$lib/toastStore.svelte';
	import '../layout.css';
	import {
		Store,
		LayoutDashboard,
		ShoppingCart,
		Package,
		Users,
		Tag,
		UserCog,
		LogOut,
		Wallet,
		Sun,
		Moon,
		Bell,
		ClipboardList
	} from '@lucide/svelte';

	let { children } = $props();

	let verificando = $state(true);
	let empleadoActual = $state<{ nombre: string; usuario: string; rol: string } | null>(null);
	let isDark = $state(false);

	let isNotificationsOpen = $state(false);
	let notificaciones = $state<any[]>([]);
	let loadingNotificaciones = $state(false);

	let currentDate = $state('');

	let pageTitle = $derived.by(() => {
		const path = $page.url.pathname;
		if (path === '/dashboard') return 'Dashboard';
		if (path.startsWith('/dashboard/ventas')) return 'Ventas';
		if (path.startsWith('/dashboard/productos')) return 'Productos';
		if (path.startsWith('/dashboard/clientes')) return 'Clientes';
		if (path.startsWith('/dashboard/promociones')) return 'Promociones';
		if (path.startsWith('/dashboard/empleados')) return 'Empleados';
		if (path.startsWith('/dashboard/cajas')) return 'Cajas';
		return 'Dashboard';
	});

	onMount(async () => {
		try {
			if (typeof window !== 'undefined') {
				const savedTheme = localStorage.getItem('theme');
				if (
					savedTheme === 'dark' ||
					(!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)
				) {
					isDark = true;
					document.documentElement.classList.add('dark');
				}
			}

			currentDate = new Intl.DateTimeFormat('es-CL', {
				weekday: 'long',
				year: 'numeric',
				month: 'long',
				day: 'numeric'
			}).format(new Date());

			empleadoActual = await checkSession();
			auth.login(empleadoActual as any);

			try {
				const res = await apiClientes.getAll();
				const clientes = Array.isArray(res) ? res : [];
				notificaciones = clientes.filter(
					(c: any) => (c.fiado_actual || 0) >= (c.fiado_maximo || 20000)
				);
			} catch (e) {}

			verificando = false;
		} catch (error) {
			goto('/login');
		}
	});

	function toggleTheme() {
		isDark = !isDark;
		if (isDark) {
			document.documentElement.classList.add('dark');
			localStorage.setItem('theme', 'dark');
		} else {
			document.documentElement.classList.remove('dark');
			localStorage.setItem('theme', 'light');
		}
	}

	async function handleLogout() {
		try {
			await logout();
			auth.logout();
		} catch (error) {
			console.error('Error al intentar cerrar sesión:', error);
		} finally {
			goto('/login');
		}
	}

	async function toggleNotifications() {
		isNotificationsOpen = !isNotificationsOpen;
		if (isNotificationsOpen) {
			loadingNotificaciones = true;
			try {
				const res = await apiClientes.getAll();
				const clientes = Array.isArray(res) ? res : [];
				notificaciones = clientes.filter(
					(c: any) => (c.fiado_actual || 0) >= (c.fiado_maximo || 20000)
				);
			} catch (err) {
				console.error('Error fetching notifications:', err);
			} finally {
				loadingNotificaciones = false;
			}
		}
	}

	function isActive(path: string) {
		if (path === '/dashboard') {
			return $page.url.pathname === '/dashboard';
		}
		return $page.url.pathname.startsWith(path);
	}
</script>

{#if verificando}
	<div class="flex h-screen items-center justify-center text-xl text-primario">
		Verificando credenciales...
	</div>
{:else}
	<div class="flex h-screen w-screen overflow-hidden">
		<aside
			class="flex h-full w-[240px] flex-col justify-between border-r border-[#2a241f] bg-[#1a1512] px-5 py-6 text-[#a39b93]"
		>
			<div class="flex flex-col gap-8">
				<div class="flex items-center gap-3 px-2">
					<div class="flex items-center justify-center rounded-lg bg-[#382a1b] p-2 text-primario">
						<Store size={22} strokeWidth={2.5} />
					</div>
					<div class="flex flex-col">
						<h2 class="text-lg font-semibold leading-tight text-white">MinimarketGo</h2>
						<p class="text-xs font-medium text-primario">Gestión comercial</p>
					</div>
				</div>

				<nav class="flex flex-col gap-1">
					<span class="mb-1 mt-2 px-3 text-xs font-bold uppercase tracking-wider text-[#7a7268]">
						Principal
					</span>

					<a
						href="/dashboard"
						class="relative flex items-center gap-3 rounded-xl border p-3 transition-colors {isActive(
							'/dashboard'
						)
							? 'border-[#4a3a28] bg-[#382a1b] text-primario'
							: 'border-transparent hover:bg-[#241e1a] hover:text-white'}"
					>
						<LayoutDashboard size={20} />
						<span class="font-medium">Dashboard</span>
						{#if isActive('/dashboard')}
							<span class="absolute right-4 h-1.5 w-1.5 rounded-full bg-primario"></span>
						{/if}
					</a>

					<a
						href="/dashboard/ventas"
						class="relative flex items-center gap-3 rounded-xl border p-3 transition-colors {isActive(
							'/dashboard/ventas'
						)
							? 'border-[#4a3a28] bg-[#382a1b] text-primario'
							: 'border-transparent hover:bg-[#241e1a] hover:text-white'}"
					>
						<ShoppingCart size={20} />
						<span class="font-medium">Ventas</span>
						{#if isActive('/dashboard/ventas')}
							<span class="absolute right-4 h-1.5 w-1.5 rounded-full bg-primario"></span>
						{/if}
					</a>

					<a
						href="/dashboard/productos"
						class="relative flex items-center gap-3 rounded-xl border p-3 transition-colors {isActive(
							'/dashboard/productos'
						)
							? 'border-[#4a3a28] bg-[#382a1b] text-primario'
							: 'border-transparent hover:bg-[#241e1a] hover:text-white'}"
					>
						<Package size={20} />
						<span class="font-medium">Productos</span>
						{#if isActive('/dashboard/productos')}
							<span class="absolute right-4 h-1.5 w-1.5 rounded-full bg-primario"></span>
						{/if}
					</a>

					<a
						href="/dashboard/clientes"
						class="relative flex items-center gap-3 rounded-xl border p-3 transition-colors {isActive(
							'/dashboard/clientes'
						)
							? 'border-[#4a3a28] bg-[#382a1b] text-primario'
							: 'border-transparent hover:bg-[#241e1a] hover:text-white'}"
					>
						<Users size={20} />
						<span class="font-medium">Clientes</span>
						{#if isActive('/dashboard/clientes')}
							<span class="absolute right-4 h-1.5 w-1.5 rounded-full bg-primario"></span>
						{/if}
					</a>

					<a
						href="/dashboard/promociones"
						class="relative flex items-center gap-3 rounded-xl border p-3 transition-colors {isActive(
							'/dashboard/promociones'
						)
							? 'border-[#4a3a28] bg-[#382a1b] text-primario'
							: 'border-transparent hover:bg-[#241e1a] hover:text-white'}"
					>
						<Tag size={20} />
						<span class="font-medium">Promociones</span>
						{#if isActive('/dashboard/promociones')}
							<span class="absolute right-4 h-1.5 w-1.5 rounded-full bg-primario"></span>
						{/if}
					</a>

					{#if empleadoActual?.rol === 'Admin'}
						<div class="mt-4 border-t border-[#2a241f] pt-4">
							<span
								class="mb-2 block px-3 text-xs font-bold uppercase tracking-wider text-[#7a7268]"
							>
								Administración
							</span>
						</div>

						<a
							href="/dashboard/empleados"
							class="relative flex items-center gap-3 rounded-xl border p-3 transition-colors {isActive(
								'/dashboard/empleados'
							)
								? 'border-[#4a3a28] bg-[#382a1b] text-primario'
								: 'border-transparent hover:bg-[#241e1a] hover:text-white'}"
						>
							<UserCog size={20} />
							<span class="font-medium">Empleados</span>
							{#if isActive('/dashboard/empleados')}
								<span class="absolute right-4 h-1.5 w-1.5 rounded-full bg-primario"></span>
							{/if}
						</a>

						<a
							href="/dashboard/cajas"
							class="relative flex items-center gap-3 rounded-xl border p-3 transition-colors {isActive(
								'/dashboard/cajas'
							)
								? 'border-[#4a3a28] bg-[#382a1b] text-primario'
								: 'border-transparent hover:bg-[#241e1a] hover:text-white'}"
						>
							<Wallet size={20} />
							<span class="font-medium">Cajas</span>
							{#if isActive('/dashboard/cajas')}
								<span class="absolute right-4 h-1.5 w-1.5 rounded-full bg-primario"></span>
							{/if}
						</a>

						<a
							href="/dashboard/turnos"
							class="relative flex items-center gap-3 rounded-xl border p-3 transition-colors {isActive(
								'/dashboard/turnos'
							)
								? 'border-[#4a3a28] bg-[#382a1b] text-primario'
								: 'border-transparent hover:bg-[#241e1a] hover:text-white'}"
						>
							<ClipboardList size={20} />
							<span class="font-medium">Turnos</span>
							{#if isActive('/dashboard/turnos')}
								<span class="absolute right-4 h-1.5 w-1.5 rounded-full bg-primario"></span>
							{/if}
						</a>
					{/if}
				</nav>
			</div>

			<div class="flex flex-col gap-6 border-t border-[#2a241f] pt-6">
				<div class="flex items-center justify-between px-2">
					<div class="flex items-center gap-3">
						<div
							class="flex h-10 w-10 items-center justify-center rounded-full bg-primario text-sm font-bold text-white"
						>
							{empleadoActual?.usuario?.substring(0, 2).toUpperCase() || 'EM'}
						</div>
						<div class="flex flex-col">
							<span class="text-sm font-semibold text-white"
								>{empleadoActual?.nombre || 'Nombre'}</span
							>
							<span class="text-xs text-[#a39b93]">{empleadoActual?.rol || 'Rol'}</span>
						</div>
					</div>
					<button
						onclick={toggleTheme}
						class="rounded-lg p-2 text-[#a39b93] hover:bg-[#241e1a] hover:text-white transition-colors"
						title="Cambiar de modo (Oscuro/Claro)"
						aria-label="Cambiar tema"
					>
						{#if isDark}
							<Sun size={18} />
						{:else}
							<Moon size={18} />
						{/if}
					</button>
				</div>

				<button
					onclick={handleLogout}
					class="flex items-center gap-3 px-2 text-left text-sm font-medium hover:text-white"
				>
					<LogOut size={20} /> Cerrar sesión
				</button>
			</div>
		</aside>

		<main class="flex-1 overflow-y-auto bg-bg-primary">
			<header
				class="flex items-center justify-between border-b border-border-color bg-bg-card px-8 py-5"
			>
				<div>
					<h1 class="text-xl font-bold text-text-primary">{pageTitle}</h1>
					<p class="text-sm text-text-secondary">{currentDate}</p>
				</div>
				<div class="relative">
					<button
						class="relative rounded-full p-2 text-text-muted transition-colors hover:bg-border-color hover:text-text-primary"
						onclick={toggleNotifications}
						aria-label="Notificaciones"
					>
						<Bell size={20} />
						{#if notificaciones.length > 0}
							<span
								class="absolute right-2 top-2 h-2 w-2 rounded-full bg-danger-color border-2 border-bg-card"
							></span>
						{/if}
					</button>
					{#if isNotificationsOpen}
						<div
							class="fixed inset-0 z-40"
							onclick={() => (isNotificationsOpen = false)}
							role="presentation"
						></div>

						<div
							class="absolute right-0 mt-2 w-80 rounded-xl border border-border-color bg-bg-card shadow-xl z-50 animate-modal-enter"
						>
							<div class="border-b border-border-color p-4">
								<h3 class="font-bold text-text-primary">Notificaciones</h3>
							</div>
							<div class="max-h-80 overflow-y-auto p-2">
								{#if loadingNotificaciones}
									<div class="p-4 text-center text-sm text-text-muted">Cargando...</div>
								{:else if notificaciones.length === 0}
									<div class="p-4 text-center text-sm text-text-muted">No hay notificaciones.</div>
								{:else}
									<div class="mb-2 rounded-lg border border-red-500/15 bg-danger-bg p-3 text-sm">
										<p class="font-semibold text-danger-color">
											{notificaciones.length} persona{notificaciones.length === 1 ? '' : 's'}
										</p>
										<p class="text-xs text-danger-color/80 mt-0.5">
											han llegado al límite de fiado, dar aviso.
										</p>
									</div>
								{/if}
							</div>
						</div>
					{/if}
				</div>
			</header>

			<div class="p-8">
				{@render children()}
			</div>
		</main>
	</div>

	<div class="fixed right-5 top-5 z-[9999] flex flex-col gap-3 pointer-events-none">
		{#each toast.toasts as t (t.id)}
			<div
				class="pointer-events-auto flex w-80 items-center gap-3 rounded-xl border bg-bg-card p-4 shadow-lg animate-modal-enter {t.type ===
				'success'
					? 'border-green-500/20 text-exito'
					: 'border-red-500/20 text-danger-color'}"
			>
				{#if t.type === 'success'}
					<div class="flex h-8 w-8 items-center justify-center rounded-full bg-exito/10">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							width="16"
							height="16"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2.5"
							stroke-linecap="round"
							stroke-linejoin="round"
							><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" /><polyline
								points="22 4 12 14.01 9 11.01"
							/></svg
						>
					</div>
				{:else}
					<div class="flex h-8 w-8 items-center justify-center rounded-full bg-danger-color/10">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							width="16"
							height="16"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2.5"
							stroke-linecap="round"
							stroke-linejoin="round"
							><circle cx="12" cy="12" r="10" /><line x1="12" y1="8" x2="12" y2="12" /><line
								x1="12"
								y1="16"
								x2="12.01"
								y2="16"
							/></svg
						>
					</div>
				{/if}
				<div class="flex-1">
					<p class="text-sm font-semibold">{t.message}</p>
				</div>
				<button
					type="button"
					class="cursor-pointer text-text-muted hover:text-text-primary"
					onclick={() => toast.dismiss(t.id)}
					aria-label="Cerrar"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="16"
						height="16"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
						><line x1="18" y1="6" x2="6" y2="18" /><line x1="6" y1="6" x2="18" y2="18" /></svg
					>
				</button>
			</div>
		{/each}
	</div>
{/if}
