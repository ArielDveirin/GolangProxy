import json
import urllib.request

def make_request():
    url = 'http://localhost:12345'
     # Create a dictionary with the string you want to send
    data = 'Request: Hello, server! This is a string from the client.'

    # Encode the dictionary as JSON
    data_json = json.dumps(data).encode('utf-8')

    # Set the Content-Type header to indicate JSON data
    headers = {'Content-Type': 'application/json'}

    # Make a POST request to the server with the JSON payload
    request = urllib.request.Request(url, data=data_json, headers=headers, method='GET')

    with urllib.request.urlopen(request) as response:
        html = response.read()
        print(html.decode())

make_request()
