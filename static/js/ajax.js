document.addEventListener('DOMContentLoaded', function () {
    const forms = document.querySelectorAll('form');

    forms.forEach(form => {
        form.addEventListener('submit', function (e) {
            e.preventDefault();

            const formData = new FormData(form);
            const formId = form.querySelector('input[name="form_id"]').value;

            fetch(form.action, {
                method: form.method,
                headers: {
                    'X-Requested-With': 'XMLHttpRequest',
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: new URLSearchParams(formData),
            })
                .then(response => response.json())
                .then(data => {
                    // Очистка предыдущих ошибок и сообщений
                    const errorSpans = form.querySelectorAll('.error');
                    errorSpans.forEach(span => span.remove());

                    const successMessages = form.querySelectorAll('.success');
                    successMessages.forEach(span => span.remove());

                    if (data.message) {
                        // Успешная регистрация
                        const successMessage = document.createElement('div');
                        successMessage.className = 'success';
                        successMessage.style.color = 'green';
                        successMessage.textContent = data.message;
                        form.appendChild(successMessage);

                        // Очистка полей формы
                        form.querySelectorAll('input[type="text"], input[type="password"]').forEach(input => {
                            input.value = '';
                        });
                    } else if (data.errors) {
                        // Отображение новых ошибок
                        Object.keys(data.errors).forEach(field => {
                            const input = form.querySelector(`[name="${formId}_${field}"]`);
                            if (input) {
                                const errorSpan = document.createElement('span');
                                errorSpan.className = 'error';
                                errorSpan.style.color = 'red';
                                errorSpan.textContent = data.errors[field];
                                input.parentNode.appendChild(errorSpan);
                            }
                        });
                    }
                })
                .catch(error => {
                    console.error('Error:', error);
                });
        });
    });
});