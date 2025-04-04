import argparse
import enum
import gzip
import socket
from http import HTTPMethod, HTTPStatus
from pathlib import Path
from threading import Thread

ACCEPTED_BUFFSIZE = 1024
"""Accepted request buffer size by the socket server."""

USER_AGENT_HEADER = "User-Agent"
ACCEPT_ENCODING_HEADER = "Accept-Encoding"
SUPPORTED_ENCODINGS = ["gzip"]


class Route(enum.StrEnum):
    """
    Enumeration of available routes.
    """

    ROOT = ""
    USER_AGENT = "user-agent"
    ECHO = "echo"
    FILES = "files"


class Request:
    def __init__(self, buffer: str) -> None:
        buffer_list = buffer.split("\r\n")

        self.body = buffer_list.pop()

        # Remove line break
        buffer_list.remove("")

        firstline, *headers = buffer_list

        self.method, self.path, self.version = firstline.split()
        self.headers = headers

    def get_header(self, header_name: str) -> str:
        try:
            header_line = next(
                header for header in self.headers if header_name in header
            )
            _name, value = header_line.split(":")
            return value.strip()
        except StopIteration:
            # Returns "" if header_name key is not found
            return ""


class Response:
    version = "HTTP/1.1"

    def __init__(
        self,
        status: HTTPStatus = HTTPStatus.OK,
        headers: dict[str, str] = {},
        data: str = "",
        bytes_data: bytes = b"",
    ) -> None:
        self.status = f"{status} {status.phrase}"

        self.headers = headers
        self.data = data
        self.content_length = len(data)
        self.bytes_data = bytes_data

        if self.data != "" and self.bytes_data != b"":
            raise ValueError("Only one of data or bytes_data can be provided")

    def __bytes__(self) -> bytes:
        response = f"{self.version} {self.status}\r\n"
        if self.data:
            self.headers["Content-Length"] = str(self.content_length)
            if self.headers.get("Content-Type") is None:
                self.headers["Content-Type"] = "text/plain"

        if self.headers:
            response += "\r\n".join(
                f"{header}: {value}" for header, value in self.headers.items()
            )
            response += "\r\n"
        response += "\r\n"

        if self.data:
            response += self.data

        encoded_response = response.encode()

        if self.bytes_data:
            encoded_response += self.bytes_data

        return encoded_response


class NotFoundResponse(Response):
    def __init__(self) -> None:
        super().__init__(HTTPStatus.NOT_FOUND)


def handle_files_route(request: Request, filename: str) -> Response:
    if not MEDIA_DIRECTORY:
        raise RuntimeError(
            "/files/ may not be requested without provide --directory argument on the server startup"
        )

    media_path = MEDIA_DIRECTORY / filename

    match request.method:
        case HTTPMethod.GET:
            if not media_path.exists():
                return NotFoundResponse()

            with open(media_path, mode="r") as file:
                return Response(
                    data=file.read(),
                    headers={"Content-Type": "application/octet-stream"},
                )

        case HTTPMethod.POST:
            with open(media_path, mode="w") as file:
                file.write(request.body)

                return Response(HTTPStatus.CREATED)

        case _:
            return Response(status=HTTPStatus.METHOD_NOT_ALLOWED)


def handle_connection(connection: socket.socket) -> None:
    # Set a timeout for the connection to prevent hanging
    connection.settimeout(5)  # 5 second timeout
    
    try:
        while True:
            try:
                buffer = connection.recv(ACCEPTED_BUFFSIZE)
                if not buffer:  # Connection closed by client
                    break
                
                request = Request(buffer.decode())
                path = request.path.split("/", 2)

                # Check for Connection: close header
                should_close = request.get_header("Connection").lower() == "close"

                # Initialize response with default headers
                response = None
                headers: dict[str, str] = {}

                match path[1]:
                    case Route.ROOT:
                        response = Response(headers={})

                    case Route.ECHO:
                        encodings = request.get_header(ACCEPT_ENCODING_HEADER)
                        if encodings != "":
                            encoding_list = encodings.split(",")
                            for encoding in encoding_list:
                                encoding = encoding.strip()
                                if encoding in SUPPORTED_ENCODINGS:
                                    headers["Content-Encoding"] = encoding
                                    break

                        if headers.get("Content-Encoding") is None:
                            response = Response(data=path[-1], headers={})
                        else:
                            body = gzip.compress(path[-1].encode())
                            headers["Content-Length"] = str(len(body))
                            response = Response(bytes_data=body, headers=headers)

                    case Route.USER_AGENT:
                        response = Response(data=request.get_header(USER_AGENT_HEADER), headers={})

                    case Route.FILES:
                        response = handle_files_route(request, filename=path[-1])

                    case _:
                        response = NotFoundResponse()

                # Add Connection: close header if requested
                if should_close:
                    response.headers["Connection"] = "close"

                connection.send(bytes(response))
                
                # Close connection if Connection: close was requested
                if should_close:
                    break
                
            except socket.timeout:
                # Connection timed out, break the loop
                break
            except Exception as e:
                # Handle any other errors and break the loop
                print(f"Error handling request: {e}")
                break
                
    finally:
        connection.close()


def main() -> None:
    print("Server listening on localhost:4221")
    server_socket = socket.create_server(("localhost", 4221), reuse_port=True)

    while True:
        connection, _client_address = server_socket.accept()  # wait for client

        thread = Thread(target=handle_connection, args=(connection,))
        thread.start()


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Run a http server")

    parser.add_argument(
        "--directory",
        metavar="path/to/directory",
        type=Path,
        nargs="?",
        help="The media directory of the server",
    )

    args = parser.parse_args()

    global MEDIA_DIRECTORY
    MEDIA_DIRECTORY: Path | None = args.directory

    main()
