<script lang="ts">
    import { page } from '$app/stores';

    const status = $page.status || 404;
    const mensaje = $page.error?.message || 'Página no encontrada';
</script>

<div style="text-align: center; margin-top: 20vh; font-family: sans-serif;">
    <h1 style="font-size: 8rem; margin: 0;">{status}</h1>
    <h2>{mensaje}</h2>
    <br>
    <a href="/" style="color: blue; text-decoration: underline;">Volver al inicio</a>
</div>
