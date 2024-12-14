import socket
import time
import argparse

def test_chat_server(username, password, commands, server_ip="127.0.0.1", server_port=7878):
    try:
        # Подключение к серверу
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
            s.connect((server_ip, server_port))
            print("Connected to the server.")

            # Получение приветственного сообщения
            welcome_message = s.recv(1024).decode("utf-8")
            print("Server:", welcome_message)
            time.sleep(1)

            # Отправка имени пользователя
            print(f"Sending username: {username}")
            s.sendall(f"{username}\n".encode("utf-8"))
            time.sleep(1)
            response = s.recv(1024).decode("utf-8")
            print("Server:", response)

            # Отправка пароля
            print(f"Sending password: {password}")
            s.sendall(f"{password}\n".encode("utf-8"))
            time.sleep(1)
            response = s.recv(1024).decode("utf-8")
            print("Server:", response)

            # Выполнение команд
            for command in commands:
                print(f"Executing command: {command}")
                s.sendall(f"{command}\n".encode("utf-8"))
                time.sleep(1)
                response = s.recv(1024).decode("utf-8")
                print("Server:", response)

            print("Test completed successfully.")
            time.sleep(1)

            # Интерактивный режим
            while True:
                user_input = input("Press Enter to disconnect or enter a command to execute: ").strip()
                if user_input == "":
                    print("Disconnecting from the server.")
                    break
                print(f"Executing command: `{user_input}`")
                s.sendall(f"{user_input}\n".encode("utf-8"))
                time.sleep(1)
                response = s.recv(1024).decode("utf-8")
                print("Server:", response)

    except Exception as e:
        print(f"An error occurred: {e}")

if __name__ == "__main__":
    # Настройка аргументов командной строки
    parser = argparse.ArgumentParser(description="Chat server tester.")
    parser.add_argument("--username", default=f"tester_{int(time.time())}", help="Username to use for testing.")
    parser.add_argument("--password", default="\n\n", help="Password to use for testing.")
    parser.add_argument("--commands", nargs="+", default=["hello"], help="List of commands to execute.")
    parser.add_argument("--server-ip", default="127.0.0.1", help="IP address of the chat server.")
    parser.add_argument("--server-port", type=int, default=7878, help="Port of the chat server.")

    args = parser.parse_args()

    # Запуск теста
    test_chat_server(args.username, args.password, args.commands, args.server_ip, args.server_port)
