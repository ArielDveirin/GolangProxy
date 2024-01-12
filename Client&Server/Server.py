import socket

# Create a socket object
server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

# Define the host and port
host = "127.0.0.1"  # localhost
port = 12345        # arbitrary port number

# Bind the socket to the address
server_socket.bind((host, port))

# Listen for incoming connections
server_socket.listen()

print("Server is listening...")

# Accept incoming connections
while True:
    client_socket, client_address = server_socket.accept()
    print(f"Connection from {client_address} has been established.")

    # Receive data from the client
    data = client_socket.recv(1024)
    print("Received:", data.decode())

    # Send a response back to the client
    client_socket.sendall(b"Hello from the server!")

# Close the connection
client_socket.close()
server_socket.close()
