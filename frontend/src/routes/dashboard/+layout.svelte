<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { checkSession, logout } from '$lib/api';
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
		Wallet
	} from '@lucide/svelte';

	let { children } = $props();

	let verificando = $state(true);
	let empleadoActual = $state<{ nombre: string; usuario: string; rol: string } | null>(null);

	onMount(async () => {
		try {
			empleadoActual = await checkSession();
			verificando = false;
		} catch (error) {
			goto('/login');
		}
	});

	async function handleLogout() {
		try {
			await logout();
		} catch (error) {
			console.error('Error al intentar cerrar sesión:', error);
		} finally {
			goto('/login');
		}
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
						<h2 class="text-lg font-semibold leading-tight text-white">GPS_2026_1</h2>
						<p class="text-xs font-medium text-primario">Gestión comercial</p>
					</div>
				</div>

				<nav class="flex flex-col gap-1">
					<span class="mb-1 mt-2 px-3 text-xs font-bold uppercase tracking-wider text-[#7a7268]">
						Principal
					</span>

					<a
						href="/dashboard"
						class="relative flex items-center gap-3 rounded-xl border border-[#4a3a28] bg-[#382a1b] p-3 text-primario transition-colors"
					>
						<LayoutDashboard size={20} />
						<span class="font-medium">Dashboard</span>
						<span class="absolute right-4 h-1.5 w-1.5 rounded-full bg-primario"></span>
					</a>

					<a
						href="/dashboard/ventas"
						class="flex items-center gap-3 rounded-xl p-3 transition-colors hover:bg-[#241e1a] hover:text-white"
					>
						<ShoppingCart size={20} />
						<span class="font-medium">Ventas</span>
					</a>
					<a
						href="/dashboard/productos"
						class="flex items-center gap-3 rounded-xl p-3 transition-colors hover:bg-[#241e1a] hover:text-white"
					>
						<Package size={20} />
						<span class="font-medium">Productos</span>
					</a>
					<a
						href="/dashboard/clientes"
						class="flex items-center gap-3 rounded-xl p-3 transition-colors hover:bg-[#241e1a] hover:text-white"
					>
						<Users size={20} />
						<span class="font-medium">Clientes</span>
					</a>
					<a
						href="/dashboard/promociones"
						class="flex items-center gap-3 rounded-xl p-3 transition-colors hover:bg-[#241e1a] hover:text-white"
					>
						<Tag size={20} />
						<span class="font-medium">Promociones</span>
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
							class="flex items-center gap-3 rounded-xl p-3 transition-colors hover:bg-[#241e1a] hover:text-white"
						>
							<UserCog size={20} />
							<span class="font-medium">Empleados</span>
						</a>
						<a
							href="/dashboard/cajas"
							class="flex items-center gap-3 rounded-xl p-3 transition-colors hover:bg-[#241e1a] hover:text-white"
						>
							<Wallet size={20} />
							<span class="font-medium">Cajas</span>
						</a>
					{/if}
				</nav>
			</div>

			<div class="flex flex-col gap-6 border-t border-[#2a241f] pt-6">
				<div class="flex items-center gap-3 px-2">
					<div
						class="flex h-10 w-10 items-center justify-center rounded-full bg-primario text-sm font-bold text-white"
					>
						{empleadoActual?.usuario?.substring(0, 2).toUpperCase() || 'EM'}
					</div>
					<div class="flex flex-col">
						<span class="text-sm font-semibold capitalize text-white"
							>{empleadoActual?.nombre || 'Nombre'}</span
						>
						<span class="text-xs">{empleadoActual?.rol || 'Rol'}</span>
					</div>
				</div>

				<button
					onclick={handleLogout}
					class="flex items-center gap-3 px-2 text-left text-sm font-medium transition-colors hover:text-white"
				>
					<LogOut size={20} /> Cerrar sesión
				</button>
			</div>
		</aside>

		<main class="flex-1 overflow-y-auto">
			{@render children()}
		</main>
	</div>
{/if}
