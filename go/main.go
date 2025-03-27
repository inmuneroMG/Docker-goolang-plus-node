package main

import (
    "github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
    "fmt"
    "gonum.org/v1/gonum/mat"
    "log"
	"time"
)
//estructura de matriz
type MatrizRequest struct {
	Matriz [][]float64 `json:"Matriz"`
}
type TestRequest struct {
	Test [][]int `json:"Test"`
}
var secretKey = []byte("clave-secreta")

func matrizToDense(data [][]float64) *mat.Dense {
	rows := len(data)
	cols := len(data[0])
	flatData := make([]float64, 0, rows*cols)
	
	//[][]array a []array
	for _, row := range data {
		flatData = append(flatData, row...)
	}

	return mat.NewDense(rows, cols, flatData)
}
func denseToMatriz(m *mat.Dense) [][]float64 {
    rows, cols := m.Dims()
    result := make([][]float64, rows)
    for i := range result {
        result[i] = make([]float64, cols)
        for j := range result[i] {
            result[i][j] = m.At(i, j)
        }
    }
    return result
}
func main() {
    // Crear aplicación Fiber
	
    app := fiber.New(fiber.Config{
        AppName: "API",
    })

    // Middlewares
    //app.Use(logger.New())
    //app.Use(cors.New())

    // Ruta de inicio
    app.Get("/", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "mensaje": "¡Bienvenido a la API de GO!",
        })
    })

    // Grupo de rutas
    appGroup := app.Group("/api", jwtMiddleware())

	app.Post("/login", login)

	appGroup.Get("/perfil", perfil)
	appGroup.Get("/datos", datos)
	// Funcion para factorizar una matriz
    appGroup.Post("/test", func(c *fiber.Ctx) error {
		var data TestRequest
		if err := c.BodyParser(&data); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Error al parsear JSON"})
		}
        return c.JSON(fiber.Map{
            "message": "Ok",
        })
	})
	appGroup.Post("/fx", func(c *fiber.Ctx) error {
		
		var data MatrizRequest
		
		//fmt.Printf("%v\n\n", mat.Formatted(&Q, mat.Prefix("    "), mat.Squeeze()))
		
		// Parsear el JSON del cuerpo
		if err := c.BodyParser(&data); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Error al parsear JSON"})
		}

		// Lista de errores de validación
		errores := make(map[string]string)

		// Validar que la matriz no esté vacía
		if len(data.Matriz) == 0 {
			errores["Matriz"] = "La matriz no puede estar vacía"
		}

		// Si hay errores, devolverlos en la respuesta
		if len(errores) > 0 {
			return c.Status(400).JSON(fiber.Map{
				"error":   "Parámetros inválidos",
				"detalles": errores,
			})
		}
		
		rows := len(data.Matriz)
		cols := len(data.Matriz[0])

		// Crear una nueva matriz transpuesta
		traspuesta := make([][]float64, cols)
		for i := range traspuesta {
			traspuesta[i] = make([]float64, rows)
		}
		
		// Llenar la matriz transpuesta
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				traspuesta[j][i] = data.Matriz[i][j]
			}
		}    
		A := matrizToDense(data.Matriz)

		// Crear matrices para almacenar Q y R
		var Q, R mat.Dense
	
		// Factorización QR
		qr := new(mat.QR)
		qr.Factorize(A)
		qr.QTo(&Q)
		qr.RTo(&R)
		
		// Imprimir resultados

		var Qm, Rm [][]float64
		Qm=denseToMatriz(&Q)
		Rm=denseToMatriz(&R)
		fmt.Println("Matriz Q:")
		fmt.Printf("%v\n\n", mat.Formatted(&Q, mat.Prefix("    "), mat.Squeeze()))
	
		fmt.Println("Matriz R:")
		fmt.Printf("%v\n", mat.Formatted(&R, mat.Prefix("    "), mat.Squeeze()))
        return c.JSON(fiber.Map{
            "Q":  Qm,
            "R":  Rm,
            "MatrizB":  traspuesta,
        })
    })

    // Iniciar servidor
    log.Fatal(app.Listen(":8085"))
}
func login(c *fiber.Ctx) error {
//Simulación de usuario autenticado (hardcoded)
	user := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	// Validación simple (reemplázalo con una base de datos)
	if user.Username != "admin" || user.Password != "1234" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Credenciales incorrectas"})
	}

//Crear token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 5).Unix(), // Expira en 1 hora
	})

//Firmar el token
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al generar el token"})
	}

//Devolver el token
	return c.JSON(fiber.Map{"token": tokenString})
}

func perfil(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Accediste al perfil", "user": "admin"})
}

func datos(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Estos son tus datos", "data": []int{1, 2, 3, 4}})
}
func jwtMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token requerido"})
		}

		// Quitar el prefijo "Bearer "
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// Validar el token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token inválido"})
		}

		// Continuar con la petición
		return c.Next()
	}
}