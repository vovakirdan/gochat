# Чат-сервер на Go

## Обзор
Этот проект реализует многофункциональный чат-сервер с поддержкой нескольких клиентов, разработанный на Go. Сервер обеспечивает удобное взаимодействие пользователей с функциями регистрации, аутентификации, приватных сообщений и общения в комнатах. Включены утилиты для генерации и проверки паролей, а также цветовые сообщения для улучшения читаемости. Сервер разработан с учётом масштабируемости, поддерживает как публичные, так и приватные комнаты, а также обеспечивает плавную работу.

---

## Особенности
- Поддержка нескольких клиентов с конкурентной обработкой с использованием горутин Go.
- Сообщения в комнатах с поддержкой публичных и защищённых паролем комнат.
- Регистрация пользователей с безопасной аутентификацией.
- Проверка силы пароля в реальном времени с предложениями по улучшению.
- Прямые приватные сообщения между пользователями.
- Динамическое создание и переключение комнат.
- Автоматическая генерация паролей, соответствующих высоким требованиям безопасности.
- Цветные системные сообщения для улучшения читаемости.
- Лёгкая встроенная база данных для управления пользователями и комнатами.
- Надёжная обработка ошибок и ведение логов.

---

## Требования
- **Go** (1.18 или новее).
- Терминальный клиент, такой как `telnet`, для подключения к серверу.

---

## Установка
1. Клонируйте репозиторий:
   ```bash
   git clone <repository-url>
   cd <repository-directory>
   ```

2. Сборка приложения с использованием Go:
   ```bash
   go build -o chat-server
   ```

3. Запустите исполняемый файл сервера:
   ```bash
   ./chat-server
   ```

По умолчанию сервер запускается на `127.0.0.1:7878`. Настройте IP и порт с помощью флагов командной строки:
```bash
./chat-server -ip=192.168.1.100 -port=9000
```

---

## Использование

### Запуск сервера
Для запуска сервера с настройками по умолчанию:
```bash
./chat-server
```
Для указания пользовательских IP и порта:
```bash
./chat-server -ip=0.0.0.0 -port=8080
```

### Подключение к серверу
Используйте любой TCP-клиент для подключения. Например, с помощью `telnet`:
```bash
telnet 127.0.0.1 7878
```

После подключения система предложит вам войти или зарегистрироваться.

---

## Команды
### Общие команды
- `/help` - Показать все доступные команды.
- `/help <command>` - Подробная информация о конкретной команде.

### Управление комнатами
- `/list rooms` - Просмотр всех существующих комнат.
- `/create room <room_name> <password>` - Создать новую комнату с опциональным паролем.
- `/switch room <room_name> <password>` - Переключиться в другую комнату. Для приватных комнат требуется пароль.

### Сообщения
- `@<username> <message>` - Отправить личное сообщение определённому пользователю.

### Управление пользователями
- `/list users` - Показать всех зарегистрированных пользователей и их текущий статус (онлайн/оффлайн).
- `/quit` - Выйти из системы и отключиться от сервера.

---

## Утилиты

### Проверка силы пароля
Сервер оценивает силу пароля во время регистрации пользователя на основе следующих критериев:
- Минимальная длина 6 символов.
- Наличие хотя бы одной цифры.
- Использование как строчных, так и прописных букв.
- Наличие хотя бы одного специального символа для повышения безопасности.

Для слабых паролей сервер предоставляет рекомендации по улучшению или предлагает автоматическую генерацию надёжного пароля, соответствующего всем критериям.

### Автоматически генерируемые пароли
Сервер может генерировать случайные пароли, соответствующие всем требованиям безопасности. Сгенерированные пароли сразу отображаются пользователю с рекомендацией сохранить их в надёжном месте.

---

## Структура файлов
- **`main.go`**: Точка входа в сервер, обработка инициализации и сетевых подключений.
- **`server.go`**: Основная логика сервера, управление клиентскими соединениями и командами.
- **`client.go`**: Определение контекста клиента и его атрибутов.
- **`database.go`**: Реализация встроенной базы данных для хранения данных пользователей и комнат.
- **`auth.go`**: Функции для проверки паролей и генерации случайных паролей.
- **`utils.go`**: Вспомогательные функции, включая раскраску текста для вывода в терминале.

---

## Рекомендации по вкладу
Мы приветствуем вклад в развитие проекта. Чтобы внести изменения:
1. Форкните репозиторий.
2. Внесите изменения в новой ветке.
3. Отправьте pull request с подробным описанием изменений.

---

## Лицензия
Этот проект является открытым и лицензирован в соответствии с лицензией MIT. Подробнее см. в файле `LICENSE`.

---

## Контакты
Для вопросов, запросов на новые функции или сообщений об ошибках откройте issue на GitHub или свяжитесь с мейнтейнером проекта через предоставленные контактные данные. Ваш вклад поможет улучшить проект для всех.
