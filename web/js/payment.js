import { renderQrCode } from './qrcode.js';
import { createHistoryTable } from './tables.js';
import { getSelectedUser } from './main.js';
import { getAllUsers } from './users.js';

export async function createPayment(event) {
	event.preventDefault();

	const amount = document.getElementById('createPayment').value;
	if (!getSelectedUser()) {
		alert('Selecione um usuário primeiro (Use o botão "User Payment Section").');
		return;
	}
	const user_id = getSelectedUser().id;

	try {
		const response = await fetch('/payment', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: amount 
				? JSON.stringify({ amount: parseFloat(amount), receiver_id: user_id })
				: JSON.stringify({ receiver_id: user_id })
		});

		if (response.ok) {
			const data = await response.json();
			renderQrCode(data);
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

export async function getAllPayments(user_id, user_type) {
	try {
		const response = await fetch(`/payments/${user_id}/${user_type}`, {
			method: 'GET',
		})

		if (response.ok) {
			const data = await response.json();
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

export async function payPayment(event) {
	event.preventDefault();

	let qrcodedata = document.getElementById('qrcodedata');
	let amount = document.getElementById('payAmount');
	if (!qrcodedata) {
		alert('Por favor, insira os dados do QR Code.');
		return;
	}
	try {
		const response = await fetch(`http://localhost:8080/payment/${getSelectedUser().id}/pay`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: amount.value
				? JSON.stringify({ amount: parseFloat(amount.value), qr_code_data: qrcodedata.value })
				: JSON.stringify({ qr_code_data: qrcodedata.value })
		})
		if (response.ok) {
			alert("Pagamento feito com sucesso!")
		} else {
			const error = await response.json();
			if (error.op === "ValidatePaymentAmount" && error.status === 412) {
				amount.hidden = false;
				return ;
			}
			alert(`${error.error}`);
		}
	} catch (err) {
		console.error("Error:", err);
		alert("Ocorreu um erro inesperado");
	}
	if (getSelectedUser()) {
		getAllPayments(getSelectedUser().id, "receiver_id");
		getAllPayments(getSelectedUser().id, "payer_id");
		getAllUsers();
        qrcodedata.value = '';
		amount.value = '';
		amount.hidden = true;

	}
}

export async function delPayment(id) {
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
	if (getSelectedUser()) {
		getAllPayments(getSelectedUser().id, "receiver_id");
		getAllPayments(getSelectedUser().id, "payer_id");
	}
}