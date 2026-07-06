<script lang="ts">
	import { onMount } from 'svelte';
	import {
		obtenerPromociones,
		crearPromocion,
		obtenerProductos,
		eliminarPromocion,
		actualizarPromocion
	} from '$lib/api';

	let promociones = $state<any[]>([]);
	let productos = $state<any[]>([]);
	let mensajeError = $state('');
	let mensajeExito = $state('');

	// Control del Modal
	let mostrarFormulario = $state(false);

	// Variables del formulario
	let categoriaSeleccionada = $state('');
	let tipoPromocion = $state('NXM');
	let productoSeleccionado = $state('');
	let lleva = $state(0);
	let paga = $state(0);
	let descuento = $state(0);
	let editandoId = $state<string | null>(null);
	
	// Variables de fecha
	let fechaInicio = $state('');
	let fechaFin = $state('');
	let hoy = $derived(new Date());
	let hoyStr = $derived(hoy.toISOString().split('T')[0]);

	let mostrarModalEliminar = $state(false);
	let idAEliminar = $state<string | null>(null);

	let categoriasUnicas = $derived([
		...new Set(productos.map((p) => p.categoria?.nombre_categoria))
	]);

	let productosFiltrados = $derived(
		categoriaSeleccionada === ''
			? productos
			: productos.filter((p) => p.categoria?.nombre_categoria === categoriaSeleccionada)
	);

	// Clasificación automática de promociones por fecha
	let activas = $derived(promociones.filter(p => {
		if (!p.fecha_inicio || !p.fecha_fin) return true; // Si aún no implementas fechas en Go, todas son activas
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

	function pedirConfirmacionEliminar(id: string) {
			idAEliminar = id;
			mostrarModalEliminar = true;
		}

		function cerrarModalEliminar() {
			idAEliminar = null;
			mostrarModalEliminar = false;
		}

	onMount(async () => {
		await cargarDatos();
	});

	async function cargarDatos() {
		try {
			productos = await obtenerProductos();
			promociones = await obtenerPromociones();
		} catch (error: any) {
			mensajeError = 'Error al cargar los datos: ' + error.message;
		}
	}

	async function handleCrearPromocion(e: Event) {
		e.preventDefault();
		mensajeError = '';
		mensajeExito = '';

		try {
			const payload = {
				producto_id: productoSeleccionado,
				tipo: tipoPromocion,
				lleva: tipoPromocion === 'NXM' ? lleva : 0,
				paga: tipoPromocion === 'NXM' ? paga : 0,
				descuento: (tipoPromocion === 'precio_fijo' || tipoPromocion === 'porcentaje') ? descuento : 0,
				fecha_inicio: fechaInicio ? new Date(fechaInicio).toISOString() : null,
				fecha_fin: fechaFin ? new Date(fechaFin).toISOString() : null
			};

			if (editandoId) {
				await actualizarPromocion(editandoId, payload);
				mensajeExito = 'Promoción actualizada exitosamente.';
			} else {
				await crearPromocion(payload);
				mensajeExito = 'Promoción creada exitosamente.';
			}
			
			mostrarFormulario = false; 
			await cargarDatos();
			setTimeout(() => { mensajeExito = ''; }, 3000);
		} catch (error: any) {
			mensajeError = error.message;
		}
	}
	async function confirmarEliminacion() {
			if (!idAEliminar) return;
			
			try {
				await eliminarPromocion(idAEliminar);
				mensajeExito = 'Promoción cancelada exitosamente.';
				await cargarDatos();
				cerrarModalEliminar();
				setTimeout(() => { mensajeExito = ''; }, 3000);
			} catch (error: any) {
				mensajeError = error.message;
				cerrarModalEliminar();
			}
		}

	function handleEditar(promo: any) {
		editandoId = promo.id_promocion;
		
		// Buscar la categoría para que el select se llene correctamente
		const prod = obtenerProducto(promo.producto_id);
		categoriaSeleccionada = prod?.categoria?.nombre_categoria || '';
		
		// Llenar el resto de datos
		productoSeleccionado = promo.producto_id;
		tipoPromocion = promo.tipo;
		lleva = promo.lleva || 0;
		paga = promo.paga || 0;
		descuento = promo.descuento || 0;

		// Extraer solo la parte YYYY-MM-DD de la fecha de la base de datos
		fechaInicio = promo.fecha_inicio ? promo.fecha_inicio.split('T')[0] : '';
		fechaFin = promo.fecha_fin ? promo.fecha_fin.split('T')[0] : '';

		mostrarFormulario = true;
	}

	// Limpia el formulario al crear una nueva
	function abrirModalNuevo() {
		editandoId = null;
		categoriaSeleccionada = '';
		productoSeleccionado = '';
		fechaInicio = '';
		fechaFin = '';
		lleva = 0;
		paga = 0;
		descuento = 0;
		mostrarFormulario = true;
	}

	async function handleEliminar(id: string) {
		if (confirm('¿Estás seguro de cancelar esta promoción?')) {
			try {
				await eliminarPromocion(id);
				await cargarDatos();
			} catch (error: any) {
				mensajeError = error.message;
			}
		}
	}

	// Helpers de la interfaz
	function obtenerProducto(id: string) {
		return productos.find((p) => p.id_producto === id);
	}

	function formatearRangoFechas(inicio: string, fin: string) {
		if (!inicio || !fin) return 'Sin límite de tiempo establecido';
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
		} else {
			// Descuento directo: Lo transformamos visualmente a porcentaje usando el precio real del producto
			if (producto && producto.precio > 0) {
				return Math.round((promo.descuento / producto.precio) * 100) + '%';
			}
			return `$${promo.descuento}`; // Fallback si el producto no tiene precio
		}
	}

	function generarTitulo(promo: any, nombreProd: string) {
		if (promo.tipo === 'NXM') return `${promo.lleva}x${promo.paga} en ${nombreProd}`;
		if (promo.tipo === 'porcentaje') return `${promo.descuento}% en ${nombreProd}`;
		return `$${promo.descuento} de descuento en ${nombreProd}`;
	}
</script>

<div class="max-w-7xl mx-auto p-6 md:p-8 w-full">
	
	<div class="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-8 gap-4">
		<div>
			<h1 class="text-2xl font-bold text-gray-900 mb-1">Promociones</h1>
			<p class="text-sm text-gray-500">{activas.length} activas · {proximas.length} próximas</p>
		</div>
		<button 
			onclick={abrirModalNuevo}
			class="bg-[#b45309] hover:bg-[#92400e] text-white px-4 py-2 rounded-lg font-semibold flex items-center gap-2 transition-colors shadow-sm"
		>
			<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
			Nueva promoción
		</button>
	</div>

	{#if mensajeExito}
		<div class="mb-6 p-4 bg-green-50 border border-green-200 text-green-700 rounded-lg font-medium">{mensajeExito}</div>
	{/if}

	{#if mensajeError && !mostrarFormulario}
		<div class="mb-6 p-4 bg-red-50 border border-red-200 text-red-700 rounded-lg font-medium">{mensajeError}</div>
	{/if}

	{#snippet promoCard(promo, estado)}
		{@const prod = obtenerProducto(promo.producto_id)}
		{@const nombreProd = prod?.nombre || 'Producto sin nombre'}
		
		<div class="bg-white border border-gray-200 rounded-xl p-5 shadow-sm hover:shadow-md transition-shadow flex flex-col justify-between {estado === 'expirada' ? 'opacity-60 grayscale' : ''}">
			
			<div>
				<div class="flex justify-between items-start mb-4">
					<div class="w-8 h-8 rounded {promo.tipo === 'NXM' ? 'bg-orange-100 text-orange-700' : 'bg-amber-100 text-amber-700'} flex items-center justify-center font-bold text-sm">
						{#if promo.tipo === 'NXM'}
							<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"></path><line x1="7" y1="7" x2="7.01" y2="7"></line></svg>
						{:else}
							%
						{/if}
					</div>
					
					{#if estado === 'activa'}
						<span class="px-3 py-1 bg-emerald-50 text-emerald-600 border border-emerald-100 text-xs rounded-full font-bold">Activa</span>
					{:else if estado === 'proxima'}
						<span class="px-3 py-1 bg-blue-50 text-blue-600 border border-blue-100 text-xs rounded-full font-bold">Próxima</span>
					{:else}
						<span class="px-3 py-1 bg-gray-100 text-gray-500 border border-gray-200 text-xs rounded-full font-bold">Expirada</span>
					{/if}
				</div>

				<h3 class="font-bold text-gray-900 text-[15px] mb-1 line-clamp-1">{generarTitulo(promo, nombreProd)}</h3>
				<p class="text-gray-500 text-xs mb-4 line-clamp-1">{nombreProd}</p>

				<div class="flex justify-between items-end mb-5">
					<div class="text-3xl font-extrabold text-[#b45309]">
						{calcularDestaque(promo, prod)}
					</div>
					<div class="bg-gray-100 text-gray-600 px-2 py-1 rounded text-[11px] font-semibold">
						{#if promo.tipo === 'NXM'}
							Lleva {promo.lleva} Paga {promo.paga}
						{:else if promo.tipo === 'porcentaje'}
							% Descuento
						{:else}
							Precio fijo
						{/if}
					</div>
				</div>
			</div>

			<div>
				<div class="flex items-center gap-2 text-gray-400 text-[11px] font-medium mb-4">
					<svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
					{formatearRangoFechas(promo.fecha_inicio, promo.fecha_fin)}
				</div>

				<div class="grid grid-cols-2 gap-3">
					<button 
						onclick={() => handleEditar(promo)}
						class="py-2 px-4 bg-gray-50 hover:bg-gray-100 text-gray-700 text-xs font-semibold rounded-lg border border-gray-200 transition-colors"
					>
						Editar
					</button>
					<button 
							onclick={() => pedirConfirmacionEliminar(promo.id_promocion)}
							class="py-2 px-4 bg-red-50 hover:bg-red-100 text-red-600 text-xs font-semibold rounded-lg border border-red-100 transition-colors"
					>
						Cancelar
					</button>
				</div>
			</div>
		</div>
	{/snippet}

	<div class="mb-10">
		<div class="flex items-center gap-2 text-emerald-600 mb-4 text-xs font-bold uppercase tracking-wider">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><polyline points="22 4 12 14.01 9 11.01"></polyline></svg>
			Activas Ahora
		</div>
		{#if activas.length === 0}
			<p class="text-gray-500 text-sm py-4">No hay promociones activas.</p>
		{:else}
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
				{#each activas as promo}
					{@render promoCard(promo, 'activa')}
				{/each}
			</div>
		{/if}
	</div>

	{#if proximas.length > 0}
		<div class="mb-10">
			<div class="flex items-center gap-2 text-blue-600 mb-4 text-xs font-bold uppercase tracking-wider">
				<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
				Próximas
			</div>
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
				{#each proximas as promo}
					{@render promoCard(promo, 'proxima')}
				{/each}
			</div>
		</div>
	{/if}

	{#if expiradas.length > 0}
		<div class="mb-10 opacity-75">
			<div class="flex items-center gap-2 text-gray-500 mb-4 text-xs font-bold uppercase tracking-wider">
				<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="15" y1="9" x2="9" y2="15"></line><line x1="9" y1="9" x2="15" y2="15"></line></svg>
				Expiradas
			</div>
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
				{#each expiradas as promo}
					{@render promoCard(promo, 'expirada')}
				{/each}
			</div>
		</div>
	{/if}

</div>

{#if mostrarFormulario}
	<div class="fixed inset-0 bg-black/40 backdrop-blur-sm flex items-center justify-center z-50 p-4">
		<div class="bg-white rounded-xl shadow-2xl p-6 max-w-2xl w-full max-h-[90vh] overflow-y-auto relative">
			
			<button 
				onclick={() => mostrarFormulario = false}
				class="absolute top-4 right-4 text-gray-400 hover:text-gray-600"
				aria-label="Cerrar formulario"
			>
				<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
			</button>

			<h2 class="text-2xl font-bold text-[#b45309] mb-6">{editandoId ? 'Editar Promoción' : 'Nueva Promoción'}</h2>

			{#if mensajeError}
				<p class="text-red-600 bg-red-50 p-3 rounded-md font-medium mb-4 text-sm border border-red-200">{mensajeError}</p>
			{/if}

			<form onsubmit={handleCrearPromocion} class="flex flex-col gap-5">
				<div class="grid grid-cols-1 sm:grid-cols-2 gap-5">
					<div>
						<label for="filtroCat" class="block text-sm font-semibold mb-1.5 text-gray-700">Categoría:</label>
						<select
							id="filtroCat"
							bind:value={categoriaSeleccionada}
							class="w-full px-3 py-2.5 bg-gray-50 border border-gray-200 rounded-lg focus:outline-none focus:border-[#b45309] focus:ring-1 focus:ring-[#b45309]"
						>
							<option value="">Todas</option>
							{#each categoriasUnicas.filter(Boolean) as cat}
								<option value={cat}>{cat}</option>
							{/each}
						</select>
					</div>

					<div>
						<label for="selectProd" class="block text-sm font-semibold mb-1.5 text-gray-700">Producto asociado:</label>
						<select
							id="selectProd"
							bind:value={productoSeleccionado}
							required
							class="w-full px-3 py-2.5 bg-gray-50 border border-gray-200 rounded-lg focus:outline-none focus:border-[#b45309] focus:ring-1 focus:ring-[#b45309] disabled:opacity-50"
							disabled={productosFiltrados.length === 0}
						>
							<option value="" disabled>Seleccione un producto...</option>
							{#each productosFiltrados as prod}
								<option value={prod.id_producto}>{prod.nombre} (${prod.precio})</option>
							{/each}
						</select>
					</div>
				</div>

				<hr class="border-gray-100 my-1">

				<div class="grid grid-cols-1 sm:grid-cols-2 gap-5">
					<div>
						<label for="fechaInicio" class="block text-sm font-semibold mb-1.5 text-gray-700">Fecha de Inicio:</label>
						<input type="date" id="fechaInicio" bind:value={fechaInicio} min={hoyStr} required class="w-full px-3 py-2.5 bg-gray-50 border border-gray-200 rounded-lg focus:outline-none focus:border-[#b45309] focus:ring-1 focus:ring-[#b45309]" />
					</div>
					<div>
						<label for="fechaFin" class="block text-sm font-semibold mb-1.5 text-gray-700">Fecha de Fin:</label>
						<input type="date" id="fechaFin" bind:value={fechaFin} min={fechaInicio || hoyStr} required class="w-full px-3 py-2.5 bg-gray-50 border border-gray-200 rounded-lg focus:outline-none focus:border-[#b45309] focus:ring-1 focus:ring-[#b45309]" />
					</div>
				</div>

				<div>
					<label for="tipoPromo" class="block text-sm font-semibold mb-1.5 text-gray-700">Tipo de Oferta:</label>
					<select
						id="tipoPromo"
						bind:value={tipoPromocion}
						class="w-full px-3 py-2.5 bg-gray-50 border border-gray-200 rounded-lg focus:outline-none focus:border-[#b45309] focus:ring-1 focus:ring-[#b45309]"
					>
						<option value="NXM">Lleva X, Paga Y (Ej. 3x2)</option>
						<option value="precio_fijo">Descuento Directo ($)</option>
						<option value="porcentaje">Descuento por Porcentaje (%)</option>
					</select>
				</div>

				<div class="bg-orange-50/50 p-4 rounded-lg border border-orange-100">
					{#if tipoPromocion === 'NXM'}
						<div class="flex gap-4">
							<div class="flex-1">
								<label for="inputLleva" class="block text-sm font-semibold mb-1.5 text-gray-700">Lleva:</label>
								<input id="inputLleva" type="number" min="2" bind:value={lleva} required class="w-full px-3 py-2.5 bg-white border border-gray-200 rounded-lg focus:outline-none focus:border-[#b45309] focus:ring-1 focus:ring-[#b45309]" />
							</div>
							<div class="flex-1">
								<label for="inputPaga" class="block text-sm font-semibold mb-1.5 text-gray-700">Paga:</label>
								<input id="inputPaga" type="number" min="1" bind:value={paga} required class="w-full px-3 py-2.5 bg-white border border-gray-200 rounded-lg focus:outline-none focus:border-[#b45309] focus:ring-1 focus:ring-[#b45309]" />
							</div>
						</div>
					{:else if tipoPromocion === 'precio_fijo'}
						<div>
							<label for="inputDesc" class="block text-sm font-semibold mb-1.5 text-gray-700">Monto a descontar ($):</label>
							<input id="inputDesc" type="number" min="1" bind:value={descuento} required class="w-full px-3 py-2.5 bg-white border border-gray-200 rounded-lg focus:outline-none focus:border-[#b45309] focus:ring-1 focus:ring-[#b45309]" />
						</div>
					{:else if tipoPromocion === 'porcentaje'}
						<div>
							<label for="inputDescPorc" class="block text-sm font-semibold mb-1.5 text-gray-700">Porcentaje de descuento (%):</label>
							<input id="inputDescPorc" type="number" min="1" max="100" bind:value={descuento} required class="w-full px-3 py-2.5 bg-white border border-gray-200 rounded-lg focus:outline-none focus:border-[#b45309] focus:ring-1 focus:ring-[#b45309]" />
						</div>
					{/if}
				</div>

				<div class="flex justify-end gap-3 mt-4">
					<button 
						type="button" 
						onclick={() => mostrarFormulario = false}
						class="px-5 py-2.5 bg-white border border-gray-300 text-gray-700 font-semibold rounded-lg hover:bg-gray-50 transition-colors"
					>
						Cancelar
					</button>
					<button
						type="submit"
						class="px-6 py-2.5 bg-[#b45309] hover:bg-[#92400e] text-white font-bold rounded-lg transition-colors shadow-sm"
					>
						{editandoId ? 'Actualizar Promoción' : 'Activar Promoción'}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
{#if mostrarModalEliminar}
	<div class="fixed inset-0 bg-black/40 backdrop-blur-sm flex items-center justify-center z-50 p-4 transition-opacity">
		<div class="bg-white rounded-xl shadow-2xl max-w-md w-full overflow-hidden transform transition-all">
			
			<div class="p-6 sm:p-8">
				<div class="mx-auto flex items-center justify-center h-14 w-14 rounded-full bg-red-100 mb-6">
					<svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-red-600"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path><line x1="12" y1="9" x2="12" y2="13"></line><line x1="12" y1="17" x2="12.01" y2="17"></line></svg>
				</div>
				
				<div class="text-center">
					<h3 class="text-xl font-bold text-gray-900 mb-2">¿Cancelar promoción?</h3>
					<p class="text-sm text-gray-500">
						Esta acción detendrá la oferta inmediatamente. La promoción se desactivará y los clientes ya no podrán usar este descuento.
					</p>
				</div>
			</div>
			
			<div class="bg-gray-50 px-6 py-4 flex flex-col-reverse sm:flex-row justify-end gap-3 border-t border-gray-100">
				<button 
					type="button" 
					onclick={cerrarModalEliminar}
					class="w-full sm:w-auto px-5 py-2.5 bg-white border border-gray-300 text-gray-700 font-semibold rounded-lg hover:bg-gray-50 transition-colors shadow-sm"
				>
					Mantener promoción
				</button>
				<button 
					type="button" 
					onclick={confirmarEliminacion}
					class="w-full sm:w-auto px-5 py-2.5 bg-red-600 hover:bg-red-700 text-white font-bold rounded-lg transition-colors shadow-sm"
				>
					Sí, cancelar ahora
				</button>
			</div>

		</div>
	</div>
{/if}