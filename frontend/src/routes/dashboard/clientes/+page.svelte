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

	let filteredClientes = $derived.by(() => {
		if (!searchQuery.trim()) return clientes;
		const query = searchQuery.toLowerCase().trim();
		const cleanQuery = query.replace(/[^0-9kK]/g, '');

		return clientes.filter((c) => {
			const matchNombre = c.nombre.toLowerCase().includes(query);
			const matchRut = cleanQuery ? c.rut.toLowerCase().replace(/[^0-9kK]/g, '').includes(cleanQuery) : false;
			return matchNombre || matchRut;
		});
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
	<title>Clientes - GPSproject</title>
	<meta name="description" content="Gestión de saldo de clientes (Fiados)" />
</svelte:head>

<!-- Stats Cards and Button Row -->
<div class="mb-6 flex flex-wrap items-start justify-between gap-6">
	<!-- Stats Cards -->
	<div class="grid flex-1 grid-cols-1 gap-6 md:grid-cols-2 lg:max-w-3xl">
		<!-- Card 1 -->
		<div class="flex flex-col justify-center rounded-xl border border-border-color bg-bg-card p-5 shadow-sm transition-all hover:border-border-color-hover">
			<span class="mb-1 text-sm font-bold uppercase tracking-wider text-text-muted">Fiado Total Pendiente</span>
			<h3 class="text-2xl font-black text-text-primary">{formatCurrency(fiadoTotalPendiente)}</h3>
			<p class="mt-2 text-xs font-medium text-text-secondary">
				<span class="text-accent">{clientesConSaldo}</span> clientes
			</p>
		</div>
		<!-- Card 2 -->
		<div class="flex flex-col justify-center rounded-xl border border-border-color bg-bg-card p-5 shadow-sm transition-all hover:border-border-color-hover">
			<span class="mb-1 text-sm font-bold uppercase tracking-wider text-text-muted">Clientes al Límite</span>
			<h3 class="text-2xl font-black text-danger-color">{clientesAlLimite}</h3>
			<p class="mt-2 text-xs font-medium text-text-secondary">Han alcanzado el monto maximo de fiado</p>
		</div>
	</div>

	<!-- Add Button -->
	<div class="flex items-start">
		{#if auth.user?.rol?.toLowerCase() === 'admin'}
			<button
				class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg bg-gradient-to-r from-accent-light to-accent px-5 py-2.5 text-sm font-semibold text-white transition-all duration-200 hover:-translate-y-[1px] hover:shadow-glow disabled:cursor-not-allowed disabled:opacity-50"
				onclick={openCreateModal}
			>
				<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
					<line x1="12" y1="5" x2="12" y2="19" />
					<line x1="5" y1="12" x2="19" y2="12" />
				</svg>
				<span>Nuevo Cliente</span>
			</button>
		{/if}
	</div>
</div>

<!-- Search Panel -->
<div class="mb-6 rounded-xl border border-border-color bg-bg-card p-6 shadow-md transition-all duration-300 hover:border-border-color-hover hover:shadow-lg">
	<div class="flex flex-wrap items-center gap-4">
		<div class="flex min-w-[280px] flex-1 overflow-hidden rounded-lg border border-[rgba(15,30,54,0.15)] bg-white focus-within:border-accent focus-within:ring-2 focus-within:ring-accent/15 dark:bg-bg-primary">
			<div class="flex items-center justify-center pl-4 pr-2 text-text-muted">
				<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
			</div>
			<input
				type="text"
				class="flex-1 border-none bg-transparent px-2 py-3 text-sm text-text-primary outline-none"
				placeholder="Buscar cliente..."
				bind:value={searchQuery}
				aria-label="Término de búsqueda"
			/>
		</div>
		<div class="flex gap-2.5">
			{#if searchQuery}
				<button
					type="button"
					class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-border-color bg-bg-secondary px-5 py-2.5 text-sm font-semibold text-text-primary transition-all duration-200 hover:bg-text-primary/5"
					onclick={clearSearch}
				>
					Limpiar
				</button>
			{/if}
		</div>
	</div>
</div>

<!-- Table View -->
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

{#if loading}
	<div class="overflow-hidden rounded-xl border border-border-color bg-bg-card shadow-md">
		<table class="w-full text-left border-collapse">
			<thead>
				<tr>
					<th class="border-b border-border-color bg-text-primary/4 p-4 text-xs font-bold uppercase tracking-wider text-text-secondary">Cliente</th>
					<th class="border-b border-border-color bg-text-primary/4 p-4 text-xs font-bold uppercase tracking-wider text-text-secondary">Teléfono</th>
					<th class="border-b border-border-color bg-text-primary/4 p-4 text-xs font-bold uppercase tracking-wider text-text-secondary">Fiado (Actual / Máx)</th>
					<th class="border-b border-border-color bg-text-primary/4 p-4 text-xs font-bold uppercase tracking-wider text-text-secondary">Última Compra</th>
					{#if auth.user?.rol?.toLowerCase() === 'admin'}
						<th class="w-[120px] border-b border-border-color bg-text-primary/4 p-4 text-right text-xs font-bold uppercase tracking-wider text-text-secondary">Acciones</th>
					{/if}
				</tr>
			</thead>
			<tbody>
				{#each Array(5) as _, i}
					<tr>
						<td class="border-b border-border-color p-4"><div class="h-4 w-40 rounded bg-gray-200 dark:bg-gray-700 animate-pulse"></div></td>
						<td class="border-b border-border-color p-4"><div class="h-4 w-28 rounded bg-gray-200 dark:bg-gray-700 animate-pulse"></div></td>
						<td class="border-b border-border-color p-4"><div class="h-4 w-32 rounded bg-gray-200 dark:bg-gray-700 animate-pulse"></div></td>
						<td class="border-b border-border-color p-4"><div class="h-4 w-24 rounded bg-gray-200 dark:bg-gray-700 animate-pulse"></div></td>
						{#if auth.user?.rol?.toLowerCase() === 'admin'}
							<td class="border-b border-border-color p-4">
								<div class="flex justify-end gap-2">
									<div class="h-8 w-8 rounded-lg bg-gray-200 dark:bg-gray-700 animate-pulse"></div>
									<div class="h-8 w-8 rounded-lg bg-gray-200 dark:bg-gray-700 animate-pulse"></div>
								</div>
							</td>
						{/if}
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
{:else if filteredClientes.length === 0}
	<div class="flex flex-col items-center justify-center rounded-xl border border-border-color bg-bg-card p-16 text-center shadow-md">
		<div class="mb-5 flex h-20 w-20 items-center justify-center rounded-full bg-text-primary/3 text-text-muted">
			<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
				<path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
				<circle cx="9" cy="7" r="4" />
				<line x1="17" y1="8" x2="22" y2="13" />
				<line x1="22" y1="8" x2="17" y2="13" />
			</svg>
		</div>
		<h3 class="mb-2 text-lg font-semibold text-text-primary">No se encontraron clientes</h3>
		<p class="max-w-[400px] text-sm text-text-secondary">Modifica el criterio de búsqueda o agrega un nuevo cliente al sistema.</p>
	</div>
{:else}
	<div class="overflow-hidden rounded-xl border border-border-color bg-bg-card shadow-md">
		<table class="w-full text-left border-collapse">
			<thead>
				<tr>
					<th class="border-b border-border-color bg-text-primary/4 p-4 text-xs font-bold uppercase tracking-wider text-text-secondary">Cliente</th>
					<th class="border-b border-border-color bg-text-primary/4 p-4 text-xs font-bold uppercase tracking-wider text-text-secondary">Teléfono</th>
					<th class="border-b border-border-color bg-text-primary/4 p-4 text-xs font-bold uppercase tracking-wider text-text-secondary min-w-[200px]">Fiado (Actual / Máx)</th>
					<th class="border-b border-border-color bg-text-primary/4 p-4 text-xs font-bold uppercase tracking-wider text-text-secondary">Última Compra</th>
					{#if auth.user?.rol?.toLowerCase() === 'admin'}
						<th class="w-[120px] border-b border-border-color bg-text-primary/4 p-4 text-right text-xs font-bold uppercase tracking-wider text-text-secondary">Acciones</th>
					{/if}
				</tr>
			</thead>
			<tbody>
				{#each filteredClientes as client (client.id_cliente)}
					{@const porcentaje = Math.min(((client.fiado_actual || 0) / (client.fiado_maximo || 20000)) * 100, 100)}
					{@const colorClase = porcentaje >= 100 ? 'bg-danger-color' : porcentaje > 80 ? 'bg-orange-500' : 'bg-exito'}
					<tr class="transition-colors hover:bg-text-primary/[0.015]">
						<td class="border-b border-border-color p-4">
							<div class="flex items-center gap-3">
								<div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-accent/10 text-accent font-bold">
									{getInitials(client.nombre)}
								</div>
								<div>
									<p class="font-semibold text-text-primary">{client.nombre}</p>
									<p class="text-xs text-text-muted">{formatRutInput(client.rut)}</p>
								</div>
							</div>
						</td>
						<td class="border-b border-border-color p-4 text-sm text-text-secondary">{client.telefono}</td>
						<td class="border-b border-border-color p-4">
							<div class="flex flex-col gap-1">
								<div class="flex justify-between text-xs font-medium">
									<span class="text-text-primary">{formatCurrency(client.fiado_actual || 0)}</span>
									<span class="text-text-muted">{formatCurrency(client.fiado_maximo || 20000)}</span>
								</div>
								<div class="h-2 w-full overflow-hidden rounded-full bg-border-color">
									<div class="h-full rounded-full {colorClase} transition-all duration-500" style="width: {porcentaje}%;"></div>
								</div>
							</div>
						</td>
						<td class="border-b border-border-color p-4 text-sm text-text-secondary">
							{formatDate(client.ultima_compra)}
						</td>
						{#if auth.user?.rol?.toLowerCase() === 'admin'}
							<td class="border-b border-border-color p-4">
								<div class="flex justify-end gap-2">
									<button
										class="inline-flex cursor-pointer items-center justify-center rounded-lg border border-border-color bg-text-primary/3 p-2 text-text-secondary transition-all duration-200 hover:border-border-color-hover hover:bg-text-primary/7 hover:text-text-primary"
										onclick={() => openEditModal(client)}
										title="Editar cliente"
									>
										<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9" /><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z" /></svg>
									</button>
									<button
										class="inline-flex cursor-pointer items-center justify-center rounded-lg border border-border-color bg-text-primary/3 p-2 text-text-secondary transition-all duration-200 hover:border-red-500/20 hover:bg-danger-bg hover:text-danger-color"
										onclick={() => openDeleteModal(client)}
										title="Eliminar cliente"
									>
										<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6" /><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" /><line x1="10" y1="11" x2="10" y2="17" /><line x1="14" y1="11" x2="14" y2="17" /></svg>
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
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-text-primary/40 p-5 backdrop-blur-[4px]" onclick={() => (showModal = false)} role="presentation">
		<div class="w-full max-w-[500px] overflow-hidden rounded-xl border border-border-color bg-bg-card shadow-lg animate-modal-enter" onclick={(e) => e.stopPropagation()} role="dialog">
			<header class="flex items-center justify-between border-b border-border-color p-5">
				<h2 class="text-lg font-bold text-text-primary">{editingCliente ? 'Editar Cliente' : 'Nuevo Cliente'}</h2>
				<button class="inline-flex cursor-pointer items-center justify-center rounded-lg border border-border-color bg-text-primary/3 p-2 text-text-secondary transition-all duration-200 hover:border-border-color-hover hover:bg-text-primary/7 hover:text-text-primary" onclick={() => (showModal = false)}>&times;</button>
			</header>
			<form onsubmit={handleSubmit}>
				<div class="p-6">
					{#if formGeneralError}
						<div class="mb-5 flex gap-3 rounded-lg border border-red-500/15 bg-danger-bg p-4 text-sm text-danger-color" role="alert">
							<span>{formGeneralError}</span>
						</div>
					{/if}
					
					<div class="mb-5 flex flex-col gap-1.5">
						<label class="text-[0.85rem] font-semibold text-text-secondary" for="formNombre">Nombre Completo</label>
						<input type="text" id="formNombre" class="rounded-lg border border-[rgba(15,30,54,0.15)] bg-white px-4 py-3 text-sm text-text-primary outline-none transition-all duration-200 focus:border-accent focus:ring-2 focus:ring-accent/15 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-bg-primary" placeholder="Ej: Juan Pérez Gómez" bind:value={formNombre} onblur={handleNombreBlur} disabled={submitLoading} required />
						{#if errNombre}<span class="mt-1 text-xs font-medium text-danger-color">{errNombre}</span>{/if}
					</div>

					<div class="mb-5 flex flex-col gap-1.5">
						<label class="text-[0.85rem] font-semibold text-text-secondary" for="formRut">RUT</label>
						<input type="text" id="formRut" class="rounded-lg border border-[rgba(15,30,54,0.15)] bg-white px-4 py-3 text-sm text-text-primary outline-none transition-all duration-200 focus:border-accent focus:ring-2 focus:ring-accent/15 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-bg-primary" placeholder="Ej: 12.345.678-K" value={formRut} oninput={handleRutInput} disabled={submitLoading || editingCliente !== null} required />
						{#if errRut}<span class="mt-1 text-xs font-medium text-danger-color">{errRut}</span>{/if}
					</div>

					<div class="mb-5 flex flex-col gap-1.5">
						<label class="text-[0.85rem] font-semibold text-text-secondary" for="formTelefono">Número Telefónico</label>
						<input type="text" id="formTelefono" class="rounded-lg border border-[rgba(15,30,54,0.15)] bg-white px-4 py-3 text-sm text-text-primary outline-none transition-all duration-200 focus:border-accent focus:ring-2 focus:ring-accent/15 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-bg-primary" placeholder="Ej: +56912345678" bind:value={formTelefono} disabled={submitLoading} required />
						{#if errTelefono}<span class="mt-1 text-xs font-medium text-danger-color">{errTelefono}</span>{/if}
					</div>

					<div class="mb-5 flex flex-col gap-1.5">
						<label class="text-[0.85rem] font-semibold text-text-secondary" for="formFiadoActual">Fiado</label>
						<div class="relative flex items-center">
							<span class="absolute left-4 text-text-muted font-medium">$</span>
							<input type="number" min="0" step="1" id="formFiadoActual" class="w-full rounded-lg border border-[rgba(15,30,54,0.15)] bg-white px-4 py-3 pl-8 text-sm text-text-primary outline-none transition-all duration-200 focus:border-accent focus:ring-2 focus:ring-accent/15 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-bg-primary [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none" bind:value={formFiadoActual} disabled={submitLoading} required />
						</div>
					</div>
				</div>
				<footer class="flex justify-end gap-3 border-t border-border-color bg-text-primary/2 p-4 px-6">
					<button type="button" class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-border-color bg-bg-secondary px-5 py-2.5 text-sm font-semibold text-text-primary transition-all duration-200 hover:bg-text-primary/5" onclick={() => (showModal = false)} disabled={submitLoading}>Cancelar</button>
					<button type="submit" class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg bg-gradient-to-r from-accent-light to-accent px-5 py-2.5 text-sm font-semibold text-white transition-all duration-200 hover:-translate-y-[1px] hover:shadow-glow disabled:cursor-not-allowed disabled:opacity-50" disabled={submitLoading}>
						{#if submitLoading} Procesando... {:else} Guardar Cambios {/if}
					</button>
				</footer>
			</form>
		</div>
	</div>
{/if}

<!-- Delete Confirmation Modal -->
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
