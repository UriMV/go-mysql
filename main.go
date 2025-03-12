package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var err error

// Configurar el almacenamiento de sesiones
var store = sessions.NewCookieStore([]byte("clave-secreta")) // Cambia "clave-secreta" por una clave segura

type User struct {
	Username string
	Password string
}

type Employee struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// Inicializar la conexión a la base de datos
func initDB() {
	// Cambia la conexión con los datos correctos de Railway
	db, err = sql.Open("mysql", "root:fFsQlmvuWNvJBtOfwFtdbRSwgFzTPduc@tcp(switchback.proxy.rlwy.net:43933)/railway")
	if err != nil {
		log.Fatal(err)
	}

	// Verificar la conexión
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Conexión a la base de datos establecida")
}

// Manejador para la página de login
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Mostrar el formulario de login
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, "No se pudo cargar la página de login", http.StatusInternalServerError)
			log.Println("Error al cargar login.html:", err)
			return
		}
		tmpl.Execute(w, nil)
	} else if r.Method == "POST" {
		// Procesar el formulario de login
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		var user User
		// Buscar el usuario en la base de datos
		err := db.QueryRow("SELECT username, password FROM users WHERE username = ?", username).Scan(&user.Username, &user.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				// Usuario no encontrado
				fmt.Fprintf(w, `<script>alert("Usuario o contraseña incorrectos"); window.location.href = "/login";</script>`)
			} else {
				log.Fatal(err)
			}
		} else {
			// Verificar la contraseña hasheada
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				// Contraseña incorrecta
				fmt.Fprintf(w, `<script>alert("Usuario o contraseña incorrectos"); window.location.href = "/login";</script>`)
			} else {
				// Crear una sesión para el usuario
				session, _ := store.Get(r, "session-name")
				session.Values["username"] = username
				session.Save(r, w)

				// Redirigir a la página principal
				http.Redirect(w, r, "/index", http.StatusFound)
			}
		}
	}
}

// Manejador para la página de registro
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Mostrar el formulario de registro
		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
			http.Error(w, "No se pudo cargar la página de registro", http.StatusInternalServerError)
			log.Println("Error al cargar register.html:", err)
			return
		}
		tmpl.Execute(w, nil)
	} else if r.Method == "POST" {
		// Procesar el formulario de registro
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")

		// Validar que las contraseñas coincidan
		if password != confirmPassword {
			fmt.Fprintf(w, `<script>alert("Las contraseñas no coinciden"); window.location.href = "/register";</script>`)
			return
		}

		// Hashear la contraseña
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error al hashear la contraseña", http.StatusInternalServerError)
			log.Println("Error al hashear la contraseña:", err)
			return
		}

		// Insertar el nuevo usuario en la base de datos
		_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, string(hashedPassword))
		if err != nil {
			// Si hay un error (por ejemplo, nombre de usuario duplicado)
			fmt.Fprintf(w, `<script>alert("Error al registrar el usuario: el nombre de usuario ya existe"); window.location.href = "/register";</script>`)
			log.Println("Error al registrar el usuario:", err)
			return
		}

		// Redirigir al login después del registro exitoso
		fmt.Fprintf(w, `<script>alert("Registro exitoso"); window.location.href = "/login";</script>`)
	}
}

// Manejador para la página principal
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener la sesión del usuario
	session, _ := store.Get(r, "session-name")

	// Verificar si el usuario está autenticado
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Mostrar la página principal
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "No se pudo cargar la página principal", http.StatusInternalServerError)
		log.Println("Error al cargar index.html:", err)
		return
	}
	tmpl.Execute(w, map[string]string{
		"Username": username,
	})
}

// Manejador para la página de perfil
func profileHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener la sesión del usuario
	session, _ := store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Mostrar la página de perfil
	tmpl, err := template.ParseFiles("templates/profile.html")
	if err != nil {
		http.Error(w, "No se pudo cargar la página de perfil", http.StatusInternalServerError)
		log.Println("Error al cargar profile.html:", err)
		return
	}
	tmpl.Execute(w, map[string]string{
		"Username": username,
	})
}

// Manejador para cerrar sesión
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener la sesión del usuario
	session, _ := store.Get(r, "session-name")

	// Eliminar la sesión
	session.Values["username"] = ""
	session.Options.MaxAge = -1 // Borrar la cookie
	session.Save(r, w)

	// Redirigir al login
	http.Redirect(w, r, "/login", http.StatusFound)
}

// Obtener todos los empleados
func getEmployees(w http.ResponseWriter, r *http.Request) {
	var employees []Employee
	rows, err := db.Query("SELECT id, name, email, role FROM employees")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var emp Employee
		err := rows.Scan(&emp.ID, &emp.Name, &emp.Email, &emp.Role)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		employees = append(employees, emp)
	}
	json.NewEncoder(w).Encode(employees)
}

// Obtener un empleado por su ID
func getEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var emp Employee
	err := db.QueryRow("SELECT id, name, email, role FROM employees WHERE id = ?", id).Scan(&emp.ID, &emp.Name, &emp.Email, &emp.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(emp)
}

// Crear un nuevo empleado
func createEmployee(w http.ResponseWriter, r *http.Request) {
	var emp Employee
	json.NewDecoder(r.Body).Decode(&emp)

	result, err := db.Exec("INSERT INTO employees (name, email, role) VALUES (?, ?, ?)", emp.Name, emp.Email, emp.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	emp.ID = int(id)
	json.NewEncoder(w).Encode(emp)
}

// Actualizar un empleado
func updateEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var emp Employee
	json.NewDecoder(r.Body).Decode(&emp)

	_, err := db.Exec("UPDATE employees SET name = ?, email = ?, role = ? WHERE id = ?", emp.Name, emp.Email, emp.Role, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	emp.ID = id
	json.NewEncoder(w).Encode(emp)
}

// Eliminar un empleado
func deleteEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	_, err := db.Exec("DELETE FROM employees WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	// Inicializar la base de datos
	initDB()
	defer db.Close()

	// Configurar el enrutador
	router := mux.NewRouter()

	// Configurar las rutas de la API
	router.HandleFunc("/api/employees", getEmployees).Methods("GET")
	router.HandleFunc("/api/employees/{id}", getEmployee).Methods("GET")
	router.HandleFunc("/api/employees", createEmployee).Methods("POST")
	router.HandleFunc("/api/employees/{id}", updateEmployee).Methods("PUT")
	router.HandleFunc("/api/employees/{id}", deleteEmployee).Methods("DELETE")

	// Configurar los manejadores para las rutas
	router.HandleFunc("/login", loginHandler)
	router.HandleFunc("/register", registerHandler)
	router.HandleFunc("/index", indexHandler)
	router.HandleFunc("/profile", profileHandler)
	router.HandleFunc("/logout", logoutHandler)
	router.HandleFunc("/employees", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/employees.html")
	}).Methods("GET")

	// Servir archivos estáticos (CSS y JS)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Iniciar el servidor
	fmt.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
