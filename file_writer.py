from http.server import HTTPServer, BaseHTTPRequestHandler
import os
import sys

class FileWriterHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        try:
            # Get the filename from the path
            filename = self.path.lstrip('/')
            if not filename:
                self.send_response(400)
                self.end_headers()
                return

            # Get the content length
            content_length = int(self.headers['Content-Length'])
            
            # Read the request body
            data = self.rfile.read(content_length)
            
            # Write to file
            filepath = os.path.join('/tmp/data/codecrafters.io/http-server-tester', filename)
            os.makedirs(os.path.dirname(filepath), exist_ok=True)
            
            with open(filepath, 'wb') as f:
                f.write(data)
            
            self.send_response(201)
            self.end_headers()
            
        except Exception as e:
            print(f"Error: {e}", file=sys.stderr)
            self.send_response(500)
            self.end_headers()

if __name__ == '__main__':
    server = HTTPServer(('localhost', 4222), FileWriterHandler)
    print("Starting file writer server on port 4222...")
    server.serve_forever() 