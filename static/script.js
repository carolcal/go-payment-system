document.getElementById('payment-form').addEventListener('submit', createPayment);

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

async function getAllPayments() {
	try {
		const response = await fetch('/payments', {
			method: 'GET',
		})

		if (response.ok) {
			const data = await response.json();
			if (Object.keys(data).length != 0) {
				createTable(data)
			} else {
				document.getElementById('tableContainer').innerHTML = "";
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

function createTable(data) {
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
		});
		payButton = `<button type="button" class="pay-btn" data-id="${item.id}">Pagar!</button>`;
		delButton = `<button type="button" class="del-btn" data-id="${item.id}">Deletar!</button>`;
		tableHTML += `<td>${payButton} ${delButton}</td>`;
		tableHTML += '</tr>';
	});
	tableHTML += '</tbody></table>';
	document.getElementById('tableContainer').innerHTML = tableHTML;

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


getAllPayments();
