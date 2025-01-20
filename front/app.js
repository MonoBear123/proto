BASE_URL = "http://client:8080"

// Навигация по кнопкам на главной странице
document.getElementById('btnLogin')?.addEventListener('click', () => {
    window.location.href = 'login.html';
});

document.getElementById('btnRegister')?.addEventListener('click', () => {
    window.location.href = 'register.html';
});

// Обработка формы входа
document.getElementById('loginForm')?.addEventListener('submit', async (e) => {
    e.preventDefault();

    const formData = new FormData(e.target);
    try {
        const {data} = await axios.post(BASE_URL + '/login',formData)
        localStorage.setItem('jwt-token', data.token); // Сохраняем токен в локальном хранилище
        window.location.href = 'authorized.html'; // Переходим на авторизованную страницу
    } catch (error) {
        alert(error.message);
    }
});

// Обработка формы регистрации
document.getElementById('registerForm')?.addEventListener('submit', async (e) => {
    e.preventDefault();

    const formData = new FormData(e.target);
    const email = formData.get('email');
    const password = formData.get('password');

    try {
        const response = await fetch(BASE_URL + '/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
                'Access-Control-Allow-Origin': '*'
            },
            body: { email, password },
            credentials: 'include'
        });

        if (!response.ok) {
            throw new Error(`Ошибка: ${response.status}`);
        }

        alert('Пожалуйста, проверьте вашу электронную почту и подтвердите регистрацию.');
        window.location.href = 'index.html'; // Возвращаемся на главную страницу
    } catch (error) {
        alert(error.message);
    }
});

// Получение списка компаний
const companySelect = document.getElementById('companySelect');
const predictButton = document.getElementById('predictButton');

if (companySelect && predictButton) {
    fetch(BASE_URL + '/search', {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('jwt-token')}`
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`Ошибка: ${response.status}`);
        }
        return response.json();
    })
    .then(data => {
        data.forEach(company => {
            const option = document.createElement('option');
            option.value = company['sec-id'];
            option.textContent = company.name;
            companySelect.appendChild(option);
        });
        predictButton.disabled = false;
    })
    .catch(error => {
        console.error('Ошибка при получении списка компаний:', error);
    });

    // Обработчик кнопки предсказания
    predictButton.addEventListener('click', () => {
        const selectedCompany = companySelect.options[companySelect.selectedIndex].value;
        window.location.href = `result.html?sec-id=${selectedCompany}`;
    });
}

// Функция для форматирования даты
function formatDate(dateString) {
    const date = new Date(dateString);
    let hours = date.getHours().toString().padStart(2, '0');
    let minutes = date.getMinutes().toString().padStart(2, '0');
    return `${hours}:${minutes}`;
}

// Получение данных предсказания
const resultText = document.getElementById('resultText');
if (resultText) {
    const searchParams = new URLSearchParams(window.location.search);
    const secId = searchParams.get('sec-id');

    fetch(`/predict?sec-id=${secId}`, {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('jwt-token')}`
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`Ошибка: ${response.status}`);
        }
        return response.json();
    })
    .then(data => {
        const priceData = data.prices; // Массив цен акций
        const labels = [];
        const prices = [];

        for (let i = 0; i < priceData.length; i++) {
            const date = new Date(); // Текущая дата
            date.setTime(date.getTime() + i * 900000); // Добавляем 15 минут для каждого элемента массива
            labels.push(formatDate(date));
            prices.push(priceData[i]);
        }

        // Построение графика
        const ctx = document.getElementById('priceChart').getContext('2d');
        new Chart(ctx, {
            type: 'line',
            data: {
                labels,
                datasets: [{
                    label: 'Цена акции',
                    data: prices,
                    fill: false,
                    borderColor: '#007bff',
                    tension: 0.1
                }]
            },
            options: {
                scales: {
                    y: {
                        beginAtZero: true
                    }
                }
            }
        });
    })
    .catch(error => {
        console.error('Ошибка при получении прогноза:', error);
        resultText.textContent = 'Не удалось получить прогноз.';
    });
}