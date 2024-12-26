package contract

// FunctionContract define el contrato que deben seguir las funciones
type FunctionContract struct {
	// Requisitos de la imagen Docker
	Requirements struct {
		// La imagen debe aceptar un único argumento string (JSON)
		Input string `json:"input"`

		// La función debe escribir su resultado en stdout
		// El resultado debe ser un string (JSON)
		Output string `json:"output"`

		// Logs y errores deben ir a stderr
		Logs string `json:"logs"`
	}
}
