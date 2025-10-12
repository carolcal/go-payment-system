export function renderQrCode(data) {
	const qrSection = document.getElementById("qrcodeSection");
	qrSection.hidden = false;
	const payload = data.qr_code_data;
	const qrDiv = document.getElementById("qrcode");
	qrDiv.innerHTML = "";

	new QRCode(qrDiv, {
		text: payload,
		width: 256,
		height: 256,
		correctLevel: QRCode.CorrectLevel.H
	});

	// Append or update a dedicated copy button without replacing existing content (avoids wiping the QR canvas)
	let copyBtn = document.getElementById('copy-qr-btn');
	if (!copyBtn) {
		copyBtn = document.createElement('button');
		copyBtn.type = 'button';
		copyBtn.id = 'copy-qr-btn';
		copyBtn.className = 'cpy-btn';
		copyBtn.textContent = 'Copy QRCode';
		copyBtn.addEventListener('click', () => {
			const text = copyBtn.getAttribute('data-qr') || '';
			if (text) navigator.clipboard.writeText(text);
		});
		qrSection.appendChild(copyBtn);
	}
	copyBtn.setAttribute('data-qr', payload);
}