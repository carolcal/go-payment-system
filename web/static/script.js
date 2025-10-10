document.getElementById('new-user-form').addEventListener('submit', createNewUser);
document.getElementById('payment-form-create').addEventListener('submit', createPayment);
document.getElementById('payment-form-pay').addEventListener('submit', payPayment);

// In-memory state to keep table row data and current selection accessible across handlers
let usersIndex = new Map(); // id -> user object
let selectedUser = null;    // currently selected user (row) for the Payment section

async function createNewUser(event) {
	event.preventDefault();

	const name = document.getElementById('new-user-name').value;
	const cpf = document.getElementById('new-user-cpf').value;
	const balance = document.getElementById('new-user-balance').value;
	const city = document.getElementById('new-user-city').value;

	try {
		const response = await fetch('/user', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ name, cpf, balance: parseFloat(balance), city })
		});

		if (response.ok) {
			alert('User created successfully!');
			getAllUsers();
		} else {
			const error = await response.json();
			alert(`Error: ${error.message}`);
		}
	} catch (err) {
		console.error("Error:", err);
		alert("An unexpected error occurred");
	}
}

async function getAllUsers() {
	try {
		const response = await fetch('/users', {
			method: 'GET',
		})

		if (response.ok) {
			const data = await response.json();
			selectedUser = null;
			document.getElementById('getUsersResult').hidden = false;
			if (Object.keys(data).length != 0) {
				createUsersTable(data)
			} else {
				document.getElementById('usersTable').innerHTML = "";
			}
		} else {
			const error = await response.json();
			alert(`Erro: ${error.error}`);
		}
	} catch (err) {
		console.error("Error:", err);
		alert("Ocorreu um erro inesperado");
	}
}

function createUsersTable(data) {
	// Rebuild the in-memory index for quick lookup by id
	usersIndex.clear();

	let tableHTML = '<table border="1"><thead><tr>';
	const firstItem = Object.values(data)[0];
	Object.keys(firstItem).forEach(key => {
		tableHTML += `<th>${key}</th>`;
	});
	tableHTML += '<th>actions</th>'
	tableHTML += '</tr></thead><tbody>';

	Object.values(data).forEach(item => {
		// Save the full row data for later usage (e.g., when clicking action buttons)
		usersIndex.set(String(item.id), item);

		tableHTML += '<tr>';
		Object.entries(item).forEach(([key, value]) => {
			if (key == "created_at" || key == "updated_at") {
				const date = new Date(value);
				value = date.toLocaleDateString('en-GB');
			} else if (key == "balance") {
				const formatter = new Intl.NumberFormat('pt-BR', {
					style: 'currency',
					currency: 'BRD',
					minimumFractionDigits: 2,
				});
				value = formatter.format(value / 100);
			}
			tableHTML += `<td>${value}</td>`;
		});

		getIdButton = `<button type="button" class="get-btn getid-btn" data="${item.id}">GET /user/:id</button>`
		delIdButton = `<button type="button" class="del-btn delid-btn" data="${item.id}">DELETE /user/:id</button>`
		viewButton = `<button type="button" class="view-btn" data="${item.id}">Open User Payment Section</button>`
		tableHTML += `<td>${getIdButton} ${delIdButton} ${viewButton}</td>`;
		tableHTML += '</tr>';
	});
	tableHTML += '</tbody></table>';
	document.getElementById('usersTable').innerHTML = tableHTML;

	document.querySelectorAll('.getid-btn').forEach(button => {
		button.addEventListener('click', (event) => {
			const id = event.target.getAttribute('data');
			document.getElementById('getUserResult').hidden = false;
			getUserById(id);
		});
	});

	document.querySelectorAll('.delid-btn').forEach(button => {
		button.addEventListener('click', (event) => {
			const id = event.target.getAttribute('data');
			delUserById(id);
		});
	});

	document.querySelectorAll('.view-btn').forEach(button => {
		button.addEventListener('click', (event) => {
			const id = event.target.getAttribute('data');

			// Set the selected user so other actions (e.g., creating payments) can use its data
			selectedUser = usersIndex.get(String(id)) || null;
			if (selectedUser) {
				// Persist selection for this session so a refresh keeps the context
				try { sessionStorage.setItem('selectedUser', JSON.stringify(selectedUser)); } catch { /* ignore */ }
			}

			document.getElementById('accountInfo').hidden = false;
			getAllPayments(id, "receiver_id");
			getAllPayments(id, "payer_id");
		});
	});
}

async function getUserById(id) {
	try {
		const response = await fetch(`/user/${id}`, {
			method: 'GET',
		})

		if (response.ok) {
			const data = await response.json();
			let userHTML = '<ul>';
			Object.entries(data).forEach(([key, value]) => {
				if (key == "created_at" || key == "updated_at") {
					userHTML += `<li><strong>${key.toUpperCase()}</strong>: ${new Date(value).toLocaleDateString('en-GB')}</li>`;
				} else if (key == "balance") {
					const formatter = new Intl.NumberFormat('pt-BR', {
						style: 'currency',
						currency: 'BRD',
						minimumFractionDigits: 2,
					});
					userHTML += `<li><strong>${key.toUpperCase()}</strong>: ${formatter.format(value / 100)}</li>`;
				} else {
					userHTML += `<li><strong>${key.toUpperCase()}</strong>: ${value}</li>`;
				}
			});
			userHTML += '</ul>';
			document.getElementById('userByIdInfo').innerHTML = userHTML;
		} else {
			const error = await response.json();
			alert(`Erro: ${error.error}`);
		}
	} catch (err) {
		console.error("Error:", err);
		alert("Ocorreu um erro inesperado");
	}
}

async function delUserById(id) {
	try {
		const response = await fetch(`/user/${id}`, {
			method: 'DELETE',
		})

		if (response.ok) {
			alert("User deleted successfully!");
			getAllUsers();
			document.getElementById('accountInfo').hidden = true;
			document.getElementById('getUserResult').hidden = true;
		} else {
			const error = await response.json();
			alert(`Erro: ${error.error}`);
		}
	} catch (err) {
		console.error("Error:", err);
		alert("Ocorreu um erro inesperado");
	}
}

async function createPayment(event) {
	event.preventDefault();

	const amount = document.getElementById('amount').value;
	// Use the selected user as the receiver of the generated payment/QR
	if (!selectedUser) {
		alert('Selecione um usuário primeiro (Use o botão "User Payment Section").');
		return;
	}
	const user_id = selectedUser.id;

	try {
		const response = await fetch('/payment', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ amount: parseFloat(amount), receiver_id: user_id })
		});

		if (response.ok) {
			const data = await response.json();
			renderQrCode(data);
			// Refresh the histories for the selected user
			getAllPayments(user_id, "receiver_id");
			getAllPayments(user_id, "payer_id");
			alert('Pagamento criado com sucesso!');
		} else {
			const error = await response.json();
			alert(`Erro: ${error.error}`);
		}
	} catch (err) {
		console.error("Error:", err);
		alert("Ocorreu um erro inesperado");
	}
}

async function getAllPayments(user_id, user_type) {
	console.log("get payments " + user_type)
	try {
		const response = await fetch(`/payments/${user_id}/${user_type}`, {
			method: 'GET',
		})

		if (response.ok) {
			const data = await response.json();
			console.log(data)
			if (Object.keys(data).length != 0) {
				createHistoryTable(data, user_type)
			} else {
				if (user_type == "receiver_id") {
					document.getElementById('accountReceiveHistory').innerHTML = "";
				} else {
					document.getElementById('accountPayHistory').innerHTML = "";
				}
			}
		} else {
			const error = await response.json();
			alert(`Erro: ${error.error}`);
		}
	} catch (err) {
		console.error("Error:", err);
		alert("Ocorreu um erro inesperado");
	}
}

function createHistoryTable(data, user_type) {
	let tableHTML = "";
	tableHTML += "<h4>GET /payments/:user_id/:user_type</h4>"
	if (user_type == "receiver_id") {
		tableHTML += "<h5>Pagamentos Gerados (user_type: receiver_id)</h5>"
	} else {
		tableHTML += "<h5>Pagamentos Efetuados (user_type: payer_id)</h5>"
	}
	tableHTML += '<table border="1"><thead><tr>';
	const firstItem = Object.values(data)[0];
	Object.keys(firstItem).forEach(key => {
		if (key != 'qr_code_data' && key != user_type) {
			tableHTML += `<th>${key}</th>`;
		}
	});
	tableHTML += '<th>actions</th>'
	tableHTML += '</tr></thead><tbody>';

	Object.values(data).forEach(item => {
		tableHTML += '<tr>';
		Object.entries(item).forEach(([key, value]) => {
			if (key != 'qr_code_data' && key != user_type) {
				if (key == "created_at" || key == "expires_at") {
					const date = new Date(value);
					value = date.toLocaleString();
				} else if (key == "amount") {
					const formatter = new Intl.NumberFormat('pt-BR', {
						style: 'currency',
						currency: 'BRD',
						minimumFractionDigits: 2,
					});
					value = formatter.format(value / 100);
				}
				tableHTML += `<td>${value}</td>`;
			}
		});
		cpyButton = `<button type="button" class="cpy-btn" data-qr="${item.qr_code_data}">Copy QRCode</button>`
		delButton = `<button type="button" class="del-btn delpay-btn" data-id="${item.id}">DELETE /payment/:id</button>`;
		tableHTML += `<td>${cpyButton} ${delButton}</td>`;
		tableHTML += '</tr>';
	});
	tableHTML += '</tbody></table>';
	if (user_type == "receiver_id") {
		document.getElementById('accountReceiveHistory').innerHTML = tableHTML;
	} else {
		document.getElementById('accountPayHistory').innerHTML = tableHTML;
	}

	document.querySelectorAll('.cpy-btn').forEach(button => {
		button.addEventListener('click', (event) => {
			const item = event.target.getAttribute('data-qr');
			navigator.clipboard.writeText(item)
		});
	});

	document.querySelectorAll('.pay-btn').forEach(button => {
		button.addEventListener('click', (event) => {
			const id = event.target.getAttribute('data-id');
			payItem(id);
		});
	});

	document.querySelectorAll('.delpay-btn').forEach(button => {
		button.addEventListener('click', (event) => {
			const id = event.target.getAttribute('data-id');
			delItem(id);
		});
	});
}

function renderQrCode(data) {
	document.getElementById("qrcode").hidden = false;
	const payload = data.qr_code_data;
	const qrDiv = document.getElementById("qrcode");
	qrDiv.innerHTML = "";

	new QRCode(qrDiv, {
		text: payload,
		width: 256,
		height: 256,
		correctLevel: QRCode.CorrectLevel.H
	});
}

async function payPayment(event) {
	event.preventDefault();

	const qrcodedata = document.getElementById('qrcodedata').value;
	if (!qrcodedata) {
		alert('Por favor, insira os dados do QR Code.');
		return;
	}
	try {
		const response = await fetch(`http://localhost:8080/payment/${selectedUser.id}/pay`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ qrcodedata })
		})
		if (response.ok) {
			alert("Pagamento feito com sucesso!")
		} else {
			const error = await response.json();
			alert(`${error.error}`);
		}
	} catch (err) {
		console.error("Error:", err);
		alert("Ocorreu um erro inesperado");
	}
	if (selectedUser) {
		getAllPayments(selectedUser.id, "receiver_id");
		getAllPayments(selectedUser.id, "payer_id");
	}
}

async function delItem(id) {
	try {
		const response = await fetch(`http://localhost:8080/payment/${id}`, {
			method: 'DELETE'
		})
		if (response.ok) {
			alert("Pagamento removido com sucesso!")
		} else {
			const error = await response.json();
			alert(`${error.error}`);
		}
	} catch (err) {
		console.error("Error:", err);
		alert("Ocorreu um erro inesperado");
	}
	if (selectedUser) {
		getAllPayments(selectedUser.id, "receiver_id");
		getAllPayments(selectedUser.id, "payer_id");
	}
}


getAllUsers()
