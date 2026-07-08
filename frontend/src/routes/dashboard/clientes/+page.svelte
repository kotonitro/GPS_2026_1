<script lang="ts">
	import { apiClientes, apiFetch } from '$lib/api';
	import { auth } from '$lib/authStore.svelte';
	import { toast } from '$lib/toastStore.svelte';
	import { onMount } from 'svelte';
	import { Edit2, Trash2, X, Coins } from '@lucide/svelte';

	interface Cliente {
		id_cliente: string;
		nombre: string;
		rut: string;
		telefono: string;
		fiado_actual: number;
		fiado_maximo: number;
		ultima_compra: string | null;
		created_at: string;
	}

	let clientes = $state<Cliente[]>([]);
	let loading = $state(true);
	let errorMsg = $state('');

	// Búsqueda Cliente-side
	let searchQuery = $state('');
    let selectedEstado = $state('Todos');

	// Abono Modal 
	let showAbonoModal = $state(false);
	let selectedAbonoCliente = $state<Cliente | null>(null);
	let abonoMonto = $state<number | ''>('');
	let abonoLoading = $state(false);
	let abonoError = $state('');

	let filteredClientes = $derived.by(() => {
		let result = clientes;
        
        if (selectedEstado === 'Con deuda') {
            result = result.filter(c => (c.fiado_actual || 0) > 0 && (c.fiado_actual || 0) < (c.fiado_maximo || 20000));
        } else if (selectedEstado === 'Al límite') {
            result = result.filter(c => (c.fiado_actual || 0) >= (c.fiado_maximo || 20000));
        }

		if (searchQuery.trim()) {
			const query = searchQuery.toLowerCase().trim();
			const cleanQuery = query.replace(/[^0-9kK]/g, '');

			result = result.filter((c) => {
				const matchNombre = c.nombre.toLowerCase().includes(query);
				const matchRut = cleanQuery ? c.rut.toLowerCase().replace(/[^0-9kK]/g, '').includes(cleanQuery) : false;
				return matchNombre || matchRut;
			});
		}
        
        return result;
	});

	// Estadísticas
	let fiadoTotalPendiente = $derived(clientes.reduce((acc, c) => acc + (c.fiado_actual || 0), 0));
	let clientesConSaldo = $derived(clientes.filter((c) => (c.fiado_actual || 0) > 0).length);
	let clientesAlLimite = $derived(
		clientes.filter((c) => (c.fiado_actual || 0) >= (c.fiado_maximo || 20000)).length
	);

	// Modales
	let showModal = $state(false);
	let showDeleteModal = $state(false);
	let editingCliente = $state<Cliente | null>(null);
	let clienteToDelete = $state<Cliente | null>(null);

	// Campos de Formulario
	let formNombre = $state('');
	let formRut = $state('');
	let formTelefono = $state('');
	let formFiadoActual = $state<number | ''>('');
	let formFiadoMaximo = $state<number>(20000);

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

	function getInitials(name: string) {
		const parts = name.trim().split(/\s+/);
		if (parts.length === 1) return parts[0].substring(0, 2).toUpperCase();
		return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
	}

	function formatDate(dateString: string | null) {
		if (!dateString) return 'Sin compras';
		const d = new Date(dateString);
		const day = String(d.getDate()).padStart(2, '0');
		const month = String(d.getMonth() + 1).padStart(2, '0');
		return `${day}/${month}/${d.getFullYear()}`;
	}

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

	function clearSearch() {
		searchQuery = '';
	}

	function formatNombreInput(val: string) {
		return val
			.toLowerCase()
			.replace(/\s+/g, ' ')
			.split(' ')
			.map(word => word.charAt(0).toUpperCase() + word.slice(1))
			.join(' ');
	}

	function handleNombreBlur() {
		formNombre = formatNombreInput(formNombre).trim();
	}

	function openCreateModal() {
		if (auth.user?.rol?.toLowerCase() !== 'admin') {
			toast.show('Acción denegada, requiere rol Administrador.', 'error');
			return;
		}
		editingCliente = null;
		formNombre = '';
		formRut = '';
		formTelefono = '';
		formFiadoActual = '';
		formFiadoMaximo = 20000;
		clearErrors();
		showModal = true;
	}

	function openEditModal(cliente: Cliente) {
		if (auth.user?.rol?.toLowerCase() !== 'admin') {
			toast.show('Acción denegada, requiere rol Administrador.', 'error');
			return;
		}
		editingCliente = cliente;
		formNombre = cliente.nombre;
		formRut = formatRutInput(cliente.rut);
		formTelefono = cliente.telefono;
		formFiadoActual = cliente.fiado_actual || 0;
		formFiadoMaximo = cliente.fiado_maximo || 20000;
		clearErrors();
		showModal = true;
	}

	function openDeleteModal(cliente: Cliente) {
		if (auth.user?.rol?.toLowerCase() !== 'admin') {
			toast.show('Acción denegada, requiere rol Administrador.', 'error');
			return;
		}
		clienteToDelete = cliente;
		showDeleteModal = true;
	}

	function openAbonoModal(cliente: Cliente) {
		selectedAbonoCliente = cliente;
		abonoMonto = '';
		abonoError = '';
		showAbonoModal = true;
	}

	async function handleAbonoSubmit(e: Event) {
		e.preventDefault();
		if (!selectedAbonoCliente || abonoMonto === '' || abonoMonto <= 0) {
			abonoError = 'Por favor, ingresa un monto válido mayor a cero.';
			return;
		}

		if (abonoMonto > selectedAbonoCliente.fiado_actual) {
			abonoError = `El abono no puede superar la deuda actual de ${formatCurrency(selectedAbonoCliente.fiado_actual)}`;
			return;
		}

		abonoLoading = true;
		abonoError = '';

		try {
			await apiFetch(`/ventas/clientes/${selectedAbonoCliente.id_cliente}/abonar`, {
				method: 'POST',
				body: JSON.stringify({ monto: Number(abonoMonto) })
			});

			toast.show(`Abono de ${formatCurrency(Number(abonoMonto))} registrado con éxito.`, 'success');
			showAbonoModal = false;
			await loadClientes();
		} catch (err: any) {
			abonoError = err.message || 'Error al procesar el abono.';
			toast.show(abonoError, 'error');
		} finally {
			abonoLoading = false;
		}
	}

	function clearErrors() {
		errNombre = '';
		errRut = '';
		errTelefono = '';
		formGeneralError = '';
	}

	async function handleSubmit(e: Event) {
		e.preventDefault();
		clearErrors();
		let isValid = true;

		if (!formNombre.trim()) {
			errNombre = 'El nombre es obligatorio.';
			isValid = false;
		} else {
			const soloLetras = /^[a-zA-ZáéíóúÁÉÍÓÚñÑ\s]+$/.test(formNombre);
			const partes = formNombre.trim().split(/\s+/);
			if (!soloLetras) {
				errNombre = 'El nombre solo debe contener letras y espacios.';
				isValid = false;
			} else if (partes.length < 2) {
				errNombre = 'Debe ingresar al menos un nombre y un apellido.';
				isValid = false;
			} else {
				formNombre = formatNombreInput(formNombre).trim();
			}
		}

		const rawRut = formRut.replace(/\./g, '');
		if (!formRut) {
			errRut = 'El RUT es obligatorio.';
			isValid = false;
		} else if (!validarRut(rawRut)) {
			errRut = 'El RUT ingresado no es válido para Chile.';
			isValid = false;
		}

		if (!formTelefono.trim()) {
			errTelefono = 'El teléfono es obligatorio.';
			isValid = false;
		} else if (!validarTelefono(formTelefono.trim())) {
			errTelefono = 'Formato no válido (ej: +56 9 1234 5678 o 912345678).';
			isValid = false;
		}

		if (formFiadoActual === '' || formFiadoActual <= 0 || formFiadoActual > 20000) {
			formGeneralError = 'El monto del fiado debe ser mayor a $0 y menor o igual a $20.000.';
			isValid = false;
		}

		if (!isValid) return;

		submitLoading = true;
		const payload = {
			nombre: formNombre.trim(),
			rut: rawRut.trim(),
			telefono: formTelefono.replace(/\s+/g, ''),
			fiado_actual: formFiadoActual === '' ? 0 : formFiadoActual,
			fiado_maximo: formFiadoMaximo
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
				if (err.errors.Nombre || err.errors.nombre)
					errNombre = err.errors.Nombre || err.errors.nombre;
				if (err.errors.Rut || err.errors.rut) errRut = err.errors.Rut || err.errors.rut;
				if (err.errors.Telefono || err.errors.telefono)
					errTelefono = err.errors.Telefono || err.errors.telefono;
			} else {
				formGeneralError = err.message || 'Error al guardar el cliente.';
			}
			toast.show('No se pudo guardar el cliente.', 'error');
		} finally {
			submitLoading = false;
		}
	}

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
	<title>Clientes - MinimarketGo</title>
	<meta name="description" content="Gestión de saldo de clientes (Fiados)" />
</svelte:head>

<div class="mb-6 flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
	<!--Búsqueda y Filtros -->
	<div class="flex flex-col gap-3 sm:flex-row sm:items-center flex-1">
		<!-- Buscador -->
		<div class="relative w-full sm:max-w-xs">
			<svg class="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-text-muted" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
			<input
				type="text"
				placeholder="Buscar cliente..."
				bind:value={searchQuery}
				class="w-full rounded-xl border border-border-color bg-bg-card py-2.5 pl-10 pr-4 text-sm text-text-primary transition-colors focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario"
			/>
		</div>

		<!-- Filtro de Estado -->
		<select
			bind:value={selectedEstado}
			class="w-full sm:w-auto rounded-xl border border-border-color bg-bg-card py-2.5 pl-4 pr-10 text-sm text-text-primary transition-colors focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario cursor-pointer"
		>
			<option value="Todos">Todos los estados</option>
			<option value="Con deuda">Con deuda</option>
			<option value="Al límite">Al límite</option>
		</select>
		
		<button
			type="button"
			title="Limpiar filtros"
			class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg text-text-muted transition-colors hover:bg-border-color hover:text-primario disabled:cursor-not-allowed disabled:opacity-50"
			onclick={() => { searchQuery = ''; selectedEstado = 'Todos'; }}
			disabled={!searchQuery && selectedEstado === 'Todos'}
		>
			<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
				<line x1="18" y1="6" x2="6" y2="18"></line>
				<line x1="6" y1="6" x2="18" y2="18"></line>
			</svg>
		</button>
	</div>

	<!-- agregar boton -->
	{#if auth.user?.rol?.toLowerCase() === 'admin'}
		<button
			onclick={openCreateModal}
			class="flex w-full items-center justify-center gap-2 rounded-xl bg-primario px-5 py-2.5 text-sm font-semibold text-white transition-transform hover:scale-[1.02] active:scale-95 sm:w-auto"
		>
			<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
				<line x1="12" y1="5" x2="12" y2="19" />
				<line x1="5" y1="12" x2="19" y2="12" />
			</svg>
			Nuevo Cliente
		</button>
	{/if}
</div>


{#if errorMsg}
	<div class="mb-5 flex gap-3 rounded-lg border border-red-500/15 bg-danger-bg p-4 text-sm text-danger-color" role="alert">
		<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
			<circle cx="12" cy="12" r="10" />
			<line x1="12" y1="8" x2="12" y2="12" />
			<line x1="12" y1="16" x2="12.01" y2="16" />
		</svg>
		<span>{errorMsg}</span>
	</div>
{/if}

	<div class="overflow-x-auto rounded-xl border border-border-color bg-bg-card shadow-sm">
		<table class="w-full whitespace-nowrap text-left text-sm text-text-primary">
			<thead class="border-b border-border-color bg-bg-primary/50 text-text-muted">
				<tr>
					<th class="px-6 py-4 font-semibold">Cliente</th>
					<th class="px-6 py-4 font-semibold">Teléfono</th>
					<th class="px-6 py-4 font-semibold min-w-[200px]">Fiado (Actual / Máx)</th>
					<th class="px-6 py-4 font-semibold">Última Compra</th>
					<th class="px-6 py-4 text-right font-semibold">Acciones</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-border-color">
				{#if loading}
					<tr>
						<td colspan="5" class="px-6 py-12 text-center text-text-muted">
							Cargando clientes...
						</td>
					</tr>
				{:else if filteredClientes.length === 0}
					<tr>
						<td colspan="5" class="px-6 py-12 text-center text-text-muted">
							No se encontraron clientes con los filtros aplicados.
						</td>
					</tr>
				{:else}
				{#each filteredClientes as client (client.id_cliente)}
					{@const porcentaje = Math.min(((client.fiado_actual || 0) / (client.fiado_maximo || 20000)) * 100, 100)}
					{@const colorClase = porcentaje >= 100 ? 'bg-danger-color' : porcentaje > 80 ? 'bg-orange-500' : 'bg-exito'}
					<tr class="group transition-colors hover:bg-bg-primary/30">
						<td class="px-6 py-4">
							<div class="flex items-center gap-3">
								<div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-border-color text-xs font-bold text-text-primary">
									{getInitials(client.nombre)}
								</div>
								<div class="flex flex-col">
									<span class="font-semibold">{client.nombre}</span>
									<span class="text-xs text-text-muted">{formatRutInput(client.rut)}</span>
								</div>
							</div>
						</td>
						<td class="px-6 py-4 font-medium">{client.telefono || '-'}</td>
						<td class="px-6 py-4">
							<div class="flex flex-col gap-1.5 w-full max-w-[200px]">
								<div class="flex justify-between text-xs">
									<span class="font-medium text-text-primary">{formatCurrency(client.fiado_actual || 0)}</span>
									<span class="text-text-muted">{formatCurrency(client.fiado_maximo || 20000)}</span>
								</div>
								<div class="h-1.5 w-full overflow-hidden rounded-full bg-border-color">
									<div class="h-full rounded-full {colorClase} transition-all duration-500" style="width: {porcentaje}%;"></div>
								</div>
							</div>
						</td>
						<td class="px-6 py-4 text-text-muted">
							{formatDate(client.ultima_compra)}
						</td>
						<td class="px-6 py-4 text-right">
							<div class="flex items-center justify-end gap-2">
								{#if (client.fiado_actual || 0) > 0}
									<button
										onclick={() => openAbonoModal(client)}
										class="inline-flex items-center gap-1 rounded-lg px-2.5 py-1.5 bg-exito/10 text-exito hover:bg-exito/20 text-xs font-bold transition-all"
										title="Abonar Deuda"
									>
										<Coins size={14} />
										<span>Abonar</span>
									</button>
								{/if}
								{#if auth.user?.rol?.toLowerCase() === 'admin'}
									<button
										onclick={() => openEditModal(client)}
										class="rounded-lg p-2 text-text-muted transition-colors hover:bg-border-color hover:text-primario"
										title="Editar"
									>
										<Edit2 size={18} />
									</button>
									<button
										onclick={() => openDeleteModal(client)}
										class="rounded-lg p-2 text-text-muted transition-colors hover:bg-danger-bg hover:text-danger-color"
										title="Eliminar"
									>
										<Trash2 size={18} />
									</button>
								{/if}
							</div>
						</td>
					</tr>
				{/each}
				{/if}
			</tbody>
		</table>
	</div>

<!-- Crear editar -->
{#if showModal}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-text-primary/40 p-4 backdrop-blur-[4px] animate-modal-enter" onclick={() => (showModal = false)} role="presentation">
		<div class="w-full max-w-lg overflow-hidden rounded-2xl border border-border-color bg-bg-card shadow-2xl" onclick={(e) => e.stopPropagation()} role="dialog">
			<header class="flex items-center justify-between border-b border-border-color px-6 py-4">
				<h3 class="text-lg font-bold text-text-primary">{editingCliente ? 'Editar Cliente' : 'Añadir Nuevo Cliente'}</h3>
				<button class="rounded-lg p-1 text-text-muted hover:bg-border-color hover:text-text-primary transition-colors" onclick={() => (showModal = false)}>
					<X size={20} />
				</button>
			</header>
			<form onsubmit={handleSubmit} autocomplete="off" class="p-6">
				{#if formGeneralError}
					<div class="mb-5 flex gap-3 rounded-lg border border-red-500/15 bg-danger-bg p-4 text-sm text-danger-color" role="alert">
						<span>{formGeneralError}</span>
					</div>
				{/if}
				
				<div class="grid grid-cols-1 gap-5">
					<div class="flex flex-col gap-1.5">
						<label class="text-sm font-semibold text-text-primary" for="formNombre">Nombre Completo</label>
						<input type="text" id="formNombre" autocomplete="off" class="rounded-xl border border-border-color bg-bg-primary px-4 py-2.5 text-sm text-text-primary focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario disabled:opacity-50" placeholder="Ej: Juan Pérez Gómez" bind:value={formNombre} onblur={handleNombreBlur} disabled={submitLoading} required />
						{#if errNombre}<span class="mt-1 text-xs font-medium text-danger-color">{errNombre}</span>{/if}
					</div>

					<div class="flex flex-col gap-1.5">
						<label class="text-sm font-semibold text-text-primary" for="formRut">RUT</label>
						<input type="text" id="formRut" class="rounded-xl border border-border-color bg-bg-primary px-4 py-2.5 text-sm text-text-primary focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario disabled:opacity-50" placeholder="Ej: 12.345.678-K" value={formRut} oninput={handleRutInput} disabled={submitLoading || editingCliente !== null} required />
						{#if errRut}<span class="mt-1 text-xs font-medium text-danger-color">{errRut}</span>{/if}
					</div>

					<div class="flex flex-col gap-1.5">
						<label class="text-sm font-semibold text-text-primary" for="formTelefono">Número Telefónico</label>
						<input type="text" id="formTelefono" class="rounded-xl border border-border-color bg-bg-primary px-4 py-2.5 text-sm text-text-primary focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario disabled:opacity-50" placeholder="Ej: +56912345678" bind:value={formTelefono} disabled={submitLoading} required />
						{#if errTelefono}<span class="mt-1 text-xs font-medium text-danger-color">{errTelefono}</span>{/if}
					</div>

					<div class="flex flex-col gap-1.5">
						<label class="text-sm font-semibold text-text-primary" for="formFiadoActual">Fiado</label>
						<div class="relative flex items-center">
							<span class="absolute left-4 text-text-muted font-medium">$</span>
							<input type="number" min="0" step="1" id="formFiadoActual" class="w-full rounded-xl border border-border-color bg-bg-primary px-4 py-2.5 pl-8 text-sm text-text-primary focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario disabled:opacity-50 [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none" bind:value={formFiadoActual} disabled={submitLoading} required />
						</div>
					</div>
				</div>
				<div class="mt-8 flex justify-end gap-4">
					<button type="button" class="font-semibold text-text-muted transition-colors hover:text-text-primary" onclick={() => (showModal = false)} disabled={submitLoading}>Cancelar</button>
					<button type="submit" class="rounded-xl bg-primario px-6 py-2.5 text-sm font-semibold text-white transition-transform hover:scale-[1.02] active:scale-95 disabled:cursor-not-allowed disabled:opacity-70" disabled={submitLoading}>
						{#if submitLoading} Procesando... {:else} {editingCliente ? 'Guardar Cambios' : 'Crear Cliente'} {/if}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<!-- borrar confirmacion del modal -->
{#if showDeleteModal && clienteToDelete}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-text-primary/40 p-5 backdrop-blur-[4px]" onclick={() => (showDeleteModal = false)} role="presentation">
		<div class="w-full max-w-[500px] overflow-hidden rounded-xl border border-border-color bg-bg-card shadow-lg animate-modal-enter" onclick={(e) => e.stopPropagation()} role="dialog">
			<header class="flex items-center justify-between border-b border-border-color p-5">
				<h2 class="text-lg font-bold text-danger-color">Confirmar Eliminación</h2>
				<button class="inline-flex cursor-pointer items-center justify-center rounded-lg border border-border-color bg-text-primary/3 p-2 text-text-secondary transition-all duration-200 hover:border-border-color-hover hover:bg-text-primary/7 hover:text-text-primary" onclick={() => (showDeleteModal = false)}>&times;</button>
			</header>
			<div class="p-6 text-text-primary">
				<p>¿Estás seguro de que deseas eliminar al cliente <strong>{clienteToDelete.nombre}</strong> de forma permanente?</p>
				<p class="mt-3 text-xs text-text-muted">No se podrá deshacer esta acción.</p>
			</div>
			<footer class="flex justify-end gap-3 border-t border-border-color bg-text-primary/2 p-4 px-6">
				<button type="button" class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-border-color bg-bg-secondary px-5 py-2.5 text-sm font-semibold text-text-primary transition-all duration-200 hover:bg-text-primary/5" onclick={() => (showDeleteModal = false)}>Cancelar</button>
				<button type="button" class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-red-500/15 bg-danger-bg px-5 py-2.5 text-sm font-semibold text-danger-color transition-all duration-200 hover:bg-danger-color hover:text-white" onclick={confirmDelete}>Eliminar Permanentemente</button>
			</footer>
		</div>
	</div>
{/if}

<!-- Abono Modal -->
{#if showAbonoModal && selectedAbonoCliente}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-text-primary/40 p-5 backdrop-blur-[4px]" onclick={() => (showAbonoModal = false)} role="presentation">
		<div class="w-full max-w-[500px] overflow-hidden rounded-xl border border-border-color bg-bg-card shadow-lg animate-modal-enter" onclick={(e) => e.stopPropagation()} role="dialog">
			<header class="flex items-center justify-between border-b border-border-color p-5">
				<h2 class="text-lg font-bold text-text-primary">Registrar Abono / Pago</h2>
				<button class="inline-flex cursor-pointer items-center justify-center rounded-lg border border-border-color bg-text-primary/3 p-2 text-text-secondary transition-all duration-200 hover:border-border-color-hover hover:bg-text-primary/7 hover:text-text-primary" onclick={() => (showAbonoModal = false)}>&times;</button>
			</header>
			<form onsubmit={handleAbonoSubmit}>
				<div class="p-6 flex flex-col gap-4 text-text-primary">
					<p class="text-sm">
						Registra un pago en efectivo o tarjeta para reducir la deuda de <strong>{selectedAbonoCliente.nombre}</strong>.
					</p>
					
					<div class="rounded-lg border border-border-color bg-bg-primary/50 p-4 text-xs flex flex-col gap-2">
						<div class="flex justify-between">
							<span>Deuda Pendiente:</span>
							<span class="font-bold text-danger-color">{formatCurrency(selectedAbonoCliente.fiado_actual)}</span>
						</div>
						<div class="flex justify-between">
							<span>Cupo Disponible:</span>
							<span class="font-bold text-exito">{formatCurrency(selectedAbonoCliente.fiado_maximo - selectedAbonoCliente.fiado_actual)}</span>
						</div>
					</div>

					<div class="flex flex-col gap-1.5 mt-2">
						<label class="text-sm font-semibold text-text-primary" for="abonoMonto">Monto del Abono</label>
						<div class="relative flex items-center">
							<span class="absolute left-4 text-text-muted font-medium">$</span>
							<input
								type="number"
								min="1"
								step="1"
								id="abonoMonto"
								class="[appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none w-full rounded-xl border border-border-color bg-bg-primary px-4 py-2.5 pl-8 text-sm text-text-primary focus:border-primario focus:outline-none disabled:opacity-50"
								placeholder="Ingresa la cantidad a abonar..."
								bind:value={abonoMonto}
								disabled={abonoLoading}
								required
								autofocus
							/>
						</div>
						{#if abonoError}
							<span class="mt-1 text-xs font-medium text-danger-color">{abonoError}</span>
						{/if}
					</div>
				</div>
				<footer class="flex justify-end gap-3 border-t border-border-color bg-text-primary/2 p-4 px-6">
					<button type="button" class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-border-color bg-bg-secondary px-5 py-2.5 text-sm font-semibold text-text-primary transition-all duration-200 hover:bg-text-primary/5" onclick={() => (showAbonoModal = false)} disabled={abonoLoading}>Cancelar</button>
					<button type="submit" class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg bg-primario px-5 py-2.5 text-sm font-semibold text-white transition-all duration-200 hover:bg-primario-hover shadow-md disabled:cursor-not-allowed disabled:opacity-70" disabled={abonoLoading}>
						{#if abonoLoading} Procesando... {:else} Registrar Pago {/if}
					</button>
				</footer>
			</form>
		</div>
	</div>
{/if}
