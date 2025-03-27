// src/index.js
const express = require('express');
const bodyParser = require('body-parser');

const app = express();
const PORT = process.env.PORT || 3000;

// Middleware para parsear JSON
app.use(bodyParser.json());

// Ruta de ejemplo
app.get('/', (req, res) => {
  res.json({ mensaje: '¡Bienvenido a mi API!' });
});

// Ruta para crear un recurso (POST)
app.post('/usuarios', (req, res) => {
  const { nombre, email } = req.body;
  
  // Aquí normalmente guardarías en una base de datos
  res.status(201).json({
    mensaje: 'Usuario creado',
    usuario: { nombre, email }
  });
});

// Ruta para obtener usuarios (GET)
app.get('/usuarios', (req, res) => {
  // En un caso real, recuperarías de una base de datos
  const usuarios = [
    { id: 1, nombre: 'Juan' },
    { id: 2, nombre: 'María' }
  ];
  
  res.json(usuarios);
});

// Iniciar servidor
app.listen(PORT, () => {
  console.log(`Servidor corriendo en puerto ${PORT}`);
});
