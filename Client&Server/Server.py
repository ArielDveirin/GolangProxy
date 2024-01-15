from http.server import SimpleHTTPRequestHandler
from socketserver import TCPServer

class MyHandler(SimpleHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()
        self.wfile.write(b'Response: Hello, this is the server response!\n')

if __name__ == "__main__":
    # Set the server address
    server_address = ('localhost', 12345)

    # Create an HTTP server
    httpd = TCPServer(server_address, MyHandler)

    print(f"Serving on {server_address[0]}:{server_address[1]}...")

    # Keep the server running indefinitely
    httpd.serve_forever()
