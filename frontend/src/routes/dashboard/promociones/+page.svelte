<script lang="ts">
	import { onMount } from 'svelte';
	import {
		obtenerPromociones,
		crearPromocion,
		obtenerProductos,
		eliminarPromocion,
		actualizarPromocion
	} from '$lib/api';
	import { auth } from '$lib/authStore.svelte'; 
	import { toast } from '$lib/toastStore.svelte'; 

	let promociones = $state<any[]>([]);
	let productos = $state<any[]>([]);
	let loading = $state(true);
	let errorMsg = $state('');

	let mostrarFormulario = $state(false);
	let mostrarModalEliminar = $state(false);

	// Variables del formulario
	let categoriaSeleccionada = $state('');
	let tipoPromocion = $state('NXM');
	let productoSeleccionado = $state<string>('');
	let comboSeleccionados = $state<string[]>([]);
	
	// Buscador de productos
	let busquedaProducto = $state('');

	let lleva = $state<number | string>('');
	let paga = $state<number | string>('');
	let descuento = $state<number | string>('');
	let editandoId = $state<string | null>(null);
	
	let fechaInicio = $state('');
	let fechaFin = $state('');
	let submitLoading = $state(false);

	let hoy = $derived(new Date());
	let hoyStr = $derived(hoy.toISOString().split('T')[0]);

	let categoriasUnicas = $derived([
		...new Set(productos.map((p) => p.categoria?.nombre_categoria))
	]);

	// buscador
	let productosFiltrados = $derived.by(() => {
		let result = productos;
		
		if (categoriaSeleccionada !== '') {
			result = result.filter((p) => p.categoria?.nombre_categoria === categoriaSeleccionada);
		}

		if (busquedaProducto.trim()) {
			const query = busquedaProducto.toLowerCase().trim();
			result = result.filter((p) => {
				const matchNombre = p.nombre.toLowerCase().includes(query);
				const matchCodigo = p.codigo_barras ? p.codigo_barras.toLowerCase().includes(query) : false;
				return matchNombre || matchCodigo;
			});
		}
		
		return result;
	});

	let activas = $derived(promociones.filter(p => {
		if (!p.fecha_inicio || !p.fecha_fin) return true;
		return new Date(p.fecha_inicio) <= hoy && new Date(p.fecha_fin) >= hoy;
	}));

	let proximas = $derived(promociones.filter(p => {
		if (!p.fecha_inicio) return false;
		return new Date(p.fecha_inicio) > hoy;
	}));

	let expiradas = $derived(promociones.filter(p => {
		if (!p.fecha_fin) return false;
		return new Date(p.fecha_fin) < hoy;
	}));

	let idAEliminar = $state<string | null>(null);

	let carruselActivas = $state<HTMLDivElement>();
	let carruselProximas = $state<HTMLDivElement>();
	let carruselExpiradas = $state<HTMLDivElement>();

	function moverCarrusel(carrusel: HTMLDivElement | undefined, direccion: 'izq' | 'der') {
		if (carrusel) {
			const desplazamiento = 340; 
			carrusel.scrollBy({ left: direccion === 'izq' ? -desplazamiento : desplazamiento, behavior: 'smooth' });
		}
	}

	onMount(async () => {
		await cargarDatos();
	});

	async function cargarDatos() {
		loading = true;
		errorMsg = '';
		try {
			productos = await obtenerProductos();
			promociones = await obtenerPromociones();
		} catch (error: any) {
			errorMsg = 'Error al cargar los datos: ' + error.message;
			if(toast) toast.show(errorMsg, 'error');
		} finally {
			loading = false;
		}
	}

	// selección del buscador
	function seleccionarProducto(id: string) {
		if (tipoPromocion === 'COMBO') {
			if (!comboSeleccionados.includes(id) && comboSeleccionados.length < 3) {
				comboSeleccionados = [...comboSeleccionados, id];
			} else if (comboSeleccionados.length >= 3) {
				if(toast) toast.show('Máximo 3 productos por combo', 'error');
			}
		} else {
			productoSeleccionado = id;
			busquedaProducto = '';
		}
	}

	function removerProductoCombo(id: string) {
		comboSeleccionados = comboSeleccionados.filter(pId => pId !== id);
	}

	// Limpiar selecciones si cambia el tipo de promoción
	$effect(() => {
		if (tipoPromocion) {
			productoSeleccionado = '';
			comboSeleccionados = [];
			busquedaProducto = '';
		}
	});

	async function handleCrearPromocion(e: Event) {
		e.preventDefault();
		
		if (tipoPromocion === 'COMBO' && comboSeleccionados.length < 2) {
			if(toast) toast.show('Un combo debe tener al menos 2 productos.', 'error');
			return;
		}
		if (tipoPromocion !== 'COMBO' && !productoSeleccionado) {
			if(toast) toast.show('Debes seleccionar un producto.', 'error');
			return;
		}

		submitLoading = true;

		try {
		const payload = {
            producto_id: tipoPromocion === 'COMBO' ? null : (productoSeleccionado || null),
            productos_combo: comboSeleccionados.length > 0 ? comboSeleccionados : [], 
            tipo: tipoPromocion,
            lleva: Number(lleva),
            paga: Number(paga),
            descuento: Number(descuento),
            fecha_inicio: fechaInicio ? new Date(fechaInicio).toISOString() : null,
            fecha_fin: fechaFin ? new Date(fechaFin).toISOString() : null
};

			if (editandoId) {
				await actualizarPromocion(editandoId, payload);
				if(toast) toast.show('Promoción actualizada exitosamente.', 'success');
			} else {
				await crearPromocion(payload);
				if(toast) toast.show('Promoción creada exitosamente.', 'success');
			}
			
			mostrarFormulario = false; 
			await cargarDatos();
		} catch (error: any) {
			if(toast) toast.show(error.message, 'error');
		} finally {
			submitLoading = false;
		}
	}

	function handleEditar(promo: any) {
		editandoId = promo.id_promocion;
		tipoPromocion = promo.tipo;
		
		if (promo.tipo === 'COMBO') {
			comboSeleccionados = promo.productos_combo || [];
		} else {
			productoSeleccionado = promo.producto_id;
			const prod = obtenerProducto(promo.producto_id);
			categoriaSeleccionada = prod?.categoria?.nombre_categoria || '';
		}
		
		lleva = promo.lleva || '';
		paga = promo.paga || '';
		descuento = promo.descuento || '';
		fechaInicio = promo.fecha_inicio ? promo.fecha_inicio.split('T')[0] : '';
		fechaFin = promo.fecha_fin ? promo.fecha_fin.split('T')[0] : '';

		mostrarFormulario = true;
	}

	function abrirModalNuevo() {
		editandoId = null;
		categoriaSeleccionada = '';
		productoSeleccionado = '';
		comboSeleccionados = [];
		busquedaProducto = '';
		tipoPromocion = 'NXM';
		fechaInicio = '';
		fechaFin = '';
		lleva = '';
		paga = '';
		descuento = '';
		mostrarFormulario = true;
	}

	function pedirConfirmacionEliminar(id: string) {
		idAEliminar = id;
		mostrarModalEliminar = true;
	}

	async function confirmarEliminacion() {
		if (!idAEliminar) return;
		try {
			await eliminarPromocion(idAEliminar);
			if(toast) toast.show('Promoción cancelada exitosamente.', 'success');
			await cargarDatos();
			mostrarModalEliminar = false;
			idAEliminar = null;
		} catch (error: any) {
			if(toast) toast.show(error.message, 'error');
			mostrarModalEliminar = false;
		}
	}

	function obtenerProducto(id: string) {
		return productos.find((p) => p.id_producto === id);
	}

	function formatearRangoFechas(inicio: string, fin: string) {
		if (!inicio || !fin) return 'Sin límite';
		const opciones: Intl.DateTimeFormatOptions = { day: '2-digit', month: 'short', year: 'numeric' };
		const fInicio = new Date(inicio).toLocaleDateString('es-ES', opciones);
		const fFin = new Date(fin).toLocaleDateString('es-ES', opciones);
		return `${fInicio} — ${fFin}`;
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

	function generarTitulo(promo: any, nombreProd: string) {
		if (promo.tipo === 'NXM') return `${promo.lleva}x${promo.paga} en ${nombreProd}`;
		if (promo.tipo === 'porcentaje') return `${promo.descuento}% en ${nombreProd}`;
		if (promo.tipo === 'COMBO') return `Combo Especial`; 
		return `$${promo.descuento} de dcto. en ${nombreProd}`;
	}

	function generarNombresCombo(promo: any) {
		if (!promo.productos_combo || promo.productos_combo.length === 0) return 'Productos del combo';
		const nombres = promo.productos_combo.map((id: string) => obtenerProducto(id)?.nombre || 'Producto').join(' + ');
		return nombres;
	}
</script>

<svelte:head>
	<title>Promociones - GPSproject</title>
</svelte:head>

<div class="mb-6 flex flex-wrap items-start justify-between gap-6">
	<div>
		<h1 class="text-2xl font-bold text-text-primary mb-1">Promociones</h1>
		<p class="text-sm text-text-secondary">{activas.length} activas · {proximas.length} próximas</p>
	</div>
	<div class="flex items-start">
		<button
			class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg bg-gradient-to-r from-accent-light to-accent px-5 py-2.5 text-sm font-semibold text-white transition-all duration-200 hover:-translate-y-[1px] hover:shadow-glow"
			onclick={abrirModalNuevo}
		>
			<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
				<line x1="12" y1="5" x2="12" y2="19" />
				<line x1="5" y1="12" x2="19" y2="12" />
			</svg>
			<span>Nueva promoción</span>
		</button>
	</div>
</div>

{#if errorMsg && !mostrarFormulario}
	<div class="mb-5 flex gap-3 rounded-lg border border-red-500/15 bg-danger-bg p-4 text-sm text-danger-color" role="alert">
		<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10" /><line x1="12" y1="8" x2="12" y2="12" /><line x1="12" y1="16" x2="12.01" y2="16" /></svg>
		<span>{errorMsg}</span>
	</div>
{/if}

{#snippet promoCard(promo, estado)}
	{@const prod = promo.tipo !== 'COMBO' ? obtenerProducto(promo.producto_id) : null}
	{@const nombreProd = promo.tipo === 'COMBO' ? generarNombresCombo(promo) : (prod?.nombre || 'Producto sin nombre')}
	
	<div class="min-w-[300px] w-[300px] sm:w-[320px] flex-none snap-start bg-bg-card border border-border-color rounded-xl p-5 shadow-sm hover:border-border-color-hover transition-all flex flex-col justify-between {estado === 'expirada' ? 'opacity-60 grayscale' : ''}">
		
		<div>
			<div class="flex justify-between items-start mb-4">
				<div class="w-8 h-8 rounded {promo.tipo === 'NXM' ? 'bg-text-primary/10 text-accent' : (promo.tipo === 'COMBO' ? 'bg-purple-500/15 text-purple-500' : 'bg-text-primary/10 text-accent')} flex items-center justify-center font-bold text-sm">
					{#if promo.tipo === 'NXM'}
						<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"></path><line x1="7" y1="7" x2="7.01" y2="7"></line></svg>
					{:else if promo.tipo === 'COMBO'}
						<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path><polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline><line x1="12" y1="22.08" x2="12" y2="12"></line></svg>
					{:else}
						%
					{/if}
				</div>
				
				{#if estado === 'activa'}
					<span class="px-3 py-1 bg-green-500/15 text-green-600 dark:text-green-400 border border-green-500/20 text-xs rounded-full font-bold">Activa</span>
				{:else if estado === 'proxima'}
					<span class="px-3 py-1 bg-blue-500/15 text-blue-600 dark:text-blue-400 border border-blue-500/20 text-xs rounded-full font-bold">Próxima</span>
				{:else}
					<span class="px-3 py-1 bg-text-primary/10 text-text-secondary border border-border-color text-xs rounded-full font-bold">Expirada</span>
				{/if}
			</div>

			<h3 class="font-bold text-text-primary text-[15px] mb-1 line-clamp-1">{generarTitulo(promo, nombreProd)}</h3>
			<p class="text-text-secondary text-xs mb-4 line-clamp-2" title={nombreProd}>{nombreProd}</p>

			<div class="flex justify-between items-end mb-5">
				<div class="text-3xl font-black text-accent">
					{calcularDestaque(promo, prod)}
				</div>
				<div class="bg-text-primary/5 text-text-secondary px-2 py-1 rounded text-[11px] font-semibold border border-border-color">
					{#if promo.tipo === 'NXM'}
						Lleva {promo.lleva} Paga {promo.paga}
					{:else if promo.tipo === 'porcentaje'}
						% Descuento
					{:else if promo.tipo === 'COMBO'}
						Precio Combo
					{:else}
						Precio fijo
					{/if}
				</div>
			</div>
		</div>

		<div>
			<div class="flex items-center gap-2 text-text-muted text-[11px] font-medium mb-4">
				<svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
				{formatearRangoFechas(promo.fecha_inicio, promo.fecha_fin)}
			</div>

			<div class="grid grid-cols-2 gap-3">
				<button 
					onclick={() => handleEditar(promo)}
					class="inline-flex cursor-pointer items-center justify-center rounded-lg border border-border-color bg-text-primary/3 p-2 text-xs font-semibold text-text-secondary transition-all duration-200 hover:border-border-color-hover hover:bg-text-primary/7 hover:text-text-primary"
				>
					Editar
				</button>
				<button 
					onclick={() => pedirConfirmacionEliminar(promo.id_promocion)}
					class="inline-flex cursor-pointer items-center justify-center rounded-lg border border-red-500/15 bg-danger-bg p-2 text-xs font-semibold text-danger-color transition-all duration-200 hover:bg-danger-color hover:text-white"
				>
					Cancelar
				</button>
			</div>
		</div>
	</div>
{/snippet}

<div class="mb-10">
	<div class="flex items-center gap-2 text-green-600 dark:text-green-500 mb-4 text-xs font-bold uppercase tracking-wider">
		<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><polyline points="22 4 12 14.01 9 11.01"></polyline></svg>
		Activas Ahora
	</div>
	
	{#if activas.length === 0}
		<p class="text-text-secondary text-sm py-4">No hay promociones activas.</p>
	{:else}
		<div class="relative group">
			<button onclick={() => moverCarrusel(carruselActivas, 'izq')} class="absolute -left-5 top-1/2 -translate-y-1/2 z-10 bg-bg-card border border-border-color shadow-md rounded-full p-2 text-text-secondary opacity-0 group-hover:opacity-100 transition-all hidden sm:block hover:bg-text-primary/5 focus:outline-none">
				<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="15 18 9 12 15 6"></polyline></svg>
			</button>
			
			<div bind:this={carruselActivas} class="flex overflow-x-auto gap-6 pb-4 snap-x snap-mandatory">
				{#each activas as promo}
					{@render promoCard(promo, 'activa')}
				{/each}
			</div>

			<button onclick={() => moverCarrusel(carruselActivas, 'der')} class="absolute -right-5 top-1/2 -translate-y-1/2 z-10 bg-bg-card border border-border-color shadow-md rounded-full p-2 text-text-secondary opacity-0 group-hover:opacity-100 transition-all hidden sm:block hover:bg-text-primary/5 focus:outline-none">
				<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"></polyline></svg>
			</button>
		</div>
	{/if}
</div>

{#if proximas.length > 0}
	<div class="mb-10">
		<div class="flex items-center gap-2 text-blue-600 dark:text-blue-500 mb-4 text-xs font-bold uppercase tracking-wider">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
			Próximas
		</div>
		
		<div class="relative group">
			<button onclick={() => moverCarrusel(carruselProximas, 'izq')} class="absolute -left-5 top-1/2 -translate-y-1/2 z-10 bg-bg-card border border-border-color shadow-md rounded-full p-2 text-text-secondary opacity-0 group-hover:opacity-100 transition-all hidden sm:block hover:bg-text-primary/5 focus:outline-none">
				<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="15 18 9 12 15 6"></polyline></svg>
			</button>
			
			<div bind:this={carruselProximas} class="flex overflow-x-auto gap-6 pb-4 snap-x snap-mandatory">
				{#each proximas as promo}
					{@render promoCard(promo, 'proxima')}
				{/each}
			</div>

			<button onclick={() => moverCarrusel(carruselProximas, 'der')} class="absolute -right-5 top-1/2 -translate-y-1/2 z-10 bg-bg-card border border-border-color shadow-md rounded-full p-2 text-text-secondary opacity-0 group-hover:opacity-100 transition-all hidden sm:block hover:bg-text-primary/5 focus:outline-none">
				<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"></polyline></svg>
			</button>
		</div>
	</div>
{/if}

{#if expiradas.length > 0}
	<div class="mb-10 opacity-75 hover:opacity-100 transition-opacity">
		<div class="flex items-center gap-2 text-text-muted mb-4 text-xs font-bold uppercase tracking-wider">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="15" y1="9" x2="9" y2="15"></line><line x1="9" y1="9" x2="15" y2="15"></line></svg>
			Expiradas
		</div>
		
		<div class="relative group">
			<button onclick={() => moverCarrusel(carruselExpiradas, 'izq')} class="absolute -left-5 top-1/2 -translate-y-1/2 z-10 bg-bg-card border border-border-color shadow-md rounded-full p-2 text-text-secondary opacity-0 group-hover:opacity-100 transition-all hidden sm:block hover:bg-text-primary/5 focus:outline-none">
				<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="15 18 9 12 15 6"></polyline></svg>
			</button>
			
			<div bind:this={carruselExpiradas} class="flex overflow-x-auto gap-6 pb-4 snap-x snap-mandatory">
				{#each expiradas as promo}
					{@render promoCard(promo, 'expirada')}
				{/each}
			</div>

			<button onclick={() => moverCarrusel(carruselExpiradas, 'der')} class="absolute -right-5 top-1/2 -translate-y-1/2 z-10 bg-bg-card border border-border-color shadow-md rounded-full p-2 text-text-secondary opacity-0 group-hover:opacity-100 transition-all hidden sm:block hover:bg-text-primary/5 focus:outline-none">
				<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"></polyline></svg>
			</button>
		</div>
	</div>
{/if}

{#if mostrarFormulario}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-text-primary/40 p-5 backdrop-blur-[4px]" onclick={() => (mostrarFormulario = false)} role="presentation">
		<div class="w-full max-w-[600px] overflow-hidden rounded-xl border border-border-color bg-bg-card shadow-lg animate-modal-enter" onclick={(e) => e.stopPropagation()} role="dialog">
			
			<header class="flex items-center justify-between border-b border-border-color p-5">
				<h2 class="text-lg font-bold text-text-primary">{editandoId ? 'Editar Promoción' : 'Nueva Promoción'}</h2>
				<button class="inline-flex cursor-pointer items-center justify-center rounded-lg border border-border-color bg-text-primary/3 p-2 text-text-secondary transition-all duration-200 hover:border-border-color-hover hover:bg-text-primary/7 hover:text-text-primary" onclick={() => (mostrarFormulario = false)}>&times;</button>
			</header>

			<form onsubmit={handleCrearPromocion}>
				<div class="p-6 overflow-y-auto max-h-[65vh] hide-scrollbar">
					
					<div class="flex flex-col gap-1.5 mb-5">
						<label for="tipoPromo" class="text-[0.85rem] font-semibold text-text-secondary">Tipo de Oferta:</label>
						<select id="tipoPromo" bind:value={tipoPromocion} class="rounded-lg border border-[rgba(15,30,54,0.15)] bg-white px-4 py-3 text-sm text-text-primary outline-none transition-all duration-200 focus:border-primario focus:ring-1 focus:ring-primario dark:bg-bg-primary">
							<option value="NXM">Lleva X, Paga Y (Ej. 3x2)</option>
							<option value="precio_fijo">Descuento Directo ($)</option>
							<option value="porcentaje">Descuento por Porcentaje (%)</option>
							<option value="COMBO">Combo (Hasta 3 productos)</option>
						</select>
					</div>

					<div class="mb-5 rounded-lg border border-border-color p-4 bg-text-primary/2">
						<label class="text-[0.85rem] font-semibold text-text-secondary mb-2 block">
							{tipoPromocion === 'COMBO' ? 'Buscar y añadir productos (Máx 3):' : 'Buscar producto asociado:'}
						</label>
						
						{#if tipoPromocion !== 'COMBO'}
							<select bind:value={categoriaSeleccionada} class="w-full mb-3 rounded-lg border border-[rgba(15,30,54,0.15)] bg-white px-4 py-2 text-sm text-text-primary outline-none transition-all duration-200 focus:border-primario focus:ring-1 focus:ring-primario dark:bg-bg-primary">
								<option value="">Todas las categorías</option>
								{#each categoriasUnicas.filter(Boolean) as cat}
									<option value={cat}>{cat}</option>
								{/each}
							</select>
						{/if}

						<div class="flex w-full overflow-hidden rounded-lg border border-[rgba(15,30,54,0.15)] bg-white transition-all duration-200 focus-within:border-primario focus-within:ring-1 focus-within:ring-primario dark:bg-bg-primary mb-3">
							<div class="flex items-center justify-center pl-3 pr-2 text-text-muted">
								<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
							</div>
							<input
								type="text"
								class="flex-1 border-none bg-transparent px-2 py-2 text-sm text-text-primary outline-none focus:ring-0"
								placeholder="Escribe el nombre o escanea el código..."
								bind:value={busquedaProducto}
							/>
						</div>

						{#if busquedaProducto.trim() && productosFiltrados.length > 0}
							<div class="max-h-32 overflow-y-auto rounded-lg border border-border-color bg-bg-card shadow-inner mb-3">
								{#each productosFiltrados as prod}
									<button type="button" onclick={() => seleccionarProducto(prod.id_producto)} class="w-full text-left px-3 py-2 text-sm border-b border-border-color hover:bg-text-primary/5 text-text-primary last:border-0 flex justify-between items-center">
										<span>{prod.nombre}</span>
										<span class="text-xs text-text-muted">{prod.codigo_barras || 'Sin código'}</span>
									</button>
								{/each}
							</div>
						{:else if busquedaProducto.trim() && productosFiltrados.length === 0}
							<p class="text-xs text-danger-color mb-3">No se encontraron productos.</p>
						{/if}

						<div class="flex flex-wrap gap-2">
							{#if tipoPromocion === 'COMBO'}
								{#each comboSeleccionados as id}
									<div class="flex items-center gap-2 rounded-full bg-accent/10 px-3 py-1 text-xs font-semibold text-accent border border-accent/20">
										{obtenerProducto(id)?.nombre || 'Producto'}
										<button type="button" onclick={() => removerProductoCombo(id)} class="text-accent hover:text-danger-color">&times;</button>
									</div>
								{/each}
								{#if comboSeleccionados.length === 0}
									<span class="text-xs text-text-muted italic">Ningún producto seleccionado.</span>
								{/if}
							{:else}
								{#if productoSeleccionado}
									<div class="flex items-center gap-2 rounded-full bg-accent/10 px-3 py-1 text-xs font-semibold text-accent border border-accent/20">
										{obtenerProducto(productoSeleccionado)?.nombre || 'Producto'}
										<button type="button" onclick={() => productoSeleccionado = ''} class="text-accent hover:text-danger-color">&times;</button>
									</div>
								{:else}
									<span class="text-xs text-text-muted italic">Ningún producto seleccionado.</span>
								{/if}
							{/if}
						</div>
					</div>

					<div class="grid grid-cols-1 sm:grid-cols-2 gap-5 mb-5">
						<div class="flex flex-col gap-1.5">
							<label for="fechaInicio" class="text-[0.85rem] font-semibold text-text-secondary">Fecha de Inicio:</label>
							<input type="date" id="fechaInicio" bind:value={fechaInicio} min={hoyStr} required class="rounded-lg border border-[rgba(15,30,54,0.15)] bg-white px-4 py-3 text-sm text-text-primary outline-none transition-all duration-200 focus:border-primario focus:ring-1 focus:ring-primario dark:bg-bg-primary" style="color-scheme: dark;" />
						</div>
						<div class="flex flex-col gap-1.5">
							<label for="fechaFin" class="text-[0.85rem] font-semibold text-text-secondary">Fecha de Fin:</label>
							<input type="date" id="fechaFin" bind:value={fechaFin} min={fechaInicio || hoyStr} required class="rounded-lg border border-[rgba(15,30,54,0.15)] bg-white px-4 py-3 text-sm text-text-primary outline-none transition-all duration-200 focus:border-primario focus:ring-1 focus:ring-primario dark:bg-bg-primary" style="color-scheme: dark;" />
						</div>
					</div>

					<div class="bg-text-primary/5 p-4 rounded-lg border border-border-color">
						{#if tipoPromocion === 'NXM'}
							<div class="flex gap-4">
								<div class="flex-1 flex flex-col gap-1.5">
									<label for="inputLleva" class="text-[0.85rem] font-semibold text-text-secondary">Lleva:</label>
									<input id="inputLleva" type="number" min="2" bind:value={lleva} required class="rounded-lg border border-[rgba(15,30,54,0.15)] bg-white px-4 py-3 text-sm text-text-primary outline-none transition-all duration-200 focus:border-primario focus:ring-1 focus:ring-primario dark:bg-bg-primary" style="color-scheme: dark;" />
								</div>
								<div class="flex-1 flex flex-col gap-1.5">
									<label for="inputPaga" class="text-[0.85rem] font-semibold text-text-secondary">Paga:</label>
									<input id="inputPaga" type="number" min="1" bind:value={paga} required class="rounded-lg border border-[rgba(15,30,54,0.15)] bg-white px-4 py-3 text-sm text-text-primary outline-none transition-all duration-200 focus:border-primario focus:ring-1 focus:ring-primario dark:bg-bg-primary" style="color-scheme: dark;" />
								</div>
							</div>
						{:else if tipoPromocion === 'precio_fijo' || tipoPromocion === 'COMBO'}
							<div class="flex flex-col gap-1.5">
								<label for="inputDesc" class="text-[0.85rem] font-semibold text-text-secondary">{tipoPromocion === 'COMBO' ? 'Precio Final del Combo ($):' : 'Monto a descontar ($):'}</label>
								<input id="inputDesc" type="number" min="1" bind:value={descuento} required class="rounded-lg border border-[rgba(15,30,54,0.15)] bg-white px-4 py-3 text-sm text-text-primary outline-none transition-all duration-200 focus:border-primario focus:ring-1 focus:ring-primario dark:bg-bg-primary" style="color-scheme: dark;" />
							</div>
						{:else if tipoPromocion === 'porcentaje'}
							<div class="flex flex-col gap-1.5">
								<label for="inputDescPorc" class="text-[0.85rem] font-semibold text-text-secondary">Porcentaje de descuento (%):</label>
								<input id="inputDescPorc" type="number" min="1" max="100" bind:value={descuento} required class="rounded-lg border border-[rgba(15,30,54,0.15)] bg-white px-4 py-3 text-sm text-text-primary outline-none transition-all duration-200 focus:border-primario focus:ring-1 focus:ring-primario dark:bg-bg-primary" style="color-scheme: dark;" />
							</div>
						{/if}
					</div>
				</div>

				<footer class="flex justify-end gap-3 border-t border-border-color bg-text-primary/2 p-4 px-6">
					<button type="button" class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-border-color bg-bg-secondary px-5 py-2.5 text-sm font-semibold text-text-primary transition-all duration-200 hover:bg-text-primary/5" onclick={() => (mostrarFormulario = false)} disabled={submitLoading}>Cancelar</button>
					<button type="submit" class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg bg-gradient-to-r from-accent-light to-accent px-5 py-2.5 text-sm font-semibold text-white transition-all duration-200 hover:-translate-y-[1px] hover:shadow-glow disabled:cursor-not-allowed disabled:opacity-50" disabled={submitLoading}>
						{#if submitLoading} Procesando... {:else} {editandoId ? 'Actualizar Promoción' : 'Activar Promoción'} {/if}
					</button>
				</footer>
			</form>
		</div>
	</div>
{/if}

{#if mostrarModalEliminar}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-text-primary/40 p-5 backdrop-blur-[4px]" onclick={() => (mostrarModalEliminar = false)} role="presentation">
		<div class="w-full max-w-[500px] overflow-hidden rounded-xl border border-border-color bg-bg-card shadow-lg animate-modal-enter" onclick={(e) => e.stopPropagation()} role="dialog">
			<header class="flex items-center justify-between border-b border-border-color p-5">
				<h2 class="text-lg font-bold text-danger-color">Confirmar Eliminación</h2>
				<button class="inline-flex cursor-pointer items-center justify-center rounded-lg border border-border-color bg-text-primary/3 p-2 text-text-secondary transition-all duration-200 hover:border-border-color-hover hover:bg-text-primary/7 hover:text-text-primary" onclick={() => (mostrarModalEliminar = false)}>&times;</button>
			</header>
			<div class="p-6 text-text-primary">
				<p>¿Estás seguro de que deseas cancelar esta promoción de forma permanente?</p>
				<p class="mt-3 text-xs text-text-muted">La promoción pasará a estar inactiva y los clientes ya no podrán usar este descuento.</p>
			</div>
			<footer class="flex justify-end gap-3 border-t border-border-color bg-text-primary/2 p-4 px-6">
				<button type="button" class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-border-color bg-bg-secondary px-5 py-2.5 text-sm font-semibold text-text-primary transition-all duration-200 hover:bg-text-primary/5" onclick={() => (mostrarModalEliminar = false)}>Mantener promoción</button>
				<button type="button" class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-red-500/15 bg-danger-bg px-5 py-2.5 text-sm font-semibold text-danger-color transition-all duration-200 hover:bg-danger-color hover:text-white" onclick={confirmarEliminacion}>Sí, cancelar ahora</button>
			</footer>
		</div>
	</div>
{/if}

<style>
	/* Scrollbar */
	::-webkit-scrollbar { width: 10px; height: 10px; }
	::-webkit-scrollbar-track { background: #f3f4f6; }
	::-webkit-scrollbar-thumb { background: #d1d5db; border-radius: 5px; border: 2px solid #f3f4f6; }
	::-webkit-scrollbar-thumb:hover { background: #9ca3af; }
	
	:global(.dark) ::-webkit-scrollbar-track { background: #18181b; }
	:global(.dark) ::-webkit-scrollbar-thumb { background: #3f3f46; border: 2px solid #18181b; }
	:global(.dark) ::-webkit-scrollbar-thumb:hover { background: #52525b; }
	
	.hide-scrollbar::-webkit-scrollbar { display: none; }
	.hide-scrollbar { -ms-overflow-style: none; scrollbar-width: none; }
</style>