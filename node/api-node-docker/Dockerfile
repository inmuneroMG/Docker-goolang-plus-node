# Usar una imagen oficial de Node.js
FROM node:16.5.0

# Establecer el directorio de trabajo en el contenedor
WORKDIR /

# Copiar package.json y package-lock.json antes de instalar dependencias
COPY package*.json ./

RUN npm config set strict-ssl false
# Instalar dependencias
RUN npm install dotenv --save-dev
RUN npm install

# Copiar el código fuente
COPY . .

# Exponer el puerto 3000
EXPOSE 3000

# Comando para ejecutar la API
CMD ["node", "index.js"]