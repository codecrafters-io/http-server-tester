import argparse
import enum
import socket
from http import HTTPMethod, HTTPStatus
from pathlib import Path
from threading import Thread

ACCEPTED_BUFFSIZE = 1024
"""Accepted request buffer size by the socket server."""

USER_AGENT_HEADER = "User-Agent"


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
        header_line = next(header for header in self.headers if header_name in header)

        _name, value = header_line.split(":")

        return value.strip()


class Response:
    version = "HTTP/1.1"

    def __init__(
        self,
        status: HTTPStatus = HTTPStatus.OK,
        data: str = "",
        content_type: str = "text/plain",
    ) -> None:
        self.status = f"{status} {status.phrase}"
        self.content_type = content_type

        self.data = data
        self.content_length = len(data)

    def __bytes__(self) -> bytes:
        response = f"{self.version} {self.status}\r\n"

        if self.data:
            response += (
                f"Content-Type: {self.content_type}\r\n"
                f"Content-Length: {self.content_length}\r\n\r\n"
                f"{self.data}"
            )
        else:
            response += "\r\n\r\n"

        return response.encode()


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
                    data=file.read(), content_type="application/octet-stream"
                )

        case HTTPMethod.POST:
            with open(media_path, mode="w") as file:
                file.write(request.body)

                return Response(HTTPStatus.CREATED)

        case _:
            return Response(status=HTTPStatus.METHOD_NOT_ALLOWED)


def handle_connection(connection: socket.socket) -> None:
    with connection:
        buffer = connection.recv(ACCEPTED_BUFFSIZE)

        request = Request(buffer.decode())

        path = request.path.split("/", 2)

        match path[1]:
            case Route.ROOT:
                response = Response()

            case Route.ECHO:
                response = Response(data=path[-1])

            case Route.USER_AGENT:
                response = Response(data=request.get_header(USER_AGENT_HEADER))

            case Route.FILES:
                response = handle_files_route(request, filename=path[-1])

            case _:
                response = NotFoundResponse()

        connection.send(bytes(response))


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
