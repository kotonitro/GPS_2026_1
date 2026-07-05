<script lang="ts">
	import { onMount } from 'svelte';
	import { obtenerPromociones, crearPromocion,  obtenerProductos, eliminarPromocion } from '$lib/api';

	let promociones = $state<any[]>([]);
	let productos = $state<any[]>([]);
	
	let mensajeError = $state('');
	let mensajeExito = $state('');

	let categoriaSeleccionada = $state('');
	let tipoPromocion = $state('NXM');
	let productoSeleccionado = $state('');
	let lleva = $state(0);
	let paga = $state(0);
	let descuento = $state(0);

	let categoriasUnicas = $derived(
		[...new Set(productos.map(p => p.categoria?.nombre_categoria))]
	);

	let productosFiltrados = $derived(
			categoriaSeleccionada === '' 
				? productos 
				: productos.filter(p => p.categoria?.nombre_categoria === categoriaSeleccionada)
		);

	onMount(async () => {
		await cargarDatos();
	});

	async function cargarDatos() {
		try {
			productos = await  obtenerProductos();
			console.log("Datos de promociones:", promociones);
			promociones = await obtenerPromociones();
		} catch (error: any) {
			mensajeError = 'Error al cargar los datos: ' + error.message;
		}
	}

	// Recibimos el evento para prevenir que recargue la página (reemplaza el |preventDefault viejo)
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
				descuento: (tipoPromocion === 'precio_fijo' || tipoPromocion === 'porcentaje') ? descuento : 0
			};

			await crearPromocion(payload);
			
			mensajeExito = 'Promoción creada y asociada al producto exitosamente.';
			
			productoSeleccionado = '';
			lleva = 0;
			paga = 0;
			descuento = 0;
			await cargarDatos();
			
		} catch (error: any) {
			mensajeError = error.message;
		}
	}

	async function handleEliminar(id: string) {
		if(confirm('¿Estás seguro de eliminar esta promoción?')) {
			try {
				await eliminarPromocion(id);
				mensajeExito = 'Promoción eliminada.';
				await cargarDatos();
			} catch (error: any) {
				mensajeError = error.message;
			}
		}
	}
</script>

<div class="max-w-6xl mx-auto p-4 grid grid-cols-1 md:grid-cols-3 gap-6">
	
	<div class="md:col-span-1 bg-white rounded-lg shadow-md p-6 h-fit">
		<h2 class="text-2xl font-bold text-primario mb-6">Nueva Promoción</h2>

		{#if mensajeError}
			<p class="text-error font-bold mb-4 text-sm">{mensajeError}</p>
		{/if}

		{#if mensajeExito}
			<p class="text-exito font-bold mb-4 text-sm">{mensajeExito}</p>
		{/if}

		<form onsubmit={handleCrearPromocion} class="flex flex-col gap-4">
			
			<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
				<div>
					<label for="filtroCat" class="block text-sm font-semibold mb-1 text-gray-700">Filtrar por Categoría:</label>
					<select
						id="filtroCat"
						bind:value={categoriaSeleccionada}
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:border-primario focus:ring-1 focus:ring-primario"
					>
						<option value="">Todas las categorías</option>
						{#each categoriasUnicas.filter(Boolean) as cat}
							<option value={cat}>{cat}</option>
						{/each}
					</select>
				</div>

				<div>
					<label for="selectProd" class="block text-sm font-semibold mb-1 text-gray-700">Producto a asociar:</label>
					<select
						id="selectProd"
						bind:value={productoSeleccionado}
						required
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:border-primario focus:ring-1 focus:ring-primario disabled:bg-gray-100 disabled:text-gray-400"
						disabled={productosFiltrados.length === 0}
					>
						<option value="" disabled>
							{productosFiltrados.length === 0 ? 'No hay productos...' : 'Seleccione un producto...'}
						</option>
						{#each productosFiltrados as prod}
							<option value={prod.id_producto}>{prod.nombre} (Stock: {prod.stock})</option>
						{/each}
					</select>
				</div>
			</div>

			<div>
							<label for="tipoPromo" class="block text-sm font-semibold mb-1 text-gray-700">Tipo de Oferta:</label>
							<select
								id="tipoPromo"
								bind:value={tipoPromocion}
								class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:border-primario focus:ring-1 focus:ring-primario"
							>
								<option value="NXM">Lleva X, Paga Y (Ej. 3x2)</option>
								<option value="precio_fijo">Descuento Directo ($)</option>
								<option value="porcentaje">Descuento por Porcentaje (%)</option>
							</select>
						</div>

						{#if tipoPromocion === 'NXM'}
							<div class="flex gap-4">
								<div class="flex-1">
									<label for="inputLleva" class="block text-sm font-semibold mb-1 text-gray-700">Lleva:</label>
									<input id="inputLleva" type="number" min="2" bind:value={lleva} required class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:border-primario focus:ring-1 focus:ring-primario" />
								</div>
								<div class="flex-1">
									<label for="inputPaga" class="block text-sm font-semibold mb-1 text-gray-700">Paga:</label>
									<input id="inputPaga" type="number" min="1" bind:value={paga} required class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:border-primario focus:ring-1 focus:ring-primario" />
								</div>
							</div>
						{:else if tipoPromocion === 'precio_fijo'}
							<div>
								<label for="inputDesc" class="block text-sm font-semibold mb-1 text-gray-700">Monto a descontar ($):</label>
								<input id="inputDesc" type="number" min="1" bind:value={descuento} required class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:border-primario focus:ring-1 focus:ring-primario" />
							</div>
						{:else if tipoPromocion === 'porcentaje'}
							<div>
								<label for="inputDescPorc" class="block text-sm font-semibold mb-1 text-gray-700">Porcentaje de descuento (%):</label>
								<input id="inputDescPorc" type="number" min="1" max="100" bind:value={descuento} required class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:border-primario focus:ring-1 focus:ring-primario" />
							</div>
						{/if}

			<button
				type="submit"
				class="w-full bg-primario hover:bg-primario-hover text-white font-bold py-2 px-4 rounded-md transition-colors mt-2"
			>
				Activar Promoción
			</button>
		</form>
	</div>

	<div class="md:col-span-2 bg-white rounded-lg shadow-md p-6 h-fit">
		<h2 class="text-xl font-bold text-primario mb-6">Promociones Activas</h2>
		
		{#if promociones.length === 0}
			<p class="text-gray-500 text-center py-8">No hay promociones activas en este momento.</p>
		{:else}
			<div class="overflow-x-auto">
				<table class="w-full text-left border-collapse">
					<thead>
						<tr class="border-b-2 border-gray-200">
							<th class="py-3 px-2 font-semibold text-gray-700">Producto</th>
							<th class="py-3 px-2 font-semibold text-gray-700">Oferta</th>
							<th class="py-3 px-2 font-semibold text-gray-700 text-right">Acción</th>
						</tr>
					</thead>
					<tbody>
					{#each promociones as promo}
						<tr class="border-b border-gray-100 hover:bg-gray-50">
						<td class="py-3 px-2 font-medium">
    {productos.find(p => p.id_producto === promo.producto_id)?.nombre || promo.producto_id || 'Sin producto asociado'}
</td>
							<td class="py-3 px-2">
								<span class="inline-block bg-primario/10 text-primario px-2 py-1 rounded text-xs font-bold">
									{#if promo.tipo === 'NXM'}
										{promo.lleva}x{promo.paga}
										{:else if promo.tipo === 'porcentaje'}
																				-{promo.descuento}%
																			{:else}
																				-${promo.descuento}
										{/if}
								</span>
							</td>
							<td class="py-3 px-2 text-right">
								<button 
									onclick={() => handleEliminar(promo.id_promocion)}
									class="text-error hover:text-red-800 text-sm font-bold"
								>
									Eliminar
								</button>
							</td>
						</tr>
{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</div>

</div>