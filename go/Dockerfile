FROM golang:1.23-alpine

# Establecer directorio de trabajo
WORKDIR /app

# Inicializar módulo
RUN go mod init api-usuarios
# Añadir dependencia
#RUN go get github.com/gofiber/fiber/v2
RUN go get gonum.org/v1/gonum/mat

# Copiar archivos de código fuente
COPY . .

# Descargar dependencias
#RUN go mod download

RUN go mod tidy -e
# Compilar la aplicación
RUN go build -o servidor

# Exponer puerto
EXPOSE 8085

# Comando para ejecutar
CMD ["./servidor"]