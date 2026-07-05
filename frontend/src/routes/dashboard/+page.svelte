<script lang="ts">
	import { ShoppingCart, Users, Wallet, Package, AlertCircle } from '@lucide/svelte';
	import { onMount } from 'svelte';
	import { apiClientes } from '$lib/api';

	// Mock Data for Dashboard Visualization (Mixed with real for Fiados)
	let ventasHoy = $state(3830);
	let clientesAtendidos = $state(47);
	let fiadosPendientes = $state(0);
	let clientesDeudores = $state(0);
	let productosAgotados = $state(2);

	onMount(async () => {
		try {
			const res = await apiClientes.getAll();
			const clientes = Array.isArray(res) ? res : [];
			let totalFiados = 0;
			let countDeudores = 0;
			
			for (const c of clientes) {
				if ((c.fiado_actual || 0) > 0) {
					totalFiados += (c.fiado_actual || 0);
					countDeudores++;
				}
			}
			
			fiadosPendientes = totalFiados;
			clientesDeudores = countDeudores;
		} catch (error) {
			console.error('Error al cargar clientes para el dashboard:', error);
		}
	});

	let ultimasVentas = $state([
		{ cliente: 'Jorge Mendoza', monto: 102, fecha: new Date(), estado: 'Efectivo', desc: 'Manzana 2kg + Leche + Pan' },
		{ cliente: 'Doña Rosa López', monto: 56, fecha: new Date(Date.now() - 3600000), estado: 'Fiado', desc: 'Jitomate 1kg + Cebolla + Refresco' },
		{ cliente: 'Público general', monto: 66, fecha: new Date(Date.now() - 7200000), estado: 'Efectivo', desc: 'Gansito x3 + Agua x2' },
		{ cliente: 'Lucía Ramírez', monto: 141, fecha: new Date(Date.now() - 10800000), estado: 'Efectivo', desc: 'Huevo + Leche + Pan Bimbo' }
	]);

	let listaPromociones = $state([
		{ nombre: 'Promo Verano 2x1', descuento: '50%', vence: new Date(Date.now() + 86400000 * 5) },
		{ nombre: 'Descuento Bebidas', descuento: '20%', vence: new Date(Date.now() + 86400000 * 2) },
		{ nombre: 'Pack Asado', descuento: '15%', vence: new Date(Date.now() + 86400000 * 10) }
	]);

	function formatCurrency(amount: number) {
		return new Intl.NumberFormat('es-CL', {
			style: 'currency',
			currency: 'CLP',
			maximumFractionDigits: 0
		}).format(amount);
	}

	function formatDate(date: Date) {
		const d = new Date(date);
		const day = String(d.getDate()).padStart(2, '0');
		const month = String(d.getMonth() + 1).padStart(2, '0');
		return `${day}/${month}/${d.getFullYear()}`;
	}
</script>

<svelte:head>
	<title>Dashboard General - GPSproject</title>
</svelte:head>

<!-- Metrics Cards -->
<div class="mb-8 grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-4">
	<!-- Card 1 -->
	<div class="flex flex-col justify-center rounded-xl border border-border-color bg-bg-card p-6 shadow-sm transition-all hover:border-border-color-hover">
		<div class="flex items-center justify-between mb-4">
			<div class="flex h-10 w-10 items-center justify-center rounded-lg bg-primario/10 text-primario">
				<ShoppingCart size={20} />
			</div>
			<span class="rounded-full bg-exito/10 px-2 py-0.5 text-xs font-bold text-exito">+12%</span>
		</div>
		<h3 class="text-3xl font-bold text-text-primary">{formatCurrency(ventasHoy)}</h3>
		<p class="mt-1 text-xs font-bold uppercase tracking-wider text-text-secondary">Ventas Hoy</p>
	</div>

	<!-- Card 2 -->
	<div class="flex flex-col justify-center rounded-xl border border-border-color bg-bg-card p-6 shadow-sm transition-all hover:border-border-color-hover">
		<div class="flex items-center justify-between mb-4">
			<div class="flex h-10 w-10 items-center justify-center rounded-lg bg-exito/10 text-exito">
				<Users size={20} />
			</div>
			<span class="rounded-full bg-exito/10 px-2 py-0.5 text-xs font-bold text-exito">+5 vs ayer</span>
		</div>
		<h3 class="text-3xl font-bold text-text-primary">{clientesAtendidos}</h3>
		<p class="mt-1 text-xs font-bold uppercase tracking-wider text-text-secondary">Clientes Atendidos</p>
	</div>

	<!-- Card 3 -->
	<div class="flex flex-col justify-center rounded-xl border border-border-color bg-bg-card p-6 shadow-sm transition-all hover:border-border-color-hover">
		<div class="flex items-center justify-between mb-4">
			<div class="flex h-10 w-10 items-center justify-center rounded-lg bg-danger-color/10 text-danger-color">
				<Wallet size={20} />
			</div>
			<span class="rounded-full bg-danger-color/10 px-2 py-0.5 text-xs font-bold text-danger-color">~ {clientesDeudores} clientes</span>
		</div>
		<h3 class="text-3xl font-bold text-text-primary">{formatCurrency(fiadosPendientes)}</h3>
		<p class="mt-1 text-xs font-bold uppercase tracking-wider text-text-secondary">Fiados Pendientes</p>
	</div>

	<!-- Card 4 -->
	<div class="flex flex-col justify-center rounded-xl border border-border-color bg-bg-card p-6 shadow-sm transition-all hover:border-border-color-hover">
		<div class="flex items-center justify-between mb-4">
			<div class="flex h-10 w-10 items-center justify-center rounded-lg bg-orange-500/10 text-orange-500">
				<AlertCircle size={20} />
			</div>
			<span class="rounded-full bg-danger-color/10 px-2 py-0.5 text-xs font-bold text-danger-color">~ requieren restock</span>
		</div>
		<h3 class="text-3xl font-bold text-text-primary">{productosAgotados}</h3>
		<p class="mt-1 text-xs font-bold uppercase tracking-wider text-text-secondary">Productos Agotados</p>
	</div>
</div>

<div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
	<!-- Últimas Ventas -->
	<div class="flex flex-col rounded-xl border border-border-color bg-bg-card shadow-md">
		<div class="border-b border-border-color p-5">
			<h2 class="flex items-center gap-2 text-lg font-bold text-text-primary">
				Últimas Ventas
			</h2>
			<a href="/dashboard/ventas" class="text-sm font-semibold text-text-secondary hover:text-primario">Ver todas ↗</a>
		</div>
		<div class="flex flex-col gap-0 p-0">
			{#each ultimasVentas as venta}
				<div class="flex items-center justify-between border-b border-border-color p-5 last:border-0 hover:bg-text-primary/[0.015] transition-colors">
					<div>
						<h4 class="font-bold text-text-primary">{venta.cliente}</h4>
						<p class="text-sm text-text-secondary">{venta.desc}</p>
					</div>
					<div class="text-right">
						<h4 class="font-bold text-text-primary">{formatCurrency(venta.monto)}</h4>
						<p class="text-xs font-bold {venta.estado === 'Fiado' ? 'text-danger-color' : 'text-text-secondary'}">{venta.estado}</p>
					</div>
				</div>
			{/each}
		</div>
	</div>

	<!-- Promociones Activas -->
	<div class="flex flex-col rounded-xl border border-border-color bg-bg-card shadow-md">
		<div class="border-b border-border-color p-5">
			<h2 class="flex items-center gap-2 text-lg font-bold text-text-primary">
				<Package size={20} class="text-primario" />
				Promociones Activas
			</h2>
		</div>
		<div class="flex flex-col gap-4 p-5">
			{#each listaPromociones as promo}
				<div class="flex items-center justify-between rounded-lg border border-border-color p-4 transition-colors hover:border-border-color-hover">
					<div>
						<h4 class="font-bold text-text-primary">{promo.nombre}</h4>
						<p class="text-sm text-text-secondary">Vence: {formatDate(promo.vence)}</p>
					</div>
					<div class="flex h-10 w-16 items-center justify-center rounded-lg bg-primario/10 font-bold text-primario">
						{promo.descuento}
					</div>
				</div>
			{/each}
		</div>
	</div>
</div>
