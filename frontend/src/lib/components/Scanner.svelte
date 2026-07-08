<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { Html5Qrcode, Html5QrcodeSupportedFormats } from 'html5-qrcode';

	// Svelte 5 props
	let { onScan, onClose } = $props<{ 
		onScan: (code: string) => void; 
		onClose?: () => void 
	}>();

	let scanner: Html5Qrcode | null = null;
	const readerElement = 'qr-reader';
	let errorMessage = $state('');

	onMount(() => {
		// Inicializar el escáner
		scanner = new Html5Qrcode(readerElement);

		scanner.start(
			{ facingMode: 'environment' }, // Restricción de cámara (trasera/por defecto)
			{
				fps: 20, // Cuadros por segundo a analizar para la detección
				qrbox: (width, height) => {
					// Área de escaneo horizontal dinámica para códigos de barra
					const w = Math.round(width * 0.75);
					const h = Math.round(height * 0.4);
					return { width: Math.max(w, 250), height: Math.max(h, 150) };
				},
				videoConstraints: {
					width: { ideal: 1280 },
					height: { ideal: 720 }
				},
				formatsToSupport: [
					Html5QrcodeSupportedFormats.EAN_13,
					Html5QrcodeSupportedFormats.EAN_8,
					Html5QrcodeSupportedFormats.CODE_128,
					Html5QrcodeSupportedFormats.QR_CODE
				]
			},
			(decodedText) => {
				// Éxito: Se detectó un código
				if (scanner && scanner.isScanning) {
					scanner.stop().then(() => {
						onScan(decodedText); // Notificar al componente padre
					}).catch(console.error);
				}
			},
			(error) => {
				// Omitir errores de escaneo de cuadros vacíos para no saturar consola
			}
		).catch((err) => {
			errorMessage = `No se pudo acceder a la cámara: ${err}`;
			console.error(err);
		});
	});

	// Limpieza al desmontar el componente
	onDestroy(() => {
		if (scanner && scanner.isScanning) {
			scanner.stop().catch(console.error);
		}
	});
</script>

<div class="scanner-container border border-border-color rounded-xl bg-bg-card p-4 shadow-md">
	{#if errorMessage}
		<p class="text-danger-color text-center font-bold text-xs mb-3">{errorMessage}</p>
	{/if}
	
	<div id={readerElement} class="rounded-lg overflow-hidden border border-border-color bg-black"></div>
	
	{#if onClose}
		<button 
			type="button" 
			onclick={onClose} 
			class="mt-4 w-full cursor-pointer bg-danger-color hover:bg-danger-color/90 text-white font-bold py-2 px-4 rounded-lg text-xs transition-colors"
		>
			Cancelar Escaneo
		</button>
	{/if}
</div>

<style>
	.scanner-container {
		width: 100%;
		max-width: 420px;
		margin: 0 auto;
	}
</style>
