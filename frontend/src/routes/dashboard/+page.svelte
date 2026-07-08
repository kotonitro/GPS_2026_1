<script lang="ts">
	import { ShoppingCart, Users, Wallet, Package, AlertCircle, Receipt } from '@lucide/svelte';
	import { onMount } from 'svelte';
	import { apiClientes, apiProductos, apiVentas, obtenerPromociones } from '$lib/api';

	let ventasHoy = $state(0);
	let clientesAtendidos = $state(0);
	let ticketPromedio = $state(0);
	let fiadosPendientes = $state(0);
	let clientesDeudores = $state(0);
	let productosAgotados = $state(0);
	let ultimasVentas = $state<any[]>([]);
	let listaPromociones = $state<any[]>([]);

	onMount(async () => {
		// 1. Cargar clientes y calcular fiados pendientes
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

		// 2. Cargar productos y contar alertas de stock
		try {
			const resProductos = await apiProductos.getAll();
			const productos = Array.isArray(resProductos) ? resProductos : [];
			let countAgotados = 0;
			
			for (const p of productos) {
				if (p.stock <= (p.stock_minimo || 0)) {
					countAgotados++;
				}
			}
			
			productosAgotados = countAgotados;
		} catch (error) {
			console.error('Error al cargar productos para el dashboard:', error);
		}

		// 3. Cargar ventas para Ventas Hoy, Clientes Atendidos, Ticket Promedio y últimas ventas
		try {
			const resVentas = await apiVentas.getAll();
			const ventas = Array.isArray(resVentas) ? resVentas : [];

			// Filtrar las de hoy (desde las 00:00:00 local)
			const inicioHoy = new Date();
			inicioHoy.setHours(0, 0, 0, 0);

			let sumaHoy = 0;
			let cantidadHoy = 0;

			for (const v of ventas) {
				const fechaVenta = new Date(v.fecha_emision);
				if (fechaVenta >= inicioHoy) {
					sumaHoy += v.monto_total || 0;
					cantidadHoy++;
				}
			}

			ventasHoy = sumaHoy;
			clientesAtendidos = cantidadHoy;
			ticketPromedio = cantidadHoy > 0 ? (sumaHoy / cantidadHoy) : 0;

			// Ordenar por fecha descendente y tomar las últimas 4
			const ordenadas = [...ventas].sort((a, b) => {
				return new Date(b.fecha_emision).getTime() - new Date(a.fecha_emision).getTime();
			});

			ultimasVentas = ordenadas.slice(0, 4).map(v => {
				return {
					cliente: v.id_empleado ? `Empleado: ${v.id_empleado.slice(0, 8)}...` : 'Público general',
					monto: v.monto_total,
					fecha: new Date(v.fecha_emision),
					estado: v.metodo_pago && v.metodo_pago.nombre_metodo ? v.metodo_pago.nombre_metodo : 'Efectivo',
					desc: `Venta #${v.id_venta.slice(0, 8)}...`
				};
			});
		} catch (error) {
			console.error('Error al cargar ventas para el dashboard:', error);
		}

		// 4. Cargar promociones activas reales
		try {
			const resPromos = await obtenerPromociones();
			const promos = Array.isArray(resPromos) ? resPromos : [];
			
			listaPromociones = promos.slice(0, 3).map(p => {
				let desc = 'Descuento';
				if (p.tipo === 'NXM') {
					desc = `${p.lleva}x${p.paga}`;
				} else if (p.tipo === 'porcentaje') {
					desc = `${p.descuento}%`;
				} else if (p.tipo === 'precio_fijo') {
					desc = `$${p.descuento}`;
				}
				
				return {
					nombre: p.tipo === 'NXM' ? `Promo NxM (${p.lleva}x${p.paga})` : `Oferta Especial`,
					descuento: desc,
					vence: p.fecha_fin ? new Date(p.fecha_fin) : new Date(Date.now() + 86400000 * 7)
				};
			});
		} catch (error) {
			console.error('Error al cargar promociones para el dashboard:', error);
		}
	});

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

<div class="mb-8 grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-4">
	<!-- Ventas Hoy -->
	<div class="flex flex-col justify-center rounded-xl border border-border-color bg-bg-card p-6 shadow-sm transition-all hover:border-border-color-hover">
		<div class="flex items-center justify-between mb-4">
			<div class="flex h-10 w-10 items-center justify-center rounded-lg bg-primario/10 text-primario">
				<ShoppingCart size={20} />
			</div>
		</div>
		<h3 class="text-3xl font-bold text-text-primary">{formatCurrency(ventasHoy)}</h3>
		<p class="mt-1 text-xs font-bold uppercase tracking-wider text-text-secondary">Ventas Hoy</p>
	</div>

	<!-- Clientes Atendidos -->
	<div class="flex flex-col justify-center rounded-xl border border-border-color bg-bg-card p-6 shadow-sm transition-all hover:border-border-color-hover">
		<div class="flex items-center justify-between mb-4">
			<div class="flex h-10 w-10 items-center justify-center rounded-lg bg-exito/10 text-exito">
				<Users size={20} />
			</div>
		</div>
		<h3 class="text-3xl font-bold text-text-primary">{clientesAtendidos}</h3>
		<p class="mt-1 text-xs font-bold uppercase tracking-wider text-text-secondary">Clientes Atendidos</p>
	</div>

	<!-- Fiados Pendientes -->
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

	<!-- Agotados / Bajo Stock -->
	<div class="flex flex-col justify-center rounded-xl border border-border-color bg-bg-card p-6 shadow-sm transition-all hover:border-border-color-hover">
		<div class="flex items-center justify-between mb-4">
			<div class="flex h-10 w-10 items-center justify-center rounded-lg bg-orange-500/10 text-orange-500">
				<AlertCircle size={20} />
			</div>
			<span class="rounded-full bg-danger-color/10 px-2 py-0.5 text-xs font-bold text-danger-color">~ requieren restock</span>
		</div>
		<h3 class="text-3xl font-bold text-text-primary">{productosAgotados}</h3>
		<p class="mt-1 text-xs font-bold uppercase tracking-wider text-text-secondary">Agotados / Bajo Stock</p>
	</div>
</div>

<div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
	
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
