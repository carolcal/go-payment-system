document.getElementById('payment-form').addEventListener('submit', createPayment);


async function getAllUsers() {
	try {
		const response = await fetch('/users', {
			method: 'GET',
		})

		if (response.ok) {
			const data = await response.json();
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
	let tableHTML = '<table border="1"><thead><tr>';
	const firstItem = Object.values(data)[0];
	Object.keys(firstItem).forEach(key => {
		tableHTML += `<th>${key}</th>`;
	});
	tableHTML += '<th>actions</th>'
	tableHTML += '</tr></thead><tbody>';

	Object.values(data).forEach(item => {
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

		viewButton = `<button type="button" class="view-btn" data="${item.id}">Visualizar</button>`
		tableHTML += `<td>${viewButton}</td>`;
		tableHTML += '</tr>';
	});
	tableHTML += '</tbody></table>';
	document.getElementById('usersTable').innerHTML = tableHTML;

	document.querySelectorAll('.view-btn').forEach(button => {
		button.addEventListener('click', (event) => {
			console.log("view infos")
			const id = event.target.getAttribute('data');
			document.getElementById('accountInfo').hidden = false;
			getAllPayments(id, "receiver_id");
			getAllPayments(id, "payer_id");
		});
	});
}


async function createPayment(event) {
	event.preventDefault();

	const amount = document.getElementById('amount').value;

	try {
		const response = await fetch('/payment', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ amount: parseFloat(amount) })
		});

		if (response.ok) {
			const data = await response.json();
			renderQrCode(data);
			getAllPayments();
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
	if (user_type == "receiver_id") {
		tableHTML += "<h5>Pagamentos Gerados</h5>"
	} else {
		tableHTML += "<h5>Pagamentos Efetuados</h5>"
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
		cpyButton = `<button type="button" class="cpy-btn" data-qr="${item.qr_code_data}">Copiar Código</button>`
		delButton = `<button type="button" class="del-btn" data-id="${item.id}">Cancelar</button>`;
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

	document.querySelectorAll('.del-btn').forEach(button => {
		button.addEventListener('click', (event) => {
			const id = event.target.getAttribute('data-id');
			delItem(id);
		});
	});
}

function renderQrCode(data) {
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

async function payItem(id) {
	try {
		const response = await fetch(`http://localhost:8080/payment/${id}/pay`, {
			method: 'POST'
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
	getAllPayments()
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
	getAllPayments()
}


getAllUsers()
