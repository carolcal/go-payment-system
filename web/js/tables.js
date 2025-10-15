import { getUserById, delUserById } from './users.js';
import { getAllPayments, payPayment, delPayment } from './payment.js';
import { getSelectedUser, setSelectedUser, usersIndex } from './main.js';

export function createUsersTable(data) {
	usersIndex.clear();

	let tableHTML = '<table border="1"><thead><tr>';
	const firstItem = Object.values(data)[0];
	Object.keys(firstItem).forEach(key => {
		tableHTML += `<th>${key}</th>`;
	});
	tableHTML += '<th>actions</th>'
	tableHTML += '</tr></thead><tbody>';

	Object.values(data).forEach(item => {
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

		const getIdButton = `<button type="button" class="get-btn getid-btn" data="${item.id}">GET /user/:id</button>`
		const updateBalanceButton = `<button type="button" class="put-btn update-btn" data="${item.id}">PUT /user/:id/balance</button>`
		const delIdButton = `<button type="button" class="del-btn delid-btn" data="${item.id}">DELETE /user/:id</button>`
		const viewButton = `<button type="button" class="view-btn" data="${item.id}">Open User Payment Section</button>`
		tableHTML += `<td>${getIdButton} ${updateBalanceButton} ${delIdButton} ${viewButton}</td>`;
		tableHTML += '</tr>';
	});

	tableHTML += '</tbody></table>';
	document.getElementById('usersTable').innerHTML = tableHTML;

	document.querySelectorAll('.getid-btn').forEach(button => {
		button.addEventListener('click', (event) => {
			const id = event.target.getAttribute('data');
			document.getElementById('getUserResult').hidden = false;
			document.getElementById('accountInfo').hidden = true;
			getUserById(id);
		});
	});

	document.querySelectorAll('.update-btn').forEach(button => {
		button.addEventListener('click', (event) => {
			const id = event.target.getAttribute('data');

			setSelectedUser(usersIndex.get(String(id)));
			document.getElementById('selectedBalanceUserName').textContent = getSelectedUser() ? getSelectedUser().name : 'N/A';
			document.getElementById('updateBalanceUser').hidden = false;
			document.getElementById('accountInfo').hidden = true;
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

			setSelectedUser(usersIndex.get(String(id)));
			document.getElementById('selectedPaymentUserName').textContent = getSelectedUser() ? getSelectedUser().name : 'N/A';
			document.getElementById('createPayment').value = '';
			document.getElementById('qrcodedata').value = '';
			document.getElementById('qrcodeSection').hidden = true;
			document.getElementById('accountInfo').hidden = false;
			document.getElementById('getUserResult').hidden = true;
			document.getElementById('updateBalanceUser').hidden = true;
			getAllPayments(id, "receiver_id");
			getAllPayments(id, "payer_id");
		});
	});
}

export function createHistoryTable(data, user_type) {
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
		const cpyButton = `<button type="button" class="cpy-btn" data-qr="${item.qr_code_data}">Copy QRCode</button>`
		const delButton = `<button type="button" class="del-btn del-${user_type}-btn" data-id="${item.id}">DELETE /payment/:id</button>`;
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
			payPayment(id);
		});
	});

	document.querySelectorAll(`.del-${user_type}-btn`).forEach(button => {
		button.addEventListener('click', (event) => {
			const id = event.target.getAttribute('data-id');
			delPayment(id);
		});
	});
}