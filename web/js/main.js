import { createNewUser, updateBalance } from './users.js';
import { createPayment, payPayment, getAllPayments } from './payment.js';
import { getAllUsers } from './users.js';

document.getElementById('new-user-form').addEventListener('submit', createNewUser);
document.getElementById('payment-form-create').addEventListener('submit', createPayment);
document.getElementById('payment-form-pay').addEventListener('submit', payPayment);
document.getElementById('update-balance-form').addEventListener('submit', updateBalance);

export let usersIndex = new Map();

let _selectedUser = null;

export function getSelectedUser() {
  return _selectedUser;
}

export function setSelectedUser(user) {
  _selectedUser = user || null;
  try {
    if (user) sessionStorage.setItem('selectedUser', JSON.stringify(user));
    else sessionStorage.removeItem('selectedUser');
  } catch { /* ignore */ }
}

export function clearSelectedUser() {
  setSelectedUser(null);
}

getAllUsers();