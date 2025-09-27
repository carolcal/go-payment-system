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
			alert('Payment created successfully!');
		} else {
			const error = await response.json();
			alert(`Error: ${error.error}`);
		}
	} catch (err) {
		console.error('Error:', err);
		alert('An unexpected error occurred.');
	}
}

async function getAllPayments() {
	try {
		const response = await fetch('/payments', {
			method: 'GET',
		})

		if (response.ok) {
			const data = await response.json();
			createTable(data)
		} else {
			const error = await response.json();
			alert(`Error: ${error.error}`);
		}
	} catch (err) {
		console.error('Error:', err);
		alert('An unexpected error occurred.');
	}
}

function createTable(data) {
	let tableHTML = '<table border="1"><thead><tr>';
	const firstItem = Object.values(data)[0];
	Object.keys(firstItem).forEach(key => {
		if (key != "qr_code_data") {
			tableHTML += `<th>${key}</th>`;
		}
	});
	tableHTML += '</tr></thead><tbody>';

	Object.values(data).forEach(item => {
		tableHTML += '<tr>';
		Object.entries(item).forEach(([key, value]) => {
			if (key != "qr_code_data") {
				tableHTML += `<td>${value}</td>`;
			}
		});
		tableHTML += '</tr>';
	});
	tableHTML += '</tbody></table>';
	 document.getElementById('tableContainer').innerHTML = tableHTML;
}

function renderQrCode(data) {
	imgHTML = `<img src=${data.qr_code_data}>`
	document.getElementById('qr-code').innerHTML = imgHTML;
}

getAllPayments();
