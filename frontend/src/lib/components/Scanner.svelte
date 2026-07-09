<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { Html5Qrcode, Html5QrcodeSupportedFormats } from 'html5-qrcode';

	// Svelte 5 props
	let { onScan, onClose } = $props<{
		onScan: (code: string) => void;
		onClose?: () => void;
	}>();

	let scanner: Html5Qrcode | null = null;
	const readerElement = 'qr-reader';
	let errorMessage = $state('');

	onMount(() => {
		// Inicializar el escáner indicando los formatos de código a soportar
		scanner = new Html5Qrcode(readerElement, {
			formatsToSupport: [
				Html5QrcodeSupportedFormats.EAN_13,
				Html5QrcodeSupportedFormats.EAN_8,
				Html5QrcodeSupportedFormats.CODE_128,
				Html5QrcodeSupportedFormats.QR_CODE
			]
		} as any);

		scanner
			.start(
				{ facingMode: 'environment' }, // Restricción de cámara trasera
				{
					fps: 30, // Mayor cantidad de cuadros para acelerar detección
					qrbox: (width: number, height: number) => {
						// Área más ajustada para asegurar enfoque central
						const w = Math.round(width * 0.8);
						const h = Math.round(height * 0.4);
						return { width: Math.max(w, 250), height: Math.max(h, 150) };
					},
					experimentalFeatures: {
						// Usar API nativa del teléfono si está disponible (muuucho más rápido en Android)
						useBarCodeDetectorIfSupported: true
					},
					videoConstraints: {
						// Resolución más alta y forzar auto-enfoque continuo
						width: { ideal: 1920 },
						height: { ideal: 1080 },
						advanced: [{ focusMode: 'continuous' }] as any
					} as any
				} as any,
				(decodedText) => {
					// Éxito: Se detectó un código
					if (scanner && scanner.isScanning) {
						scanner
							.stop()
							.then(() => {
								onScan(decodedText); // Notificar al componente padre
							})
							.catch(console.error);
					}
				},
				(error) => {
					// Omitir errores de escaneo de cuadros vacíos para no saturar consola
				}
			)
			.catch((err) => {
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

	<div
		id={readerElement}
		class="rounded-lg overflow-hidden border border-border-color bg-black"
	></div>

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
