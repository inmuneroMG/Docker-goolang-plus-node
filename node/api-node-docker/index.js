require("dotenv").config();
const express = require('express');
const jwt = require("jsonwebtoken");
const bodyParser = require("body-parser");
const cors = require("cors");
const axios = require('axios');
//import axios from 'axios';
const app = express();
//const appaxios = axios();
const PORT = process.env.PORT || 3000;
const SECRET_KEY = process.env.SECRET_KEY || "clave_secreta";

app.use(express.json());

app.get('/', (req, res) => {
    res.json({ message: 'API en Node.js con Express dentro de Docker!' });
});
//publico login
app.post("/login", (req, res) => {
    const { username, password } = req.body;

    // 🔹 Validación simple (esto debería verificarse en una base de datos)
    if (username !== "admin" || password !== "1234") {
        return res.status(401).json({ error: "Credenciales incorrectas" });
    }

    // 🔹 Generar token JWT válido por 1 hora
    const token = jwt.sign({ username, password }, SECRET_KEY, { expiresIn: "5h" });

    res.json({ token });
});
//middleware
const verifyToken = (req, res, next) => {
    const token = req.headers["authorization"];

    if (!token) {
        return res.status(403).json({ error: "Token requerido" });
    }

    try {
        const decoded = jwt.verify(token.split(" ")[1], SECRET_KEY);
        req.user = decoded;
        next();
    } catch (error) {
        return res.status(401).json({ error: "Token inválido" });
    }
};
app.get("/perfil", verifyToken, (req, res) => {
    res.json({ message: "Accediste al perfil", user: req.user.username });
});
app.post('/api/matriz', verifyToken, async (req, res) => {
    try {
        const { matriz } = req.body; // Extraemos la matriz del body
        //var Matriz ={"Matriz":matriz}
        var Matriz=JSON.stringify({"Matriz":matriz});
        console.log("Matriz recibida:", matriz);
        // Verificamos que sea un array de arrays de números
        if (!Array.isArray(matriz) || !matriz.every(row => Array.isArray(row) && row.every(num => typeof num === 'number'))) {
            return res.status(400).json({ error: 'El formato debe ser un array de arrays de números' });
        }
        console.log("tokenGOLANG INI");
        console.log("username", req.user.username);
        console.log("password", req.user.password);
        var loginJWT=JSON.stringify({"username": req.user.username,"password": req.user.password});
        
        console.log("loginJWTs", loginJWT);

        const responseJWT = await axios.post('http://host.docker.internal:8085/login',loginJWT, {
            headers: { 'Content-Type': 'application/json' }});
        
            const { token } = responseJWT.data;
            console.log("tokenGOLANG FIN",token);
        /*const {data} = axios.post('http://localhost:8085/api/test',JSON.stringify(jsonData), {
            headers: { "Content-Type": "application/json" }
          })*/
            const response = await axios.post('http://host.docker.internal:8085/api/fx',Matriz, {
                headers: { 'Content-Type': 'application/json',"Authorization": `Bearer ${token}` }});
            //.then(response => res.json(response.data))
            //.catch(error => res.status(501).json({ error: 'Error procesando la solicitud' }));
            
            const { Q, R, MatrizB } = response.data;

            //console.log("message",message);
            console.log("response",response.data);
            console.log("Q",Q,"|",esDiagonal(Q));
            console.log("R",R,"|",esDiagonal(R));
            console.log("MatrizB",MatrizB,"|",esDiagonal(MatrizB));
    
            const matrices=Q.concat(R).concat(MatrizB)
    var responseJson={
        maxValor:obtenerMaximo(matrices),
        minValor:obtenerMinimo(matrices),
        promedioValor:calcularPromedio(matrices),
        sumaValores:calcularSuma(matrices),
        unaOMasMatricesDiagonal:(esDiagonal(Q)||esDiagonal(R)||esDiagonal(MatrizB) ? "Sí" : "No"),
        MatrizInvertida:MatrizB,
        MatrizQ:Q,
        MatrizR:R
    }

    console.log("maxValor",responseJson);
    // Verificar si alguna de las matrices es diagonal
    //console.log("¿Alguna matriz es diagonal?", matrices.some(esDiagonal) ? "Sí" : "No");
            res.json(responseJson);

    } catch (error) {
        res.status(500).json({ error: 'Error procesando la solicitud' });
    }
});

// Función para obtener el valor máximo de todas las matrices
function obtenerMaximo(matrices) {
    return Math.max(...matrices.flat(2));
  }
  
// Función para obtener el valor mínimo de todas las matrices
  function obtenerMinimo(matrices) {
    return Math.min(...matrices.flat(2));
  }
  
// Función para calcular el promedio de todos los valores
  function calcularPromedio(matrices) {
    const valores = matrices.flat(2);
    return valores.reduce((sum, val) => sum + val, 0) / valores.length;
  }
  
// Función para calcular la suma total de los valores
  function calcularSuma(matrices) {
    return matrices.flat(2).reduce((sum, val) => sum + val, 0);
  }
  
// Función para verificar si una matriz es diagonal
  function esDiagonal(matriz) {
    return matriz.every((fila, i) => fila.every((val, j) => (i === j ? true : val === 0)));
  }
app.listen(PORT, () => {
    console.log(`Servidor corriendo en http://localhost:${PORT}`);
});