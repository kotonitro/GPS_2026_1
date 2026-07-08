<script lang="ts">
	import { ShoppingCart, Users, Wallet, Package, AlertCircle, Receipt } from '@lucide/svelte';
	import { onMount } from 'svelte';
	import {
		apiClientes,
		apiProductos,
		apiVentas,
		obtenerPromociones,
		apiCajas,
		apiEmpleados
	} from '$lib/api';

	let ventasHoy = $state(0);
	let clientesAtendidos = $state(0);
	let ticketPromedio = $state(0);
	let fiadosPendientes = $state(0);
	let clientesDeudores = $state(0);
	let productosAgotados = $state(0);

	let ultimasVentas = $state<any[]>([]);
	let productos = $state<any[]>([]);
	let promociones = $state<any[]>([]);
	let hoy = $derived(new Date());

	let promocionesActivas = $derived(
		promociones
			.filter((p) => {
				if (!p.fecha_inicio || !p.fecha_fin) return true;
				return new Date(p.fecha_inicio) <= hoy && new Date(p.fecha_fin) >= hoy;
			})
			.slice(0, 3)
	);

	function obtenerProducto(id: string) {
		return productos.find((p) => p.id_producto === id);
	}

	function generarNombresCombo(promo: any) {
		if (!promo.productos_combo || promo.productos_combo.length === 0) return 'Productos del combo';
		const nombres = promo.productos_combo
			.map((id: string) => obtenerProducto(id)?.nombre || 'Producto')
			.join(' + ');
		return nombres;
	}

	function generarTitulo(promo: any, nombreProd: string) {
		if (promo.tipo === 'NXM') return `${promo.lleva}x${promo.paga} en ${nombreProd}`;
		if (promo.tipo === 'porcentaje') return `${promo.descuento}% en ${nombreProd}`;
		if (promo.tipo === 'COMBO') return `Combo Especial`;
		return `$${promo.descuento} dcto. en ${nombreProd}`;
	}

	function calcularDestaque(promo: any, producto: any) {
		if (promo.tipo === 'NXM') {
			return `${promo.lleva}x${promo.paga}`;
		} else if (promo.tipo === 'porcentaje') {
			return `${promo.descuento}%`;
		} else if (promo.tipo === 'COMBO') {
			return `$${promo.descuento}`;
		} else {
			if (producto && producto.precio > 0) {
				return Math.round((promo.descuento / producto.precio) * 100) + '%';
			}
			return `$${promo.descuento}`;
		}
	}

	onMount(async () => {
		// 1. Cargar clientes y calcular fiados pendientes
		try {
			const res = await apiClientes.getAll();
			const clientes = Array.isArray(res) ? res : [];
			let totalFiados = 0;
			let countDeudores = 0;

			for (const c of clientes) {
				if ((c.fiado_actual || 0) > 0) {
					totalFiados += c.fiado_actual || 0;
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
			productos = Array.isArray(resProductos) ? resProductos : [];
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
			// Cargar empleados para mapear ID -> Nombre
			let mapaEmpleados = new Map<string, string>();
			try {
				const resEmp = await apiEmpleados.getAll();
				const emps = Array.isArray(resEmp) ? resEmp : [];
				for (const e of emps) {
					mapaEmpleados.set(e.id_empleado, `${e.nombre} ${e.apellido || ''}`.trim());
				}
			} catch (error) {
				console.error('Error al cargar empleados para mapear nombres:', error);
			}

			// Cargar cajas para mapear ID -> Nombre
			let mapaCajas = new Map<string, string>();
			try {
				const resCajas = await apiCajas.getAll();
				const cjs = Array.isArray(resCajas) ? resCajas : [];
				for (const c of cjs) {
					mapaCajas.set(c.id_caja, c.nombre);
				}
			} catch (error) {
				console.error('Error al cargar cajas para mapear nombres:', error);
			}

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
			ticketPromedio = cantidadHoy > 0 ? sumaHoy / cantidadHoy : 0;

			// Ordenar por fecha descendente y tomar las últimas 4
			const ordenadas = [...ventas].sort((a, b) => {
				return new Date(b.fecha_emision).getTime() - new Date(a.fecha_emision).getTime();
			});

			ultimasVentas = ordenadas.slice(0, 4).map((v) => {
				const empNombre = mapaEmpleados.get(v.id_empleado) || 'Público general';
				const cajaNombre = mapaCajas.get(v.id_caja) || 'Caja';
				return {
					empleado: empNombre,
					caja: cajaNombre,
					monto: v.monto_total,
					fecha: new Date(v.fecha_emision),
					estado:
						v.metodo_pago && v.metodo_pago.nombre_metodo ? v.metodo_pago.nombre_metodo : 'Efectivo',
					desc: `Venta #${v.id_venta.slice(0, 8)}...`
				};
			});
		} catch (error) {
			console.error('Error al cargar ventas para el dashboard:', error);
		}

		// 4. Cargar promociones activas reales
		try {
			const resPromos = await obtenerPromociones();
			promociones = Array.isArray(resPromos) ? resPromos : [];
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

	function formatDate(date: Date | string) {
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
	<div
		class="flex flex-col justify-center rounded-xl border border-border-color bg-bg-card p-6 shadow-sm transition-all hover:border-border-color-hover"
	>
		<div class="flex items-center justify-between mb-4">
			<div
				class="flex h-10 w-10 items-center justify-center rounded-lg bg-primario/10 text-primario"
			>
				<ShoppingCart size={20} />
			</div>
		</div>
		<h3 class="text-3xl font-bold text-text-primary">{formatCurrency(ventasHoy)}</h3>
		<p class="mt-1 text-xs font-bold uppercase tracking-wider text-text-secondary">Ventas Hoy</p>
	</div>

	<!-- Clientes Atendidos -->
	<div
		class="flex flex-col justify-center rounded-xl border border-border-color bg-bg-card p-6 shadow-sm transition-all hover:border-border-color-hover"
	>
		<div class="flex items-center justify-between mb-4">
			<div class="flex h-10 w-10 items-center justify-center rounded-lg bg-exito/10 text-exito">
				<Users size={20} />
			</div>
		</div>
		<h3 class="text-3xl font-bold text-text-primary">{clientesAtendidos}</h3>
		<p class="mt-1 text-xs font-bold uppercase tracking-wider text-text-secondary">
			Clientes Atendidos
		</p>
	</div>

	<!-- Fiados Pendientes -->
	<div
		class="flex flex-col justify-center rounded-xl border border-border-color bg-bg-card p-6 shadow-sm transition-all hover:border-border-color-hover"
	>
		<div class="flex items-center justify-between mb-4">
			<div
				class="flex h-10 w-10 items-center justify-center rounded-lg bg-danger-color/10 text-danger-color"
			>
				<Wallet size={20} />
			</div>
			<span class="rounded-full bg-danger-color/10 px-2 py-0.5 text-xs font-bold text-danger-color"
				>{clientesDeudores} clientes</span
			>
		</div>
		<h3 class="text-3xl font-bold text-text-primary">{formatCurrency(fiadosPendientes)}</h3>
		<p class="mt-1 text-xs font-bold uppercase tracking-wider text-text-secondary">
			Fiados Pendientes
		</p>
	</div>

	<!-- Agotados / Bajo Stock -->
	<div
		class="flex flex-col justify-center rounded-xl border border-border-color bg-bg-card p-6 shadow-sm transition-all hover:border-border-color-hover"
	>
		<div class="flex items-center justify-between mb-4">
			<div
				class="flex h-10 w-10 items-center justify-center rounded-lg bg-orange-500/10 text-orange-500"
			>
				<AlertCircle size={20} />
			</div>
			<span class="rounded-full bg-danger-color/10 px-2 py-0.5 text-xs font-bold text-danger-color"
				>requieren restock</span
			>
		</div>
		<h3 class="text-3xl font-bold text-text-primary">{productosAgotados}</h3>
		<p class="mt-1 text-xs font-bold uppercase tracking-wider text-text-secondary">
			Productos Agotados / Bajo Stock
		</p>
	</div>
</div>

<div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
	<div class="flex flex-col rounded-xl border border-border-color bg-bg-card shadow-md">
		<div class="border-b border-border-color p-5">
			<h2 class="flex items-center gap-2 text-lg font-bold text-text-primary">
				<ShoppingCart size={20} class="text-primario" />
				Últimas Ventas
			</h2>
			<a
				href="/dashboard/ventas"
				class="text-sm font-semibold text-text-secondary hover:text-primario">Ver todas</a
			>
		</div>
		<div class="flex flex-col gap-0 p-0">
			{#each ultimasVentas as venta}
				<div
					class="flex items-center justify-between border-b border-border-color p-5 last:border-0 hover:bg-text-primary/[0.015] transition-colors"
				>
					<div>
						<h4 class="font-bold text-text-primary">{venta.empleado}</h4>
						<p class="text-sm text-text-secondary">{venta.caja} • {venta.desc}</p>
					</div>
					<div class="text-right">
						<h4 class="font-bold text-text-primary">{formatCurrency(venta.monto)}</h4>
						<p
							class="text-xs font-bold {venta.estado === 'Fiado'
								? 'text-danger-color'
								: 'text-text-secondary'}"
						>
							{venta.estado}
						</p>
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
			<a
				href="/dashboard/promociones"
				class="text-sm font-semibold text-text-secondary hover:text-primario">Ver todas</a
			>
		</div>
		<div class="flex flex-col gap-4 p-5">
			{#each promocionesActivas as promo}
				{@const nombreProd =
					promo.tipo === 'COMBO'
						? generarNombresCombo(promo)
						: obtenerProducto(promo.producto_id)?.nombre || 'Producto'}
				<div
					class="flex items-center justify-between rounded-lg border border-border-color p-4 transition-colors hover:border-border-color-hover"
				>
					<div>
						<h4
							class="font-bold text-text-primary line-clamp-1"
							title={generarTitulo(promo, nombreProd)}
						>
							{generarTitulo(promo, nombreProd)}
						</h4>
						<p class="text-sm text-text-secondary">
							Vence: {promo.fecha_fin ? formatDate(promo.fecha_fin) : 'Sin límite'}
						</p>
					</div>
					<div
						class="flex h-10 min-w-[64px] px-3 items-center justify-center rounded-lg bg-primario/10 font-bold text-primario"
					>
						{calcularDestaque(
							promo,
							promo.tipo !== 'COMBO' ? obtenerProducto(promo.producto_id) : null
						)}
					</div>
				</div>
			{:else}
				<p class="text-sm text-text-secondary py-2 text-center">No hay promociones activas.</p>
			{/each}
		</div>
	</div>
</div>
