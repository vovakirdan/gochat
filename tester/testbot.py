import socket
import time
import argparse

def test_chat_server():
    server_ip = "127.0.0.1"
    server_port = 7878

    username = "test_user"
    password = "Test123!"
    message = "hello"

    try:
        # Подключение к серверу
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
            s.connect((server_ip, server_port))
            print("Connected to the server.")

            # Получение приветственного сообщения
            welcome_message = s.recv(1024).decode("utf-8")
            print("Server:", welcome_message)

            # Отправка имени пользователя
            s.sendall(f"{username}\n".encode("utf-8"))
            time.sleep(0.5)
            response = s.recv(1024).decode("utf-8")
            print("Server:", response)

            # Отправка пароля
            s.sendall(f"{password}\n".encode("utf-8"))
            time.sleep(0.5)
            response = s.recv(1024).decode("utf-8")
            print("Server:", response)

            # Отправка сообщения
            s.sendall(f"{message}\n".encode("utf-8"))
            time.sleep(0.5)
            response = s.recv(1024).decode("utf-8")
            print("Server:", response)

            print("Test completed successfully.")

    except Exception as e:
        print(f"An error occurred: {e}")

if __name__ == "__main__":
    test_chat_server()
