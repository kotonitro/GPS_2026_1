package clientes

type Cliente struct {
	ID       string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id_cliente"`
	Nombre   string `json:"nombre"`
	Rut      string `json:"rut"`
	Telefono string `json:"telefono"`
}
