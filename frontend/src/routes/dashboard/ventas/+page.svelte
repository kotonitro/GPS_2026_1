<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { apiFetch, apiClientes, apiVentas, apiCajas, obtenerPromociones } from '$lib/api';
	import { auth } from '$lib/authStore.svelte';
	import { toast } from '$lib/toastStore.svelte';
	import {
		Search,
		Plus,
		Minus,
		Trash2,
		User,
		CreditCard,
		Coins,
		FileText,
		History,
		ShoppingBag,
		ShoppingCart,
		ChevronRight,
		Percent,
		AlertCircle,
		Keyboard,
		Laptop,
		Sparkles,
		Camera
	} from '@lucide/svelte';
	import Scanner from '$lib/components/Scanner.svelte';

	interface Producto {
		id_producto: string;
		nombre: string;
		descripcion: string;
		stock: number;
		stock_minimo: number;
		precio: number;
		marca: string;
		codigo_barras: string;
		estado: boolean;
		id_categoria: string;
	}

	interface Cliente {
		id_cliente: string;
		nombre: string;
		rut: string;
		telefono: string;
		fiado_actual: number;
		fiado_maximo: number;
	}

	interface Caja {
		id_caja: string;
		nombre: string;
		ubicacion: string;
		activo: boolean;
	}

	interface DetalleVenta {
		id_detalle?: string;
		id_venta?: string;
		id_producto: string;
		cantidad: number;
		monto_final: number;
		producto?: Producto;
	}

	interface Venta {
		id_venta: string;
		id_caja: string;
		id_metodo: string;
		id_empleado: string;
		fecha_emision: string;
		pago: number;
		vuelto: number;
		monto_total: number;
		monto_descuento: number;
		estado_sync: string;
		detalles?: DetalleVenta[];
		metodo_pago?: {
			id_metodo: string;
			nombre_metodo: string;
		};
	}

	interface Promocion {
		id_promocion: string;
		tipo: string;
		lleva: number;
		paga: number;
		descuento: number;
		producto_id: string;
		fecha_inicio: string | null;
		fecha_fin: string | null;
	}

	// State variables
	let activeTab = $state<'pos' | 'history'>('pos');
	let loading = $state(true);
	let submitting = $state(false);

	// Data from API
	let productos = $state<Producto[]>([]);
	let clientes = $state<Cliente[]>([]);
	let promociones = $state<Promocion[]>([]);
	let cajas = $state<Caja[]>([]);
	let ventas = $state<Venta[]>([]);

	// POS Cart & Register State
	let selectedCajaId = $state<string>('');
	let selectedMetodoId = $state<string>('11111111-1111-1111-1111-111111111111'); // Default to Efectivo
	let selectedClienteId = $state<string>('');
	let discountPercent = $state<number>(0);
	let cashReceived = $state<number | ''>('');

	// Modals
	let showCajaConfigModal = $state(false);
	let showManualCodeModal = $state(false);
	let manualBarcodeValue = $state('');
	let modoEscaneo = $state(false);

	// Cart Items
	interface CartItem {
		producto: Producto;
		cantidad: number;
	}
	let cart = $state<CartItem[]>([]);

	// Last scanned item for visual feedback
	let lastScannedProduct = $state<Producto | null>(null);

	// Global barcode scan variables
	let barcodeBuffer = '';
	let lastKeyTime = 0;

	// Search query for history
	let searchHistoryQuery = $state('');

	// Selected sale for detail modal
	let selectedSale = $state<Venta | null>(null);
	let showDetailModal = $state(false);

	// Filtered list for history
	let filteredSales = $derived.by(() => {
		const query = searchHistoryQuery.toLowerCase().trim();
		return ventas.filter((v) => {
			if (!query) return true;
			return (
				v.id_venta.toLowerCase().includes(query) ||
				v.id_empleado.toLowerCase().includes(query)
			);
		});
	});

	function isPromotionActive(p: Promocion): boolean {
		const now = new Date();
		if (p.fecha_inicio) {
			const start = new Date(p.fecha_inicio);
			if (now < start) return false;
		}
		if (p.fecha_fin) {
			const end = new Date(p.fecha_fin);
			if (now > end) return false;
		}
		return true;
	}

	function calculateItemDiscount(item: CartItem, promo: Promocion): number {
		const price = item.producto.precio;
		const qty = item.cantidad;

		if (promo.tipo === 'NXM') {
			if (promo.lleva > 0 && promo.paga > 0 && promo.lleva > promo.paga) {
				const sets = Math.floor(qty / promo.lleva);
				const discountQty = sets * (promo.lleva - promo.paga);
				return discountQty * price;
			}
		} else if (promo.tipo === 'porcentaje') {
			if (promo.descuento > 0) {
				return Math.round(qty * price * (promo.descuento / 100));
			}
		} else if (promo.tipo === 'precio_fijo') {
			if (promo.descuento > 0 && price > promo.descuento) {
				const unitDiscount = price - promo.descuento;
				return qty * unitDiscount;
			}
		}
		return 0;
	}

	function getBestDiscountForItem(item: CartItem): { discount: number; promo: Promocion | null } {
		const activePromos = promociones.filter(
			(p) => p.producto_id === item.producto.id_producto && isPromotionActive(p)
		);
		let maxDiscount = 0;
		let bestPromo: Promocion | null = null;

		for (const promo of activePromos) {
			const disc = calculateItemDiscount(item, promo);
			if (disc > maxDiscount) {
				maxDiscount = disc;
				bestPromo = promo;
			}
		}

		return { discount: maxDiscount, promo: bestPromo };
	}

	// Totals computations
	let subtotal = $derived(cart.reduce((acc, item) => acc + item.producto.precio * item.cantidad, 0));
	let discountAmount = $derived(
		cart.reduce((acc, item) => acc + getBestDiscountForItem(item).discount, 0)
	);
	let total = $derived(Math.max(subtotal - discountAmount, 0));
	let change = $derived(
		selectedMetodoId === '11111111-1111-1111-1111-111111111111' && typeof cashReceived === 'number'
			? Math.max(cashReceived - total, 0)
			: 0
	);

	// History statistics
	let totalVendidoHoy = $derived(
		ventas
			.filter((v) => {
				const date = new Date(v.fecha_emision);
				const today = new Date();
				return (
					date.getDate() === today.getDate() &&
					date.getMonth() === today.getMonth() &&
					date.getFullYear() === today.getFullYear()
				);
			})
			.reduce((acc, v) => acc + v.monto_total, 0)
	);
	let cantidadVentasHoy = $derived(
		ventas.filter((v) => {
			const date = new Date(v.fecha_emision);
			const today = new Date();
			return (
				date.getDate() === today.getDate() &&
				date.getMonth() === today.getMonth() &&
				date.getFullYear() === today.getFullYear()
			);
		}).length
	);
	let ticketPromedio = $derived(cantidadVentasHoy > 0 ? Math.round(totalVendidoHoy / cantidadVentasHoy) : 0);

	// Selected client details
	let selectedClientData = $derived(clientes.find((c) => c.id_cliente === selectedClienteId));
	let isFiadoLimitExceeded = $derived.by(() => {
		if (selectedMetodoId !== 'fiado' || !selectedClientData) return false;
		return selectedClientData.fiado_actual + total > selectedClientData.fiado_maximo;
	});

	// Global Keyboard Listener for barcode scanning
	function handleGlobalKeydown(e: KeyboardEvent) {
		if (activeTab !== 'pos') return;

		// Ignore keys if focus is inside any text input or selection elements
		const target = e.target as HTMLElement;
		if (
			target.tagName === 'INPUT' ||
			target.tagName === 'SELECT' ||
			target.tagName === 'TEXTAREA' ||
			target.isContentEditable
		) {
			return;
		}

		const currentTime = Date.now();

		// If delay is > 100ms, assume human typed and reset buffer
		if (currentTime - lastKeyTime > 100) {
			barcodeBuffer = '';
		}
		lastKeyTime = currentTime;

		if (e.key === 'Enter') {
			const code = barcodeBuffer.trim();
			if (code) {
				processBarcodeScan(code);
			}
			barcodeBuffer = '';
			e.preventDefault();
		} else if (e.key.length === 1) {
			barcodeBuffer += e.key;
		}
	}

	onMount(async () => {
		window.addEventListener('keydown', handleGlobalKeydown);
		await loadData();
	});

	onDestroy(() => {
		if (typeof window !== 'undefined') {
			window.removeEventListener('keydown', handleGlobalKeydown);
		}
	});

	async function loadData() {
		loading = true;
		try {
			const [resProducts, resClients, resCajas, resSales, resPromociones] = await Promise.all([
				apiFetch('/inventario/productos').catch(() => []),
				apiClientes.getAll().catch(() => []),
				apiCajas.getAll().catch(() => []),
				apiVentas.getAll().catch(() => []),
				obtenerPromociones().catch(() => [])
			]);

			productos = Array.isArray(resProducts) ? resProducts : [];
			clientes = Array.isArray(resClients) ? resClients : [];
			cajas = Array.isArray(resCajas) ? resCajas.filter((c) => c.activo) : [];
			ventas = Array.isArray(resSales) ? resSales : [];
			promociones = Array.isArray(resPromociones) ? resPromociones : [];

			// Detect Caja using Local Storage
			detectCaja();
		} catch (err: any) {
			toast.show(err.message || 'Error al cargar los datos del módulo de ventas.', 'error');
		} finally {
			loading = false;
		}
	}

	function detectCaja() {
		const savedCajaId = localStorage.getItem('id_caja');
		const activeCajaExists = cajas.some((c) => c.id_caja === savedCajaId);

		if (savedCajaId && activeCajaExists) {
			selectedCajaId = savedCajaId;
		} else {
			// If not set, or selected caja is no longer active/exists
			selectedCajaId = '';
			showCajaConfigModal = true;
		}
	}

	function handleCajaSelection(cajaId: string) {
		selectedCajaId = cajaId;
		localStorage.setItem('id_caja', cajaId);
		showCajaConfigModal = false;
		const cajaName = cajas.find((c) => c.id_caja === cajaId)?.nombre || 'Caja';
		toast.show(`Equipo asociado correctamente a: ${cajaName}`, 'success');
	}

	function changeCaja() {
		showCajaConfigModal = true;
	}

	function formatCurrency(amount: number) {
		return new Intl.NumberFormat('es-CL', {
			style: 'currency',
			currency: 'CLP',
			maximumFractionDigits: 0
		}).format(amount);
	}

	function processBarcodeScan(code: string) {
		const product = productos.find((p) => p.codigo_barras === code);
		if (product) {
			if (product.stock <= 0) {
				toast.show(`El producto "${product.nombre}" no tiene stock disponible.`, 'error');
			} else {
				addToCart(product);
				lastScannedProduct = product;
				toast.show(`Añadido: ${product.nombre}`, 'success');
			}
		} else {
			toast.show(`Código de barras "${code}" no encontrado.`, 'error');
		}
	}

	function handleManualBarcodeSubmit(e: Event) {
		e.preventDefault();
		const code = manualBarcodeValue.trim();
		if (!code) return;
		processBarcodeScan(code);
		manualBarcodeValue = '';
		showManualCodeModal = false;
	}

	function manejarEscaneo(codigo: string) {
		modoEscaneo = false;
		processBarcodeScan(codigo);
	}

	function addToCart(producto: Producto) {
		const existingItemIndex = cart.findIndex((item) => item.producto.id_producto === producto.id_producto);
		if (existingItemIndex > -1) {
			const item = cart[existingItemIndex];
			if (item.cantidad + 1 > producto.stock) {
				toast.show(`Stock máximo alcanzado para ${producto.nombre}`, 'error');
				return;
			}
			cart[existingItemIndex].cantidad += 1;
		} else {
			cart.push({ producto, cantidad: 1 });
		}
		cashReceived = '';
	}

	function updateQuantity(productId: string, delta: number) {
		const index = cart.findIndex((item) => item.producto.id_producto === productId);
		if (index > -1) {
			const item = cart[index];
			const newQty = item.cantidad + delta;
			if (newQty <= 0) {
				cart.splice(index, 1);
			} else if (newQty > item.producto.stock) {
				toast.show(`Solo hay ${item.producto.stock} unidades disponibles de ${item.producto.nombre}`, 'error');
			} else {
				cart[index].cantidad = newQty;
			}
		}
		cashReceived = '';
	}

	function removeFromCart(productId: string) {
		cart = cart.filter((item) => item.producto.id_producto !== productId);
		cashReceived = '';
	}

	function clearCart() {
		cart = [];
		discountPercent = 0;
		cashReceived = '';
		selectedClienteId = '';
		selectedMetodoId = '11111111-1111-1111-1111-111111111111';
		lastScannedProduct = null;
	}

	async function handleCheckout(e: Event) {
		e.preventDefault();
		if (cart.length === 0) {
			toast.show('El carrito está vacío.', 'error');
			return;
		}

		if (!selectedCajaId) {
			toast.show('Por favor, asocia el equipo a una Caja activa.', 'error');
			showCajaConfigModal = true;
			return;
		}

		if (selectedMetodoId === '11111111-1111-1111-1111-111111111111') {
			if (cashReceived === '' || cashReceived < total) {
				toast.show('El pago ingresado es insuficiente.', 'error');
				return;
			}
		}

		if (selectedMetodoId === 'fiado') {
			if (!selectedClienteId) {
				toast.show('Por favor, selecciona un cliente para registrar el Fiado.', 'error');
				return;
			}
			if (isFiadoLimitExceeded) {
				toast.show('El monto excede el fiado máximo permitido para este cliente.', 'error');
				return;
			}
		}

		submitting = true;
		
		const detallesPayload = cart.map((item) => {
			const { discount } = getBestDiscountForItem(item);
			return {
				id_producto: item.producto.id_producto,
				cantidad: item.cantidad,
				monto_final: item.producto.precio * item.cantidad - discount
			};
		});

		const metodoIdForBackend = selectedMetodoId === 'fiado' 
			? '11111111-1111-1111-1111-111111111111' 
			: selectedMetodoId;

		const payload = {
			id_caja: selectedCajaId,
			id_metodo: metodoIdForBackend,
			pago: selectedMetodoId === '11111111-1111-1111-1111-111111111111' ? (cashReceived as number) : total,
			vuelto: selectedMetodoId === '11111111-1111-1111-1111-111111111111' ? change : 0,
			monto_total: total,
			monto_descuento: discountAmount,
			detalles: detallesPayload
		};

		try {
			await apiVentas.create(payload);

			if (selectedMetodoId === 'fiado' && selectedClientData) {
				const nuevoFiado = selectedClientData.fiado_actual + total;
				await apiClientes.update(selectedClienteId, {
					fiado_actual: nuevoFiado
				} as any);
				toast.show(`Fiado registrado por ${formatCurrency(total)} para ${selectedClientData.nombre}`, 'success');
			}

			toast.show('Venta registrada con éxito.', 'success');
			clearCart();
			await loadData();
		} catch (err: any) {
			toast.show(err.message || 'Error al procesar la venta.', 'error');
		} finally {
			submitting = false;
		}
	}

	function openSaleDetail(sale: Venta) {
		selectedSale = sale;
		showDetailModal = true;
	}

	function getProductDetailName(productId: string) {
		const prod = productos.find((p) => p.id_producto === productId);
		return prod ? prod.nombre : 'Producto desconocido';
	}

	// Get name of current caja
	let activeCajaName = $derived(cajas.find((c) => c.id_caja === selectedCajaId)?.nombre || 'Ninguna');
</script>

<svelte:head>
	<title>Ventas - GPSproject</title>
	<meta name="description" content="Registro y control de ventas" />
</svelte:head>

<!-- Tabs Navigation -->
<div class="mb-6 flex border-b border-border-color justify-between items-center">
	<div class="flex">
		<button
			onclick={() => (activeTab = 'pos')}
			class="inline-flex items-center gap-2 border-b-2 px-6 py-3 text-sm font-semibold transition-all duration-200 {activeTab === 'pos' ? 'border-primario text-primario' : 'border-transparent text-text-secondary hover:text-text-primary'}"
		>
			<ShoppingBag size={18} />
			<span>Registrar Venta (Scanner POS)</span>
		</button>
		<button
			onclick={() => (activeTab = 'history')}
			class="inline-flex items-center gap-2 border-b-2 px-6 py-3 text-sm font-semibold transition-all duration-200 {activeTab === 'history' ? 'border-primario text-primario' : 'border-transparent text-text-secondary hover:text-text-primary'}"
		>
			<History size={18} />
			<span>Historial de Ventas</span>
		</button>
	</div>

	<!-- Caja Detection Indicator -->
	<div class="flex items-center gap-3 pr-4">
		<div class="flex items-center gap-2 rounded-lg border border-border-color bg-bg-card px-3 py-1.5 text-xs text-text-secondary">
			<Laptop size={14} class="text-primario" />
			<span>Terminal asociado a: <strong class="text-text-primary">{activeCajaName}</strong></span>
		</div>
		<button
			onclick={changeCaja}
			class="text-xs text-primario hover:underline font-semibold"
		>
			Cambiar Caja
		</button>
	</div>
</div>

{#if loading}
	<div class="flex h-64 items-center justify-center text-xl text-primario">
		Cargando datos de ventas...
	</div>
{:else}
	<!-- POS TAB -->
	{#if activeTab === 'pos'}
		<div class="grid grid-cols-1 gap-6 lg:grid-cols-12">
			
			<!-- POS Scanner Information Panel (Left side) -->
			<div class="flex flex-col gap-6 lg:col-span-7 justify-center items-center rounded-xl border border-dashed border-border-color bg-bg-card/50 p-8 text-center min-h-[400px]">
				
				<!-- Bouncing Scan Circle -->
				<div class="relative flex h-28 w-28 items-center justify-center rounded-full bg-primario/10 text-primario animate-pulse">
					<svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
						<path stroke-linecap="round" stroke-linejoin="round" d="M3 5v14c0 1.1.9 2 2 2h14a2 2 0 0 0 2-2V5a2 2 0 0 0-2-2H5a2 2 0 0 0-2 2z"/>
						<path stroke-linecap="round" stroke-linejoin="round" d="M7 7v10M10 7v10M13 7v10M17 7v10"/>
					</svg>
					<!-- Glowing Green Status Indicator -->
					<span class="absolute right-1 bottom-1 flex h-4 w-4">
						<span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-exito opacity-75"></span>
						<span class="relative inline-flex rounded-full h-4 w-4 bg-exito border-2 border-bg-card"></span>
					</span>
				</div>

				<div>
					<h3 class="text-xl font-bold text-text-primary flex items-center justify-center gap-2">
						<span>Lector de Barra Activo</span>
					</h3>
					<p class="text-sm text-text-secondary max-w-sm mt-2">
						El sistema está listo. Simplemente pasa la etiqueta de código de barras del producto por el lector físico para registrarlo en el carrito.
					</p>
					<p class="text-xs text-text-muted mt-1">
						(No es necesario hacer clic en ningún campo de texto)
					</p>
				</div>

				<!-- Visual Feedback: Last Scanned Item -->
				{#if lastScannedProduct}
					<div class="w-full max-w-md rounded-xl border border-primario/30 bg-primario/5 p-4 text-left flex justify-between items-center animate-modal-enter">
						<div class="flex-1">
							<span class="text-[10px] uppercase font-bold tracking-wider text-primario flex items-center gap-1">
								<Sparkles size={10} /> Escaneado Recientemente
							</span>
							<p class="font-bold text-text-primary text-sm mt-1">{lastScannedProduct.nombre}</p>
							<p class="text-xs text-text-muted mt-0.5">{lastScannedProduct.marca} • Cód: {lastScannedProduct.codigo_barras}</p>
						</div>
						<div class="text-right">
							<span class="text-base font-black text-primario">{formatCurrency(lastScannedProduct.precio)}</span>
							<p class="text-[10px] text-text-muted mt-0.5">Stock restante: {lastScannedProduct.stock}</p>
						</div>
					</div>
				{/if}

				<!-- Manual Code & Camera Scan Buttons -->
				<div class="mt-4 flex flex-wrap gap-3">
					<button
						type="button"
						onclick={() => (showManualCodeModal = true)}
						class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-border-color bg-bg-card px-4 py-2.5 text-xs font-semibold text-text-secondary transition-all hover:bg-text-primary/5"
					>
						<Keyboard size={14} />
						<span>Ingresar Código a Mano</span>
					</button>
					<button
						type="button"
						onclick={() => (modoEscaneo = true)}
						class="inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-border-color bg-bg-card px-4 py-2.5 text-xs font-semibold text-text-secondary transition-all hover:bg-text-primary/5"
					>
						<Camera size={14} />
						<span>Escanear con Cámara</span>
					</button>
				</div>
			</div>

			<!-- Cart & Checkout (Right side) -->
			<div class="lg:col-span-5">
				<form onsubmit={handleCheckout} class="flex flex-col rounded-xl border border-border-color bg-bg-card shadow-lg overflow-hidden">
					<header class="border-b border-border-color bg-text-primary/[0.02] p-4 flex justify-between items-center">
						<h3 class="font-bold text-text-primary flex items-center gap-2">
							<ShoppingCart size={18} class="text-primario" />
							<span>Carrito de Ventas</span>
						</h3>
						{#if cart.length > 0}
							<button type="button" onclick={clearCart} class="text-xs text-danger-color hover:underline">
								Vaciar
							</button>
						{/if}
					</header>

					<!-- Cart Items List -->
					<div class="p-4 flex flex-col gap-3 min-h-[220px] max-h-[300px] overflow-y-auto border-b border-border-color">
						{#if cart.length === 0}
							<div class="flex flex-col items-center justify-center h-full text-center text-text-muted py-8">
								<ShoppingBag size={32} class="mb-2 opacity-50" />
								<p class="text-sm">Escanea productos para agregarlos al carrito</p>
							</div>
						{:else}
							{#each cart as item (item.producto.id_producto)}
								{@const { discount, promo } = getBestDiscountForItem(item)}
								<div class="flex items-center justify-between gap-3 p-2 rounded-lg bg-text-primary/[0.015] border border-border-color/50">
									<div class="flex-1">
										<p class="text-sm font-semibold text-text-primary line-clamp-1">{item.producto.nombre}</p>
										<p class="text-xs text-primario font-bold">{formatCurrency(item.producto.precio)} c/u</p>
									</div>
									
									<!-- Quantity Control -->
									<div class="flex items-center gap-2">
										<button
											type="button"
											onclick={() => updateQuantity(item.producto.id_producto, -1)}
											class="flex h-7 w-7 items-center justify-center rounded-lg border border-border-color bg-bg-secondary text-text-primary hover:bg-text-primary/5"
										>
											<Minus size={12} />
										</button>
										<span class="text-sm font-bold w-6 text-center text-text-primary">{item.cantidad}</span>
										<button
											type="button"
											onclick={() => updateQuantity(item.producto.id_producto, 1)}
											class="flex h-7 w-7 items-center justify-center rounded-lg border border-border-color bg-bg-secondary text-text-primary hover:bg-text-primary/5"
										>
											<Plus size={12} />
										</button>
									</div>

									<!-- Item Total & Delete -->
									<div class="text-right flex items-center gap-3">
										<div class="flex flex-col min-w-[80px]">
											{#if discount > 0}
												<span class="text-xs text-text-muted line-through">
													{formatCurrency(item.producto.precio * item.cantidad)}
												</span>
												<span class="text-sm font-bold text-danger-color">
													{formatCurrency(item.producto.precio * item.cantidad - discount)}
												</span>
												{#if promo}
													<span class="text-[9px] text-accent font-bold">
														{promo.tipo === 'NXM' ? `${promo.lleva}x${promo.paga}` : promo.tipo === 'porcentaje' ? `-${promo.descuento}%` : 'Promo'}
													</span>
												{/if}
											{:else}
												<span class="text-sm font-bold text-text-primary">
													{formatCurrency(item.producto.precio * item.cantidad)}
												</span>
											{/if}
										</div>
										<button
											type="button"
											onclick={() => removeFromCart(item.producto.id_producto)}
											class="text-text-muted hover:text-danger-color p-1"
										>
											<Trash2 size={14} />
										</button>
									</div>
								</div>
							{/each}
						{/if}
					</div>

					<!-- Checkout Options -->
					<div class="p-4 flex flex-col gap-4 bg-text-primary/[0.005]">

						<!-- Payment Method Selection -->
						<div class="flex flex-col gap-1">
							<span class="text-xs font-bold uppercase tracking-wider text-text-secondary mb-1">Método de Pago</span>
							<div class="grid grid-cols-3 gap-2">
								<button
									type="button"
									onclick={() => { selectedMetodoId = '11111111-1111-1111-1111-111111111111'; cashReceived = ''; }}
									class="flex flex-col items-center justify-center p-2 rounded-lg border text-xs font-semibold gap-1 transition-all {selectedMetodoId === '11111111-1111-1111-1111-111111111111' ? 'border-primario bg-primario/10 text-primario' : 'border-border-color text-text-secondary hover:bg-text-primary/5'}"
								>
									<Coins size={16} />
									<span>Efectivo</span>
								</button>
								<button
									type="button"
									onclick={() => { selectedMetodoId = '22222222-2222-2222-2222-222222222222'; cashReceived = ''; selectedClienteId = ''; }}
									class="flex flex-col items-center justify-center p-2 rounded-lg border text-xs font-semibold gap-1 transition-all {selectedMetodoId === '22222222-2222-2222-2222-222222222222' ? 'border-primario bg-primario/10 text-primario' : 'border-border-color text-text-secondary hover:bg-text-primary/5'}"
								>
									<CreditCard size={16} />
									<span>Tarjeta</span>
								</button>
								<button
									type="button"
									onclick={() => { selectedMetodoId = 'fiado'; cashReceived = ''; }}
									class="flex flex-col items-center justify-center p-2 rounded-lg border text-xs font-semibold gap-1 transition-all {selectedMetodoId === 'fiado' ? 'border-primario bg-primario/10 text-primario' : 'border-border-color text-text-secondary hover:bg-text-primary/5'}"
								>
									<User size={16} />
									<span>Fiado</span>
								</button>
							</div>
						</div>

						<!-- Conditional Fields -->
						{#if selectedMetodoId === '11111111-1111-1111-1111-111111111111'}
							<!-- Cash Received -->
							<div class="flex gap-4">
								<div class="flex-1 flex flex-col gap-1">
									<label for="cashReceived" class="text-xs font-bold uppercase tracking-wider text-text-secondary">Efectivo Recibido</label>
									<div class="relative flex items-center">
										<span class="absolute left-3 text-sm text-text-muted font-bold">$</span>
										<input id="cashReceived" type="number" min="0" placeholder="0" bind:value={cashReceived} class="w-full rounded-lg border border-border-color bg-bg-card py-2 pl-7 pr-3 text-sm text-text-primary outline-none focus:border-accent" required />
									</div>
								</div>
								<div class="flex-1 flex flex-col justify-end pb-1">
									<span class="text-xs font-bold uppercase tracking-wider text-text-muted">Vuelto:</span>
									<span class="text-lg font-black text-exito mt-1">{formatCurrency(change)}</span>
								</div>
							</div>
						{:else if selectedMetodoId === 'fiado'}
							<!-- Client Selection -->
							<div class="flex flex-col gap-1">
								<label for="clientSelect" class="text-xs font-bold uppercase tracking-wider text-text-secondary">Cliente Asoc.</label>
								<select id="clientSelect" bind:value={selectedClienteId} class="w-full rounded-lg border border-border-color bg-bg-card p-2.5 text-sm text-text-primary outline-none focus:border-accent" required>
									<option value="" disabled>-- Selecciona un cliente --</option>
									{#each clientes as cli}
										<option value={cli.id_cliente}>{cli.nombre} (Deuda: {formatCurrency(cli.fiado_actual)})</option>
									{/each}
								</select>
								{#if selectedClientData}
									<div class="mt-2 p-3 rounded-lg border text-xs flex flex-col gap-1 {isFiadoLimitExceeded ? 'border-red-500/15 bg-danger-bg text-danger-color' : 'border-border-color bg-text-primary/[0.01]' }">
										<div class="flex justify-between">
											<span>Fiado Actual:</span>
											<span class="font-bold">{formatCurrency(selectedClientData.fiado_actual)}</span>
										</div>
										<div class="flex justify-between">
											<span>Cupo Máximo:</span>
											<span class="font-bold">{formatCurrency(selectedClientData.fiado_maximo)}</span>
										</div>
										<div class="flex justify-between border-t border-dashed pt-1 mt-1 font-semibold">
											<span>Nueva Deuda:</span>
											<span>{formatCurrency(selectedClientData.fiado_actual + total)}</span>
										</div>
										{#if isFiadoLimitExceeded}
											<p class="font-bold text-[10px] uppercase mt-1">⚠️ Excede el saldo máximo permitido</p>
										{/if}
									</div>
								{/if}
							</div>
						{/if}

						<!-- Checkout Summary -->
						<div class="border-t border-border-color pt-4 flex flex-col gap-1">
							<div class="flex justify-between text-xs text-text-secondary">
								<span>Subtotal:</span>
								<span>{formatCurrency(subtotal)}</span>
							</div>
							{#if discountAmount > 0}
								<div class="flex justify-between text-xs text-danger-color">
									<span>Descuento por Promoción:</span>
									<span>-{formatCurrency(discountAmount)}</span>
								</div>
							{/if}
							<div class="flex justify-between items-center text-text-primary border-t border-dashed border-border-color pt-2 mt-1">
								<span class="font-black text-sm">TOTAL:</span>
								<span class="font-black text-2xl text-primario">{formatCurrency(total)}</span>
							</div>
						</div>
					</div>

					<footer class="p-4 bg-text-primary/[0.02] border-t border-border-color">
						<button
							type="submit"
							disabled={cart.length === 0 || submitting || isFiadoLimitExceeded || !selectedCajaId}
							class="w-full inline-flex cursor-pointer items-center justify-center gap-2 rounded-lg bg-gradient-to-r from-accent-light to-accent py-3 font-bold text-white shadow-md transition-all hover:bg-primario-hover hover:shadow-glow disabled:cursor-not-allowed disabled:opacity-50"
						>
							{#if submitting}
								Procesando Venta...
							{:else}
								Confirmar Venta ({formatCurrency(total)})
							{/if}
						</button>
					</footer>
				</form>
			</div>
		</div>
	{/if}

	<!-- HISTORY TAB -->
	{#if activeTab === 'history'}
		<!-- Stats -->
		<div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
			<div class="flex flex-col justify-center rounded-xl border border-border-color bg-bg-card p-5 shadow-sm">
				<span class="mb-1 text-sm font-bold uppercase tracking-wider text-text-muted">Total Recaudado Hoy</span>
				<h3 class="text-2xl font-black text-text-primary">{formatCurrency(totalVendidoHoy)}</h3>
			</div>
			<div class="flex flex-col justify-center rounded-xl border border-border-color bg-bg-card p-5 shadow-sm">
				<span class="mb-1 text-sm font-bold uppercase tracking-wider text-text-muted">Ventas Registradas Hoy</span>
				<h3 class="text-2xl font-black text-text-primary">{cantidadVentasHoy}</h3>
			</div>
			<div class="flex flex-col justify-center rounded-xl border border-border-color bg-bg-card p-5 shadow-sm">
				<span class="mb-1 text-sm font-bold uppercase tracking-wider text-text-muted">Ticket Promedio Hoy</span>
				<h3 class="text-2xl font-black text-primario">{formatCurrency(ticketPromedio)}</h3>
			</div>
		</div>

		<!-- Search Sales -->
		<div class="mb-6 rounded-xl border border-border-color bg-bg-card p-4 shadow-sm">
			<div class="relative flex overflow-hidden rounded-lg border border-[rgba(15,30,54,0.15)] bg-white focus-within:border-accent dark:bg-bg-primary">
				<div class="flex items-center pl-4 text-text-muted">
					<Search size={16} />
				</div>
				<input
					type="text"
					bind:value={searchHistoryQuery}
					placeholder="Buscar venta por ID..."
					class="w-full border-none bg-transparent px-3 py-2.5 text-sm text-text-primary outline-none"
				/>
			</div>
		</div>

		<!-- Sales Table -->
		{#if filteredSales.length === 0}
			<div class="flex flex-col items-center justify-center rounded-xl border border-border-color bg-bg-card p-16 text-center shadow-md">
				<FileText size={48} class="mb-3 text-text-muted" />
				<h3 class="font-semibold text-text-primary">No se encontraron ventas</h3>
				<p class="text-sm text-text-secondary max-w-sm mt-1">
					No hay ventas que coincidan con la búsqueda o no se han realizado ventas en el sistema.
				</p>
			</div>
		{:else}
			<div class="overflow-hidden rounded-xl border border-border-color bg-bg-card shadow-md">
				<table class="w-full text-left border-collapse">
					<thead>
						<tr>
							<th class="border-b border-border-color bg-text-primary/4 p-4 text-xs font-bold uppercase tracking-wider text-text-secondary">ID Venta</th>
							<th class="border-b border-border-color bg-text-primary/4 p-4 text-xs font-bold uppercase tracking-wider text-text-secondary">Fecha y Hora</th>
							<th class="border-b border-border-color bg-text-primary/4 p-4 text-xs font-bold uppercase tracking-wider text-text-secondary">Método de Pago</th>
							<th class="border-b border-border-color bg-text-primary/4 p-4 text-xs font-bold uppercase tracking-wider text-text-secondary">Descuento</th>
							<th class="border-b border-border-color bg-text-primary/4 p-4 text-xs font-bold uppercase tracking-wider text-text-secondary">Total</th>
							<th class="w-[100px] border-b border-border-color bg-text-primary/4 p-4 text-right text-xs font-bold uppercase tracking-wider text-text-secondary">Detalles</th>
						</tr>
					</thead>
					<tbody>
						{#each filteredSales as sale (sale.id_venta)}
							<tr class="transition-colors hover:bg-text-primary/[0.01]">
								<td class="border-b border-border-color p-4 font-mono text-xs text-text-primary">
									{sale.id_venta.substring(0, 8)}...
								</td>
								<td class="border-b border-border-color p-4 text-sm text-text-secondary">
									{new Date(sale.fecha_emision).toLocaleString('es-CL')}
								</td>
								<td class="border-b border-border-color p-4">
									<span class="text-xs font-semibold px-2 py-0.5 rounded-full {sale.metodo_pago?.nombre_metodo === 'Tarjeta' ? 'bg-blue-100 text-blue-800' : 'bg-green-100 text-green-800'}">
										{sale.metodo_pago?.nombre_metodo || 'Efectivo'}
									</span>
								</td>
								<td class="border-b border-border-color p-4 text-sm text-danger-color">
									{sale.monto_descuento > 0 ? `-${formatCurrency(sale.monto_descuento)}` : 'Ninguno'}
								</td>
								<td class="border-b border-border-color p-4 text-sm font-bold text-text-primary">
									{formatCurrency(sale.monto_total)}
								</td>
								<td class="border-b border-border-color p-4 text-right">
									<button
										onclick={() => openSaleDetail(sale)}
										class="inline-flex cursor-pointer items-center justify-center rounded-lg border border-border-color bg-text-primary/3 p-2 text-text-secondary transition-all duration-200 hover:border-border-color-hover hover:bg-text-primary/7 hover:text-text-primary"
										title="Ver detalles de venta"
									>
										<ChevronRight size={16} />
									</button>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	{/if}
{/if}

<!-- Detail Modal -->
{#if showDetailModal && selectedSale}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-text-primary/40 p-5 backdrop-blur-[4px]" onclick={() => (showDetailModal = false)} role="presentation">
		<div class="w-full max-w-[600px] overflow-hidden rounded-xl border border-border-color bg-bg-card shadow-lg animate-modal-enter" onclick={(e) => e.stopPropagation()} role="dialog">
			<header class="flex items-center justify-between border-b border-border-color p-5">
				<div>
					<h2 class="text-lg font-bold text-text-primary">Detalle de Venta</h2>
					<p class="text-xs text-text-muted mt-0.5">ID: {selectedSale.id_venta}</p>
				</div>
				<button class="inline-flex cursor-pointer items-center justify-center rounded-lg border border-border-color bg-text-primary/3 p-2 text-text-secondary transition-all duration-200 hover:border-border-color-hover hover:bg-text-primary/7 hover:text-text-primary" onclick={() => (showDetailModal = false)}>&times;</button>
			</header>
			<div class="p-6">
				<!-- Meta info -->
				<div class="grid grid-cols-2 gap-4 mb-6 text-sm border-b border-border-color pb-4">
					<div>
						<p class="text-text-muted font-medium text-xs uppercase">Fecha y Hora</p>
						<p class="font-semibold text-text-primary mt-0.5">
							{new Date(selectedSale.fecha_emision).toLocaleString('es-CL')}
						</p>
					</div>
					<div>
						<p class="text-text-muted font-medium text-xs uppercase">Método de Pago</p>
						<p class="font-semibold text-text-primary mt-0.5">
							{selectedSale.metodo_pago?.nombre_metodo || 'Efectivo'}
						</p>
					</div>
					<div>
						<p class="text-text-muted font-medium text-xs uppercase">ID Empleado</p>
						<p class="font-semibold text-text-primary mt-0.5">{selectedSale.id_empleado}</p>
					</div>
					<div>
						<p class="text-text-muted font-medium text-xs uppercase">ID Caja</p>
						<p class="font-semibold text-text-primary mt-0.5">{selectedSale.id_caja}</p>
					</div>
				</div>

				<!-- Products table -->
				<h4 class="font-bold text-xs uppercase text-text-secondary mb-3">Productos Vendidos</h4>
				<div class="max-h-60 overflow-y-auto mb-6 border border-border-color rounded-lg">
					<table class="w-full text-left border-collapse text-sm">
						<thead>
							<tr class="bg-text-primary/2">
								<th class="p-3 font-semibold text-xs text-text-secondary">Producto</th>
								<th class="p-3 font-semibold text-xs text-text-secondary text-center">Cant.</th>
								<th class="p-3 font-semibold text-xs text-text-secondary text-right font-bold">Total</th>
							</tr>
						</thead>
						<tbody>
							{#if selectedSale.detalles}
								{#each selectedSale.detalles as det}
									<tr>
										<td class="p-3 border-b border-border-color font-medium text-text-primary">
											{getProductDetailName(det.id_producto)}
										</td>
										<td class="p-3 border-b border-border-color text-center text-text-primary">
											{det.cantidad}
										</td>
										<td class="p-3 border-b border-border-color text-right font-semibold text-text-primary">
											{formatCurrency(det.monto_final)}
										</td>
									</tr>
								{/each}
							{:else}
								<tr>
									<td colspan="3" class="p-4 text-center text-text-muted">
										No hay detalles de productos para esta venta
									</td>
								</tr>
							{/if}
						</tbody>
					</table>
				</div>

				<!-- Money Breakdown -->
				<div class="flex flex-col gap-2 border-t border-border-color pt-4 text-sm font-semibold">
					{#if selectedSale.monto_descuento > 0}
						<div class="flex justify-between text-text-secondary">
							<span>Descuento Aplicado:</span>
							<span class="text-danger-color">-{formatCurrency(selectedSale.monto_descuento)}</span>
						</div>
					{/if}
					{#if selectedSale.pago > 0}
						<div class="flex justify-between text-text-secondary">
							<span>Pago recibido:</span>
							<span>{formatCurrency(selectedSale.pago)}</span>
						</div>
						<div class="flex justify-between text-text-secondary">
							<span>Vuelto entregado:</span>
							<span>{formatCurrency(selectedSale.vuelto)}</span>
						</div>
					{/if}
					<div class="flex justify-between text-base font-black border-t border-dashed pt-2 mt-2">
						<span class="text-text-primary">Monto Total de Venta:</span>
						<span class="text-primario">{formatCurrency(selectedSale.monto_total)}</span>
					</div>
				</div>
			</div>
			<footer class="flex justify-end border-t border-border-color bg-text-primary/2 p-4 px-6">
				<button type="button" class="inline-flex cursor-pointer items-center justify-center rounded-lg border border-border-color bg-bg-secondary px-5 py-2.5 text-sm font-semibold text-text-primary transition-all duration-200 hover:bg-text-primary/5" onclick={() => (showDetailModal = false)}>
					Cerrar
				</button>
			</footer>
		</div>
	</div>
{/if}

<!-- Caja Config Configuration Modal -->
{#if showCajaConfigModal}
	<div class="fixed inset-0 z-[100] flex items-center justify-center bg-black/60 p-5 backdrop-blur-sm">
		<div class="w-full max-w-[450px] overflow-hidden rounded-xl border border-border-color bg-bg-card p-6 shadow-2xl animate-modal-enter text-center">
			<div class="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-full bg-primario/10 text-primario">
				<Laptop size={28} />
			</div>
			<h3 class="text-lg font-bold text-text-primary">Asociar Caja a este Terminal</h3>
			<p class="text-xs text-text-secondary mt-1.5 max-w-xs mx-auto">
				Por seguridad y control de arqueo, debes indicar a qué caja corresponde este equipo antes de realizar ventas. Esta selección quedará guardada en este dispositivo.
			</p>
			
			<div class="my-6">
				<label for="cajaConfigSelect" class="block text-left text-xs font-bold uppercase tracking-wider text-text-secondary mb-2">Selecciona la caja activa</label>
				<select id="cajaConfigSelect" class="w-full rounded-lg border border-border-color bg-bg-card p-3 text-sm text-text-primary outline-none focus:border-primario" onchange={(e) => handleCajaSelection((e.target as HTMLSelectElement).value)}>
					<option value="" disabled selected={!selectedCajaId}>-- Selecciona una caja --</option>
					{#each cajas as c}
						<option value={c.id_caja}>{c.nombre} ({c.ubicacion})</option>
					{/each}
				</select>
			</div>

			<p class="text-[10px] text-text-muted">
				Podrás cambiar esta asignación más tarde desde la esquina superior derecha si es necesario.
			</p>
		</div>
	</div>
{/if}

<!-- Manual Code Input Modal (For damaged labels) -->
{#if showManualCodeModal}
	<div class="fixed inset-0 z-[100] flex items-center justify-center bg-black/50 p-5 backdrop-blur-sm" onclick={() => (showManualCodeModal = false)} role="presentation">
		<div class="w-full max-w-[400px] overflow-hidden rounded-xl border border-border-color bg-bg-card shadow-2xl animate-modal-enter" onclick={(e) => e.stopPropagation()} role="dialog">
			<header class="flex items-center justify-between border-b border-border-color p-4">
				<h3 class="font-bold text-text-primary text-sm">Ingreso Manual de Código</h3>
				<button class="text-text-muted hover:text-text-primary text-lg" onclick={() => (showManualCodeModal = false)}>&times;</button>
			</header>
			<form onsubmit={handleManualBarcodeSubmit}>
				<div class="p-5">
					<label for="manualBarcode" class="block text-xs font-bold uppercase tracking-wider text-text-secondary mb-2">Código de barras</label>
					<input
						id="manualBarcode"
						type="text"
						bind:value={manualBarcodeValue}
						placeholder="Escribe el código numérico..."
						class="w-full rounded-lg border border-border-color bg-bg-card p-3 text-sm text-text-primary outline-none focus:border-primario"
						required
						autofocus
					/>
				</div>
				<footer class="flex justify-end gap-3 border-t border-border-color bg-text-primary/2 p-3 px-5">
					<button type="button" class="px-4 py-2 rounded-lg border border-border-color text-xs font-semibold text-text-secondary bg-bg-secondary hover:bg-text-primary/5" onclick={() => (showManualCodeModal = false)}>Cancelar</button>
					<button type="submit" class="px-4 py-2 rounded-lg bg-primario text-white text-xs font-bold hover:bg-primario-hover shadow-sm">Agregar Producto</button>
				</footer>
			</form>
		</div>
	</div>
{/if}

<!-- Scanner Modal -->
{#if modoEscaneo}
	<div class="fixed inset-0 z-[110] flex items-center justify-center bg-black/50 p-5 backdrop-blur-sm" onclick={() => (modoEscaneo = false)} role="presentation">
		<div class="w-full max-w-[460px] overflow-hidden rounded-xl border border-border-color bg-bg-card shadow-2xl animate-modal-enter" onclick={(e) => e.stopPropagation()} role="dialog">
			<header class="flex items-center justify-between border-b border-border-color p-4">
				<h3 class="font-bold text-text-primary text-sm">Escanear Código de Barras</h3>
				<button class="text-text-muted hover:text-text-primary text-lg" onclick={() => (modoEscaneo = false)}>&times;</button>
			</header>
			<div class="p-6 flex flex-col items-center justify-center">
				<Scanner onScan={manejarEscaneo} onClose={() => (modoEscaneo = false)} />
			</div>
		</div>
	</div>
{/if}
