const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

export interface Empleado {
    id?: string;
    id_empleado?: string;
    nombre: string;
    usuario: string;
    rol: string;
    es_admin: boolean;
}

export interface Categoria {
    id_categoria: string;
    nombre_categoria: string;
    descripcion?: string;
}

export interface Producto {
    id_producto: string;
    nombre: string;
    descripcion: string;
    stock: number;
    stock_minimo: number;
    precio: number;
    marca: string;
    unidad: string;
    codigo_barras: string;
    estado: boolean;
    id_categoria: string;
    categoria?: Categoria;
}

export interface Promocion {
    id_promocion?: string;
    producto_id?: string | null;
    productos_combo?: string[];
    tipo: string;
    lleva?: number;
    paga?: number;
    descuento?: number;
    fecha_inicio?: string | null;
    fecha_fin?: string | null;
}

export async function apiFetch(endpoint: string, options: RequestInit = {}) {
    const headers = {
        'Content-Type': 'application/json',
        ...options.headers
    };

    const response = await fetch(`${API_URL}${endpoint}`, {
        ...options,
        headers,
        credentials: 'include'
    });

    if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        
        const mensajePrincipal = errorData.detalle || errorData.error || errorData.Error || errorData.mensaje || 'Error de conexión con el servidor';
        
        const errorObj = new Error(mensajePrincipal);
        
        Object.assign(errorObj, errorData);
        
        if (errorData.errores) {
            (errorObj as any).errors = errorData.errores;
        }

        throw errorObj;
    }

    return response.json();
}

export async function login(usuario: string, contrasena: string): Promise<Empleado> {
	return apiFetch('/auth/login', {
		method: 'POST',
		body: JSON.stringify({ usuario, contrasena })
	});
}

export async function logout(): Promise<void> {
	return apiFetch('/auth/logout', {
		method: 'POST'
	});
}

export async function checkSession(): Promise<Empleado> {
	return apiFetch('/auth/me');
}

//promociones 
export async function obtenerPromociones(): Promise<Promocion[]> {
	return apiFetch('/promociones/');
}

export async function crearPromocion(data: Partial<Promocion>) {
	return apiFetch('/promociones', {
		method: 'POST',
		body: JSON.stringify(data)
	});
}

export async function eliminarPromocion(id: string) {
	return apiFetch(`/promociones/${id}`, {
		method: 'DELETE'
	});
}

export async function actualizarPromocion(id: string, payload: Partial<Promocion>) {
	return apiFetch(`/promociones/${id}`, {
		method: 'PATCH',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify(payload)
	});
}

export async function obtenerProductos() {
	return apiFetch('/inventario/productos');
}

export const apiClientes = {
	getAll: () => apiFetch('/clientes/'),
	getById: (id: string) => apiFetch(`/clientes/${id}`),
	getByRut: (rut: string) => apiFetch(`/clientes/rut/${rut}`),
	searchByNombre: (nombre: string) => apiFetch(`/clientes/search/nombre?nombre=${encodeURIComponent(nombre)}`),
	create: (data: { nombre: string; rut: string; telefono: string }) => apiFetch('/clientes', {
		method: 'POST',
		body: JSON.stringify(data)
	}),
	update: (id: string, data: { nombre?: string; rut?: string; telefono?: string }) => apiFetch(`/clientes/${id}`, {
		method: 'PATCH',
		body: JSON.stringify(data)
	}),
	delete: (id: string) => apiFetch(`/clientes/${id}`, {
		method: 'DELETE'
	})
};

export const apiVentas = {
	getAll: () => apiFetch('/ventas'),
	getById: (id: string) => apiFetch(`/ventas/${id}`),
	create: (data: {
		id_caja: string;
		id_metodo: string;
		pago: number;
		vuelto: number;
		monto_total: number;
		monto_descuento: number;
		detalles: Array<{
			id_producto: string;
			cantidad: number;
			monto_final: number;
		}>;
	}) => apiFetch('/ventas', {
		method: 'POST',
		body: JSON.stringify(data)
	})
};


export const apiCategorias = {
    getAll: () => apiFetch('/inventario/categorias'),
    create: (data: { nombre_categoria: string }) => apiFetch('/inventario/categorias', {
        method: 'POST',
        body: JSON.stringify(data)
    }),
    update: (id: string, data: { nombre_categoria: string }) => apiFetch(`/inventario/categorias/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data)
    }),
    delete: (id: string) => apiFetch(`/inventario/categorias/${id}`, {
        method: 'DELETE'
    })
};

export const apiProductos = {
    getAll: () => apiFetch('/inventario/productos'),
    getById: (id: string) => apiFetch(`/inventario/productos/${id}`),
    create: (data: Partial<Producto>) => apiFetch('/inventario/productos', {
        method: 'POST',
        body: JSON.stringify(data)
    }),
    update: (id: string, data: Partial<Producto>) => apiFetch(`/inventario/productos/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data)
    }),
    delete: (id: string) => apiFetch(`/inventario/productos/${id}`, {
        method: 'DELETE'
    })
};

export const apiEmpleados = {
    getAll: () => apiFetch('/empleados'),
    getById: (id: string) => apiFetch(`/empleados/${id}`),
    create: (data: Partial<Producto>) => apiFetch('/empleados', {
        method: 'POST',
        body: JSON.stringify(data)
    }),
    update: (id: string, data: Partial<Producto>) => apiFetch(`/empleados/${id}`, {
        method: 'PATCH',
        body: JSON.stringify(data)
    }),
    delete: (id: string) => apiFetch(`/empleados/${id}`, {
        method: 'DELETE'
    })
};

export const apiRoles = {
    getAll: () => apiFetch('/empleados/roles'),
    create: (data: any) => apiFetch('/empleados/roles', { method: 'POST', body: JSON.stringify(data) }),
    update: (id: string, data: any) => apiFetch(`/empleados/roles/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
    delete: (id: string) => apiFetch(`/empleados/roles/${id}`, { method: 'DELETE' })
};

export const apiCajas = {
    getAll: () => apiFetch('/cajas'),
    getById: (id: string) => apiFetch(`/cajas/${id}`),
    create: (data: any) => apiFetch('/cajas', { method: 'POST', body: JSON.stringify(data) }),
    update: (id: string, data: any) => apiFetch(`/cajas/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
    delete: (id: string) => apiFetch(`/cajas/${id}`, { method: 'DELETE' })
};