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
		
		// Añadir puntos al cuerpo si tiene suficiente longitud
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
			toast.show('Acción denegada, requiere rol admin.', 'error');
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
			toast.show('Acción denegada, requiere rol admin.', 'error');
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
			toast.show('Acción denegada, requiere rol admin', 'error');
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
			errRut = 'El RUT ingresado no es válido.';
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
<div class="page-header">
	<div>
		<h1 class="page-title">Gestión de Clientes</h1>
		<p class="page-subtitle">Admin info clientes</p>
	</div>
	{#if auth.user?.rol === 'Admin'}
		<button class="btn btn-primary" onclick={openCreateModal} id="btn-nuevo-cliente">
			<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
				<line x1="12" y1="5" x2="12" y2="19"/>
				<line x1="5" y1="12" x2="19" y2="12"/>
			</svg>
			<span>Nuevo Cliente</span>
		</button>
	{/if}
</div>

<!-- Search Panel -->
<div class="card search-card">
	<form onsubmit={handleSearch} class="search-form">
		<div class="search-input-group">
			<div class="select-wrapper">
				<select class="form-input search-select" bind:value={searchType} aria-label="Tipo de búsqueda">
					<option value="nombre">Por Nombre</option>
					<option value="rut">Por RUT</option>
				</select>
			</div>
			<input
				type="text"
				class="form-input search-input"
				placeholder={searchType === 'nombre' ? 'Ej: Juan Pérez...' : 'Ej: 12345678-9...'}
				bind:value={searchQuery}
				aria-label="Término de búsqueda"
			/>
		</div>
		<div class="search-actions">
			<button type="submit" class="btn btn-primary btn-search" id="btn-search-submit">
				<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
					<circle cx="11" cy="11" r="8"/>
					<line x1="21" y1="21" x2="16.65" y2="16.65"/>
				</svg>
				<span>Buscar</span>
			</button>
			{#if searchQuery}
				<button type="button" class="btn btn-secondary" onclick={clearSearch} id="btn-search-clear">
					Limpiar
				</button>
			{/if}
		</div>
	</form>
</div>

<!-- Table View -->
{#if errorMsg}
	<div class="alert alert-error" role="alert">
		<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
			<circle cx="12" cy="12" r="10"/>
			<line x1="12" y1="8" x2="12" y2="12"/>
			<line x1="12" y1="16" x2="12.01" y2="16"/>
		</svg>
		<span>{errorMsg}</span>
	</div>
{/if}

{#if loading}
	<div class="table-container">
		<table class="data-table">
			<thead>
				<tr>
					<th>Nombre</th>
					<th>RUT</th>
					<th>Teléfono</th>
					{#if auth.user?.rol === 'Admin'}
						<th class="actions-header">Acciones</th>
					{/if}
				</tr>
			</thead>
			<tbody>
				{#each Array(5) as _, i}
					<tr>
						<td><div class="skeleton skeleton-text" style="width: 150px;"></div></td>
						<td><div class="skeleton skeleton-text" style="width: 100px;"></div></td>
						<td><div class="skeleton skeleton-text" style="width: 120px;"></div></td>
						{#if auth.user?.rol === 'Admin'}
							<td>
								<div style="display: flex; gap: 8px;">
									<div class="skeleton skeleton-btn" style="width: 32px; height: 32px; border-radius: 6px;"></div>
									<div class="skeleton skeleton-btn" style="width: 32px; height: 32px; border-radius: 6px;"></div>
								</div>
							</td>
						{/if}
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
{:else if clientes.length === 0}
	<div class="card empty-card">
		<div class="empty-icon">
			<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
				<path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
				<circle cx="9" cy="7" r="4"/>
				<line x1="17" y1="8" x2="22" y2="13"/>
				<line x1="22" y1="8" x2="17" y2="13"/>
			</svg>
		</div>
		<h3>No se encontraron clientes</h3>
		<p>Modifica el criterio de búsqueda o agrega un nuevo cliente al sistema.</p>
	</div>
{:else}
	<div class="table-container">
		<table class="data-table">
			<thead>
				<tr>
					<th>Nombre</th>
					<th>RUT</th>
					<th>Teléfono</th>
					{#if auth.user?.rol === 'Admin'}
						<th class="actions-header">Acciones</th>
					{/if}
				</tr>
			</thead>
			<tbody>
				{#each clientes as client (client.id_cliente)}
					<tr>
						<td class="client-name">{client.nombre}</td>
						<td class="client-rut">{formatRutInput(client.rut)}</td>
						<td class="client-tel">{client.telefono}</td>
						{#if auth.user?.rol === 'Admin'}
							<td>
								<div class="actions-cell">
									<button
										class="btn-icon"
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
										class="btn-icon btn-delete-icon"
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
	<div class="modal-backdrop" onclick={() => showModal = false} role="presentation">
		<div class="modal-content" onclick={(e) => e.stopPropagation()} role="dialog" aria-modal="true" aria-labelledby="modal-title">
			<header class="modal-header">
				<h2 id="modal-title">{editingCliente ? 'Editar Cliente' : 'Nuevo Cliente'}</h2>
				<button class="btn-icon" onclick={() => showModal = false} aria-label="Cerrar modal">&times;</button>
			</header>
			<form onsubmit={handleSubmit}>
				<div class="modal-body">
					{#if formGeneralError}
						<div class="alert alert-error" role="alert">
							<span>{formGeneralError}</span>
						</div>
					{/if}

					<div class="form-group">
						<label class="form-label" for="formNombre">Nombre Completo</label>
						<input
							type="text"
							id="formNombre"
							class="form-input"
							placeholder="Ej: Juan Pérez Gómez"
							bind:value={formNombre}
							disabled={submitLoading}
							required
						/>
						{#if errNombre}
							<span class="form-error">{errNombre}</span>
						{/if}
					</div>

					<div class="form-group">
						<label class="form-label" for="formRut">RUT</label>
						<input
							type="text"
							id="formRut"
							class="form-input"
							placeholder="Ej: 12.345.678-K"
							value={formRut}
							oninput={handleRutInput}
							disabled={submitLoading || editingCliente !== null}
							required
						/>
						{#if errRut}
							<span class="form-error">{errRut}</span>
						{/if}
					</div>

					<div class="form-group">
						<label class="form-label" for="formTelefono">Número Telefónico</label>
						<input
							type="text"
							id="formTelefono"
							class="form-input"
							placeholder="Ej: +56912345678"
							bind:value={formTelefono}
							disabled={submitLoading}
							required
						/>
						{#if errTelefono}
							<span class="form-error">{errTelefono}</span>
						{/if}
					</div>
				</div>
				<footer class="modal-footer">
					<button type="button" class="btn btn-secondary" onclick={() => showModal = false} disabled={submitLoading}>
						Cancelar
					</button>
					<button type="submit" class="btn btn-primary" disabled={submitLoading} id="btn-modal-save">
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
	<div class="modal-backdrop" onclick={() => showDeleteModal = false} role="presentation">
		<div class="modal-content" onclick={(e) => e.stopPropagation()} role="dialog" aria-modal="true" aria-labelledby="delete-modal-title">
			<header class="modal-header">
				<h2 id="delete-modal-title" style="color: var(--danger-color);">Confirmar Eliminación</h2>
				<button class="btn-icon" onclick={() => showDeleteModal = false} aria-label="Cerrar modal">&times;</button>
			</header>
			<div class="modal-body">
				<p>¿eliminar al cliente <strong>{clienteToDelete.nombre}</strong> de forma permanente?</p>
				<p style="margin-top: 12px; font-size: 0.85rem; color: var(--text-muted);">No se puede deshacer, ademas se pierde informacion.</p>
			</div>
			<footer class="modal-footer">
				<button type="button" class="btn btn-secondary" onclick={() => showDeleteModal = false}>
					Cancelar
				</button>
				<button type="button" class="btn btn-danger" onclick={confirmDelete} id="btn-modal-confirm-delete">
					Eliminar Permanentemente
				</button>
			</footer>
		</div>
	</div>
{/if}

<style>
	.search-card {
		margin-bottom: 24px;
		background: var(--bg-card);
	}

	.search-form {
		display: flex;
		flex-wrap: wrap;
		gap: 16px;
		align-items: center;
	}

	.search-input-group {
		display: flex;
		flex: 1;
		min-width: 280px;
		border: 1px solid rgba(15, 30, 54, 0.15);
		border-radius: 8px;
		overflow: hidden;
		background: #ffffff;
	}

	.search-input-group:focus-within {
		border-color: var(--accent-color);
		box-shadow: 0 0 0 2px rgba(217, 119, 6, 0.15);
	}

	.select-wrapper {
		border-right: 1px solid var(--border-color);
	}

	.search-select {
		border: none;
		background: transparent;
		padding-right: 8px;
		cursor: pointer;
		height: 100%;
		border-radius: 0;
		outline: none !important;
		box-shadow: none !important;
	}

	.search-input {
		flex: 1;
		border: none;
		background: transparent;
		border-radius: 0;
		outline: none !important;
		box-shadow: none !important;
	}

	.search-actions {
		display: flex;
		gap: 10px;
	}

	.btn-search {
		padding-left: 16px;
		padding-right: 16px;
	}

	.empty-card {
		padding: 60px 20px;
		text-align: center;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		background: var(--bg-card);
		border: 1px solid var(--border-color);
	}

	.empty-icon {
		width: 80px;
		height: 80px;
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--text-muted);
		background: rgba(15, 30, 54, 0.03);
		border-radius: 50%;
		margin-bottom: 20px;
	}

	.empty-card h3 {
		font-size: 1.2rem;
		font-weight: 600;
		margin-bottom: 8px;
	}

	.empty-card p {
		color: var(--text-secondary);
		font-size: 0.9rem;
		max-width: 400px;
	}

	.client-name {
		font-weight: 600;
		color: var(--text-primary);
	}

	.client-rut {
		font-family: monospace;
		font-size: 0.9rem;
		color: var(--text-secondary);
	}

	.client-tel {
		color: var(--text-secondary);
	}

	.actions-header {
		text-align: right;
		width: 120px;
	}

	.actions-cell {
		display: flex;
		justify-content: flex-end;
		gap: 8px;
	}

	.btn-delete-icon:hover {
		color: var(--danger-color);
		background: var(--danger-bg);
		border-color: rgba(239, 68, 68, 0.2);
	}
</style>
