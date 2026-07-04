<script lang="ts">
	import { apiPromociones } from '$lib/api'; // Asegúrate de crear esto en tu $lib/api
	import { auth } from '$lib/authStore.svelte';
	import { toast } from '$lib/toastStore.svelte';
	import { onMount } from 'svelte';

	interface Promocion {
		id_promocion: string;
		tipo: 'NXM' | 'DESCUENTO';
		lleva: number;
		paga: number;
		descuento: number;
		created_at?: string;
	}

	let promociones = $state<Promocion[]>([]);
	let loading = $state(true);
	let errorMsg = $state('');

	// Búsqueda / Filtros
	let searchTipo = $state<'TODOS' | 'NXM' | 'DESCUENTO'>('TODOS');

	// Modales
	let showModal = $state(false);
	let showDeleteModal = $state(false);
	let editingPromo = $state<Promocion | null>(null);
	let promoToDelete = $state<Promocion | null>(null);

	// Campos de Formulario
	let formTipo = $state<'NXM' | 'DESCUENTO'>('NXM');
	let formLleva = $state<number>(0);
	let formPaga = $state<number>(0);
	let formDescuento = $state<number>(0);

	// Errores de Formulario
	let errLlevaPaga = $state('');
	let errDescuento = $state('');
	let formGeneralError = $state('');
	let submitLoading = $state(false);

	onMount(loadPromociones);

	async function loadPromociones() {
		loading = true;
		errorMsg = '';
		try {
			const res = await apiPromociones.getAll();
			promociones = Array.isArray(res) ? res : [];
		} catch (err: any) {
			errorMsg = err.message || 'Error al conectar con la base de datos de promociones.';
		} finally {
			loading = false;
		}
	}

	// Manejo de Búsqueda / Filtro
	async function handleSearch(e: Event) {
		e.preventDefault();
		loading = true;
		errorMsg = '';

		try {
			if (searchTipo === 'TODOS') {
				await loadPromociones();
			} else {
				const res = await apiPromociones.getByTipo(searchTipo);
				promociones = Array.isArray(res) ? res : [];
			}
		} catch (err: any) {
			if (err.status === 404) {
				promociones = [];
			} else {
				errorMsg = err.message || 'Error al filtrar promociones.';
			}
		} finally {
			loading = false;
		}
	}

	function clearSearch() {
		searchTipo = 'TODOS';
		loadPromociones();
	}

	// Modales
	function openCreateModal() {
		if (auth.user?.rol !== 'Admin') {
			toast.show('Acción denegada, requiere rol Administrador.', 'error');
			return;
		}
		editingPromo = null;
		formTipo = 'NXM';
		formLleva = 2;
		formPaga = 1;
		formDescuento = 0;
		clearErrors();
		showModal = true;
	}

	function openEditModal(promo: Promocion) {
		if (auth.user?.rol !== 'Admin') {
			toast.show('Acción denegada, requiere rol Administrador.', 'error');
			return;
		}
		editingPromo = promo;
		formTipo = promo.tipo;
		formLleva = promo.lleva || 0;
		formPaga = promo.paga || 0;
		formDescuento = promo.descuento || 0;
		clearErrors();
		showModal = true;
	}

	function openDeleteModal(promo: Promocion) {
		if (auth.user?.rol !== 'Admin') {
			toast.show('Acción denegada, requiere rol Administrador.', 'error');
			return;
		}
		promoToDelete = promo;
		showDeleteModal = true;
	}

	function clearErrors() {
		errLlevaPaga = '';
		errDescuento = '';
		formGeneralError = '';
	}

	// Submit Formulario
	async function handleSubmit(e: Event) {
		e.preventDefault();
		clearErrors();
		let isValid = true;

		// Validaciones lógicas según el tipo de promoción
		if (formTipo === 'NXM') {
			if (formLleva <= 0 || formPaga <= 0) {
				errLlevaPaga = 'Las cantidades deben ser mayores a 0.';
				isValid = false;
			} else if (formLleva <= formPaga) {
				errLlevaPaga = 'La cantidad que "Lleva" debe ser mayor a la que "Paga" (ej: Lleva 3, Paga 2).';
				isValid = false;
			}
			formDescuento = 0; // Limpiamos el otro campo por seguridad
		} else if (formTipo === 'DESCUENTO') {
			if (formDescuento <= 0) {
				errDescuento = 'El descuento debe ser un valor válido mayor a 0.';
				isValid = false;
			}
			formLleva = 0;
			formPaga = 0;
		}

		if (!isValid) return;

		submitLoading = true;
		const payload = {
			tipo: formTipo,
			lleva: formTipo === 'NXM' ? formLleva : 0,
			paga: formTipo === 'NXM' ? formPaga : 0,
			descuento: formTipo === 'DESCUENTO' ? formDescuento : 0
		};

		try {
			if (editingPromo) {
				await apiPromociones.update(editingPromo.id_promocion, payload);
				toast.show('Promoción modificada exitosamente.', 'success');
			} else {
				await apiPromociones.create(payload);
				toast.show('Promoción creada exitosamente.', 'success');
			}
			showModal = false;
			loadPromociones();
		} catch (err: any) {
			formGeneralError = err.message || 'Error al guardar la promoción.';
		} finally {
			submitLoading = false;
		}
	}

	// Confirmar Borrado
	async function confirmDelete() {
		if (!promoToDelete) return;
		try {
			await apiPromociones.delete(promoToDelete.id_promocion);
			toast.show('Promoción eliminada con éxito.', 'success');
			showDeleteModal = false;
			promoToDelete = null;
			loadPromociones();
		} catch (err: any) {
			toast.show(err.message || 'No se pudo eliminar la promoción.', 'error');
		}
	}
</script>

<svelte:head>
	<title>Promociones - GPSproject</title>
	<meta name="description" content="Gestión de promociones y descuentos" />
</svelte:head>

<div class="flex justify-between items-center mb-8 flex-wrap gap-4">
	<div>
		<h1 class="text-3xl font-bold tracking-tight text-text-primary">Gestión de Promociones</h1>
		<p class="text-text-secondary text-sm mt-1">Configura las ofertas NXM (ej. 2x1) y descuentos directos del inventario.</p>
	</div>
	{#if auth.user?.rol === 'Admin'}
		<button class="inline-flex items-center justify-center gap-2 font-semibold text-sm px-5 py-2.5 rounded-lg cursor-pointer bg-gradient-to-r from-accent-light to-accent text-white hover:shadow-glow hover:-translate-y-[1px] disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200" onclick={openCreateModal}>
			<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
				<line x1="12" y1="5" x2="12" y2="19"/>
				<line x1="5" y1="12" x2="19" y2="12"/>
			</svg>
			<span>Nueva Promoción</span>
		</button>
	{/if}
</div>

<div class="bg-bg-card border border-border-color rounded-xl p-6 shadow-md hover:border-border-color-hover hover:shadow-lg transition-all duration-300 mb-6">
	<form onsubmit={handleSearch} class="flex flex-wrap gap-4 items-center">
		<div class="flex flex-1 min-w-[280px] border border-[rgba(15,30,54,0.15)] rounded-lg overflow-hidden bg-white focus-within:border-accent focus-within:ring-2 focus-within:ring-accent/15">
			<select class="w-full border-none bg-transparent px-4 py-3 cursor-pointer outline-none text-sm text-text-primary" bind:value={searchTipo}>
				<option value="TODOS">Todas las Promociones</option>
				<option value="NXM">Solo Formato NXM (Lleva/Paga)</option>
				<option value="DESCUENTO">Solo Descuentos Directos</option>
			</select>
		</div>
		<div class="flex gap-2.5">
			<button type="submit" class="inline-flex items-center justify-center gap-2 font-semibold text-sm px-5 py-2.5 rounded-lg cursor-pointer bg-gradient-to-r from-accent-light to-accent text-white hover:shadow-glow hover:-translate-y-[1px] transition-all duration-200">
				<span>Filtrar</span>
			</button>
			{#if searchTipo !== 'TODOS'}
				<button type="button" class="inline-flex items-center justify-center gap-2 font-semibold text-sm px-5 py-2.5 rounded-lg cursor-pointer bg-bg-secondary text-text-primary border border-border-color hover:bg-text-primary/5 transition-all duration-200" onclick={clearSearch}>
					Limpiar
				</button>
			{/if}
		</div>
	</form>
</div>

{#if errorMsg}
	<div class="p-4 rounded-lg flex gap-3 text-sm mb-5 bg-danger-bg text-danger-color border border-red-500/15" role="alert">
		<span>{errorMsg}</span>
	</div>
{/if}

{#if loading}
	<div class="bg-bg-card border border-border-color rounded-xl overflow-hidden shadow-md">
		<table class="w-full border-collapse text-left">
			<thead>
				<tr>
					<th class="bg-text-primary/4 p-4 font-bold text-xs text-text-secondary uppercase tracking-wider border-b border-border-color">Tipo</th>
					<th class="bg-text-primary/4 p-4 font-bold text-xs text-text-secondary uppercase tracking-wider border-b border-border-color">Detalle</th>
					<th class="bg-text-primary/4 p-4 font-bold text-xs text-text-secondary uppercase tracking-wider border-b border-border-color">Descuento</th>
					{#if auth.user?.rol === 'Admin'}
						<th class="bg-text-primary/4 p-4 font-bold text-xs text-text-secondary uppercase tracking-wider border-b border-border-color text-right w-[120px]">Acciones</th>
					{/if}
				</tr>
			</thead>
			<tbody>
				{#each Array(3) as _}
					<tr>
						<td class="p-4 border-b border-border-color"><div class="skeleton-loader h-4 w-20 rounded"></div></td>
						<td class="p-4 border-b border-border-color"><div class="skeleton-loader h-4 w-32 rounded"></div></td>
						<td class="p-4 border-b border-border-color"><div class="skeleton-loader h-4 w-24 rounded"></div></td>
						{#if auth.user?.rol === 'Admin'}
							<td class="p-4 border-b border-border-color">
								<div class="flex justify-end gap-2">
									<div class="skeleton-loader w-8 h-8 rounded-lg"></div>
									<div class="skeleton-loader w-8 h-8 rounded-lg"></div>
								</div>
							</td>
						{/if}
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
{:else if promociones.length === 0}
	<div class="bg-bg-card border border-border-color rounded-xl p-16 text-center flex flex-col items-center justify-center shadow-md">
		<h3 class="text-lg font-semibold text-text-primary mb-2">No hay promociones configuradas</h3>
		<p class="text-text-secondary text-sm">Crea tu primera oferta para aplicarla en la caja.</p>
	</div>
{:else}
	<div class="bg-bg-card border border-border-color rounded-xl overflow-hidden shadow-md">
		<table class="w-full border-collapse text-left">
			<thead>
				<tr>
					<th class="bg-text-primary/4 p-4 font-bold text-xs text-text-secondary uppercase tracking-wider border-b border-border-color">Tipo</th>
					<th class="bg-text-primary/4 p-4 font-bold text-xs text-text-secondary uppercase tracking-wider border-b border-border-color">Detalle (Lleva/Paga)</th>
					<th class="bg-text-primary/4 p-4 font-bold text-xs text-text-secondary uppercase tracking-wider border-b border-border-color">Descuento Directo</th>
					{#if auth.user?.rol === 'Admin'}
						<th class="bg-text-primary/4 p-4 font-bold text-xs text-text-secondary uppercase tracking-wider border-b border-border-color text-right w-[120px]">Acciones</th>
					{/if}
				</tr>
			</thead>
			<tbody>
				{#each promociones as promo (promo.id_promocion)}
					<tr class="hover:bg-text-primary/[0.015] transition-colors">
						<td class="p-4 border-b border-border-color font-semibold text-text-primary">
							<span class="bg-text-primary/10 text-text-primary px-2 py-1 rounded text-xs font-bold tracking-wider">
								{promo.tipo}
							</span>
						</td>
						<td class="p-4 border-b border-border-color text-text-secondary">
							{#if promo.tipo === 'NXM'}
								Lleva {promo.lleva}, Paga {promo.paga}
							{:else}
								-
							{/if}
						</td>
						<td class="p-4 border-b border-border-color text-text-secondary">
							{#if promo.tipo === 'DESCUENTO'}
								${promo.descuento}
							{:else}
								-
							{/if}
						</td>
						{#if auth.user?.rol === 'Admin'}
							<td class="p-4 border-b border-border-color">
								<div class="flex justify-end gap-2">
									<button class="p-2 bg-text-primary/3 text-text-secondary rounded-lg border border-border-color hover:text-text-primary transition-all" onclick={() => openEditModal(promo)}>
										<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg>
									</button>
									<button class="p-2 bg-text-primary/3 text-text-secondary rounded-lg border border-border-color hover:text-danger-color hover:bg-danger-bg transition-all" onclick={() => openDeleteModal(promo)}>
										<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/><line x1="10" y1="11" x2="10" y2="17"/><line x1="14" y1="11" x2="14" y2="17"/></svg>
									</button>
								</div>
							</td>
						{/if}
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
{/if}

{#if showModal}
	<div class="fixed inset-0 bg-text-primary/40 backdrop-blur-[4px] flex items-center justify-center z-50 p-5" onclick={() => showModal = false} role="presentation">
		<div class="bg-bg-card border border-border-color rounded-xl w-full max-w-[500px] shadow-lg animate-modal-enter overflow-hidden" onclick={(e) => e.stopPropagation()}>
			<header class="p-5 border-b border-border-color flex justify-between items-center">
				<h2 class="text-lg font-bold text-text-primary">{editingPromo ? 'Editar Promoción' : 'Nueva Promoción'}</h2>
				<button class="p-2 text-text-secondary hover:text-text-primary" onclick={() => showModal = false}>&times;</button>
			</header>
			<form onsubmit={handleSubmit}>
				<div class="p-6">
					{#if formGeneralError}
						<div class="p-4 rounded-lg flex gap-3 text-sm mb-5 bg-danger-bg text-danger-color border border-red-500/15">{formGeneralError}</div>
					{/if}

					<div class="flex flex-col gap-1.5 mb-5">
						<label class="text-[0.85rem] font-semibold text-text-secondary" for="formTipo">Tipo de Oferta</label>
						<select id="formTipo" class="bg-white border border-[rgba(15,30,54,0.15)] rounded-lg px-4 py-3 text-text-primary text-sm outline-none focus:border-accent focus:ring-2 focus:ring-accent/15" bind:value={formTipo} disabled={submitLoading || editingPromo !== null}>
							<option value="NXM">Formato NXM (Ej. 3x2)</option>
							<option value="DESCUENTO">Descuento Directo</option>
						</select>
					</div>

					{#if formTipo === 'NXM'}
						<div class="grid grid-cols-2 gap-4 mb-5">
							<div class="flex flex-col gap-1.5">
								<label class="text-[0.85rem] font-semibold text-text-secondary">Cantidad que LLEVA</label>
								<input type="number" min="1" class="bg-white border border-[rgba(15,30,54,0.15)] rounded-lg px-4 py-3 text-text-primary" bind:value={formLleva} disabled={submitLoading} />
							</div>
							<div class="flex flex-col gap-1.5">
								<label class="text-[0.85rem] font-semibold text-text-secondary">Cantidad que PAGA</label>
								<input type="number" min="1" class="bg-white border border-[rgba(15,30,54,0.15)] rounded-lg px-4 py-3 text-text-primary" bind:value={formPaga} disabled={submitLoading} />
							</div>
						</div>
						{#if errLlevaPaga}
							<span class="text-danger-color text-xs font-medium mt-[-10px] mb-5 block">{errLlevaPaga}</span>
						{/if}
					{:else}
						<div class="flex flex-col gap-1.5 mb-5">
							<label class="text-[0.85rem] font-semibold text-text-secondary">Monto de Descuento ($)</label>
							<input type="number" min="1" class="bg-white border border-[rgba(15,30,54,0.15)] rounded-lg px-4 py-3 text-text-primary" placeholder="Ej: 1500" bind:value={formDescuento} disabled={submitLoading} />
							{#if errDescuento}
								<span class="text-danger-color text-xs font-medium mt-1">{errDescuento}</span>
							{/if}
						</div>
					{/if}
				</div>
				<footer class="p-4 px-6 bg-text-primary/2 border-t border-border-color flex justify-end gap-3">
					<button type="button" class="px-5 py-2.5 rounded-lg border border-border-color text-sm" onclick={() => showModal = false}>Cancelar</button>
					<button type="submit" class="px-5 py-2.5 rounded-lg bg-gradient-to-r from-accent-light to-accent text-white text-sm" disabled={submitLoading}>Guardar</button>
				</footer>
			</form>
		</div>
	</div>
{/if}

{#if showDeleteModal && promoToDelete}
	<div class="fixed inset-0 bg-text-primary/40 backdrop-blur-[4px] flex items-center justify-center z-50 p-5" onclick={() => showDeleteModal = false}>
		<div class="bg-bg-card border border-border-color rounded-xl w-full max-w-[500px] shadow-lg overflow-hidden" onclick={(e) => e.stopPropagation()}>
			<header class="p-5 border-b border-border-color flex justify-between items-center">
				<h2 class="text-lg font-bold text-danger-color">Confirmar Eliminación</h2>
			</header>
			<div class="p-6 text-text-primary">
				<p>¿Estás seguro de eliminar la promoción 
					<strong>
						{#if promoToDelete.tipo === 'NXM'}
							{promoToDelete.lleva}x{promoToDelete.paga}
						{:else}
							Descuento ${promoToDelete.descuento}
						{/if}
					</strong>?
				</p>
			</div>
			<footer class="p-4 px-6 bg-text-primary/2 border-t flex justify-end gap-3">
				<button type="button" class="px-5 py-2.5 rounded-lg border border-border-color" onclick={() => showDeleteModal = false}>Cancelar</button>
				<button type="button" class="px-5 py-2.5 rounded-lg bg-danger-bg text-danger-color border border-red-500/15" onclick={confirmDelete}>Eliminar</button>
			</footer>
		</div>
	</div>
{/if}