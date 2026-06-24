<script lang="ts">
	import { apiClientes } from '$lib/api';
	import { auth } from '$lib/authStore.svelte';
	import { toast } from '$lib/toastStore.svelte';
	import { onMount } from 'svelte';

	interface Cliente {
		id_cliente: string;
		nombre: string;
		rut: string;
		telefono: string;
		created_at: string;
	}

	let clientes = $state<Cliente[]>([]);
	let loading = $state(true);
	let errorMsg = $state('');

	// Búsqueda
	let searchQuery = $state('');
	let searchType = $state<'nombre' | 'rut'>('nombre');

	// Modales
	let showModal = $state(false);
	let showDeleteModal = $state(false);
	let editingCliente = $state<Cliente | null>(null);
	let clienteToDelete = $state<Cliente | null>(null);

	// Campos de Formulario
	let formNombre = $state('');
	let formRut = $state('');
	let formTelefono = $state('');

	// Errores de Formulario
	let errNombre = $state('');
	let errRut = $state('');
	let errTelefono = $state('');
	let formGeneralError = $state('');
	let submitLoading = $state(false);

	onMount(loadClientes);

	async function loadClientes() {
		loading = true;
		errorMsg = '';
		try {
			const res = await apiClientes.getAll();
			clientes = Array.isArray(res) ? res : [];
		} catch (err: any) {
			errorMsg = err.message || 'Error al conectar con la base de datos de clientes.';
		} finally {
			loading = false;
		}
	}

	// Validadores de frontend para RUT chileno y teléfono móvil
	function validarRut(rutOriginal: string): boolean {
		const rutLimpio = rutOriginal.replace(/\./g, '').replace(/-/g, '').trim().toUpperCase();
		if (rutLimpio.length < 8) return false;
		const cuerpo = rutLimpio.slice(0, -1);
		const dvIngresado = rutLimpio.slice(-1);

		let suma = 0;
		let multiplicador = 2;
		for (let i = cuerpo.length - 1; i >= 0; i--) {
			const digito = parseInt(cuerpo[i], 10);
			if (isNaN(digito)) return false;
			suma += digito * multiplicador;
			multiplicador++;
			if (multiplicador > 7) multiplicador = 2;
		}

		const resto = suma % 11;
		const resultado = 11 - resto;
		let dvCalculado = '';
		if (resultado === 11) dvCalculado = '0';
		else if (resultado === 10) dvCalculado = 'K';
		else dvCalculado = resultado.toString();

		return dvCalculado === dvIngresado;
	}

	function validarTelefono(tel: string): boolean {
		const patron = /^(\+56)?[\s.-]?9[\s.-]?\d{4}[\s.-]?\d{4}$/;
		return patron.test(tel);
	}

	function formatRutInput(val: string) {
		let clean = val.replace(/[^0-9kK]/g, '');
		if (clean.length === 0) return '';
		
		const dv = clean.slice(-1);
		let cuerpo = clean.slice(0, -1);
		
		if (cuerpo.length > 0) {
			let formattedCuerpo = '';
			let j = 0;
			for (let i = cuerpo.length - 1; i >= 0; i--) {
				formattedCuerpo = cuerpo[i] + formattedCuerpo;
				j++;
				if (j === 3 && i > 0) {
					formattedCuerpo = '.' + formattedCuerpo;
					j = 0;
				}
			}
			return `${formattedCuerpo}-${dv}`;
		}
		return dv;
	}

	function handleRutInput(e: Event) {
		const target = e.target as HTMLInputElement;
		const formatted = formatRutInput(target.value);
		formRut = formatted;
		target.value = formatted;
	}

	// Manejo de Búsqueda
	async function handleSearch(e: Event) {
		e.preventDefault();
		if (!searchQuery.trim()) {
			loadClientes();
			return;
		}

		loading = true;
		errorMsg = '';
		const query = searchQuery.trim();

		try {
			if (searchType === 'rut') {
				const cleanedRut = query.replace(/\./g, '');
				const res = await apiClientes.getByRut(cleanedRut);
				if (res && res.id_cliente) {
					clientes = [res];
				} else {
					clientes = [];
				}
			} else {
				const res = await apiClientes.searchByNombre(query);
				clientes = Array.isArray(res.resultados) ? res.resultados : [];
			}
		} catch (err: any) {
			if (err.status === 404) {
				clientes = [];
			} else {
				errorMsg = err.message || 'Error al buscar clientes.';
			}
		} finally {
			loading = false;
		}
	}

	function clearSearch() {
		searchQuery = '';
		loadClientes();
	}

	// Modales
	function openCreateModal() {
		if (auth.user?.rol !== 'Admin') {
			toast.show('Acción denegada, requiere rol Administrador.', 'error');
			return;
		}
		editingCliente = null;
		formNombre = '';
		formRut = '';
		formTelefono = '';
		clearErrors();
		showModal = true;
	}

	function openEditModal(cliente: Cliente) {
		if (auth.user?.rol !== 'Admin') {
			toast.show('Acción denegada, requiere rol Administrador.', 'error');
			return;
		}
		editingCliente = cliente;
		formNombre = cliente.nombre;
		formRut = formatRutInput(cliente.rut);
		formTelefono = cliente.telefono;
		clearErrors();
		showModal = true;
	}

	function openDeleteModal(cliente: Cliente) {
		if (auth.user?.rol !== 'Admin') {
			toast.show('Acción denegada, requiere rol Administrador.', 'error');
			return;
		}
		clienteToDelete = cliente;
		showDeleteModal = true;
	}

	function clearErrors() {
		errNombre = '';
		errRut = '';
		errTelefono = '';
		formGeneralError = '';
	}

	// Submit Formulario
	async function handleSubmit(e: Event) {
		e.preventDefault();
		clearErrors();
		let isValid = true;

		// Nombre
		if (!formNombre.trim()) {
			errNombre = 'El nombre es obligatorio.';
			isValid = false;
		}

		// RUT
		const rawRut = formRut.replace(/\./g, '');
		if (!formRut) {
			errRut = 'El RUT es obligatorio.';
			isValid = false;
		} else if (!validarRut(rawRut)) {
			errRut = 'El RUT ingresado no es válido para Chile.';
			isValid = false;
		}

		// Teléfono
		if (!formTelefono.trim()) {
			errTelefono = 'El teléfono es obligatorio.';
			isValid = false;
		} else if (!validarTelefono(formTelefono.trim())) {
			errTelefono = 'Formato no válido (ej: +56 9 1234 5678 o 912345678).';
			isValid = false;
		}

		if (!isValid) return;

		submitLoading = true;
		const payload = {
			nombre: formNombre.trim(),
			rut: rawRut.trim(),
			telefono: formTelefono.replace(/\s+/g, '')
		};

		try {
			if (editingCliente) {
				await apiClientes.update(editingCliente.id_cliente, payload);
				toast.show('Cliente modificado exitosamente.', 'success');
			} else {
				await apiClientes.create(payload);
				toast.show('Cliente creado exitosamente.', 'success');
			}
			showModal = false;
			loadClientes();
		} catch (err: any) {
			if (err.errors) {
				if (err.errors.Nombre || err.errors.nombre) errNombre = err.errors.Nombre || err.errors.nombre;
				if (err.errors.Rut || err.errors.rut) errRut = err.errors.Rut || err.errors.rut;
				if (err.errors.Telefono || err.errors.telefono) errTelefono = err.errors.Telefono || err.errors.telefono;
			} else {
				formGeneralError = err.message || 'Error al guardar el cliente.';
			}
		} finally {
			submitLoading = false;
		}
	}

	// Confirmar Borrado
	async function confirmDelete() {
		if (!clienteToDelete) return;
		try {
			await apiClientes.delete(clienteToDelete.id_cliente);
			toast.show('Cliente eliminado con éxito.', 'success');
			showDeleteModal = false;
			clienteToDelete = null;
			loadClientes();
		} catch (err: any) {
			toast.show(err.message || 'No se pudo eliminar el cliente.', 'error');
		}
	}
</script>

<svelte:head>
	<title>Clientes - GPSproject</title>
	<meta name="description" content="gestión de salgo de clientes" />
</svelte:head>

<!-- Page Header -->
<div class="flex justify-between items-center mb-8 flex-wrap gap-4">
	<div>
		<h1 class="text-3xl font-bold tracking-tight text-text-primary">Gestión de Clientes</h1>
		<p class="text-text-secondary text-sm mt-1">Visualiza, busca y administra la información de clientes registrados en el sistema.</p>
	</div>
	{#if auth.user?.rol === 'Admin'}
		<button class="inline-flex items-center justify-center gap-2 font-semibold text-sm px-5 py-2.5 rounded-lg cursor-pointer bg-gradient-to-r from-accent-light to-accent text-white hover:shadow-glow hover:-translate-y-[1px] disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200" onclick={openCreateModal} id="btn-nuevo-cliente">
			<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
				<line x1="12" y1="5" x2="12" y2="19"/>
				<line x1="5" y1="12" x2="19" y2="12"/>
			</svg>
			<span>Nuevo Cliente</span>
		</button>
	{/if}
</div>

<!-- Search Panel -->
<div class="bg-bg-card border border-border-color rounded-xl p-6 shadow-md hover:border-border-color-hover hover:shadow-lg transition-all duration-300 mb-6">
	<form onsubmit={handleSearch} class="flex flex-wrap gap-4 items-center">
		<div class="flex flex-1 min-w-[280px] border border-[rgba(15,30,54,0.15)] rounded-lg overflow-hidden bg-white focus-within:border-accent focus-within:ring-2 focus-within:ring-accent/15">
			<div class="border-r border-border-color">
				<select class="border-none bg-transparent pr-2 pl-4 py-3 cursor-pointer h-full outline-none text-sm text-text-primary" bind:value={searchType} aria-label="Tipo de búsqueda">
					<option value="nombre">Por Nombre</option>
					<option value="rut">Por RUT</option>
				</select>
			</div>
			<input
				type="text"
				class="flex-1 border-none bg-transparent px-4 py-3 outline-none text-sm text-text-primary"
				placeholder={searchType === 'nombre' ? 'Ej: Juan Pérez...' : 'Ej: 12345678-9...'}
				bind:value={searchQuery}
				aria-label="Término de búsqueda"
			/>
		</div>
		<div class="flex gap-2.5">
			<button type="submit" class="inline-flex items-center justify-center gap-2 font-semibold text-sm px-5 py-2.5 rounded-lg cursor-pointer bg-gradient-to-r from-accent-light to-accent text-white hover:shadow-glow hover:-translate-y-[1px] disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200" id="btn-search-submit">
				<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
					<circle cx="11" cy="11" r="8"/>
					<line x1="21" y1="21" x2="16.65" y2="16.65"/>
				</svg>
				<span>Buscar</span>
			</button>
			{#if searchQuery}
				<button type="button" class="inline-flex items-center justify-center gap-2 font-semibold text-sm px-5 py-2.5 rounded-lg cursor-pointer bg-bg-secondary text-text-primary border border-border-color hover:bg-text-primary/5 transition-all duration-200" onclick={clearSearch} id="btn-search-clear">
					Limpiar
				</button>
			{/if}
		</div>
	</form>
</div>

<!-- Table View -->
{#if errorMsg}
	<div class="p-4 rounded-lg flex gap-3 text-sm mb-5 bg-danger-bg text-danger-color border border-red-500/15" role="alert">
		<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
			<circle cx="12" cy="12" r="10"/>
			<line x1="12" y1="8" x2="12" y2="12"/>
			<line x1="12" y1="16" x2="12.01" y2="16"/>
		</svg>
		<span>{errorMsg}</span>
	</div>
{/if}

{#if loading}
	<div class="bg-bg-card border border-border-color rounded-xl overflow-hidden shadow-md">
		<table class="w-full border-collapse text-left">
			<thead>
				<tr>
					<th class="bg-text-primary/4 p-4 font-bold text-xs text-text-secondary uppercase tracking-wider border-b border-border-color">Nombre</th>
					<th class="bg-text-primary/4 p-4 font-bold text-xs text-text-secondary uppercase tracking-wider border-b border-border-color">RUT</th>
					<th class="bg-text-primary/4 p-4 font-bold text-xs text-text-secondary uppercase tracking-wider border-b border-border-color">Teléfono</th>
					{#if auth.user?.rol === 'Admin'}
						<th class="bg-text-primary/4 p-4 font-bold text-xs text-text-secondary uppercase tracking-wider border-b border-border-color text-right w-[120px]">Acciones</th>
					{/if}
				</tr>
			</thead>
			<tbody>
				{#each Array(5) as _, i}
					<tr>
						<td class="p-4 border-b border-border-color"><div class="skeleton-loader h-4 w-40 rounded"></div></td>
						<td class="p-4 border-b border-border-color"><div class="skeleton-loader h-4 w-28 rounded"></div></td>
						<td class="p-4 border-b border-border-color"><div class="skeleton-loader h-4 w-32 rounded"></div></td>
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
{:else if clientes.length === 0}
	<div class="bg-bg-card border border-border-color rounded-xl p-16 text-center flex flex-col items-center justify-center shadow-md">
		<div class="w-20 h-20 flex items-center justify-center text-text-muted bg-text-primary/3 rounded-full mb-5">
			<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
				<path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
				<circle cx="9" cy="7" r="4"/>
				<line x1="17" y1="8" x2="22" y2="13"/>
				<line x1="22" y1="8" x2="17" y2="13"/>
			</svg>
		</div>
		<h3 class="text-lg font-semibold text-text-primary mb-2">No se encontraron clientes</h3>
		<p class="text-text-secondary text-sm max-w-[400px]">Modifica el criterio de búsqueda o agrega un nuevo cliente al sistema.</p>
	</div>
{:else}
	<div class="bg-bg-card border border-border-color rounded-xl overflow-hidden shadow-md">
		<table class="w-full border-collapse text-left">
			<thead>
				<tr>
					<th class="bg-text-primary/4 p-4 font-bold text-xs text-text-secondary uppercase tracking-wider border-b border-border-color">Nombre</th>
					<th class="bg-text-primary/4 p-4 font-bold text-xs text-text-secondary uppercase tracking-wider border-b border-border-color">RUT</th>
					<th class="bg-text-primary/4 p-4 font-bold text-xs text-text-secondary uppercase tracking-wider border-b border-border-color">Teléfono</th>
					{#if auth.user?.rol === 'Admin'}
						<th class="bg-text-primary/4 p-4 font-bold text-xs text-text-secondary uppercase tracking-wider border-b border-border-color text-right w-[120px]">Acciones</th>
					{/if}
				</tr>
			</thead>
			<tbody>
				{#each clientes as client (client.id_cliente)}
					<tr class="hover:bg-text-primary/[0.015] transition-colors">
						<td class="p-4 border-b border-border-color font-semibold text-text-primary">{client.nombre}</td>
						<td class="p-4 border-b border-border-color font-mono text-sm text-text-secondary">{formatRutInput(client.rut)}</td>
						<td class="p-4 border-b border-border-color text-text-secondary">{client.telefono}</td>
						{#if auth.user?.rol === 'Admin'}
							<td class="p-4 border-b border-border-color">
								<div class="flex justify-end gap-2">
									<button
										class="p-2 bg-text-primary/3 text-text-secondary rounded-lg border border-border-color cursor-pointer inline-flex items-center justify-center transition-all duration-200 hover:text-text-primary hover:bg-text-primary/7 hover:border-border-color-hover"
										onclick={() => openEditModal(client)}
										title="Editar cliente"
										aria-label="Editar"
									>
										<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
											<path d="M12 20h9"/>
											<path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/>
										</svg>
									</button>
									<button
										class="p-2 bg-text-primary/3 text-text-secondary rounded-lg border border-border-color cursor-pointer inline-flex items-center justify-center transition-all duration-200 hover:text-danger-color hover:bg-danger-bg hover:border-red-500/20"
										onclick={() => openDeleteModal(client)}
										title="Eliminar cliente"
										aria-label="Eliminar"
									>
										<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
											<polyline points="3 6 5 6 21 6"/>
											<path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
											<line x1="10" y1="11" x2="10" y2="17"/>
											<line x1="14" y1="11" x2="14" y2="17"/>
										</svg>
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

<!-- Form Dialog Modal (Create / Edit) -->
{#if showModal}
	<div class="fixed inset-0 bg-text-primary/40 backdrop-blur-[4px] flex items-center justify-center z-50 p-5" onclick={() => showModal = false} role="presentation">
		<div class="bg-bg-card border border-border-color rounded-xl w-full max-w-[500px] shadow-lg animate-modal-enter overflow-hidden" onclick={(e) => e.stopPropagation()} role="dialog" aria-modal="true" aria-labelledby="modal-title">
			<header class="p-5 border-b border-border-color flex justify-between items-center">
				<h2 id="modal-title" class="text-lg font-bold text-text-primary">{editingCliente ? 'Editar Cliente' : 'Nuevo Cliente'}</h2>
				<button class="p-2 bg-text-primary/3 text-text-secondary rounded-lg border border-border-color cursor-pointer inline-flex items-center justify-center transition-all duration-200 hover:text-text-primary hover:bg-text-primary/7 hover:border-border-color-hover" onclick={() => showModal = false} aria-label="Cerrar modal">&times;</button>
			</header>
			<form onsubmit={handleSubmit}>
				<div class="p-6">
					{#if formGeneralError}
						<div class="p-4 rounded-lg flex gap-3 text-sm mb-5 bg-danger-bg text-danger-color border border-red-500/15" role="alert">
							<span>{formGeneralError}</span>
						</div>
					{/if}

					<div class="flex flex-col gap-1.5 mb-5">
						<label class="text-[0.85rem] font-semibold text-text-secondary" for="formNombre">Nombre Completo</label>
						<input
							type="text"
							id="formNombre"
							class="bg-white border border-[rgba(15,30,54,0.15)] rounded-lg px-4 py-3 text-text-primary text-sm outline-none transition-all duration-200 focus:border-accent focus:ring-2 focus:ring-accent/15 placeholder:text-text-muted/80 disabled:opacity-50 disabled:cursor-not-allowed"
							placeholder="Ej: Juan Pérez Gómez"
							bind:value={formNombre}
							disabled={submitLoading}
							required
						/>
						{#if errNombre}
							<span class="text-danger-color text-xs font-medium mt-1">{errNombre}</span>
						{/if}
					</div>

					<div class="flex flex-col gap-1.5 mb-5">
						<label class="text-[0.85rem] font-semibold text-text-secondary" for="formRut">RUT</label>
						<input
							type="text"
							id="formRut"
							class="bg-white border border-[rgba(15,30,54,0.15)] rounded-lg px-4 py-3 text-text-primary text-sm outline-none transition-all duration-200 focus:border-accent focus:ring-2 focus:ring-accent/15 placeholder:text-text-muted/80 disabled:opacity-50 disabled:cursor-not-allowed"
							placeholder="Ej: 12.345.678-K"
							value={formRut}
							oninput={handleRutInput}
							disabled={submitLoading || editingCliente !== null}
							required
						/>
						{#if errRut}
							<span class="text-danger-color text-xs font-medium mt-1">{errRut}</span>
						{/if}
					</div>

					<div class="flex flex-col gap-1.5 mb-5">
						<label class="text-[0.85rem] font-semibold text-text-secondary" for="formTelefono">Número Telefónico</label>
						<input
							type="text"
							id="formTelefono"
							class="bg-white border border-[rgba(15,30,54,0.15)] rounded-lg px-4 py-3 text-text-primary text-sm outline-none transition-all duration-200 focus:border-accent focus:ring-2 focus:ring-accent/15 placeholder:text-text-muted/80 disabled:opacity-50 disabled:cursor-not-allowed"
							placeholder="Ej: +56912345678"
							bind:value={formTelefono}
							disabled={submitLoading}
							required
						/>
						{#if errTelefono}
							<span class="text-danger-color text-xs font-medium mt-1">{errTelefono}</span>
						{/if}
					</div>
				</div>
				<footer class="p-4 px-6 bg-text-primary/2 border-t border-border-color flex justify-end gap-3">
					<button type="button" class="inline-flex items-center justify-center gap-2 font-semibold text-sm px-5 py-2.5 rounded-lg cursor-pointer bg-bg-secondary text-text-primary border border-border-color hover:bg-text-primary/5 transition-all duration-200" onclick={() => showModal = false} disabled={submitLoading}>
						Cancelar
					</button>
					<button type="submit" class="inline-flex items-center justify-center gap-2 font-semibold text-sm px-5 py-2.5 rounded-lg cursor-pointer bg-gradient-to-r from-accent-light to-accent text-white hover:shadow-glow hover:-translate-y-[1px] disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200" disabled={submitLoading} id="btn-modal-save">
						{#if submitLoading}
							Procesando...
						{:else}
							Guardar Cambios
						{/if}
					</button>
				</footer>
			</form>
		</div>
	</div>
{/if}

<!-- Delete Confirmation Modal -->
{#if showDeleteModal && clienteToDelete}
	<div class="fixed inset-0 bg-text-primary/40 backdrop-blur-[4px] flex items-center justify-center z-50 p-5" onclick={() => showDeleteModal = false} role="presentation">
		<div class="bg-bg-card border border-border-color rounded-xl w-full max-w-[500px] shadow-lg animate-modal-enter overflow-hidden" onclick={(e) => e.stopPropagation()} role="dialog" aria-modal="true" aria-labelledby="delete-modal-title">
			<header class="p-5 border-b border-border-color flex justify-between items-center">
				<h2 id="delete-modal-title" class="text-lg font-bold text-danger-color">Confirmar Eliminación</h2>
				<button class="p-2 bg-text-primary/3 text-text-secondary rounded-lg border border-border-color cursor-pointer inline-flex items-center justify-center transition-all duration-200 hover:text-text-primary hover:bg-text-primary/7 hover:border-border-color-hover" onclick={() => showDeleteModal = false} aria-label="Cerrar modal">&times;</button>
			</header>
			<div class="p-6 text-text-primary">
				<p>¿Estás seguro de que deseas eliminar al cliente <strong>{clienteToDelete.nombre}</strong> (RUT: {formatRutInput(clienteToDelete.rut)}) de forma permanente?</p>
				<p class="mt-3 text-xs text-text-muted">Esta acción no se puede deshacer y puede afectar a los reportes de ventas vinculados.</p>
			</div>
			<footer class="p-4 px-6 bg-text-primary/2 border-t border-border-color flex justify-end gap-3">
				<button type="button" class="inline-flex items-center justify-center gap-2 font-semibold text-sm px-5 py-2.5 rounded-lg cursor-pointer bg-bg-secondary text-text-primary border border-border-color hover:bg-text-primary/5 transition-all duration-200" onclick={() => showDeleteModal = false}>
					Cancelar
				</button>
				<button type="button" class="inline-flex items-center justify-center gap-2 font-semibold text-sm px-5 py-2.5 rounded-lg cursor-pointer bg-danger-bg text-danger-color border border-red-500/15 hover:bg-danger-color hover:text-white transition-all duration-200" onclick={confirmDelete} id="btn-modal-confirm-delete">
					Eliminar Permanentemente
				</button>
			</footer>
		</div>
	</div>
{/if}
