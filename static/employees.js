document.addEventListener('DOMContentLoaded', function () {
    const addEmployeeBtn = document.getElementById('addEmployeeBtn');
    const modal = document.getElementById('employeeModal');
    const closeModal = document.querySelector('.close');
    const form = document.getElementById('employeeForm');
    const tableBody = document.querySelector('#employeeTable tbody');

    // Mostrar modal al hacer clic en el botón de agregar empleado
    addEmployeeBtn.addEventListener('click', function () {
        document.getElementById('modalTitle').textContent = 'Agregar Empleado';
        form.reset();
        modal.style.display = 'block';
    });

    // Ocultar modal al hacer clic en la 'x'
    closeModal.addEventListener('click', function () {
        modal.style.display = 'none';
    });

    // Ocultar modal al hacer clic fuera del modal
    window.addEventListener('click', function (event) {
        if (event.target === modal) {
            modal.style.display = 'none';
        }
    });

    // Manejar el envío del formulario
    form.addEventListener('submit', function (e) {
        e.preventDefault();
        const id = document.getElementById('employeeId').value;
        const employee = {
            name: document.getElementById('name').value,
            email: document.getElementById('email').value,
            role: document.getElementById('role').value,
        };

        if (id) {
            updateEmployee(id, employee);
        } else {
            createEmployee(employee);
        }
    });

    // Cargar empleados al iniciar
    loadEmployees();

    // Cargar empleados
    function loadEmployees() {
        fetch('/api/employees')
            .then(response => {
                if (!response.ok) {
                    throw new Error('Error al cargar empleados');
                }
                return response.json();
            })
            .then(data => {
                tableBody.innerHTML = '';
                data.forEach(emp => {
                    const row = document.createElement('tr');
                    row.innerHTML = `
                        <td>${emp.id}</td>
                        <td>${emp.name}</td>
                        <td>${emp.email}</td>
                        <td>${emp.role}</td>
                        <td class="actions">
                            <button class="edit" onclick="editEmployee(${emp.id})">Editar</button>
                            <button class="delete" onclick="deleteEmployee(${emp.id})">Eliminar</button>
                        </td>
                    `;
                    tableBody.appendChild(row);
                });
            })
            .catch(error => {
                console.error('Error:', error);
            });
    }

    // Crear empleado
    function createEmployee(employee) {
        fetch('/api/employees', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(employee),
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Error al crear empleado');
                }
                return response.json();
            })
            .then(() => {
                loadEmployees();
                form.reset();
                modal.style.display = 'none';
            })
            .catch(error => {
                console.error('Error:', error);
            });
    }

    // Editar empleado
    window.editEmployee = function (id) {
        fetch(`/api/employees/${id}`)
            .then(response => {
                if (!response.ok) {
                    throw new Error('Error al obtener empleado');
                }
                return response.json();
            })
            .then(data => {
                document.getElementById('modalTitle').textContent = 'Editar Empleado';
                document.getElementById('employeeId').value = data.id;
                document.getElementById('name').value = data.name;
                document.getElementById('email').value = data.email;
                document.getElementById('role').value = data.role;
                modal.style.display = 'block';
            })
            .catch(error => {
                console.error('Error:', error);
            });
    }

    // Actualizar empleado
    function updateEmployee(id, employee) {
        fetch(`/api/employees/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(employee),
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Error al actualizar empleado');
                }
                return response.json();
            })
            .then(() => {
                loadEmployees();
                form.reset();
                modal.style.display = 'none';
            })
            .catch(error => {
                console.error('Error:', error);
            });
    }

    // Eliminar empleado
    window.deleteEmployee = function (id) {
        if (confirm('¿Estás seguro de que deseas eliminar este empleado?')) {
            fetch(`/api/employees/${id}`, { method: 'DELETE' })
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Error al eliminar empleado');
                    }
                    loadEmployees();
                })
                .catch(error => {
                    console.error('Error:', error);
                });
        }
    }
});