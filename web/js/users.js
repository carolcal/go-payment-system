import { getSelectedUser, clearSelectedUser } from './main.js';
import { createUsersTable } from './tables.js';

export async function createNewUser(event) {
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
			alert(`Error: ${error.error}`);
		}
	} catch (err) {
		console.error("Error:", err);
		alert("An unexpected error occurred");
	}
}

export async function getAllUsers() {
	try {
		const response = await fetch('/users', {
			method: 'GET',
		})

		if (response.ok) {
			const data = await response.json();
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



export async function getUserById(id) {
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

export async function updateBalance(event) {
	event.preventDefault();

	const user_id = getSelectedUser().id;
	let diff = document.getElementById('updateBalanceAmount').value;
	const depositOrWithdraw = document.getElementById('updateBalanceType').value;

	if (depositOrWithdraw === "withdraw") {
		diff = -diff;
	}

	try {
		const response = await fetch(`/user/${user_id}/balance`, {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ diff: parseFloat(diff) })
		});

		if (response.ok) {
			alert("Balance updated successfully!");
			getAllUsers();
			document.getElementById('updateBalanceAmount').value = '';
			document.getElementById('updateBalanceUser').hidden = true;
			document.getElementById('accountInfo').hidden = true;
			document.getElementById('getUserResult').hidden = true;
			clearSelectedUser();
		} else {
			const error = await response.json();
			alert(`Erro: ${error.error}`);
		}
	} catch (err) {
		console.error("Error:", err);
		alert("Ocorreu um erro inesperado");
	}
}

export async function delUserById(id) {
	try {
		const response = await fetch(`/user/${id}`, {
			method: 'DELETE',
		})

		if (response.ok) {
			alert("User deleted successfully!");
			getAllUsers();
			document.getElementById('updateBalanceAmount').value = '';
			document.getElementById('updateBalanceUser').hidden = true;
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