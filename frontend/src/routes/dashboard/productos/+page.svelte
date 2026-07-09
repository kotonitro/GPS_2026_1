<script lang="ts">
	import { apiProductos, apiCategorias, type Producto, type Categoria } from '$lib/api';
	import { auth } from '$lib/authStore.svelte';
	import { toast } from '$lib/toastStore.svelte';
	import { onMount } from 'svelte';
	import {
		Search,
		Plus,
		Edit2,
		Trash2,
		X,
		AlertTriangle,
		ArrowUpDown,
		ArrowUp,
		ArrowDown,
		Settings,
		Camera
	} from '@lucide/svelte';
	import Scanner from '$lib/components/Scanner.svelte';

	let productos = $state<Producto[]>([]);
	let categorias = $state<Categoria[]>([]);
	let loading = $state(true);
	let errorMsg = $state('');

	// Gestión de Categorías Modal
	let showCategoriesModal = $state(false);
	let manageCatLoading = $state(false);
	let editingCatId = $state<string | null>(null);
	let editingCatName = $state('');
	let catGeneralError = $state('');
	let catSearchQuery = $state('');

	let filteredCategorias = $derived(
		categorias.filter((c) =>
			c.nombre_categoria.toLowerCase().includes(catSearchQuery.toLowerCase())
		)
	);

	// Filtros y Búsqueda
	let searchQuery = $state('');
	let modoEscaneo = $state(false);
	let selectedCategoria = $state('Todas');
	let selectedEstado = $state('Todas');

	// Ordenamiento
	let sortStock = $state<'none' | 'asc' | 'desc'>('none');

	// Lógica de filtrado y ordenamiento
	let filteredProductos = $derived.by(() => {
		let result = [...productos];

		// Búsqueda por texto (nombre o código barras)
		if (searchQuery.trim()) {
			const query = searchQuery.toLowerCase().trim();
			result = result.filter(
				(p) =>
					p.nombre.toLowerCase().includes(query) ||
					(p.codigo_barras && p.codigo_barras.toLowerCase().includes(query))
			);
		}

		// Filtro Categoría
		if (selectedCategoria !== 'Todas') {
			result = result.filter((p) => p.id_categoria === selectedCategoria);
		}

		// Filtro Estado
		if (selectedEstado !== 'Todas') {
			result = result.filter((p) => {
				if (selectedEstado === 'Agotado') return p.stock <= 0;
				if (selectedEstado === 'Bajo stock') return p.stock > 0 && p.stock <= p.stock_minimo;
				if (selectedEstado === 'Disponible') return p.stock > p.stock_minimo;
				return true;
			});
		}

		// Ordenamiento
		if (sortStock !== 'none') {
			result.sort((a, b) => {
				if (sortStock === 'asc') return a.stock - b.stock;
				return b.stock - a.stock;
			});
		}

		return result;
	});

	// Estadísticas y Alertas

	// Modales
	let showModal = $state(false);
	let showDeleteModal = $state(false);
	let editingProducto = $state<Producto | null>(null);
	let productoToDelete = $state<Producto | null>(null);

	// Campos de Formulario
	let formNombre = $state('');
	let formDescripcion = $state('');
	let formPrecio = $state<number | ''>('');
	let formStock = $state<number | ''>('');
	let formStockMinimo = $state<number | ''>('');
	let formMarca = $state('');
	let formUnidad = $state('unidades');
	let formCodigoBarras = $state('');
	let formCategoriaNombre = $state('');
	let inlineSuggestion = $derived.by(() => {
		if (!formCategoriaNombre) return '';
		const lowerInput = formCategoriaNombre.toLowerCase();
		const match = categorias.find((c) => c.nombre_categoria.toLowerCase().startsWith(lowerInput));
		if (match) {
			return formCategoriaNombre + match.nombre_categoria.substring(formCategoriaNombre.length);
		}
		return '';
	});
	// Errores de Formulario
	let errNombre = $state('');
	let errPrecio = $state('');
	let errCategoria = $state('');
	let errCodigoBarras = $state('');
	let formGeneralError = $state('');
	let submitLoading = $state(false);

	onMount(loadData);

	function manejarEscaneo(codigo: string) {
		searchQuery = codigo;
		modoEscaneo = false;
		toast.show(`Código escaneado: ${codigo}`, 'success');
	}

	async function loadData() {
		loading = true;
		errorMsg = '';
		try {
			const [resProductos, resCategorias] = await Promise.all([
				apiProductos.getAll(),
				apiCategorias.getAll()
			]);
			productos = Array.isArray(resProductos) ? resProductos : [];
			categorias = Array.isArray(resCategorias) ? resCategorias : [];
		} catch (err: any) {
			errorMsg = err.message || 'Error al conectar con la base de datos.';
			toast.show(errorMsg, 'error');
		} finally {
			loading = false;
		}
	}

	function formatCurrency(amount: number) {
		return new Intl.NumberFormat('es-CL', {
			style: 'currency',
			currency: 'CLP',
			maximumFractionDigits: 0
		}).format(amount);
	}

	function getEstadoBadge(producto: Producto) {
		if (producto.stock <= 0) {
			return {
				text: 'Agotado',
				classes: 'border-red-500/20 bg-red-500/10 text-red-600 dark:text-red-400'
			};
		} else if (producto.stock <= producto.stock_minimo) {
			return {
				text: 'Bajo stock',
				classes: 'border-orange-500/20 bg-orange-500/10 text-orange-600 dark:text-orange-400'
			};
		} else {
			return {
				text: 'Disponible',
				classes: 'border-green-500/20 bg-green-500/10 text-green-600 dark:text-green-400'
			};
		}
	}

	function toggleSortStock() {
		if (sortStock === 'none') sortStock = 'asc';
		else if (sortStock === 'asc') sortStock = 'desc';
		else sortStock = 'none';
	}

	function openCreateModal() {
		if (auth.user?.rol?.toLowerCase() !== 'admin') {
			toast.show('Acción denegada, requiere rol Administrador.', 'error');
			return;
		}
		editingProducto = null;
		formNombre = '';
		formDescripcion = '';
		formPrecio = '';
		formStock = '';
		formStockMinimo = '';
		formMarca = '';
		formUnidad = 'unidades';
		formCodigoBarras = '';
		formCategoriaNombre = '';
		clearErrors();
		showModal = true;
	}

	function handleNombreBlur() {
		if (formNombre) {
			formNombre = formNombre.replace(/\b\w/g, (c) => c.toUpperCase());
		}
	}

	function handleCategoriaBlur() {
		if (formCategoriaNombre) {
			formCategoriaNombre = formCategoriaNombre.replace(/\b\w/g, (c) => c.toUpperCase());
		}
	}

	function handleCategoriaKeyDown(e: KeyboardEvent) {
		if (e.key === 'Tab' && inlineSuggestion && inlineSuggestion !== formCategoriaNombre) {
			e.preventDefault();
			formCategoriaNombre = inlineSuggestion;
		}
	}

	function openEditModal(producto: Producto) {
		if (auth.user?.rol?.toLowerCase() !== 'admin') {
			toast.show('Acción denegada, requiere rol Administrador.', 'error');
			return;
		}
		editingProducto = producto;
		formNombre = producto.nombre;
		formDescripcion = producto.descripcion;
		formPrecio = producto.precio;
		formStock = producto.stock;
		formStockMinimo = producto.stock_minimo;
		formMarca = producto.marca || '';
		formUnidad = producto.unidad || 'unidades';
		formCodigoBarras = producto.codigo_barras || '';
		formCategoriaNombre = producto.categoria?.nombre_categoria || '';
		clearErrors();
		showModal = true;
	}

	function closeModal() {
		showModal = false;
		editingProducto = null;
	}

	function confirmDelete(producto: Producto) {
		if (auth.user?.rol?.toLowerCase() !== 'admin') {
			toast.show('Acción denegada, requiere rol Administrador.', 'error');
			return;
		}
		productoToDelete = producto;
		showDeleteModal = true;
	}

	function closeDeleteModal() {
		showDeleteModal = false;
		productoToDelete = null;
	}

	function openCategoriesModal() {
		if (auth.user?.rol?.toLowerCase() !== 'admin') {
			toast.show('Acción denegada, requiere rol Administrador.', 'error');
			return;
		}
		showCategoriesModal = true;
		catGeneralError = '';
		editingCatId = null;
		editingCatName = '';
		catSearchQuery = '';
	}

	function closeCategoriesModal() {
		showCategoriesModal = false;
	}

	function startEditCat(cat: Categoria) {
		editingCatId = cat.id_categoria;
		editingCatName = cat.nombre_categoria;
		catGeneralError = '';
	}

	function cancelEditCat() {
		editingCatId = null;
		editingCatName = '';
	}

	async function saveEditCat(id: string) {
		if (!editingCatName.trim()) {
			catGeneralError = 'El nombre de la categoría no puede estar vacío.';
			return;
		}
		manageCatLoading = true;
		catGeneralError = '';
		try {
			const finalName = editingCatName.trim().replace(/\b\w/g, (c) => c.toUpperCase());
			const res = await apiCategorias.update(id, { nombre_categoria: finalName });

			// Actualizamos el array local manualmente ya que el backend solo devuelve un mensaje
			categorias = categorias.map((c) =>
				c.id_categoria === id ? { ...c, nombre_categoria: finalName } : c
			);
			toast.show('Categoría actualizada con éxito.', 'success');
			editingCatId = null;
		} catch (err: any) {
			catGeneralError = err.message || 'Error al actualizar la categoría.';
		} finally {
			manageCatLoading = false;
		}
	}

	async function deleteCat(id: string) {
		manageCatLoading = true;
		catGeneralError = '';
		try {
			await apiCategorias.delete(id);
			categorias = categorias.filter((c) => c.id_categoria !== id);
			if (selectedCategoria === id) {
				selectedCategoria = 'Todas';
			}
			toast.show('Categoría eliminada con éxito.', 'success');
		} catch (err: any) {
			catGeneralError =
				err.message || 'Error al eliminar la categoría. Probablemente tenga productos asociados.';
		} finally {
			manageCatLoading = false;
		}
	}

	function clearErrors() {
		errNombre = '';
		errPrecio = '';
		errCategoria = '';
		errCodigoBarras = '';
		formGeneralError = '';
	}

	async function handleSubmit(e: Event) {
		e.preventDefault();
		clearErrors();
		let isValid = true;

		if (!formNombre.trim()) {
			errNombre = 'El nombre es obligatorio.';
			isValid = false;
		}
		if (Number(formPrecio) < 0) {
			errPrecio = 'El precio no puede ser negativo.';
			isValid = false;
		}
		if (!formCategoriaNombre.trim()) {
			errCategoria = 'Debe especificar una categoría.';
			isValid = false;
		}

		const barcodeRegex = /^\d{13}$/;
		if (!barcodeRegex.test(formCodigoBarras.trim())) {
			errCodigoBarras = 'El código de barras debe tener exactamente 13 dígitos numéricos.';
			isValid = false;
		}

		if (!isValid) return;

		submitLoading = true;
		let selectedCatId = '';

		if (formCategoriaNombre.trim()) {
			const existingCat = categorias.find(
				(c) => c.nombre_categoria.toLowerCase() === formCategoriaNombre.trim().toLowerCase()
			);
			if (existingCat) {
				selectedCatId = existingCat.id_categoria;
			} else {
				try {
					const resCat = await apiCategorias.create({
						nombre_categoria: formCategoriaNombre.trim()
					});
					if (resCat && resCat.categoria) {
						selectedCatId = resCat.categoria.id_categoria;
						categorias = [...categorias, resCat.categoria];
					}
				} catch (err) {
					errCategoria = 'Error al crear la nueva categoría.';
					submitLoading = false;
					return;
				}
			}
		}

		if (!selectedCatId) {
			errCategoria = 'Debe especificar una categoría.';
			submitLoading = false;
			return;
		}

		const payload: Partial<Producto> = {
			nombre: formNombre.trim(),
			descripcion: formDescripcion.trim(),
			precio: Number(formPrecio),
			stock: Number(formStock),
			stock_minimo: Number(formStockMinimo),
			marca: formMarca.trim(),
			unidad: formUnidad,
			codigo_barras: formCodigoBarras.trim(),
			id_categoria: selectedCatId
		};

		try {
			if (editingProducto) {
				await apiProductos.update(editingProducto.id_producto, payload);
				toast.show('Producto modificado exitosamente.', 'success');
			} else {
				await apiProductos.create(payload);
				toast.show('Producto creado exitosamente.', 'success');
			}
			closeModal();
			await loadData();
		} catch (error: any) {
			formGeneralError = error.message || 'Error al procesar la solicitud.';
		} finally {
			submitLoading = false;
		}
	}

	async function handleDelete() {
		if (!productoToDelete) return;
		submitLoading = true;
		try {
			await apiProductos.delete(productoToDelete.id_producto);
			toast.show('Producto eliminado exitosamente.', 'success');
			closeDeleteModal();
			await loadData();
		} catch (error: any) {
			toast.show(error.message || 'Error al eliminar el producto.', 'error');
		} finally {
			submitLoading = false;
		}
	}
</script>

<svelte:head>
	<title>Productos - MinimarketGo</title>
</svelte:head>

<div class="h-full">
	<!-- busqueda y acciones -->
	<div class="mb-6 flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
		<!-- Bloque de Búsqueda y Filtros -->
		<div class="flex flex-col gap-3 sm:flex-row sm:items-center flex-1">
			<!-- Buscador -->
			<div class="relative w-full sm:max-w-xs">
				<Search
					class="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-text-muted"
					size={20}
				/>
				<input
					type="text"
					placeholder="Buscar producto..."
					bind:value={searchQuery}
					class="w-full rounded-xl border border-border-color bg-bg-card py-2.5 pl-10 pr-12 text-sm text-text-primary focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario"
				/>
				<button
					type="button"
					onclick={() => (modoEscaneo = true)}
					class="absolute right-3 top-1/2 -translate-y-1/2 cursor-pointer text-text-muted hover:text-primario flex items-center justify-center"
					title="Escanear con cámara"
				>
					<Camera size={18} />
				</button>
			</div>

			<!-- filtro categorias -->
			<div class="flex items-center gap-2">
				<select
					bind:value={selectedCategoria}
					class="w-full sm:w-auto rounded-xl border border-border-color bg-bg-card py-2.5 pl-4 pr-10 text-sm text-text-primary focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario cursor-pointer"
				>
					<option value="Todas">Todas las categorías</option>
					{#each categorias as cat}
						<option value={cat.id_categoria}>{cat.nombre_categoria}</option>
					{/each}
				</select>
				{#if auth.user?.rol?.toLowerCase() === 'admin'}
					<button
						type="button"
						title="Gestionar categorías"
						onclick={openCategoriesModal}
						class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg text-text-muted hover:bg-border-color hover:text-primario disabled:cursor-not-allowed disabled:opacity-50"
					>
						<Settings size={14} strokeWidth={2.5} />
					</button>
				{/if}
			</div>

			<!-- filtro estados -->
			<select
				bind:value={selectedEstado}
				class="w-full sm:w-auto rounded-xl border border-border-color bg-bg-card py-2.5 pl-4 pr-10 text-sm text-text-primary focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario cursor-pointer"
			>
				<option value="Todas">Todos los estados</option>
				<option value="Disponible">Disponible</option>
				<option value="Bajo stock">Bajo stock</option>
				<option value="Agotado">Agotado</option>
			</select>

			<button
				type="button"
				title="Limpiar filtros"
				class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg text-text-muted hover:bg-border-color hover:text-primario disabled:cursor-not-allowed disabled:opacity-50"
				onclick={() => {
					searchQuery = '';
					selectedCategoria = 'Todas';
					selectedEstado = 'Todas';
				}}
				disabled={!searchQuery && selectedCategoria === 'Todas' && selectedEstado === 'Todas'}
			>
				<X size={14} strokeWidth={2.5} />
			</button>
		</div>

		<!-- boton  -->
		{#if auth.user?.rol?.toLowerCase() === 'admin'}
			<button
				onclick={openCreateModal}
				class="flex w-full items-center justify-center gap-2 rounded-xl bg-primario px-5 py-2.5 text-sm font-semibold text-white transition-transform hover:scale-[1.02] active:scale-95 sm:w-auto"
			>
				<Plus size={18} />
				Nuevo Producto
			</button>
		{/if}
	</div>

	<!-- tabla de informacion -->
	<div class="overflow-x-auto rounded-xl border border-border-color bg-bg-card shadow-sm">
		<table class="w-full whitespace-nowrap text-left text-sm text-text-primary">
			<thead class="border-b border-border-color bg-bg-primary/50 text-text-muted">
				<tr>
					<th class="px-6 py-4 font-semibold">Producto</th>
					<th class="px-6 py-4 font-semibold">Categoría</th>
					<th class="px-6 py-4 font-semibold">Precio</th>
					<th class="px-6 py-4 font-semibold">
						<button
							class="flex items-center gap-1 hover:text-text-primary"
							onclick={toggleSortStock}
						>
							STOCK
							{#if sortStock === 'asc'}
								<ArrowUp size={12} />
							{:else if sortStock === 'desc'}
								<ArrowDown size={12} />
							{:else}
								<ArrowUpDown size={12} class="opacity-50" />
							{/if}
						</button>
					</th>
					<th class="px-6 py-4 font-semibold">Unidad</th>
					<th class="px-6 py-4 font-semibold">Estado</th>
					{#if auth.user?.rol?.toLowerCase() === 'admin'}
						<th class="px-6 py-4 text-right font-semibold">Acciones</th>
					{/if}
				</tr>
			</thead>
			<tbody class="divide-y divide-border-color">
				{#if loading}
					<tr>
						<td colspan="7" class="px-6 py-12 text-center text-text-muted">
							Cargando productos...
						</td>
					</tr>
				{:else if filteredProductos.length === 0}
					<tr>
						<td colspan="7" class="px-6 py-12 text-center text-text-muted">
							No se encontraron productos que coincidan con la búsqueda.
						</td>
					</tr>
				{:else}
					{#each filteredProductos as producto (producto.id_producto)}
						{@const estado = getEstadoBadge(producto)}
						<tr class="group hover:bg-bg-primary/30">
							<td class="px-6 py-4">
								<div class="flex flex-col">
									<span class="font-bold text-text-primary">{producto.nombre}</span>
									{#if producto.marca}
										<span class="text-xs text-text-muted">{producto.marca}</span>
									{/if}
								</div>
							</td>
							<td class="px-6 py-4">
								<span
									class="rounded-full border border-accent-light/40 bg-accent-light/10 px-3 py-1 text-xs font-bold text-accent"
								>
									{producto.categoria?.nombre_categoria || 'Sin categoría'}
								</span>
							</td>
							<td class="px-6 py-4 font-bold text-text-primary">
								{formatCurrency(producto.precio)}
							</td>
							<td class="px-6 py-4 font-semibold text-text-primary">
								{producto.stock}
							</td>
							<td class="px-6 py-4 text-text-secondary capitalize">
								{producto.unidad || 'unidades'}
							</td>
							<td class="px-6 py-4">
								<span
									class="rounded-full border px-2.5 py-1 text-[0.70rem] font-bold {estado.classes}"
								>
									{estado.text}
								</span>
							</td>
							{#if auth.user?.rol?.toLowerCase() === 'admin'}
								<td class="px-6 py-4 text-right">
									<div class="flex items-center justify-end gap-2">
										<button
											class="rounded-lg p-2 text-text-muted hover:bg-border-color hover:text-primario"
											onclick={() => openEditModal(producto)}
											title="Editar"
										>
											<Edit2 size={18} />
										</button>
										<button
											class="rounded-lg p-2 text-text-muted hover:bg-danger-bg hover:text-danger-color"
											onclick={() => confirmDelete(producto)}
											title="Eliminar"
										>
											<Trash2 size={18} />
										</button>
									</div>
								</td>
							{/if}
						</tr>
					{/each}
				{/if}
			</tbody>
		</table>
	</div>
</div>

<!-- Crear/Editar -->
{#if showModal}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-text-primary/40 p-4 backdrop-blur-[4px] animate-modal-enter"
		onclick={closeModal}
		role="presentation"
	>
		<div
			class="w-full max-w-2xl overflow-hidden rounded-2xl border border-border-color bg-bg-card shadow-2xl"
			onclick={(e) => e.stopPropagation()}
			role="dialog"
		>
			<header class="flex items-center justify-between border-b border-border-color px-6 py-4">
				<h3 class="text-lg font-bold text-text-primary">
					{editingProducto ? 'Editar Producto' : 'Añadir Nuevo Producto'}
				</h3>
				<button
					class="rounded-lg p-1 text-text-muted hover:bg-border-color hover:text-text-primary"
					onclick={closeModal}
				>
					<X size={20} />
				</button>
			</header>

			<form onsubmit={handleSubmit} autocomplete="off" class="p-6 max-h-[80vh] overflow-y-auto">
				{#if formGeneralError}
					<div
						class="mb-5 flex gap-3 rounded-lg border border-red-500/15 bg-danger-bg p-4 text-sm text-danger-color"
						role="alert"
					>
						<span>{formGeneralError}</span>
					</div>
				{/if}

				<div class="grid grid-cols-1 gap-5 md:grid-cols-2">
					<div class="flex flex-col gap-1.5 md:col-span-2">
						<label class="text-sm font-semibold text-text-primary" for="formNombre"
							>Nombre del Producto *</label
						>
						<input
							type="text"
							id="formNombre"
							autocomplete="off"
							class="rounded-xl border border-border-color bg-bg-primary px-4 py-2.5 text-sm text-text-primary focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario disabled:opacity-50"
							placeholder="Ej: Manzana Roja"
							bind:value={formNombre}
							onblur={handleNombreBlur}
							disabled={submitLoading}
							required
						/>
						{#if errNombre}<span class="mt-1 text-xs font-medium text-danger-color"
								>{errNombre}</span
							>{/if}
					</div>

					<div class="flex flex-col gap-1.5 md:col-span-2">
						<label class="text-sm font-semibold text-text-primary" for="formDescripcion"
							>Descripción</label
						>
						<textarea
							id="formDescripcion"
							rows="2"
							class="rounded-xl border border-border-color bg-bg-primary px-4 py-2.5 text-sm text-text-primary focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario disabled:opacity-50 resize-none"
							placeholder="Breve descripción del producto..."
							bind:value={formDescripcion}
							disabled={submitLoading}></textarea>
					</div>

					<div class="flex flex-col gap-1.5 md:col-span-1">
						<label class="text-sm font-semibold text-text-primary" for="formUnidad">Unidad *</label>
						<select
							id="formUnidad"
							class="rounded-xl border border-border-color bg-bg-primary px-4 py-2.5 text-sm text-text-primary focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario disabled:opacity-50"
							bind:value={formUnidad}
							disabled={submitLoading}
							required
						>
							<option value="unidades">Unidades</option>
							<option value="kg">Kilogramos (kg)</option>
						</select>
					</div>

					<div class="flex flex-col gap-1.5 md:col-span-1">
						<label class="text-sm font-semibold text-text-primary" for="formCategoria"
							>Categoría *</label
						>
						<div class="relative flex items-center">
							<input
								type="text"
								class="absolute inset-0 z-0 w-full rounded-xl border border-transparent bg-transparent px-4 py-2.5 text-sm text-text-secondary/40 outline-none"
								value={inlineSuggestion}
								disabled
							/>
							<input
								type="text"
								id="formCategoria"
								autocomplete="off"
								class="relative z-10 w-full rounded-xl border border-border-color bg-transparent px-4 py-2.5 text-sm text-text-primary focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario disabled:opacity-50"
								placeholder="Ej: Frutas, Abarrotes..."
								bind:value={formCategoriaNombre}
								onkeydown={handleCategoriaKeyDown}
								onblur={handleCategoriaBlur}
								disabled={submitLoading}
								required
							/>
						</div>
						{#if errCategoria}<span class="mt-1 text-xs font-medium text-danger-color"
								>{errCategoria}</span
							>{/if}
					</div>

					<div class="flex flex-col gap-1.5">
						<label class="text-sm font-semibold text-text-primary" for="formPrecio"
							>Precio (CLP) *</label
						>
						<input
							type="number"
							min="0"
							id="formPrecio"
							class="[appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none rounded-xl border border-border-color bg-bg-primary px-4 py-2.5 text-sm text-text-primary focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario disabled:opacity-50"
							bind:value={formPrecio}
							disabled={submitLoading}
							required
							onfocus={(e) => e.currentTarget.select()}
						/>
						{#if errPrecio}<span class="mt-1 text-xs font-medium text-danger-color"
								>{errPrecio}</span
							>{/if}
					</div>

					<div class="flex flex-col gap-1.5">
						<label class="text-sm font-semibold text-text-primary" for="formStock"
							>Stock Actual</label
						>
						<input
							type="number"
							id="formStock"
							class="[appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none rounded-xl border border-border-color bg-bg-primary px-4 py-2.5 text-sm text-text-primary focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario disabled:opacity-50"
							bind:value={formStock}
							disabled={submitLoading}
							required
							onfocus={(e) => e.currentTarget.select()}
						/>
					</div>

					<div class="flex flex-col gap-1.5">
						<label class="text-sm font-semibold text-text-primary" for="formStockMinimo"
							>Alerta Stock Mínimo</label
						>
						<input
							type="number"
							id="formStockMinimo"
							min="0"
							class="[appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none rounded-xl border border-border-color bg-bg-primary px-4 py-2.5 text-sm text-text-primary focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario disabled:opacity-50"
							bind:value={formStockMinimo}
							disabled={submitLoading || !!editingProducto}
							onfocus={(e) => e.currentTarget.select()}
						/>
					</div>

					<div class="flex flex-col gap-1.5">
						<label class="text-sm font-semibold text-text-primary" for="formMarca">Marca</label>
						<input
							type="text"
							id="formMarca"
							class="rounded-xl border border-border-color bg-bg-primary px-4 py-2.5 text-sm text-text-primary focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario disabled:opacity-50"
							placeholder="Ej: Coca-Cola"
							bind:value={formMarca}
							disabled={submitLoading}
						/>
					</div>

					<div class="flex flex-col gap-1.5">
						<label class="text-sm font-semibold text-text-primary" for="formCodigoBarras"
							>Código de Barras *</label
						>
						<input
							type="text"
							id="formCodigoBarras"
							autocomplete="off"
							class="rounded-xl border border-border-color bg-bg-primary px-4 py-2.5 text-sm text-text-primary focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario disabled:opacity-50"
							placeholder="Ej: 7801234567890"
							bind:value={formCodigoBarras}
							disabled={submitLoading}
							required
						/>
						{#if errCodigoBarras}<span class="mt-1 text-xs font-medium text-danger-color"
								>{errCodigoBarras}</span
							>{/if}
					</div>
				</div>

				<div class="mt-8 flex justify-end gap-4">
					<button
						type="button"
						class="font-semibold text-text-muted hover:text-text-primary"
						onclick={closeModal}
						disabled={submitLoading}>Cancelar</button
					>
					<button
						type="submit"
						class="rounded-xl bg-primario px-6 py-2.5 text-sm font-semibold text-white transition-transform hover:scale-[1.02] active:scale-95 disabled:cursor-not-allowed disabled:opacity-70"
						disabled={submitLoading}
					>
						{#if submitLoading}
							Procesando...
						{:else}
							{editingProducto ? 'Guardar Cambios' : 'Crear Producto'}
						{/if}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<!-- borrar confirmacion modeal -->
{#if showDeleteModal && productoToDelete}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-text-primary/40 p-5 backdrop-blur-[4px]"
		onclick={closeDeleteModal}
		role="presentation"
	>
		<div
			class="w-full max-w-[500px] overflow-hidden rounded-xl border border-border-color bg-bg-card shadow-lg animate-modal-enter"
			onclick={(e) => e.stopPropagation()}
			role="dialog"
		>
			<header class="flex items-center justify-between border-b border-border-color p-5">
				<h2 class="text-lg font-bold text-danger-color">Confirmar Eliminación</h2>
				<button
					class="inline-flex cursor-pointer items-center justify-center rounded-lg border border-border-color bg-text-primary/3 p-2 text-text-secondary hover:border-border-color-hover hover:bg-text-primary/7 hover:text-text-primary"
					onclick={closeDeleteModal}>&times;</button
				>
			</header>
			<div class="p-6 text-sm text-text-primary">
				<p class="mb-3 text-base">
					¿Estás seguro de que deseas eliminar el producto <strong class="font-bold"
						>{productoToDelete.nombre}</strong
					> de forma permanente?
				</p>
				<p class="text-text-muted">No se podrá deshacer esta acción.</p>
			</div>
			<footer
				class="flex justify-end gap-3 border-t border-border-color bg-text-primary/2 p-4 px-6"
			>
				<button
					class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-border-color bg-bg-secondary px-5 py-2.5 text-sm font-semibold text-text-primary hover:bg-text-primary/5"
					onclick={closeDeleteModal}
					disabled={submitLoading}>Cancelar</button
				>
				<button
					class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-red-500/20 bg-danger-bg px-5 py-2.5 text-sm font-semibold text-danger-color hover:bg-red-500/20 disabled:cursor-not-allowed disabled:opacity-50"
					onclick={handleDelete}
					disabled={submitLoading}
				>
					{#if submitLoading}
						Procesando...
					{:else}
						Eliminar Permanentemente
					{/if}
				</button>
			</footer>
		</div>
	</div>
{/if}

<!-- abrir la interfaz de categorias -->
{#if showCategoriesModal}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-text-primary/40 p-4 backdrop-blur-[4px] animate-modal-enter"
		onclick={closeCategoriesModal}
		role="presentation"
	>
		<div
			class="w-full max-w-md overflow-hidden rounded-2xl border border-border-color bg-bg-card shadow-2xl"
			onclick={(e) => e.stopPropagation()}
			role="dialog"
		>
			<header class="flex items-center justify-between border-b border-border-color px-6 py-4">
				<h3 class="text-lg font-bold text-text-primary">Gestionar Categorías</h3>
				<button
					class="rounded-lg p-1 text-text-muted hover:bg-border-color hover:text-text-primary"
					onclick={closeCategoriesModal}
				>
					<X size={20} />
				</button>
			</header>

			<div class="p-6">
				{#if catGeneralError}
					<div
						class="mb-5 flex gap-3 rounded-lg border border-red-500/15 bg-danger-bg p-4 text-sm text-danger-color"
						role="alert"
					>
						<span>{catGeneralError}</span>
					</div>
				{/if}

				<div class="mb-4 relative">
					<Search size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-text-muted" />
					<input
						type="text"
						class="w-full rounded-xl border border-border-color bg-bg-primary py-2.5 pl-9 pr-4 text-sm text-text-primary focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario"
						placeholder="Buscar categoría..."
						bind:value={catSearchQuery}
					/>
				</div>

				{#if filteredCategorias.length === 0}
					<p class="text-center text-sm text-text-muted py-4">
						{catSearchQuery ? 'No se encontraron categorías.' : 'No hay categorías registradas.'}
					</p>
				{:else}
					<ul class="flex flex-col gap-3 max-h-60 overflow-y-auto pr-2">
						{#each filteredCategorias as cat (cat.id_categoria)}
							<li
								class="flex items-center justify-between rounded-xl border border-border-color bg-bg-primary px-4 py-3"
							>
								{#if editingCatId === cat.id_categoria}
									<div class="flex flex-1 items-center gap-2 mr-2">
										<input
											type="text"
											class="w-full rounded-lg border border-border-color bg-bg-card px-3 py-1.5 text-sm text-text-primary focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario"
											bind:value={editingCatName}
											disabled={manageCatLoading}
										/>
										<button
											class="shrink-0 rounded-lg bg-primario px-3 py-1.5 text-xs font-semibold text-white hover:bg-primario-hover disabled:opacity-50"
											onclick={() => saveEditCat(cat.id_categoria)}
											disabled={manageCatLoading}>Guardar</button
										>
										<button
											class="shrink-0 rounded-lg px-3 py-1.5 text-xs font-semibold text-text-muted hover:text-text-primary disabled:opacity-50"
											onclick={cancelEditCat}
											disabled={manageCatLoading}>Cancelar</button
										>
									</div>
								{:else}
									<span class="text-sm font-medium text-text-primary">{cat.nombre_categoria}</span>
									<div class="flex items-center gap-1">
										<button
											class="rounded-lg p-2 text-text-muted hover:bg-border-color hover:text-primario disabled:opacity-50"
											title="Editar"
											onclick={() => startEditCat(cat)}
											disabled={manageCatLoading || editingCatId !== null}
										>
											<Edit2 size={16} />
										</button>
										<button
											class="rounded-lg p-2 text-text-muted hover:bg-danger-bg hover:text-danger-color disabled:opacity-50"
											title="Eliminar"
											onclick={() => deleteCat(cat.id_categoria)}
											disabled={manageCatLoading || editingCatId !== null}
										>
											<Trash2 size={16} />
										</button>
									</div>
								{/if}
							</li>
						{/each}
					</ul>
				{/if}
			</div>
		</div>
	</div>
{/if}

<!-- Scanner Modal para la cam -->
{#if modoEscaneo}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-text-primary/40 p-5 backdrop-blur-[4px]"
		onclick={() => (modoEscaneo = false)}
		role="presentation"
	>
		<div
			class="w-full max-w-[460px] overflow-hidden rounded-xl border border-border-color bg-bg-card shadow-lg animate-modal-enter"
			onclick={(e) => e.stopPropagation()}
			role="dialog"
		>
			<header class="flex items-center justify-between border-b border-border-color p-5">
				<h2 class="text-base font-bold text-text-primary">Escanear Código de Barras</h2>
				<button
					class="inline-flex cursor-pointer items-center justify-center rounded-lg border border-border-color bg-text-primary/3 p-2 text-text-secondary hover:border-border-color-hover hover:bg-text-primary/7 hover:text-text-primary"
					onclick={() => (modoEscaneo = false)}>&times;</button
				>
			</header>
			<div class="p-6 flex flex-col items-center justify-center">
				<Scanner onScan={manejarEscaneo} onClose={() => (modoEscaneo = false)} />
			</div>
		</div>
	</div>
{/if}
