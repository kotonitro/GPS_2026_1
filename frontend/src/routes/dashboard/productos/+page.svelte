<script lang="ts">
	import { apiProductos, apiCategorias, type Producto, type Categoria } from '$lib/api';
	import { auth } from '$lib/authStore.svelte';
	import { toast } from '$lib/toastStore.svelte';
	import { onMount } from 'svelte';
	import { Search, Plus, Edit2, Trash2, X, AlertTriangle, ArrowUpDown, ArrowUp, ArrowDown } from '@lucide/svelte';

	let productos = $state<Producto[]>([]);
	let categorias = $state<Categoria[]>([]);
	let loading = $state(true);
	let errorMsg = $state('');

	// Filtros y Búsqueda
	let searchQuery = $state('');
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
            result = result.filter(p => 
                p.nombre.toLowerCase().includes(query) || 
                (p.codigo_barras && p.codigo_barras.toLowerCase().includes(query))
            );
        }
        
        // Filtro Categoría
        if (selectedCategoria !== 'Todas') {
            result = result.filter(p => p.id_categoria === selectedCategoria);
        }
        
        // Filtro Estado
        if (selectedEstado !== 'Todas') {
            result = result.filter(p => {
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
        const match = categorias.find(c => c.nombre_categoria.toLowerCase().startsWith(lowerInput));
        if (match) {
            return formCategoriaNombre + match.nombre_categoria.substring(formCategoriaNombre.length);
        }
        return '';
    });
	// Errores de Formulario
	let errNombre = $state('');
	let errPrecio = $state('');
    let errCategoria = $state('');
	let formGeneralError = $state('');
	let submitLoading = $state(false);

	onMount(loadData);

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
            return { text: 'Agotado', classes: 'border-red-500/20 bg-red-500/10 text-red-600 dark:text-red-400' };
        } else if (producto.stock <= producto.stock_minimo) {
            return { text: 'Bajo stock', classes: 'border-orange-500/20 bg-orange-500/10 text-orange-600 dark:text-orange-400' };
        } else {
            return { text: 'Disponible', classes: 'border-green-500/20 bg-green-500/10 text-green-600 dark:text-green-400' };
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
			formNombre = formNombre.replace(/\b\w/g, c => c.toUpperCase());
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

	function clearErrors() {
		errNombre = '';
		errPrecio = '';
        errCategoria = '';
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
        if (formPrecio < 0) {
            errPrecio = 'El precio no puede ser negativo.';
            isValid = false;
        }
        if (!formCategoriaNombre.trim()) {
            errCategoria = 'Debe especificar una categoría.';
            isValid = false;
        }

		if (!isValid) return;

		submitLoading = true;
        let selectedCatId = '';

        if (formCategoriaNombre.trim()) {
            const existingCat = categorias.find(c => c.nombre_categoria.toLowerCase() === formCategoriaNombre.trim().toLowerCase());
            if (existingCat) {
                selectedCatId = existingCat.id_categoria;
            } else {
                try {
                    const resCat = await apiCategorias.create({ nombre_categoria: formCategoriaNombre.trim() });
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
	<!-- Page Header & Stats -->

    <div class="mb-8 flex flex-wrap items-center justify-end gap-4">
		{#if auth.user?.rol?.toLowerCase() === 'admin'}
            <button
                onclick={openCreateModal}
                class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg bg-gradient-to-r from-accent-light to-accent px-5 py-2.5 text-sm font-semibold text-white transition-all duration-200 hover:-translate-y-[1px] hover:shadow-glow disabled:cursor-not-allowed disabled:opacity-50"
            >
                <Plus size={18} strokeWidth={2.5} />
                <span>Agregar producto</span>
            </button>
        {/if}
	</div>

	<!-- Main Content Area -->
	<div class="flex flex-col gap-6">
		<!-- Search and Filters Bar -->
		<div class="mb-6 w-full md:w-fit rounded-xl border border-border-color bg-bg-card p-6 shadow-md transition-all duration-300 hover:border-border-color-hover hover:shadow-lg">
			<div class="flex flex-col sm:flex-row items-center gap-4">
				<div class="flex min-w-[280px] w-full sm:w-80 overflow-hidden rounded-lg border border-[rgba(15,30,54,0.15)] bg-white focus-within:border-accent focus-within:ring-2 focus-within:ring-accent/15 dark:bg-bg-primary">
					<div class="flex items-center justify-center pl-4 pr-2 text-text-muted">
                        <Search size={18} />
					</div>
					<input
						type="text"
						class="flex-1 border-none bg-transparent px-2 py-3 text-sm text-text-primary outline-none focus:ring-0"
						placeholder="Buscar producto..."
						bind:value={searchQuery}
					/>
				</div>
                <div class="flex gap-4 w-full sm:w-auto">
                    <!-- Category Filter -->
                    <select 
                        bind:value={selectedCategoria}
                        class="w-full sm:w-48 rounded-lg border border-[rgba(15,30,54,0.15)] bg-white px-4 py-3 text-sm font-medium text-text-primary outline-none transition-all focus:border-accent focus:ring-2 focus:ring-accent/15 dark:bg-bg-primary"
                    >
                        <option value="Todas">Todas las categorías</option>
                        {#each categorias as cat (cat.id_categoria)}
                            <option value={cat.id_categoria}>{cat.nombre_categoria}</option>
                        {/each}
                    </select>
                    
                    <!-- Status Filter -->
                    <select 
                        bind:value={selectedEstado}
                        class="w-full sm:w-48 rounded-lg border border-[rgba(15,30,54,0.15)] bg-white px-4 py-3 text-sm font-medium text-text-primary outline-none transition-all focus:border-accent focus:ring-2 focus:ring-accent/15 dark:bg-bg-primary"
                    >
                        <option value="Todas">Todos los estados</option>
                        <option value="Disponible">Disponible</option>
                        <option value="Bajo stock">Bajo stock</option>
                        <option value="Agotado">Agotado</option>
                    </select>

                    <button
                        type="button"
                        title="Limpiar filtros"
                        class="flex h-[46px] w-[46px] shrink-0 items-center justify-center rounded-lg border border-border-color bg-bg-secondary text-text-muted transition-all duration-200 hover:bg-text-primary/10 hover:text-text-primary focus:border-accent focus:outline-none focus:ring-2 focus:ring-accent/15 disabled:cursor-not-allowed disabled:opacity-30"
                        onclick={() => { searchQuery = ''; selectedCategoria = 'Todas'; selectedEstado = 'Todas'; }}
                        disabled={!searchQuery && selectedCategoria === 'Todas' && selectedEstado === 'Todas'}
                    >
                        <X size={18} />
                    </button>
                </div>
			</div>
		</div>

		<!-- Data Table -->
		<div class="overflow-hidden rounded-xl border border-[rgba(15,30,54,0.08)] bg-bg-card shadow-sm">
			<div class="overflow-x-auto">
				<table class="w-full text-left text-sm text-text-secondary">
					<thead class="border-b border-[rgba(15,30,54,0.08)] bg-gray-50/50 text-[0.65rem] font-bold uppercase tracking-wider text-text-muted dark:bg-bg-primary/20">
						<tr>
							<th class="px-6 py-4">Producto</th>
							<th class="px-6 py-4">Categoría</th>
							<th class="px-6 py-4">Precio</th>
                            <th class="px-6 py-4">
                                <button class="flex items-center gap-1 hover:text-text-primary" onclick={toggleSortStock}>
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
                            <th class="px-6 py-4">Unidad</th>
                            <th class="px-6 py-4">Estado</th>
                            {#if auth.user?.rol?.toLowerCase() === 'admin'}
							    <th class="px-6 py-4 text-right">Acciones</th>
                            {/if}
						</tr>
					</thead>
					<tbody>
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
								<tr class="group border-b border-[rgba(15,30,54,0.04)] bg-bg-card transition-colors hover:bg-gray-50/80 last:border-0 dark:hover:bg-bg-primary/40">
									<td class="px-6 py-4">
                                        <div class="flex flex-col">
                                            <span class="font-bold text-text-primary">{producto.nombre}</span>
                                            {#if producto.marca}
                                                <span class="text-xs text-text-muted">{producto.marca}</span>
                                            {/if}
                                        </div>
									</td>
                                    <td class="px-6 py-4">
                                        <span class="rounded-full border border-accent-light/40 bg-accent-light/10 px-3 py-1 text-xs font-bold text-accent">
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
                                        <span class="rounded-full border px-2.5 py-1 text-[0.70rem] font-bold {estado.classes}">
                                            {estado.text}
                                        </span>
                                    </td>
                                    {#if auth.user?.rol?.toLowerCase() === 'admin'}
                                        <td class="px-6 py-4 text-right">
                                            <div class="flex justify-end gap-2">
                                                <button
                                                    onclick={() => openEditModal(producto)}
                                                    class="inline-flex cursor-pointer items-center justify-center rounded-lg border border-border-color bg-text-primary/3 p-2 text-text-secondary transition-all duration-200 hover:border-border-color-hover hover:bg-text-primary/7 hover:text-text-primary"
                                                    aria-label="Editar"
                                                    title="Editar"
                                                >
                                                    <Edit2 size={16} />
                                                </button>
                                                <button
                                                    onclick={() => confirmDelete(producto)}
                                                    class="inline-flex cursor-pointer items-center justify-center rounded-lg border border-border-color bg-text-primary/3 p-2 text-text-secondary transition-all duration-200 hover:border-red-500/20 hover:bg-danger-bg hover:text-danger-color"
                                                    aria-label="Eliminar"
                                                    title="Eliminar"
                                                >
                                                    <Trash2 size={16} />
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
	</div>
</div>

<!-- Modal Crear/Editar -->
{#if showModal}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-text-primary/40 p-5 backdrop-blur-[4px]" onclick={closeModal} role="presentation">
		<div class="w-full max-w-2xl overflow-hidden rounded-xl border border-border-color bg-bg-card shadow-lg animate-modal-enter" onclick={(e) => e.stopPropagation()} role="dialog">
			<header class="flex items-center justify-between border-b border-border-color p-5">
				<h2 class="text-lg font-bold text-text-primary">{editingProducto ? 'Editar Producto' : 'Nuevo Producto'}</h2>
				<button class="inline-flex cursor-pointer items-center justify-center rounded-lg border border-border-color bg-text-primary/3 p-2 text-text-secondary transition-all duration-200 hover:border-border-color-hover hover:bg-text-primary/7 hover:text-text-primary" onclick={closeModal}>&times;</button>
			</header>

			<form onsubmit={handleSubmit} autocomplete="off" class="p-6 max-h-[80vh] overflow-y-auto">
					{#if formGeneralError}
						<div class="mb-6 rounded-xl border border-red-500/20 bg-danger-bg p-4 text-sm text-danger-color">
							{formGeneralError}
						</div>
					{/if}
					
                    <div class="grid grid-cols-1 gap-5 md:grid-cols-2 mb-5">
                        <div class="flex flex-col gap-1.5 md:col-span-2">
                            <label class="text-[0.85rem] font-semibold text-text-secondary" for="formNombre">Nombre del Producto *</label>
                            <input type="text" id="formNombre" autocomplete="off" class="rounded-lg border border-[rgba(15,30,54,0.15)] bg-white px-4 py-3 text-sm text-text-primary outline-none transition-all duration-200 focus:border-accent focus:ring-2 focus:ring-accent/15 disabled:opacity-50 dark:bg-bg-primary" placeholder="Ej: Manzana Roja" bind:value={formNombre} onblur={handleNombreBlur} disabled={submitLoading} required />
                            {#if errNombre}<span class="mt-1 text-xs font-medium text-danger-color">{errNombre}</span>{/if}
                        </div>

                        <div class="flex flex-col gap-1.5 md:col-span-2">
                            <label class="text-[0.85rem] font-semibold text-text-secondary" for="formDescripcion">Descripción</label>
                            <textarea id="formDescripcion" rows="2" class="rounded-lg border border-[rgba(15,30,54,0.15)] bg-white px-4 py-3 text-sm text-text-primary outline-none transition-all duration-200 focus:border-accent focus:ring-2 focus:ring-accent/15 disabled:opacity-50 dark:bg-bg-primary resize-none" placeholder="Breve descripción del producto..." bind:value={formDescripcion} disabled={submitLoading}></textarea>
                        </div>
                        
                        <div class="flex flex-col gap-1.5 md:col-span-1">
                            <label class="text-[0.85rem] font-semibold text-text-secondary" for="formUnidad">Unidad *</label>
                            <select id="formUnidad" class="rounded-lg border border-[rgba(15,30,54,0.15)] bg-white px-4 py-3 text-sm text-text-primary outline-none transition-all duration-200 focus:border-accent focus:ring-2 focus:ring-accent/15 disabled:opacity-50 dark:bg-bg-primary" bind:value={formUnidad} disabled={submitLoading} required>
                                <option value="unidades">Unidades</option>
                                <option value="kg">Kilogramos (kg)</option>
                            </select>
                        </div>
                        
                        <div class="flex flex-col gap-1.5 md:col-span-1">
                            <label class="text-[0.85rem] font-semibold text-text-secondary" for="formCategoria">Categoría *</label>
                            <div class="relative flex items-center">
                                <input type="text" class="absolute inset-0 z-0 w-full rounded-lg border border-transparent bg-white px-4 py-3 text-sm text-text-secondary/40 outline-none dark:bg-bg-primary" value={inlineSuggestion} disabled />
                                <input type="text" id="formCategoria" autocomplete="off" class="relative z-10 w-full rounded-lg border border-[rgba(15,30,54,0.15)] bg-transparent px-4 py-3 text-sm text-text-primary outline-none transition-all duration-200 focus:border-accent focus:ring-2 focus:ring-accent/15 disabled:opacity-50" placeholder="Ej: Frutas, Abarrotes..." bind:value={formCategoriaNombre} onkeydown={handleCategoriaKeyDown} disabled={submitLoading} required />
                            </div>
                            {#if errCategoria}<span class="mt-1 text-xs font-medium text-danger-color">{errCategoria}</span>{/if}
                        </div>
                        
                        <div class="flex flex-col gap-1.5">
                            <label class="text-[0.85rem] font-semibold text-text-secondary" for="formPrecio">Precio (CLP) *</label>
                            <input type="number" min="0" id="formPrecio" class="[appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none rounded-lg border border-[rgba(15,30,54,0.15)] bg-white px-4 py-3 text-sm text-text-primary outline-none transition-all duration-200 focus:border-accent focus:ring-2 focus:ring-accent/15 disabled:opacity-50 dark:bg-bg-primary" bind:value={formPrecio} disabled={submitLoading} required onfocus={(e) => e.currentTarget.select()} />
                            {#if errPrecio}<span class="mt-1 text-xs font-medium text-danger-color">{errPrecio}</span>{/if}
                        </div>

                        <div class="flex flex-col gap-1.5">
                            <label class="text-[0.85rem] font-semibold text-text-secondary" for="formStock">Stock Actual</label>
                            <input type="number" id="formStock" class="[appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none rounded-lg border border-[rgba(15,30,54,0.15)] bg-white px-4 py-3 text-sm text-text-primary outline-none transition-all duration-200 focus:border-accent focus:ring-2 focus:ring-accent/15 disabled:opacity-50 dark:bg-bg-primary" bind:value={formStock} disabled={submitLoading} required onfocus={(e) => e.currentTarget.select()} />
                        </div>
                        
                        <div class="flex flex-col gap-1.5">
                            <label class="text-[0.85rem] font-semibold text-text-secondary" for="formStockMinimo">Alerta Stock Mínimo</label>
                            <input type="number" id="formStockMinimo" min="0" class="[appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none rounded-lg border border-[rgba(15,30,54,0.15)] bg-white px-4 py-3 text-sm text-text-primary outline-none transition-all duration-200 focus:border-accent focus:ring-2 focus:ring-accent/15 disabled:opacity-50 dark:bg-bg-primary" bind:value={formStockMinimo} disabled={submitLoading || !!editingProducto} onfocus={(e) => e.currentTarget.select()} />
                        </div>

                        <div class="flex flex-col gap-1.5">
                            <label class="text-[0.85rem] font-semibold text-text-secondary" for="formMarca">Marca</label>
                            <input type="text" id="formMarca" class="rounded-lg border border-[rgba(15,30,54,0.15)] bg-white px-4 py-3 text-sm text-text-primary outline-none transition-all duration-200 focus:border-accent focus:ring-2 focus:ring-accent/15 disabled:opacity-50 dark:bg-bg-primary" placeholder="Ej: Coca-Cola" bind:value={formMarca} disabled={submitLoading} />
                        </div>
                        
                        <div class="flex flex-col gap-1.5">
                            <label class="text-[0.85rem] font-semibold text-text-secondary" for="formCodigoBarras">Código de Barras *</label>
                            <input type="text" id="formCodigoBarras" autocomplete="off" class="rounded-lg border border-[rgba(15,30,54,0.15)] bg-white px-4 py-3 text-sm text-text-primary outline-none transition-all duration-200 focus:border-accent focus:ring-2 focus:ring-accent/15 disabled:opacity-50 dark:bg-bg-primary" placeholder="Ej: 780123456789" bind:value={formCodigoBarras} disabled={submitLoading} required />
                        </div>
                    </div>
                    
					<footer class="flex justify-end gap-3 border-t border-border-color bg-text-primary/2 p-4 px-6 mt-6">
						<button type="button" class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-border-color bg-bg-secondary px-5 py-2.5 text-sm font-semibold text-text-primary transition-all duration-200 hover:bg-text-primary/5" onclick={closeModal} disabled={submitLoading}>Cancelar</button>
						<button type="submit" class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg bg-gradient-to-r from-accent-light to-accent px-5 py-2.5 text-sm font-semibold text-white transition-all duration-200 hover:-translate-y-[1px] hover:shadow-glow disabled:cursor-not-allowed disabled:opacity-50" disabled={submitLoading}>
							{#if submitLoading} Procesando... {:else} Guardar Cambios {/if}
						</button>
					</footer>
				</form>
		</div>
	</div>
{/if}

<!-- Delete Confirmation Modal -->
{#if showDeleteModal && productoToDelete}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-text-primary/40 p-5 backdrop-blur-[4px]" onclick={closeDeleteModal} role="presentation">
		<div class="w-full max-w-[500px] overflow-hidden rounded-xl border border-border-color bg-bg-card shadow-lg animate-modal-enter" onclick={(e) => e.stopPropagation()} role="dialog">
			<header class="flex items-center justify-between border-b border-border-color p-5">
				<h2 class="text-lg font-bold text-danger-color">Confirmar Eliminación</h2>
				<button class="inline-flex cursor-pointer items-center justify-center rounded-lg border border-border-color bg-text-primary/3 p-2 text-text-secondary transition-all duration-200 hover:border-border-color-hover hover:bg-text-primary/7 hover:text-text-primary" onclick={closeDeleteModal}>&times;</button>
			</header>
			<div class="p-6 text-sm text-text-primary">
				<p class="mb-3 text-base">¿Estás seguro de que deseas eliminar el producto <strong class="font-bold">{productoToDelete.nombre}</strong> de forma permanente?</p>
				<p class="text-text-muted">No se podrá deshacer esta acción.</p>
			</div>
			<footer class="flex justify-end gap-3 border-t border-border-color bg-text-primary/2 p-4 px-6">
				<button class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-border-color bg-bg-secondary px-5 py-2.5 text-sm font-semibold text-text-primary transition-all duration-200 hover:bg-text-primary/5" onclick={closeDeleteModal} disabled={submitLoading}>Cancelar</button>
				<button class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-red-500/20 bg-danger-bg px-5 py-2.5 text-sm font-semibold text-danger-color transition-all duration-200 hover:bg-red-500/20 disabled:cursor-not-allowed disabled:opacity-50" onclick={handleDelete} disabled={submitLoading}>
					{#if submitLoading} Procesando... {:else} Eliminar Permanentemente {/if}
				</button>
			</footer>
		</div>
	</div>
{/if}
