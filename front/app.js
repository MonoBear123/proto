BASE_URL = "http://localhost:8080"

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

    try {
        const {data} = await axios.post(BASE_URL + '/register',formData)

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
    axios.get(`${BASE_URL}/search`)

        .then(({ data }) => {
            console.log(data)
            // Проверяем, является ли data объектом с ключами
            if (typeof data === 'object' && !Array.isArray(data)) {
                for (const companyName in data) {
                    const company = data[companyName];
                    const option = document.createElement('option');
                    option.value = company;  // Применяем значение из словаря
                    option.textContent = companyName;  // Используем имя компании
                    companySelect.appendChild(option);
                }
            }
            predictButton.disabled = false;
        })
        .catch(error => {
            console.error('Ошибка при получении списка компаний:', error);
        });


    // Обработчик кнопки предсказания
    predictButton.addEventListener('click', () => {
        if (companySelect.selectedIndex > 0) {  // Проверка, что была выбрана компания
            const selectedCompany = companySelect.options[companySelect.selectedIndex].value;
            console.log(selectedCompany)
            window.location.href = `result.html?sec-id=${encodeURIComponent(selectedCompany)}`;  // Используем encodeURIComponent
        } else {
            alert('Пожалуйста, выберите компанию для предсказания.');
        }
    });
}

// Функция для форматирования даты
function formatDate(dateString) {
    const date = new Date(dateString);
    // Форматирование числа (например, 20 января будет "20")
    let day = date.getDate().toString().padStart(2, '0');
    let month = (date.getMonth() + 1).toString().padStart(2, '0');  // Добавляем 1 к месяцу, так как индексация начинается с 0

    // Форматирование времени
    let hours = date.getHours().toString().padStart(2, '0');
    return `${day}/${month} ${hours}`;
}

// Получение данных предсказания
// Получение данных предсказания
const resultText = document.getElementById('resultText');
if (resultText) {
    const searchParams = new URLSearchParams(window.location.search);
    const secId = searchParams.get('sec-id');
    const token = localStorage.getItem("jwt-token");
    console.log(secId);

    if (secId) {
        const formSec = new FormData();
        formSec.append("secid", "SBER"); //secId
        formSec.append("token", token);

        axios.post(BASE_URL + '/predict', formSec)
            .then(response => {
                const data = response.data;  // Убедитесь, что здесь получаем правильное поле `data`
                console.log(data);
                if (data.result && Array.isArray(data.result)) {  // Проверка, что есть поле `prices`
                    const priceData = data.result;
                    const labels = [];
                    const prices = [];

                    for (let i = 0; i < priceData.length; i++) {
                        const oneWeekAgo = new Date();
                        oneWeekAgo.setDate(oneWeekAgo.getDate() - 6);  // Начало недели назад

                        const date = new Date(oneWeekAgo);  // Текущая дата начинается с "недели назад"
                        date.setTime(date.getTime() + i * 3600000);  // Добавляем 1 час для каждого элемента массива

                        labels.push(formatDate(date));  // Форматируем дату
                        prices.push(priceData[i]);  // Добавляем цену из массива
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
                                borderColor: (context) => {
                                    const index = context.dataIndex;
                                    return index >= prices.length - 24 ? 'red' : '#007bff';
                                },
                                tension: 0.1
                            }]
                        },
                        options: {
                            maintainAspectRatio: false,
                            scales: {
                                y: {
                                    beginAtZero: false,
                                }
                            }
                        }
                    });
                } else {
                    resultText.textContent = 'Нет данных для отображения.';
                }
            })
            .catch(error => {
                console.error('Ошибка при получении прогноза:', error);
                resultText.textContent = 'Не удалось получить прогноз.';
            });
    } else {
        console.warn('Нет sec-id в URL.');
    }
}