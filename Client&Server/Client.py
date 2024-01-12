import socket

# Create a socket object
client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

# Define the host and port to connect
host = "127.0.0.1"  # localhost
port = 12345        # same port as server

# Connect to the server
client_socket.connect((host, port))

# Send data to the server
message = "Hello from the client!"
client_socket.sendall(message.encode())

# Receive a response from the server
data = client_socket.recv(1024)
print("Received:", data.decode())

# Close the connection
client_socket.close()
